package network

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"

	"github.com/ignite/cli/ignite/pkg/cosmoserror"
	"github.com/ignite/cli/ignite/pkg/cosmosutil"
	"github.com/ignite/cli/ignite/pkg/events"
	"github.com/ignite/cli/ignite/pkg/xurl"
	"github.com/ignite/cli/ignite/services/network/networkchain"
	"github.com/ignite/cli/ignite/services/network/networktypes"
)

type joinOptions struct {
	accountAmount sdk.Coins
	gentxPath     string
	publicAddress string
}

type JoinOption func(*joinOptions)

func WithAccountRequest(amount sdk.Coins) JoinOption {
	return func(o *joinOptions) {
		o.accountAmount = amount
	}
}

// TODO accept struct not file path
func WithCustomGentxPath(path string) JoinOption {
	return func(o *joinOptions) {
		o.gentxPath = path
	}
}

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
	options ...JoinOption,
) error {
	o := joinOptions{}
	for _, apply := range options {
		apply(&o)
	}

	isCustomGentx := o.gentxPath != ""
	var (
		nodeID string
		peer   launchtypes.Peer
		err    error
	)

	// if the custom gentx is not provided, get the chain default from the chain home folder.
	if !isCustomGentx {
		if nodeID, err = c.NodeID(ctx); err != nil {
			return err
		}

		if xurl.IsHTTP(o.publicAddress) {
			peer = launchtypes.NewPeerTunnel(nodeID, networkchain.HTTPTunnelChisel, o.publicAddress)
		} else {
			peer = launchtypes.NewPeerConn(nodeID, o.publicAddress)

		}

		if o.gentxPath, err = c.DefaultGentxPath(); err != nil {
			return err
		}
	}

	// parse the gentx content
	gentxInfo, gentx, err := cosmosutil.GentxFromPath(o.gentxPath)
	if err != nil {
		return err
	}

	if isCustomGentx {
		if peer, err = ParsePeerAddress(gentxInfo.Memo); err != nil {
			return err
		}
	}

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
		if err := n.ensureAccount(
			ctx,
			genesisPath,
			isCustomGentx,
			launchID,
			accountAddress,
			o.accountAmount,
		); err != nil {
			return err
		}
	}

	return n.sendValidatorRequest(ctx, launchID, peer, accountAddress, gentx, gentxInfo)
}

// ensureAccount creates an add AddAccount request message.
func (n Network) ensureAccount(
	ctx context.Context,
	genesisPath string,
	isCustomGentx bool,
	launchID uint64,
	address string,
	amount sdk.Coins,
) (err error) {
	n.ev.Send(events.New(events.StatusOngoing, "Verifying account already exists "+address))

	// if is custom gentx path, avoid to check account into genesis from the home folder
	var accExist bool
	if !isCustomGentx {
		accExist, err = cosmosutil.CheckGenesisContainsAddress(genesisPath, address)
		if err != nil {
			return err
		}
		if accExist {
			return fmt.Errorf("account %s already exist", address)
		}
	}
	// check if account exists as a genesis account in SPN chain launch information
	hasAccount, err := n.hasAccount(ctx, launchID, address)
	if err != nil {
		return err
	}
	if hasAccount {
		return fmt.Errorf("account %s already exist", address)
	}

	return n.sendAccountRequest(launchID, address, amount)
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
	// Check if the validator request already exist
	hasValidator, err := n.hasValidator(ctx, launchID, valAddress)
	if err != nil {
		return err
	}
	if hasValidator {
		return fmt.Errorf("validator %s already exist", valAddress)
	}

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

// hasValidator verify if the validator already exist into the SPN store
func (n Network) hasValidator(ctx context.Context, launchID uint64, address string) (bool, error) {
	_, err := n.launchQuery.GenesisValidator(ctx, &launchtypes.QueryGetGenesisValidatorRequest{
		LaunchID: launchID,
		Address:  address,
	})
	if cosmoserror.Unwrap(err) == cosmoserror.ErrNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// hasAccount verify if the account already exist into the SPN store
func (n Network) hasAccount(ctx context.Context, launchID uint64, address string) (bool, error) {
	_, err := n.launchQuery.VestingAccount(ctx, &launchtypes.QueryGetVestingAccountRequest{
		LaunchID: launchID,
		Address:  address,
	})
	if cosmoserror.Unwrap(err) == cosmoserror.ErrNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	_, err = n.launchQuery.GenesisAccount(ctx, &launchtypes.QueryGetGenesisAccountRequest{
		LaunchID: launchID,
		Address:  address,
	})
	if cosmoserror.Unwrap(err) == cosmoserror.ErrNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
