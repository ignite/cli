package network

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/events"
)

// Join creates the RequestAddValidator message into the SPN
func (b *Builder) Join(
	ctx context.Context,
	launchID uint64,
	peer,
	valKeyName,
	valAddress string,
	gentx,
	consPubKey []byte,
	selfDelegation sdk.Coin,
) (uint64, error) {
	// Change the chain address prefix to spn
	spnValAddress, err := SetSPNPrefix(valAddress)
	if err != nil {
		return 0, err
	}

	// Check if the validator request already exist
	exist, err := b.CheckValidatorExist(ctx, launchID, spnValAddress)
	if err != nil {
		return 0, err
	}
	if exist {
		return 0, errors.New("validator already exist: " + spnValAddress)
	}

	msg := launchtypes.NewMsgRequestAddValidator(
		spnValAddress,
		launchID,
		gentx,
		consPubKey,
		selfDelegation,
		peer,
	)

	b.ev.Send(events.New(events.StatusOngoing, "Broadcasting validator transaction"))
	res, err := b.cosmos.BroadcastTx(valKeyName, msg)
	if err != nil {
		return 0, err
	}

	var requestRes launchtypes.MsgRequestResponse
	if err := res.Decode(&requestRes); err != nil {
		return 0, err
	}
	b.ev.Send(events.New(events.StatusDone, "AddValidator transaction sent"))

	return requestRes.RequestID, err
}

// CreateAccountRequestMsg creates an add AddAccount request message
func (b *Builder) CreateAccountRequestMsg(
	ctx context.Context,
	chainHome string,
	customGentx bool,
	launchID uint64,
	amount sdk.Coin,
) (uint64, error) {
	address := b.account.Address(SPNAddressPrefix)
	spnAddress, err := SetSPNPrefix(address)
	if err != nil {
		return 0, err
	}

	b.ev.Send(events.New(events.StatusOngoing, "Verifying account already exists "+spnAddress))

	shouldCreateAcc := false
	if !customGentx {
		exist, err := CheckGenesisAddress(chainHome, spnAddress)
		if err != nil {
			return 0, err
		}
		if !exist {
			exist, err = b.CheckAccountExist(ctx, launchID, spnAddress)
			if err != nil {
				return 0, err
			}
		}
		shouldCreateAcc = !exist
	}
	if shouldCreateAcc || customGentx {
		msg := launchtypes.NewMsgRequestAddAccount(
			spnAddress,
			launchID,
			sdk.NewCoins(amount),
		)

		b.ev.Send(events.New(events.StatusOngoing, "Broadcasting account transactions"))
		res, err := b.cosmos.BroadcastTx(b.account.Name, msg)
		if err != nil {
			return 0, err
		}

		var requestRes launchtypes.MsgRequestResponse
		if err := res.Decode(&requestRes); err != nil {
			return 0, err
		}
		b.ev.Send(events.New(events.StatusDone, "AddAccount transactions sent"))

		return requestRes.RequestID, nil
	}

	b.ev.Send(events.New(events.StatusDone, "Account already exist"))
	return 0, nil
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

// CheckAccountExist check if the account already exists or is pending approval
func (b *Builder) CheckAccountExist(ctx context.Context, launchID uint64, address string) (bool, error) {
	if b.hasAccount(ctx, launchID, address) {
		return true, nil
	}
	// verify if the account is pending approval
	requests, err := b.fetchRequests(ctx, launchID)
	if err != nil {
		return false, err
	}
	for _, request := range requests {
		switch req := request.Content.Content.(type) {
		case *launchtypes.RequestContent_GenesisAccount:
			if req.GenesisAccount.Address == address {
				return true, nil
			}
		case *launchtypes.RequestContent_VestingAccount:
			if req.VestingAccount.Address == address {
				return true, nil
			}
		}
	}
	return false, nil
}

// CheckValidatorExist check if the validator already exists or is pending approval
func (b *Builder) CheckValidatorExist(ctx context.Context, launchID uint64, address string) (bool, error) {
	if b.hasValidator(ctx, launchID, address) {
		return true, nil
	}
	// verify if the validator is pending approval
	requests, err := b.fetchRequests(ctx, launchID)
	if err != nil {
		return false, err
	}
	for _, request := range requests {
		genesisVal := request.Content.GetGenesisValidator()
		if genesisVal == nil {
			continue
		}
		if genesisVal.Address == address {
			return true, nil
		}
	}
	return false, nil
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

// fetchRequests fetches the chain requests from SPN by launch id
func (b *Builder) fetchRequests(ctx context.Context, launchID uint64) ([]launchtypes.Request, error) {
	res, err := launchtypes.NewQueryClient(b.cosmos.Context).RequestAll(ctx, &launchtypes.QueryAllRequestRequest{
		LaunchID: launchID,
	})
	if err != nil {
		return nil, err
	}

	return res.Request, err
}
