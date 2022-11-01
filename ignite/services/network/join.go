package network

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"

	"github.com/ignite/cli/ignite/pkg/cosmosutil"
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
		if err := n.SendAccountRequest(ctx, launchID, accountAddress, o.accountAmount); err != nil {
			return err
		}
	}

	return n.SendValidatorRequest(ctx, launchID, peer, accountAddress, gentx, gentxInfo)
}

func (n Network) SendAccountRequestForCoordinator(ctx context.Context, launchID uint64, amount sdk.Coins) error {
	addr, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return err
	}

	return n.SendAccountRequest(ctx, launchID, addr, amount)
}
