package network

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/events"
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

	var requestRes launchtypes.MsgSettleRequestResponse
	if err := res.Decode(&requestRes); err != nil {
		return err
	}
	b.ev.Send(events.New(events.StatusDone, "Settle request transactions sent"))
	return nil
}
