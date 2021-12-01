package network

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network/networkchain"
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

// SubmitRequest submits reviewals for proposals in batch for chain.
func (n Network) SubmitRequest(launchID uint64, reviewal ...Reviewal) error {
	n.ev.Send(events.New(events.StatusOngoing, "Submitting requests..."))

	messages := make([]sdk.Msg, len(reviewal))
	for i, reviewal := range reviewal {
		messages[i] = launchtypes.NewMsgSettleRequest(
			n.account.Address(networkchain.SPN),
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
