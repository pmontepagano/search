package searcherrors

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrChannelNotFound = errors.New("channel not found")

// https://grpc.io/docs/guides/error/#richer-error-model
// https://jbrandhorst.com/post/grpc-errors/
func BrokerageFailedError(failedParticipants []string) error {
	st := status.New(codes.NotFound, "CHANNEL_BROKERAGE_FAILED")
	ei := &errdetails.ErrorInfo{
		Reason: "CHANNEL_BROKERAGE_FAILED",
		Domain: "github.com/pmontepagano/search", // TODO: is this valid?
		Metadata: map[string]string{
			"failed_participants": strings.Join(failedParticipants[:], ","),
		},
	}
	st, err := st.WithDetails(ei)
	if err != nil {
		// If this errored, it will always error
		// here, so better panic so we can figure
		// out why than have this silently passing.
		panic(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
	}
	return st.Err()
}
