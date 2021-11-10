package network

import (
	"context"
	"errors"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
)

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

// CheckRequestAccount check if the add account request already exist
func (b *Blockchain) CheckRequestAccount(ctx context.Context, launchID uint64, addr string) (bool, error) {
	requests, err := b.builder.fetchRequests(ctx, launchID)
	if err != nil {
		return false, err
	}
	for _, request := range requests {
		genesisAcc := request.Content.GetGenesisAccount()
		if genesisAcc == nil {
			continue
		}
		if genesisAcc.Address == addr {
			return true, nil
		}
	}
	return false, nil
}

// checkRequestValidator check if the add validator request already exist
func (b *Blockchain) checkRequestValidator(ctx context.Context, launchID uint64, addr string) (bool, uint64, error) {
	requests, err := b.builder.fetchRequests(ctx, launchID)
	if err != nil {
		return false, 0, err
	}
	for _, request := range requests {
		genesisVal := request.Content.GetGenesisValidator()
		if genesisVal == nil {
			continue
		}
		if genesisVal.Address == addr {
			return true, genesisVal.LaunchID, nil
		}
	}
	return false, 0, nil
}

// Join creates the RequestAddValidator message into the SPN
func (b *Blockchain) Join(
	ctx context.Context,
	launchID uint64,
	valAddress, peer string,
	gentx, consPubKey []byte,
	selfDelegation sdk.Coin,
) (string, error) {
	// Check if the validator request already exist
	exist, launchID, err := b.checkRequestValidator(ctx, launchID, valAddress)
	if err != nil {
		return "", err
	}
	if exist {
		return strconv.Itoa(int(launchID)), nil
	}

	msgCreateChain := launchtypes.NewMsgRequestAddValidator(
		valAddress,
		launchID,
		gentx,
		consPubKey,
		selfDelegation,
		peer,
	)

	response, err := b.builder.cosmos.BroadcastTx(b.builder.account.Name, msgCreateChain)
	if err != nil {
		return "", err
	}

	out, err := b.builder.cosmos.Context.Codec.MarshalJSON(response)
	if err != nil {
		return "", err
	}

	return string(out), err
}

// CreateAccount creates an add AddAccount request
func (b *Blockchain) CreateAccount(launchID uint64, coins sdk.Coins) (string, error) {
	msgCreateChain := launchtypes.NewMsgRequestAddAccount(
		b.builder.account.Address(SPNAddressPrefix),
		launchID,
		coins,
	)

	response, err := b.builder.cosmos.BroadcastTx(b.builder.account.Name, msgCreateChain)
	if err != nil {
		return "", err
	}

	out, err := b.builder.cosmos.Context.Codec.MarshalJSON(response)
	if err != nil {
		return "", err
	}

	return string(out), err
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
