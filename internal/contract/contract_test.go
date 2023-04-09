package contract

import (
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

	sys, err := ParseFSAFile(exampleFile)
	if err != nil {
		t.Errorf("Error parsing FSA: %s", err)
	}

	if len(sys.CFSMs) != 2 {
		t.Errorf("Parsed FSA has %d CFSMs. Expected 2.", len(sys.CFSMs))
	}
	firstCFSM := sys.CFSMs[0]
	if firstCFSM.Comment != "A" {
		t.Errorf("Expected fist CFSM to have Comment 'A' but instead had %s", firstCFSM.Comment)
	}
	if len(firstCFSM.States()) != 2 {
		t.Errorf("Expected 2 states in first CFSM. Found %d", len(firstCFSM.States()))
	}
	if firstCFSM.Start.Label != "q0" {
		t.Errorf("Expected start state of first CFSM to be 'q0'. Found %s instead", firstCFSM.Start.Label)
	}

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

	sys, err := ParseFSAFile(exampleFile)
	if err != nil {
		t.Errorf("Error parsing FSA: %s", err)
	}

	if len(sys.CFSMs) != 2 {
		t.Errorf("Parsed FSA has %d CFSMs. Expected 2.", len(sys.CFSMs))
	}
	require.ElementsMatch(t, []string{"Ping", "Pong"}, sys.GetParticipants())

	firstCFSM := sys.CFSMs[0]
	if firstCFSM.Comment != "Ping" {
		t.Errorf("Expected fist CFSM to have Comment 'Ping' but instead had %s", firstCFSM.Comment)
	}
	if len(firstCFSM.States()) != 6 {
		t.Errorf("Expected 6 states in first CFSM. Found %d", len(firstCFSM.States()))
	}

}
