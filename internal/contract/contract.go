package contract

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/nickng/cfsm"

	pb "github.com/clpombo/search/gen/go/search/v1"
)

type Contract interface {
	GetParticipants() []string
	GetRemoteParticipantNames() []string
	GetLocalParticipantName() string // This returns the name of the local participant of this contract.
	// TODO: add GetNextState()
}

type CFSMContract struct {
	*cfsm.System
	localParticipant *cfsm.CFSM
}

func (s *CFSMContract) GetParticipants() []string {
	return s.getParticipants(true)
}

func (c *CFSMContract) GetLocalParticipantName() string {
	return strings.Clone(c.localParticipant.Comment)
}

func (c *CFSMContract) GetRemoteParticipantNames() (ret []string) {
	return c.getParticipants(false)
}

func (c *CFSMContract) getParticipants(includeLocal bool) (ret []string) {
	for _, m := range c.CFSMs {
		if !includeLocal && m == c.localParticipant {
			continue
		}
		// TODO: maybe instead of using Comment to save each CFSMs name, fork the library and change attribute.
		if m.Comment != "" {
			ret = append(ret, m.Comment)
		} else {
			ret = append(ret, strconv.Itoa(m.ID))
		}
	}
	return
}

func ConvertPBContract(pbContract *pb.Contract) (Contract, error) {
	if pbContract.Format == pb.ContractFormat_CONTRACT_FORMAT_FSA {
		cfsmSystem, err := ParseFSAFile(bytes.NewReader(pbContract.Contract))
		if err != nil {
			return nil, err
		}
		var localCFSM *cfsm.CFSM
		for _, m := range cfsmSystem.CFSMs {
			if m.Comment != "" && m.Comment == pbContract.LocalParticipant {
				localCFSM = m
				break
			} else {
				if strconv.Itoa(m.ID) == pbContract.LocalParticipant {
					localCFSM = m
					break
				}
			}
		}
		if localCFSM == nil {
			return nil, fmt.Errorf("invalid contract. local_participant not present in FSA.")
		}
		contract := CFSMContract{
			System:           cfsmSystem,
			localParticipant: localCFSM,
		}
		return &contract, nil
	}
	return nil, fmt.Errorf("not implemented")
}

func ParseFSAFile(reader io.Reader) (*cfsm.System, error) {
	// f, err := os.Open("/tmp/dat")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()

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
		Send TransitionType = iota
		Recv
	)
	type Transition struct {
		FromState    *cfsm.State
		OtherMachine string
		Action       TransitionType
		Message      string
		NextState    *cfsm.State
	}

	var namesOfCFSMs []string            // array of names of CFSMs as we find in order. If they don't have name, we use str(int) of the machine order (starts with 0).
	numberOfCFSM := make(map[string]int) // inverse of the latter
	var status FSAParserStatus = Initial
	sys := cfsm.NewSystem() // This is what we'll return.

	// While reading transitions we can find transitions that refer to CFSMs that we haven't yet parsed.
	// So we'll save all transitions we find in this map, and add them all after consuming the entire fsa file.
	transitionsToBackfill := make(map[*cfsm.CFSM][]Transition)

	var stateNames map[string]*cfsm.State
	var currentMachine *cfsm.CFSM

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
			if matches[2] != "" {
				num, err := strconv.Atoi(matches[2])
				if err == nil && num != currentMachineNumber {
					return nil, fmt.Errorf("fsa file invalid: machine number %d is named %d", currentMachineNumber, num)
				}
				namesOfCFSMs = append(namesOfCFSMs, matches[2])
				numberOfCFSM[matches[2]] = currentMachineNumber
			} else {
				namesOfCFSMs = append(namesOfCFSMs, strconv.Itoa(len(namesOfCFSMs)))
				numberOfCFSM[strconv.Itoa(currentMachineNumber)] = currentMachineNumber
			}
			currentMachine = sys.NewMachine()
			currentMachine.Comment = namesOfCFSMs[currentMachineNumber]
			stateNames = make(map[string]*cfsm.State)
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
				action = Recv
			case "!":
				action = Send
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
			// TODO: We could also see transitions?
			if currentLine == ".end" {
				status = Initial
				continue
			} else {
				return nil, errors.New("fsa file invalid. Expected '.end'")
			}
		}

	}

	for _, transitions := range transitionsToBackfill {
		for _, t := range transitions {
			otherMachineNum, err := strconv.Atoi(t.OtherMachine)
			if err != nil {
				var ok bool
				otherMachineNum, ok = numberOfCFSM[t.OtherMachine]
				if !ok {
					return nil, errors.New("non existant CFSM referenced")
				}
			}

			otherMachine := sys.CFSMs[otherMachineNum]
			switch t.Action {
			case Send:
				e := cfsm.NewSend(otherMachine, t.Message)
				e.SetNext(t.NextState)
				t.FromState.AddTransition(e)
			case Recv:
				e := cfsm.NewRecv(otherMachine, t.Message)
				e.SetNext(t.NextState)
				t.FromState.AddTransition(e)
			}

		}
	}

	return sys, nil
}
