package cfsm

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Parse a single CFSM in fsa format.
func ParseSingleCFSMFSA(reader io.Reader) (*CFSM, error) {
	sys, err := parseFSA(reader, true)
	if err != nil {
		return nil, err
	}
	return sys.CFSMs[0], nil
}

// Parse a system of CFSMs in fsa format.
func ParseSystemCFSMsFSA(reader io.Reader) (*System, error) {
	return parseFSA(reader, false)
}

func parseFSA(reader io.Reader, singleCFSM bool) (*System, error) {
	// Regular expressions for parser
	lineCommentRe := regexp.MustCompile(`^--.*`)
	machineStartRe := regexp.MustCompile(`^\.outputs(\s+(\S+))?$`)
	startStateRe := regexp.MustCompile(`^\.marking\s+(\S+)$`)
	messageRe := regexp.MustCompile(`^(\S+)\s(\S+)\s([\?!])\s(\S+)\s(\S+)$`)

	// Possible states of the parser while consuming input.
	type FSAParserStatus byte
	const (
		Initial         FSAParserStatus = iota
		MachineStarting                 // we've consumed ".outputs"
		Transitions                     // we've consumed ".state graph". We can loop here reading transitions.
		StartStateRead                  // we've consumed ".marking X" and can still read more transitions.
	)
	type TransitionType byte
	const (
		SendType TransitionType = iota
		RecvType
	)
	type Transition struct {
		FromState    *State
		OtherMachine string
		Action       TransitionType
		Message      string
		NextState    *State
	}

	var namesOfCFSMs []string            // array of names of CFSMs as we find in order. If they don't have name, we use str(int) of the machine order (starts with 0).
	numberOfCFSM := make(map[string]int) // inverse of the latter
	var status FSAParserStatus = Initial
	sys := NewSystem() // This is what we'll return.

	// While reading transitions we can find transitions that refer to CFSMs that we haven't yet parsed.
	// So we'll save all transitions we find in this map, and add them all after consuming the entire fsa file.
	transitionsToBackfill := make(map[*CFSM][]Transition)

	var stateNames map[string]*State
	var currentMachine *CFSM

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		currentLine := scanner.Text()
		currentLine = strings.TrimSpace(currentLine)
		if lineCommentRe.MatchString(currentLine) || len(currentLine) == 0 {
			// Skip comments or empty lines.
			continue
		}
		switch status {
		case Initial:
			// We only expect the start marker.
			matches := machineStartRe.FindStringSubmatch(currentLine)
			if matches == nil {
				return nil, errors.New("fsa file invalid. Expected .outputs")
			}
			currentMachineNumber := len(namesOfCFSMs)
			machineName := matches[2]
			if machineName != "" {
				num, err := strconv.Atoi(machineName)
				if err == nil && num != currentMachineNumber {
					return nil, fmt.Errorf("fsa file invalid: machine number %d is named %d", currentMachineNumber, num)
				}
				namesOfCFSMs = append(namesOfCFSMs, machineName)
				numberOfCFSM[machineName] = currentMachineNumber
			} else {
				machineName = strconv.Itoa(currentMachineNumber)
				namesOfCFSMs = append(namesOfCFSMs, strconv.Itoa(len(namesOfCFSMs)))
				numberOfCFSM[machineName] = currentMachineNumber
			}
			var err error
			currentMachine, err = sys.NewNamedMachine(machineName)
			if err != nil {
				return nil, fmt.Errorf("failed to create machine %s", machineName)
			}
			// currentMachine.Comment = namesOfCFSMs[currentMachineNumber]
			stateNames = make(map[string]*State)
			transitionsToBackfill[currentMachine] = make([]Transition, 0)
			status = MachineStarting
		case MachineStarting:
			// We can only expect now a (useless) ".state graph" line.
			if currentLine == ".state graph" {
				status = Transitions
				continue
			} else {
				return nil, errors.New("fsa file invalid: expected '.state graph' after CFSM declaration")
			}
		case Transitions:
			// Here we can see transition lines or the start state marker.
			startStateMatch := startStateRe.FindStringSubmatch(currentLine)
			if len(startStateMatch) == 2 {
				startStateName := startStateMatch[1]
				val, ok := stateNames[startStateName]
				if !ok {
					val := currentMachine.NewState()
					val.Label = startStateName
					stateNames[startStateName] = val
				}
				currentMachine.Start = val
				status = StartStateRead
				continue
			}
			transitionMatches := messageRe.FindStringSubmatch(currentLine)
			if transitionMatches == nil {
				return nil, errors.New("expected transition")
			}

			// Create fromState and nextState if they don't exist already.
			val, ok := stateNames[transitionMatches[1]]
			if !ok {
				val = currentMachine.NewState()
				val.Label = transitionMatches[1]
				stateNames[transitionMatches[1]] = val
			}
			fromState := val

			val, ok = stateNames[transitionMatches[5]]
			if !ok {
				val = currentMachine.NewState()
				val.Label = transitionMatches[5]
				stateNames[transitionMatches[5]] = val
			}
			nextState := val

			// Parse action
			var action TransitionType
			switch a := transitionMatches[3]; a {
			case "?":
				action = RecvType
			case "!":
				action = SendType
			default:
				return nil, errors.New("fsa file invalid: invalid action")
			}

			msg := transitionMatches[4]
			otherCFSM := transitionMatches[2]

			// Add transition to be backfilled.
			transitionsToBackfill[currentMachine] = append(transitionsToBackfill[currentMachine], Transition{
				FromState:    fromState,
				OtherMachine: otherCFSM,
				Action:       action,
				Message:      msg,
				NextState:    nextState,
			})

		case StartStateRead:
			// Here we can only see the end marker.
			// TODO: We could also see transitions? As-is, the line ".marking <START_STATE>" has to come after
			//   all transitions and must be followed by an ".end" line.
			if currentLine == ".end" {
				status = Initial
				if singleCFSM {
					break
				}
				continue
			} else {
				return nil, errors.New("fsa file invalid. Expected '.end'")
			}
		}

	}
	if status != Initial {
		return nil, errors.New("fsa file invalid: unexpected EOF")
	}

	for _, transitions := range transitionsToBackfill {
		for _, t := range transitions {
			var otherMachine *CFSM
			if !singleCFSM {
				// t.OtherMachine can be a name or a string representing a number.
				otherMachineNum, err := strconv.Atoi(t.OtherMachine)
				if err != nil {
					var ok bool
					otherMachineNum, ok = numberOfCFSM[t.OtherMachine]
					if !ok {
						return nil, errors.New("non existant CFSM referenced")
					}
				}
				otherMachine = sys.CFSMs[otherMachineNum]
			}

			switch t.Action {
			case SendType:
				var e *Send
				if !singleCFSM {
					e = NewSend(otherMachine, t.Message)
				} else {
					e = NewSendToName(t.OtherMachine, t.Message)
				}
				e.SetNext(t.NextState)
				t.FromState.AddTransition(e)
			case RecvType:
				var e *Recv
				if !singleCFSM {
					e = NewRecv(otherMachine, t.Message)
				} else {
					e = NewRecvFromName(t.OtherMachine, t.Message)
				}
				e.SetNext(t.NextState)
				t.FromState.AddTransition(e)
			}

		}
	}

	return sys, nil
}
