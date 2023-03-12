package contract

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/nickng/cfsm"
)

// GlobalContract represents a requirements contract that characterizes a channel
type GlobalContract struct {
	contract     string
	participants []string
}

// PartialContract represents a projection of a GlobalContract to a participant
type PartialContract struct {
	contract     string
	participants []string
}

// LocalContract represents a provides contract
type LocalContract struct {
	contract     string
	participants []string
}

// Projects GC into PartialContract. For now, this is a dummy implementation that only
// considers contracts of two participants, so we simply copy contract and participant fields
func (gc *GlobalContract) ProjectPartialContract(participant string) PartialContract {
	return PartialContract{
		contract:     gc.contract,
		participants: gc.participants,
	}
}

// Determines if candidate LocalContract satisfies the PartialContract. For now, this is
// a dummy implementation that always returns TRUE
// TODO: do the check, this is a dummy implementation
func (pc *PartialContract) IsSatisfiedBy(candidate LocalContract) bool {
	return true
}

func ParseFSAFile(filename string) (cfsm.System, err error) {
	f, err := os.Open("/tmp/dat")
	if err != nil {
		panic(err)
	}

	// Regular expressions for parser
	lineCommentRe := regexp.MustCompile(`^--.*`)
	machineStartRe := regexp.MustCompile(`^\.outputs\s+(\S)?$`)
	startStateRe := regexp.MustCompile(`^\.marking\s+(\S)$`)
	machineEndRe := regexp.MustCompile(`^\.end$`)
	messageRe := regexp.MustCompile(`^(\S)\s(\S)\s([\?!])\s(\S)\s(\S)$`)

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
		Send	TransitionType = iota
		Recv
	)
	type Transition struct {
		FromState *cfsm.State
		OtherMachine	string
		Action TransitionType
		Message string
		NextState *cfsm.State
	}

	var namesOfCFSMs []string  // array of names of CFSMs as we find in order. If they don't have name, we use str(int) of the machine order (starts with 0).
	numberOfCFSM := make(map[string]int)  // inverse of the latter
	var status FSAParserStatus = Initial
	sys := cfsm.NewSystem()  // This is what we'll return.
	
	// While parsing we may find transitions that refer to a CFSM that we haven't yet parsed.
	// In that case, we add the machine name to this set. We'll have to check that this set is empty
	// when we finish parsing.
	machinesNotYetDeclared := make(map[string]struct{})
	
	// While reading transitions we can find transitions that refer to CFSMs that we haven't yet parsed.
	// In that case, we add the transition to this set and backfill the transitions after reading the
	// entire file.
	transitionsToBackfill := make(map[*cfsm.CFSM][]Transition)

	var stateNames map[string]*cfsm.State
	var currentMachine *cfsm.CFSM

	
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		currentLine := scanner.Text()
		currentLine = strings.TrimSpace(currentLine)
		if lineCommentRe.MatchString(currentLine) || len(currentLine == 0) {
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
			if len(matches) == 2 {
				num, err := strconv.Atoi(matches[1])
				if err == nil && num != currentMachineNumber {
					return nil, errors.New(fmt.Sprintf("fsa file invalid. Machine number %d is named %d", currentMachineNumber, num))
				}
				namesOfCFSMs = append(namesOfCFSMs, matches[1])
				numberOfCFSM[matches[1]] = currentMachineNumber
			} else {
				namesOfCFSMs = append(namesOfCFSMs, len(namesOfCFSMs))
				numberOfCFSM[currentMachineNumber] = currentMachineNumber
			}
			currentMachine = sys.NewMachine()
			stateNames = make(map[string]*cfsm.State)
			transitionsToBackfill[&currentMachine] = make([]Transition)
			status = MachineStarting
		case MachineStarting:
			// We can only expect now a (useless) ".state graph" line.
			if currentLine == ".state graph"{
				status = Transitions
				continue
			} else {
				return nil, errors.New("fsa file invalid. Expected '.state graph' after CFSM declaration.")
			}
		case Transitions:
			// Here we can see transition lines or the start state marker.
			startStateMatch := startStateRe.FindStringSubmatch(currentLine)
			if startStateMatch != nil && len(startStateMatch) == 2 {
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
				return nil, errors.New("Expected transition.")
			}
			// _, fromStateStr, otherCFSM, action, msg, nextStateStr := transitionMatches...
			
			// Create fromState and nextState if they don't exist already.
			val, ok := stateNames[transitionMatches[1]]
			if !ok {
				val := currentMachine.NewState()
				val.Label = transitionMatches[1]
				stateNames[transitionMatches[1]] = val
			}
			fromState := val
			
			val, ok = stateNames[transitionMatches[5]]
			if !ok {
				val := currentMachine.NewState()
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
				return nil, errors.New("Invalid action.")
			}
			
			msg := transitionMatches[4]

			// Add transition to be backfilled.

			transitionsToBackfill[currentMachine] = append(transitionsToBackfill[currentMachine], Transition{
				FromState: fromState,
				OtherMachine: otherCFSM,
				Action: action,
				Message: msg,
				NextState: nextState,
			})

		case StartStateRead:
			// Here we can only see the end marker.
			// TODO: We could also see transitions?
			if currentLine == ".end"{
				status = Initial
				continue
			} else {
				return nil, errors.New("fsa file invalid. Expected '.end'.")
			}
		}

		
	}

	for machine, transitions := range transitionsToBackfill {
		for _, t := range transitions {
			otherMachineNum, ok := numberOfCFSM[t.OtherMachine]
			if !ok {
				return nil, errors.New("Non existant CFSM referenced.")
			}
			otherMachine := sys.CFSMs[otherMachineNum]
			switch t.Action {
			case Send:
				e := cfsm.NewSend(otherMachine, t.Message)
			case Recv:
				e := cfsm.NewRecv(otherMachine, t.Message)
			}
			e.SetNext(t.NextState)
			t.FromState.AddTransition(e)
		}
	}

	if len(machinesNotYetDeclared) > 0 {
		return nil, errors.New("Invalid fsa. Some transition refers to an undeclared CFSM.")
	}

	return sys, nil
}
