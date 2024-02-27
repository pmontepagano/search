package cfsm

import (
	"bytes"
	"reflect"
	"testing"
)

func TestConvertCFSMToPythonBisimulationFormat(t *testing.T) {
	// Create a sample CFSM for testing
	contract, err := ParseSingleCFSMFSA(bytes.NewReader([]byte(`
	.outputs self
	.state graph
	q0 Server ! login() q1
	q1 Server ? accept() q2
	.marking q0
	.end
	`)))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedOutput := []byte(`from cfsm_bisimulation import CommunicatingFiniteStateMachine

cfsm = CommunicatingFiniteStateMachine(["Self","Server"])

cfsm.add_states('0')
cfsm.add_states('1')
cfsm.add_states('2')

cfsm.set_as_initial('0')

cfsm.add_transition_between('0', '1', 'Self Server ! message0()')
cfsm.add_transition_between('1', '2', 'Self Server ? message1()')

`)

	output, err := ConvertCFSMToPythonBisimulationFormat(contract)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("Output does not match expected value.\nExpected: %s\nGot: %s", expectedOutput, output)
	}
}
