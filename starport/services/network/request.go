package network

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/cosmosutil"
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

// Requests fetches the chain requests from SPN by launch id
func (n Network) Requests(ctx context.Context, launchID uint64) ([]launchtypes.Request, error) {
	res, err := launchtypes.NewQueryClient(n.cosmos.Context).RequestAll(ctx, &launchtypes.QueryAllRequestRequest{
		LaunchID: launchID,
	})
	if err != nil {
		return nil, err
	}

	return res.Request, err
}

// Request fetches the chain request from SPN by launch and request id
func (n Network) Request(ctx context.Context, launchID, requestID uint64) (launchtypes.Request, error) {
	res, err := launchtypes.NewQueryClient(n.cosmos.Context).Request(ctx, &launchtypes.QueryGetRequestRequest{
		LaunchID:  launchID,
		RequestID: requestID,
	})
	if err != nil {
		return launchtypes.Request{}, err
	}
	return res.Request, err
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

// verifyAddValidatorRequest verify the validator request parameters
func (Network) verifyAddValidatorRequest(req *launchtypes.RequestContent_GenesisValidator) error {
	// If this is an add validator request
	var (
		peer           = req.GenesisValidator.Peer
		valAddress     = req.GenesisValidator.Address
		consPubKey     = req.GenesisValidator.ConsPubKey
		selfDelegation = req.GenesisValidator.SelfDelegation
	)

	// Check values inside the gentx are correct
	info, _, err := cosmosutil.ParseGentx(req.GenesisValidator.GenTx)
	if err != nil {
		return fmt.Errorf("cannot parse gentx %s", err.Error())
	}

	// Change the address prefix fetched from the gentx to the one used on SPN
	// Because all on-chain stored address on SPN uses the SPN prefix
	spnFetchedAddress, err := cosmosutil.ChangeAddressPrefix(info.DelegatorAddress, networkchain.SPN)
	if err != nil {
		return err
	}

	// Check validator address
	if valAddress != spnFetchedAddress {
		return fmt.Errorf(
			"the validator address %s doesn't match the one inside the gentx %s",
			valAddress,
			spnFetchedAddress,
		)
	}

	// Check validator address
	if !info.PubKey.Equal(consPubKey) {
		return fmt.Errorf(
			"the consensus pub key %s doesn't match the one inside the gentx %s",
			string(consPubKey),
			string(info.PubKey),
		)
	}

	// Check self delegation
	if selfDelegation.Denom != info.SelfDelegation.Denom ||
		!selfDelegation.IsEqual(info.SelfDelegation) {
		return fmt.Errorf(
			"the self delegation %s doesn't match the one inside the gentx %s",
			selfDelegation.String(),
			info.SelfDelegation.String(),
		)
	}

	// Check the format of the peer
	if !cosmosutil.VerifyPeerFormat(peer) {
		return fmt.Errorf(
			"the peer %s doesn't match the peer format <node-id>@<host>",
			peer,
		)
	}
	return nil
}

// VerifyRequests if the requests are correct and simulate them with the current launch information
// Correctness means checks that have to be performed off-chain
func (n Network) VerifyRequests(ctx context.Context, launchID uint64, requests ...uint64) error {
	n.ev.Send(events.New(events.StatusOngoing, "Verifying requests..."))
	// Check all request
	for _, id := range requests {
		request, err := n.Request(ctx, launchID, id)
		if err != nil {
			return err
		}

		req, ok := request.Content.Content.(*launchtypes.RequestContent_GenesisValidator)
		if ok {
			err := n.verifyAddValidatorRequest(req)
			if err != nil {
				return fmt.Errorf("request %d error: %s", id, err.Error())
			}
		}
	}
	n.ev.Send(events.New(events.StatusDone, "Requests verified"))

	// TODO simulate the requests
	// If all requests are correct, simulate them
	// return n.SimulateRequests(ctx, launchID, requests)

	return nil
}
