package network

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strings"

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
func (b *Builder) SubmitRequest(launchID uint64, reviewal ...Reviewal) error {
	b.ev.Send(events.New(events.StatusOngoing, "Submitting requests..."))

	messages := make([]sdk.Msg, len(reviewal))
	for i, reviewal := range reviewal {
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

// verifyAddValidatorRequest verify the validator request parameters
func (*Builder) verifyAddValidatorRequest(req *launchtypes.RequestContent_GenesisValidator) error {
	// If this is an add validator request
	peer := req.GenesisValidator.Peer
	valAddress := req.GenesisValidator.Address
	consPubKey := req.GenesisValidator.ConsPubKey
	selfDelegation := req.GenesisValidator.SelfDelegation

	// Check values inside the gentx are correct
	info, _, err := gentx.ParseGentx(req.GenesisValidator.GenTx)
	if err != nil {
		return fmt.Errorf("cannot parse gentx %s", err.Error())
	}

	// Check validator address
	if valAddress != info.DelegatorAddress {
		return fmt.Errorf(
			"the validator address %s doesn't match the one inside the gentx %s",
			valAddress,
			info.DelegatorAddress,
		)
	}

	// Check validator address
	if res := bytes.Compare(consPubKey, info.PubKey); res != 0 {
		return fmt.Errorf(
			"the consensus pub key %s doesn't match the one inside the gentx %s",
			string(consPubKey),
			string(info.PubKey),
		)
	}

	// Check self delegation
	if !selfDelegation.IsEqual(info.SelfDelegation) {
		return fmt.Errorf(
			"the self delegation %s doesn't match the one inside the gentx %s",
			selfDelegation.String(),
			info.SelfDelegation.String(),
		)
	}

	// Check the format of the peer
	nodeHost := strings.Split(peer, "@")
	if len(nodeHost) < 2 {
		return fmt.Errorf(
			"the peer %s doesn't match the peer format <node-id>@<host>",
			peer,
		)
	}
	nodeID := nodeHost[0]
	if len(nodeID) == 0 {
		return fmt.Errorf("empty peer node id")
	}

	// Looks up the given host
	host := nodeHost[1]
	hostName, err := net.LookupHost(host)
	if err != nil || len(hostName) == 0 {
		return fmt.Errorf(
			"the peer host %s contains an invalid host",
			host,
		)
	}
	return nil
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
			err := b.verifyAddValidatorRequest(req)
			if err != nil {
				return fmt.Errorf("request %d error: %s", id, err.Error())
			}
		}
	}
	b.ev.Send(events.New(events.StatusDone, "Requests verified"))

	// TODO simulate the proposals
	// If all proposals are correct, simulate them
	// return b.SimulateProposals(ctx, chainID, proposals, commandOut)

	return nil
}
