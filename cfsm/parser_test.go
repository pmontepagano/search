package cfsm

import (
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestHelloWorldFSAParse(t *testing.T) {
	const exampleFSAContent = `
.outputs A
.state graph
q0 1 ! hello q1
q0 1 ! world q1
.marking q0
.end

.outputs
.state graph
q0 0 ? hello q1
q0 0 ? world q1
.marking q0
.end`
	fs := fstest.MapFS{
		"example.fsa": {Data: []byte(exampleFSAContent)},
	}

	exampleFile, err := fs.Open("example.fsa")
	if err != nil {
		t.Error("Failed reading example.fsa.")
	}
	defer exampleFile.Close()

	sys, err := ParseSystemCFSMsFSA(exampleFile)
	if err != nil {
		t.Errorf("Error parsing FSA: %s", err)
	}

	if len(sys.CFSMs) != 2 {
		t.Errorf("Parsed FSA has %d CFSMs. Expected 2.", len(sys.CFSMs))
	}
	firstCFSM := sys.CFSMs[0]
	if firstCFSM.Name != "A" {
		t.Errorf("Expected fist CFSM to have Name 'A' but instead had %s", firstCFSM.Name)
	}
	if len(firstCFSM.States()) != 2 {
		t.Errorf("Expected 2 states in first CFSM. Found %d", len(firstCFSM.States()))
	}
	if firstCFSM.Start.Label != "q0" {
		t.Errorf("Expected start state of first CFSM to be 'q0'. Found %s instead", firstCFSM.Start.Label)
	}

	// TODO: we should add a test to check output format.
	// outString := sys.String()
	// require.Equal(t, exampleFSAContent, outString)

}

func TestPingPongFSAParser(t *testing.T) {
	const exampleFSAContent = `
.outputs Ping
.state graph
0 1 ! ping 5
2 1 ? bye 1
3 1 ! bye 2
3 1 ! finished 2
4 1 ! *<1 0
4 1 ! >*1 3
5 1 ? pong 4
.marking 0
.end

.outputs Pong
.state graph
0 0 ? ping 5
2 0 ! bye 1
3 0 ? bye 2
3 0 ? finished 2
4 0 ? *<1 0
4 0 ? >*1 3
5 0 ! pong 4
.marking 0
.end
`
	fs := fstest.MapFS{
		"example.fsa": {Data: []byte(exampleFSAContent)},
	}

	exampleFile, err := fs.Open("example.fsa")
	if err != nil {
		t.Error("Failed reading example.fsa.")
	}
	defer exampleFile.Close()

	sys, err := ParseSystemCFSMsFSA(exampleFile)
	if err != nil {
		t.Errorf("Error parsing FSA: %s", err)
	}

	if len(sys.CFSMs) != 2 {
		t.Errorf("Parsed FSA has %d CFSMs. Expected 2.", len(sys.CFSMs))
	}

	require.ElementsMatch(t, []string{"Ping", "Pong"}, sys.GetAllMachineNames())

	firstCFSM := sys.CFSMs[0]
	if firstCFSM.Name != "Ping" {
		t.Errorf("Expected first CFSM to have Name 'Ping' but instead had %s", firstCFSM.Name)
	}
	if len(firstCFSM.States()) != 6 {
		t.Errorf("Expected 6 states in first CFSM. Found %d", len(firstCFSM.States()))
	}

}

func TestPingPongFSALocalContractParser(t *testing.T) {
	const exampleFSAContent = `
.outputs Ping
.state graph
0 1 ! ping 5
2 1 ? bye 1
3 1 ! bye 2
3 1 ! finished 2
4 1 ! *<1 0
4 1 ! >*1 3
5 1 ? pong 4
.marking 0
.end
`
	fs := fstest.MapFS{
		"example.fsa": {Data: []byte(exampleFSAContent)},
	}

	exampleFile, err := fs.Open("example.fsa")
	if err != nil {
		t.Error("Failed reading example.fsa.")
	}
	defer exampleFile.Close()

	machine, err := ParseSingleCFSMFSA(exampleFile)
	if err != nil {
		t.Errorf("Error parsing FSA: %s", err)
	}

	if machine.Name != "Ping" {
		t.Errorf("Expected CFSM to have Name 'Ping' but instead had %s", machine.Name)
	}
	if len(machine.States()) != 6 {
		t.Errorf("Expected 6 states in first CFSM. Found %d", len(machine.States()))
	}

}

func TestCFSMSerialization(t *testing.T) {
	const pingPongFSA = `
.outputs Ping
.state graph
0 Pong ! ping 1
1 Pong ? pong 0
0 Pong ! bye 2
0 Pong ! finished 3
2 Pong ? bye 3
.marking 0
.end

.outputs Pong
.state graph
0 Ping ? ping 1
1 Ping ! pong 0
0 Ping ? bye 2
2 Ping ! bye 3
0 Ping ? finished 3
.marking 0
.end
`
	const cfsmStringPingPong = `
-- Machine #0
.outputs Ping
.state graph
q00 Pong ! ping q01
q00 Pong ! bye q02
q00 Pong ! finished q03
q01 Pong ? pong q00
q02 Pong ? bye q03
.marking q00
.end

-- Machine #1
.outputs Pong
.state graph
q10 Ping ? ping q11
q10 Ping ? bye q12
q10 Ping ? finished q13
q11 Ping ! pong q10
q12 Ping ! bye q13
.marking q10
.end
`
	pingPongSystem, err := ParseSystemCFSMsFSA(strings.NewReader(pingPongFSA))
	if err != nil {
		t.Errorf("Error parsing FSA: %s", err)
	}

	require.Equal(t, cfsmStringPingPong, pingPongSystem.String())
}

func TestSingleCFSMSerialization(t *testing.T) {
	const cfsmPongFSA = `.outputs Pong
	.state graph
	0 Ping ? ping 1
	1 Ping ! pong 0
	0 Ping ? bye 2
	2 Ping ! bye 3
	0 Ping ? finished 3
	.marking 0
	.end`
	const normalizedPongString = `
-- Machine #0
.outputs Pong
.state graph
q00 Ping ? ping q01
q00 Ping ? bye q02
q00 Ping ? finished q03
q01 Ping ! pong q00
q02 Ping ! bye q03
.marking q00
.end
`
	pongCFSM, err := ParseSingleCFSMFSA(strings.NewReader(cfsmPongFSA))
	if err != nil {
		t.Errorf("Error parsing FSA: %s", err)
	}

	require.Equal(t, normalizedPongString, pongCFSM.String())
}
