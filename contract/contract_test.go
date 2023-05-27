package contract

import (
	"testing"

	pb "github.com/clpombo/search/gen/go/search/v1"
	"github.com/stretchr/testify/require"
)

func TestProjectionID(t *testing.T) {
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
	gc, err := ConvertPBGlobalContract(&pb.GlobalContract{
		Contract:      []byte(pingPongFSA),
		Format:        pb.GlobalContractFormat_GLOBAL_CONTRACT_FORMAT_FSA,
		InitiatorName: "Ping",
	})
	if err != nil {
		t.Errorf("Error converting FSA: %s", err)
	}

	require.Equal(t, "0d6ebde80b89baf1b76d33b636c49985c0792cdd93b9fa204932347ee1b5eabaa21c796f20161d6f0aba42e0b7172ef07566b56f8a3840512df44c77f9c84d42", gc.GetContractID())
}
