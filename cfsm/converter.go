package cfsm

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
)

// This function converts a CFSM to a string in the format of the Python Bisimulation library
// https://github.com/diegosenarruzza/bisimulation/
func ConvertCFSMToPythonBisimulationFormat(contract *CFSM) ([]byte, error) {
	// Import statement
	const importStatement = "from cfsm_bisimulation import CommunicatingFiniteStateMachine\n\n"

	// Create CFSM
	// Diego's format expects also the name for the CFSM we're defining, so we'll have to pick
	// a name that is unused. We'll use 'self' unless it's already used. If it is, we'll generate
	// a random string and use that (verifying it's also unused).
	otherCFSMs := contract.OtherCFSMs()
	selfname := "self"
	restart := false
	for restart {
		restart = false
		for _, name := range otherCFSMs {
			if name == selfname {
				selfname = "self_" + generateRandomString(5)
				restart = true
				break
			}
		}
	}
	allCFSMs := append(otherCFSMs, selfname)
	allCFSMsJSON, err := json.Marshal(allCFSMs)
	if err != nil {
		return nil, err
	}
	createCFSM := "cfsm = CommunicatingFiniteStateMachine(" + string(allCFSMsJSON) + ")\n\n"

	// Add states
	addStates := ""
	for _, state := range contract.States() {
		addStates += fmt.Sprintf("cfsm.add_states('%d')\n", state.ID)
	}
	addStates += "\n"

	// Set initial state
	setInitialState := fmt.Sprintf("cfsm.set_as_initial('%d')\n\n", contract.Start.ID)

	// Add transitions
	addTransitions := ""
	for _, state := range contract.States() {
		for _, transition := range state.Transitions() {
			// cfsm_bisimulation format expects the message to be a string with the format
			// is SenderReceiver[!|?]function([typed parameters]). e.g "ClientServer!f(int x)" or "ServerClient?g()"

			msgRegex := `\w+\(.*\)`
			matched, err := regexp.MatchString(msgRegex, transition.Message())
			if err != nil {
				return nil, err
			}
			if !matched {
				return nil, fmt.Errorf("the message in transition %s does not match the format expected by cfsm_simulation", transition.Label())
			}

			var transitionTypeMarker string
			if transition.IsSend() {
				transitionTypeMarker = "!"
			} else {
				transitionTypeMarker = "?"
			}
			action := fmt.Sprintf("%s%s%s%s", selfname, transition.NameOfOtherCFSM(), transitionTypeMarker, transition.Message())
			addTransitions += fmt.Sprintf("cfsm.add_transition_between('%d', '%d', '%s', True)\n", state.ID, transition.State().ID, action)
		}
	}

	// Combine all code blocks
	code := importStatement + createCFSM + addStates + setInitialState + addTransitions + "\n"

	return []byte(code), nil
}

func generateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
