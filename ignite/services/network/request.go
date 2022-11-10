package network

import (
	"context"
	"fmt"

	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/cosmosutil"

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
func (n Network) SubmitRequest(ctx context.Context, launchID uint64, reviewal ...Reviewal) error {
	n.ev.Send("Submitting requests...", events.ProgressStart())

	addr, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return err
	}

	messages := make([]sdk.Msg, len(reviewal))
	for i, reviewal := range reviewal {
		messages[i] = launchtypes.NewMsgSettleRequest(
			addr,
			launchID,
			reviewal.RequestID,
			reviewal.IsApproved,
		)
	}

	res, err := n.cosmos.BroadcastTx(ctx, n.account, messages...)
	if err != nil {
		return err
	}

	var requestRes launchtypes.MsgSettleRequestResponse
	return res.Decode(&requestRes)
}

// SendAccountRequest creates an add AddAccount request message.
func (n Network) SendAccountRequest(
	ctx context.Context,
	launchID uint64,
	address string,
	amount sdk.Coins,
) error {
	addr, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return err
	}

	msg := launchtypes.NewMsgSendRequest(
		addr,
		launchID,
		launchtypes.NewGenesisAccount(
			launchID,
			address,
			amount,
		),
	)

	n.ev.Send("Broadcasting account transactions", events.ProgressStart())

	res, err := n.cosmos.BroadcastTx(ctx, n.account, msg)
	if err != nil {
		return err
	}

	var requestRes launchtypes.MsgSendRequestResponse
	if err := res.Decode(&requestRes); err != nil {
		return err
	}

	if requestRes.AutoApproved {
		n.ev.Send(
			"Account added to the network by the coordinator!",
			events.Icon(icons.Bullet),
			events.ProgressFinish(),
		)
	} else {
		n.ev.Send(
			fmt.Sprintf("Request %d to add account to the network has been submitted!", requestRes.RequestID),
			events.Icon(icons.Bullet),
			events.ProgressFinish(),
		)
	}
	return nil
}

// SendValidatorRequest creates the RequestAddValidator message into the SPN
func (n Network) SendValidatorRequest(
	ctx context.Context,
	launchID uint64,
	peer launchtypes.Peer,
	valAddress string,
	gentx []byte,
	gentxInfo cosmosutil.GentxInfo,
) error {
	addr, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return err
	}

	msg := launchtypes.NewMsgSendRequest(
		addr,
		launchID,
		launchtypes.NewGenesisValidator(
			launchID,
			valAddress,
			gentx,
			gentxInfo.PubKey,
			gentxInfo.SelfDelegation,
			peer,
		),
	)

	n.ev.Send("Broadcasting validator transaction", events.ProgressStart())

	res, err := n.cosmos.BroadcastTx(ctx, n.account, msg)
	if err != nil {
		return err
	}

	var requestRes launchtypes.MsgSendRequestResponse
	if err := res.Decode(&requestRes); err != nil {
		return err
	}

	if requestRes.AutoApproved {
		n.ev.Send("Validator added to the network by the coordinator!", events.ProgressFinish())
	} else {
		n.ev.Send(
			fmt.Sprintf("Request %d to join the network as a validator has been submitted!", requestRes.RequestID),
			events.ProgressFinish(),
		)
	}
	return nil
}

// SendValidatorRemoveRequest creates the RequestRemoveValidator message to SPN
func (n Network) SendValidatorRemoveRequest(
	ctx context.Context,
	launchID uint64,
	valAddress string,
) error {
	addr, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return err
	}

	msg := launchtypes.NewMsgSendRequest(
		addr,
		launchID,
		launchtypes.NewValidatorRemoval(
			valAddress,
		),
	)

	n.ev.Send("Broadcasting transaction", events.ProgressStart())

	res, err := n.cosmos.BroadcastTx(ctx, n.account, msg)
	if err != nil {
		return err
	}

	var requestRes launchtypes.MsgSendRequestResponse
	if err := res.Decode(&requestRes); err != nil {
		return err
	}

	if requestRes.AutoApproved {
		n.ev.Send("Validator removed from network by the coordinator!", events.ProgressFinish())
	} else {
		n.ev.Send(
			fmt.Sprintf(
				"Request %d to remove validator from the network has been submitted!", requestRes.RequestID,
			),
			events.ProgressFinish(),
		)
	}
	return nil
}

// SendAccountRemoveRequest creates the RequestRemoveAccount message to SPN
func (n Network) SendAccountRemoveRequest(
	ctx context.Context,
	launchID uint64,
	address string,
) error {
	addr, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return err
	}

	msg := launchtypes.NewMsgSendRequest(
		addr,
		launchID,
		launchtypes.NewAccountRemoval(
			address,
		),
	)

	n.ev.Send("Broadcasting transaction", events.ProgressStart())

	res, err := n.cosmos.BroadcastTx(ctx, n.account, msg)
	if err != nil {
		return err
	}

	var requestRes launchtypes.MsgSendRequestResponse
	if err := res.Decode(&requestRes); err != nil {
		return err
	}

	if requestRes.AutoApproved {
		n.ev.Send("Account removed from network by the coordinator!", events.ProgressFinish())
	} else {
		n.ev.Send(
			fmt.Sprintf(
				"Request %d to remove account from the network has been submitted!", requestRes.RequestID,
			),
			events.ProgressFinish(),
		)
	}
	return nil
}
