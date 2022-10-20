package network

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"

	"github.com/ignite/cli/ignite/pkg/cosmosutil"
	"github.com/ignite/cli/ignite/pkg/events"
	"github.com/ignite/cli/ignite/pkg/xurl"
	"github.com/ignite/cli/ignite/services/network/networkchain"
	"github.com/ignite/cli/ignite/services/network/networktypes"
)

type joinOptions struct {
	accountAmount sdk.Coins
	publicAddress string
}

type JoinOption func(*joinOptions)

// WithAccountRequest allows to join the chain by requesting a genesis account with the specified amount of tokens
func WithAccountRequest(amount sdk.Coins) JoinOption {
	return func(o *joinOptions) {
		o.accountAmount = amount
	}
}

// WithPublicAddress allows to specify a peer public address for the node
func WithPublicAddress(addr string) JoinOption {
	return func(o *joinOptions) {
		o.publicAddress = addr
	}
}

// Join to the network.
func (n Network) Join(
	ctx context.Context,
	c Chain,
	launchID uint64,
	gentxPath string,
	options ...JoinOption,
) error {
	o := joinOptions{}
	for _, apply := range options {
		apply(&o)
	}

	var (
		nodeID string
		peer   launchtypes.Peer
		err    error
	)

	// parse the gentx content
	gentxInfo, gentx, err := cosmosutil.GentxFromPath(gentxPath)
	if err != nil {
		return err
	}

	// get the peer address
	if o.publicAddress != "" {
		if nodeID, err = c.NodeID(ctx); err != nil {
			return err
		}

		if xurl.IsHTTP(o.publicAddress) {
			peer = launchtypes.NewPeerTunnel(nodeID, networkchain.HTTPTunnelChisel, o.publicAddress)
		} else {
			peer = launchtypes.NewPeerConn(nodeID, o.publicAddress)
		}
	} else {
		// if the peer address is not specified, we parse it from the gentx memo
		if peer, err = ParsePeerAddress(gentxInfo.Memo); err != nil {
			return err
		}
	}

	// change the chain address prefix to spn
	accountAddress, err := cosmosutil.ChangeAddressPrefix(gentxInfo.DelegatorAddress, networktypes.SPN)
	if err != nil {
		return err
	}

	if !o.accountAmount.IsZero() {
		if err := n.sendAccountRequest(ctx, launchID, accountAddress, o.accountAmount); err != nil {
			return err
		}
	}

	return n.sendValidatorRequest(ctx, launchID, peer, accountAddress, gentx, gentxInfo)
}

// sendValidatorRequest creates the RequestAddValidator message into the SPN
func (n Network) sendValidatorRequest(
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

	n.ev.Send("Broadcasting validator transaction", events.ProgressStarted())

	res, err := n.cosmos.BroadcastTx(ctx, n.account, msg)
	if err != nil {
		return err
	}

	var requestRes launchtypes.MsgSendRequestResponse
	if err := res.Decode(&requestRes); err != nil {
		return err
	}

	if requestRes.AutoApproved {
		n.ev.Send("Validator added to the network by the coordinator!", events.ProgressFinished())
	} else {
		n.ev.Send(
			fmt.Sprintf("Request %d to join the network as a validator has been submitted!", requestRes.RequestID),
			events.ProgressFinished(),
		)
	}
	return nil
}
