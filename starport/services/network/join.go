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

// Join to the network.
func (n Network) Join(
	ctx context.Context,
	c Chain,
	launchID uint64,
	amount sdk.Coin,
	publicAddress string,
	gentxPath string,
) error {
	peer, err := c.Peer(ctx, publicAddress)
	if err != nil {
		return err
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

	if err := n.sendAccountRequest(
		ctx,
		genesisPath,
		isCustomGentx,
		launchID,
		amount,
	); err != nil {
		return err
	}

	return n.sendValidatorRequest(ctx, launchID, peer, gentxPath)
}

// sendAccountRequest creates an add AddAccount request message.
func (n Network) sendAccountRequest(
	ctx context.Context,
	genesisPath string,
	isCustomGentx bool,
	launchID uint64,
	amount sdk.Coin,
) error {
	address := n.account.Address(networkchain.SPN)
	spnAddress, err := cosmosutil.ChangeAddressPrefix(address, networkchain.SPN)
	if err != nil {
		return err
	}

	n.ev.Send(events.New(events.StatusOngoing, "Verifying account already exists "+spnAddress))

	// if is custom gentx path, avoid to check account into genesis from the home folder
	var accExist bool

	if !isCustomGentx {
		accExist, err = cosmosutil.CheckGenesisContainsAddress(genesisPath, spnAddress)
		if err != nil {
			return err
		}
	}
	// check if account exists as a genesis account in SPN chain launch information
	if !accExist && !n.hasAccount(ctx, launchID, spnAddress) {
		msg := launchtypes.NewMsgRequestAddAccount(
			spnAddress,
			launchID,
			sdk.NewCoins(amount),
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

	n.ev.Send(events.New(events.StatusDone, "Account already exist"))
	return nil
}

// sendValidatorRequest creates the RequestAddValidator message into the SPN
func (n Network) sendValidatorRequest(
	ctx context.Context,
	launchID uint64,
	peer string,
	gentxPath string,
) error {
	// Parse the gentx content
	gentxInfo, gentx, err := cosmosutil.GentxFromPath(gentxPath)
	if err != nil {
		return err
	}

	// Change the chain address prefix to spn
	spnValAddress, err := cosmosutil.ChangeAddressPrefix(gentxInfo.DelegatorAddress, networkchain.SPN)
	if err != nil {
		return err
	}

	// Check if the validator request already exist
	if n.hasValidator(ctx, launchID, spnValAddress) {
		return fmt.Errorf("validator %s already exist", spnValAddress)
	}

	msg := launchtypes.NewMsgRequestAddValidator(
		spnValAddress,
		launchID,
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
func (n Network) hasValidator(ctx context.Context, launchID uint64, address string) bool {
	_, err := launchtypes.NewQueryClient(n.cosmos.Context).GenesisValidator(ctx, &launchtypes.QueryGetGenesisValidatorRequest{
		LaunchID: launchID,
		Address:  address,
	})
	return err == nil
}

// hasAccount verify if the account already exist into the SPN store
func (n Network) hasAccount(ctx context.Context, launchID uint64, address string) bool {
	_, err := launchtypes.NewQueryClient(n.cosmos.Context).VestingAccount(ctx, &launchtypes.QueryGetVestingAccountRequest{
		LaunchID: launchID,
		Address:  address,
	})
	if err == nil {
		return true
	}
	_, err = launchtypes.NewQueryClient(n.cosmos.Context).GenesisAccount(ctx, &launchtypes.QueryGetGenesisAccountRequest{
		LaunchID: launchID,
		Address:  address,
	})
	return err == nil
}
