package network

import (
	"context"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/gentx"
)

func (b *Builder) Join(
	ctx context.Context,
	chainHome string,
	launchID uint64,
	customGentx bool,
	amount sdk.Coin,
	peer,
	valKeyName string,
	gentx []byte,
	gentxInfo gentx.Info) error {
	if err := b.SendAccountRequestMsg(ctx,
		chainHome,
		customGentx,
		launchID,
		amount); err != nil {
		return err
	}
	return b.SendValidatorRequestMsg(ctx, launchID, peer, valKeyName, gentx, gentxInfo)
}

// SendAccountRequestMsg creates an add AddAccount request message
func (b *Builder) SendAccountRequestMsg(
	ctx context.Context,
	chainHome string,
	isCustomGentx bool,
	launchID uint64,
	amount sdk.Coin,
) error {
	address := b.account.Address(SPNAddressPrefix)
	spnAddress, err := SetSPNPrefix(address)
	if err != nil {
		return err
	}

	b.ev.Send(events.New(events.StatusOngoing, "Verifying account already exists "+spnAddress))

	// if is custom gentx path, avoid to check account into genesis from the home folder
	accExist := false
	if !isCustomGentx {
		accExist, err = CheckGenesisContainsAddress(chainHome, spnAddress)
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

		var requestRes launchtypes.MsgRequestResponse
		if err := res.Decode(&requestRes); err != nil {
			return err
		}
		b.ev.Send(events.New(events.StatusDone, "AddAccount transactions sent"))

		if requestRes.AutoApproved {
			b.ev.Send(events.New(events.StatusDone, "Account added to the network by the coordinator!"))
		} else {
			b.ev.Send(events.New(events.StatusDone,
				fmt.Sprintf("%s Request %d to add account to the network has been submitted!\n",
					clispinner.OK, requestRes.RequestID),
			))
		}
		return nil
	}

	b.ev.Send(events.New(events.StatusDone, "Account already exist"))
	return err
}

// SendValidatorRequestMsg creates the RequestAddValidator message into the SPN
func (b *Builder) SendValidatorRequestMsg(
	ctx context.Context,
	launchID uint64,
	peer,
	valKeyName string,
	gentx []byte,
	gentxInfo gentx.Info,
) error {
	// Change the chain address prefix to spn
	spnValAddress, err := SetSPNPrefix(gentxInfo.DelegatorAddress)
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
	res, err := b.cosmos.BroadcastTx(valKeyName, msg)
	if err != nil {
		return err
	}

	var requestRes launchtypes.MsgRequestResponse
	if err := res.Decode(&requestRes); err != nil {
		return err
	}
	b.ev.Send(events.New(events.StatusDone, "MsgRequestAddValidator transaction sent"))

	if requestRes.AutoApproved {
		b.ev.Send(events.New(events.StatusDone, "Validator added to the network by the coordinator!\n"))
	} else {
		b.ev.Send(events.New(events.StatusDone,
			fmt.Sprintf("Request %d to join the network as a validator has been submitted!\n",
				requestRes.RequestID),
		))
	}

	return nil
}

// GetAccountAddress return an account address for the blockchain by name
func (b *Blockchain) GetAccountAddress(ctx context.Context, accountName string) (string, error) {
	if !b.isInitialized {
		return "", errors.New("the blockchain must be initialized to show an account")
	}

	chainCmd, err := b.chain.Commands(ctx)
	if err != nil {
		return "", err
	}
	acc, err := chainCmd.ShowAccount(ctx, accountName)
	if err != nil {
		return "", err
	}
	return acc.Address, nil
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
