package network

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"

	"github.com/ignite/cli/ignite/pkg/events"
	"github.com/ignite/cli/ignite/services/network/networktypes"
)

// Reviewal keeps a request's reviewal.
type Reviewal struct {
	RequestID  uint64
	IsApproved bool
}

// ApproveRequest returns approval for a request with id.
func ApproveRequest(requestID uint64) Reviewal {
	return Reviewal{
		RequestID:  requestID,
		IsApproved: true,
	}
}

// RejectRequest returns rejection for a request with id.
func RejectRequest(requestID uint64) Reviewal {
	return Reviewal{
		RequestID:  requestID,
		IsApproved: false,
	}
}

// Requests fetches all the chain requests from SPN by launch id
func (n Network) Requests(ctx context.Context, launchID uint64) ([]networktypes.Request, error) {
	res, err := n.launchQuery.RequestAll(ctx, &launchtypes.QueryAllRequestRequest{
		LaunchID: launchID,
	})
	if err != nil {
		return nil, err
	}
	requests := make([]networktypes.Request, len(res.Request))
	for i, req := range res.Request {
		requests[i] = networktypes.ToRequest(req)
	}
	return requests, nil
}

// Request fetches the chain request from SPN by launch and request id
func (n Network) Request(ctx context.Context, launchID, requestID uint64) (networktypes.Request, error) {
	res, err := n.launchQuery.Request(ctx, &launchtypes.QueryGetRequestRequest{
		LaunchID:  launchID,
		RequestID: requestID,
	})
	if err != nil {
		return networktypes.Request{}, err
	}
	return networktypes.ToRequest(res.Request), nil
}

// RequestFromIDs fetches the chain requested from SPN by launch and provided request IDs
// TODO: once implemented, use the SPN query from https://github.com/tendermint/spn/issues/420
func (n Network) RequestFromIDs(ctx context.Context, launchID uint64, requestIDs ...uint64) (reqs []networktypes.Request, err error) {
	for _, id := range requestIDs {
		req, err := n.Request(ctx, launchID, id)
		if err != nil {
			return reqs, err
		}
		reqs = append(reqs, req)
	}
	return reqs, nil
}

// SubmitRequest submits reviewals for proposals in batch for chain.
func (n Network) SubmitRequest(launchID uint64, reviewal ...Reviewal) error {
	n.ev.Send(events.New(events.StatusOngoing, "Submitting requests..."))

	messages := make([]sdk.Msg, len(reviewal))
	for i, reviewal := range reviewal {
		messages[i] = launchtypes.NewMsgSettleRequest(
			n.account.Address(networktypes.SPN),
			launchID,
			reviewal.RequestID,
			reviewal.IsApproved,
		)
	}

	res, err := n.cosmos.BroadcastTx(n.account.Name, messages...)
	if err != nil {
		return err
	}

	var requestRes launchtypes.MsgSettleRequestResponse
	return res.Decode(&requestRes)
}
