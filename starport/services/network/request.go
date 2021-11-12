package network

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/gentx"
)

// Reviewal keeps a request's reviewal.
type Reviewal struct {
	RequestID  uint64
	IsApproved bool
}

// ApproveProposal returns approval for a proposal with id.
func ApproveProposal(requestID uint64) Reviewal {
	return Reviewal{
		RequestID:  requestID,
		IsApproved: true,
	}
}

// RejectProposal returns rejection for a proposals with id.
func RejectProposal(requestID uint64) Reviewal {
	return Reviewal{
		RequestID:  requestID,
		IsApproved: false,
	}
}

// SubmitRequest submits reviewals for proposals in batch for chain.
func (b *Builder) SubmitRequest(launchID uint64, reviewals ...Reviewal) error {
	b.ev.Send(events.New(events.StatusOngoing, "Approving requests..."))

	messages := make([]sdk.Msg, len(reviewals))
	for i, reviewal := range reviewals {
		messages[i] = launchtypes.NewMsgSettleRequest(
			b.account.Address(SPNAddressPrefix),
			launchID,
			reviewal.RequestID,
			reviewal.IsApproved,
		)
	}

	res, err := b.cosmos.BroadcastTx(b.account.Name, messages...)
	if err != nil {
		return err
	}

	var requestRes launchtypes.MsgRequestResponse
	if err := res.Decode(&requestRes); err != nil {
		return err
	}
	b.ev.Send(events.New(events.StatusDone, "Settle request transactions sent"))
	return nil
}

// fetchRequest fetches the chain request from SPN by launch and request id
func (b *Builder) fetchRequest(ctx context.Context, launchID, requestID uint64) (launchtypes.Request, error) {
	res, err := launchtypes.NewQueryClient(b.cosmos.Context).Request(ctx, &launchtypes.QueryGetRequestRequest{
		LaunchID:  launchID,
		RequestID: requestID,
	})
	if err != nil {
		return launchtypes.Request{}, err
	}
	return res.Request, err
}

// VerifyRequests if the requests are correct and simulate them with the current launch information
// Correctness means checks that have to be performed off-chain
func (b *Builder) VerifyRequests(ctx context.Context, launchID uint64, requests []uint64) error {
	b.ev.Send(events.New(events.StatusOngoing, "Verifying requests..."))
	// Check all request
	for _, id := range requests {
		request, err := b.fetchRequest(ctx, launchID, id)
		if err != nil {
			return err
		}

		req, ok := request.Content.Content.(*launchtypes.RequestContent_GenesisValidator)
		if ok {
			// If this is an add validator request
			valAddress := req.GenesisValidator.Address
			selfDelegation := req.GenesisValidator.SelfDelegation

			// Check values inside the gentx are correct
			info, _, err := gentx.ParseGentx(req.GenesisValidator.GenTx)
			if err != nil {
				return fmt.Errorf("cannot parse request %v gentx: %v", id, err.Error())
			}

			// Check validator address
			if valAddress != info.DelegatorAddress {
				return fmt.Errorf(
					"request %v contains a validator address %v that doesn't match the one inside the gentx: %v",
					id,
					valAddress,
					info.DelegatorAddress,
				)
			}

			// Check self delegation
			if !selfDelegation.IsEqual(info.SelfDelegation) {
				return fmt.Errorf(
					"request %v contains a self delegation %v that doesn't match the one inside the gentx: %v",
					id,
					selfDelegation.String(),
					info.SelfDelegation.String(),
				)
			}
		}
	}
	b.ev.Send(events.New(events.StatusDone, "Requests verified"))

	// TODO simulate the proposals
	// If all proposals are correct, simulate them
	// return b.SimulateProposals(ctx, chainID, proposals, commandOut)

	return nil
}
