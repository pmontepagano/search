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
	q0 OurServer ! login() q1
	q1 OurServer ? accept() q2
	.marking q0
	.end
	`)))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedOutput := []byte(`{{.MachineName}} = CommunicatingFiniteStateMachine(["Self","Courserver"])

{{.MachineName}}.add_states('0')
{{.MachineName}}.add_states('1')
{{.MachineName}}.add_states('2')

{{.MachineName}}.set_as_initial('0')

{{.MachineName}}.add_transition_between('0', '1', 'Self Courserver ! message0()')
{{.MachineName}}.add_transition_between('1', '2', 'Self Courserver ? message1()')

`)

	output, _, _, err := ConvertCFSMToPythonBisimulationFormat(contract)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("Output does not match expected value.\nExpected: %s\nGot: %s", expectedOutput, output)
	}
}
