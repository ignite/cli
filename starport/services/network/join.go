package network

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/cosmosaddress"
	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"github.com/tendermint/starport/starport/pkg/events"
)

func (b *Blockchain) Join(
	ctx context.Context,
	launchID uint64,
	amount sdk.Coin,
	publicAddress string,
	gentxPath string) error {

	// Get the chain node id to build the peer string `<nodeID>@<host>`
	cmd, err := b.chain.Commands(ctx)
	if err != nil {
		return err
	}
	nodeID, err := cmd.ShowNodeID(ctx)
	if err != nil {
		return err
	}
	peer := fmt.Sprintf("%s@%s", nodeID, publicAddress)

	// Check if a custom gentx is provided
	isCustomGentx := gentxPath != ""
	// If the custom gentx is not provided, get the chain
	// default from the chain home folder
	if !isCustomGentx {
		gentxPath, err = b.chain.DefaultGentxPath()
		if err != nil {
			return err
		}
	}

	// Get the chain genesis path from the home folder
	genesisPath, err := b.chain.GenesisPath()
	if err != nil {
		return err
	}

	if err := b.builder.sendAccountRequest(ctx,
		genesisPath,
		isCustomGentx,
		launchID,
		amount); err != nil {
		return err
	}
	return b.builder.sendValidatorRequest(ctx, launchID, peer, gentxPath)
}

// sendAccountRequest creates an add AddAccount request message
func (b *Builder) sendAccountRequest(
	ctx context.Context,
	genesisPath string,
	isCustomGentx bool,
	launchID uint64,
	amount sdk.Coin,
) error {
	address := b.account.Address(SPNAddressPrefix)
	spnAddress, err := cosmosaddress.ChangePrefix(address, SPNAddressPrefix)
	if err != nil {
		return err
	}

	b.ev.Send(events.New(events.StatusOngoing, "Verifying account already exists "+spnAddress))

	// if is custom gentx path, avoid to check account into genesis from the home folder
	accExist := false
	if !isCustomGentx {
		accExist, err = cosmosutil.CheckGenesisContainsAddress(genesisPath, spnAddress)
		if err != nil {
			return err
		}
	}
	// check if account exists as a genesis account in SPN chain launch information
	if !accExist && !b.hasAccount(ctx, launchID, spnAddress) {
		msg := launchtypes.NewMsgRequestAddAccount(
			spnAddress,
			launchID,
			sdk.NewCoins(amount),
		)

		b.ev.Send(events.New(events.StatusOngoing, "Broadcasting account transactions"))
		res, err := b.cosmos.BroadcastTx(b.account.Name, msg)
		if err != nil {
			return err
		}

		var requestRes launchtypes.MsgRequestAddAccountResponse
		if err := res.Decode(&requestRes); err != nil {
			return err
		}
		b.ev.Send(events.New(events.StatusDone, "MsgRequestAddAccount transactions sent"))

		if requestRes.AutoApproved {
			b.ev.Send(events.New(events.StatusDone, "Account added to the network by the coordinator!"))
		} else {
			b.ev.Send(events.New(events.StatusDone,
				fmt.Sprintf("Request %d to add account to the network has been submitted!",
					requestRes.RequestID),
			))
		}
		return nil
	}

	b.ev.Send(events.New(events.StatusDone, "Account already exist"))
	return nil
}

// sendValidatorRequest creates the RequestAddValidator message into the SPN
func (b *Builder) sendValidatorRequest(
	ctx context.Context,
	launchID uint64,
	peer string,
	gentxPath string,
) error {
	// Parse the gentx content
	gentxInfo, gentx, err := cosmosutil.ParseGentx(gentxPath)
	if err != nil {
		return err
	}

	// Change the chain address prefix to spn
	spnValAddress, err := cosmosaddress.ChangePrefix(gentxInfo.DelegatorAddress, SPNAddressPrefix)
	if err != nil {
		return err
	}

	// Check if the validator request already exist
	if b.hasValidator(ctx, launchID, spnValAddress) {
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

	b.ev.Send(events.New(events.StatusOngoing, "Broadcasting validator transaction"))
	res, err := b.cosmos.BroadcastTx(b.account.Name, msg)
	if err != nil {
		return err
	}

	var requestRes launchtypes.MsgRequestAddValidatorResponse
	if err := res.Decode(&requestRes); err != nil {
		return err
	}
	b.ev.Send(events.New(events.StatusDone, "MsgRequestAddValidator transaction sent"))

	if requestRes.AutoApproved {
		b.ev.Send(events.New(events.StatusDone, "Validator added to the network by the coordinator!"))
	} else {
		b.ev.Send(events.New(events.StatusDone,
			fmt.Sprintf("Request %d to join the network as a validator has been submitted!",
				requestRes.RequestID),
		))
	}
	return nil
}

// hasValidator verify if the validator already exist into the SPN store
func (b *Builder) hasValidator(ctx context.Context, launchID uint64, address string) bool {
	_, err := launchtypes.NewQueryClient(b.cosmos.Context).GenesisValidator(ctx, &launchtypes.QueryGetGenesisValidatorRequest{
		LaunchID: launchID,
		Address:  address,
	})
	return err == nil
}

// hasAccount verify if the account already exist into the SPN store
func (b *Builder) hasAccount(ctx context.Context, launchID uint64, address string) bool {
	_, err := launchtypes.NewQueryClient(b.cosmos.Context).VestingAccount(ctx, &launchtypes.QueryGetVestingAccountRequest{
		LaunchID: launchID,
		Address:  address,
	})
	if err == nil {
		return true
	}
	_, err = launchtypes.NewQueryClient(b.cosmos.Context).GenesisAccount(ctx, &launchtypes.QueryGetGenesisAccountRequest{
		LaunchID: launchID,
		Address:  address,
	})
	return err == nil
}
