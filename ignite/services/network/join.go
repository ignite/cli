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
	nodeID        string
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

func WithNodeID(nodeID string) JoinOption {
	return func(o *joinOptions) {
		o.nodeID = nodeID
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
		nodeID = o.nodeID
		peer   launchtypes.Peer
		err    error
	)

	// parse the gentx content
	gentxInfo, gentx, err := cosmosutil.GentxFromPath(gentxPath)
	if err != nil {
		return err
	}

	// get the peer address
	if nodeID == "" {
		if nodeID, err = c.NodeID(ctx); err != nil {
			return err
		}
	}

	switch {
	case xurl.IsHTTP(o.publicAddress):
		peer = launchtypes.NewPeerTunnel(nodeID, networkchain.HTTPTunnelChisel, o.publicAddress)
	case xurl.IsTCP(o.publicAddress):
		peer = launchtypes.NewPeerConn(nodeID, o.publicAddress)
	case o.publicAddress == "":
		peer = launchtypes.NewPeerEmpty(nodeID)
	default:
		return fmt.Errorf("unsupported public address format: %s", o.publicAddress)
	}

	// TODO: support peer parsing form memo
	// peer, err = ParsePeerAddress(gentxInfo.Memo)

	// get the chain genesis path from the home folder
	genesisPath, err := c.GenesisPath()
	if err != nil {
		return err
	}

	// change the chain address prefix to spn
	accountAddress, err := cosmosutil.ChangeAddressPrefix(gentxInfo.DelegatorAddress, networktypes.SPN)
	if err != nil {
		return err
	}

	if !o.accountAmount.IsZero() {
		if err := n.ensureAccount(genesisPath, launchID, accountAddress, o.accountAmount); err != nil {
			return err
		}
	}

	return n.sendValidatorRequest(launchID, peer, accountAddress, gentx, gentxInfo)
}

// ensureAccount creates an add AddAccount request message.
func (n Network) ensureAccount(
	genesisPath string,
	launchID uint64,
	address string,
	amount sdk.Coins,
) (err error) {
	n.ev.Send(events.New(events.StatusOngoing, "Verifying account already exists "+address))

	// the account may already exist in the initial genesis, we check it from the generated genesis
	accExist, err := cosmosutil.CheckGenesisContainsAddress(genesisPath, address)
	if err != nil {
		return err
	}
	if accExist {
		return fmt.Errorf("account %s already exist in the initial genesis", address)
	}

	return n.sendAccountRequest(launchID, address, amount)
}

// sendValidatorRequest creates the RequestAddValidator message into the SPN
func (n Network) sendValidatorRequest(
	launchID uint64,
	peer launchtypes.Peer,
	valAddress string,
	gentx []byte,
	gentxInfo cosmosutil.GentxInfo,
) error {
	msg := launchtypes.NewMsgRequestAddValidator(
		n.account.Address(networktypes.SPN),
		launchID,
		valAddress,
		gentx,
		gentxInfo.PubKey,
		gentxInfo.SelfDelegation,
		peer,
	)

	n.ev.Send(events.New(events.StatusOngoing, "Broadcasting validator transaction"))

	res, err := n.cosmos.BroadcastTx(n.account.Name, msg)
	if err != nil {
		return err
	}

	var requestRes launchtypes.MsgRequestAddValidatorResponse
	if err := res.Decode(&requestRes); err != nil {
		return err
	}

	if requestRes.AutoApproved {
		n.ev.Send(events.New(events.StatusDone, "Validator added to the network by the coordinator!"))
	} else {
		n.ev.Send(events.New(events.StatusDone,
			fmt.Sprintf("Request %d to join the network as a validator has been submitted!",
				requestRes.RequestID),
		))
	}
	return nil
}
