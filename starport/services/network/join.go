package network

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/cosmoserror"
	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/xurl"
	"github.com/tendermint/starport/starport/services/network/networkchain"
	"github.com/tendermint/starport/starport/services/network/networktypes"
)

type joinOptions struct {
	accountAmount sdk.Coins
}

type JoinOption func(*joinOptions)

func WithAccountRequest(amount sdk.Coins) JoinOption {
	return func(o *joinOptions) {
		o.accountAmount = amount
	}
}

// Join to the network.
func (n Network) Join(
	ctx context.Context,
	c Chain,
	launchID uint64,
	publicAddress,
	gentxPath string,
	options ...JoinOption,
) error {
	o := joinOptions{}
	for _, apply := range options {
		apply(&o)
	}

	nodeID, err := c.NodeID(ctx)
	if err != nil {
		return err
	}

	var peer launchtypes.Peer
	if xurl.IsHTTP(publicAddress) {
		peer = launchtypes.NewPeerTunnel(nodeID, networkchain.HTTPTunnelChisel, publicAddress)
	} else {
		peer = launchtypes.NewPeerConn(nodeID, publicAddress)

	}

	isCustomGentx := gentxPath != ""

	// if the custom gentx is not provided, get the chain default from the chain home folder.
	if !isCustomGentx {
		gentxPath, err = c.DefaultGentxPath()
		if err != nil {
			return err
		}
	}

	// get the chain genesis path from the home folder
	genesisPath, err := c.GenesisPath()
	if err != nil {
		return err
	}

	// parse the gentx content
	gentxInfo, gentx, err := cosmosutil.GentxFromPath(gentxPath)
	if err != nil {
		return err
	}

	// change the chain address prefix to spn
	accountAddress, err := cosmosutil.ChangeAddressPrefix(gentxInfo.DelegatorAddress, networktypes.SPN)
	if err != nil {
		return err
	}

	if !o.accountAmount.IsZero() {
		if err := n.sendAccountRequest(
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

// sendAccountRequest creates an add AddAccount request message.
func (n Network) sendAccountRequest(
	ctx context.Context,
	genesisPath string,
	isCustomGentx bool,
	launchID uint64,
	accountAddress string,
	amount sdk.Coins,
) (err error) {
	address := n.account.Address(networktypes.SPN)
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

	msg := launchtypes.NewMsgRequestAddAccount(
		n.account.Address(networktypes.SPN),
		launchID,
		accountAddress,
		amount,
	)

	n.ev.Send(events.New(events.StatusOngoing, "Broadcasting account transactions"))
	res, err := n.cosmos.BroadcastTx(n.account.Name, msg)
	if err != nil {
		return err
	}

	var requestRes launchtypes.MsgRequestAddAccountResponse
	if err := res.Decode(&requestRes); err != nil {
		return err
	}

	if requestRes.AutoApproved {
		n.ev.Send(events.New(events.StatusDone, "Account added to the network by the coordinator!"))
	} else {
		n.ev.Send(events.New(events.StatusDone,
			fmt.Sprintf("Request %d to add account to the network has been submitted!",
				requestRes.RequestID),
		))
	}
	return nil
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
