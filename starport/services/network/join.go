package network

import (
	"context"
	"errors"

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

func (b *Blockchain) CheckRequestAccount(ctx context.Context, launchID uint64, addr string) (bool, error) {
	requests, err := b.builder.fetchRequests(ctx, launchID)
	if err != nil {
		return false, err
	}

	for _, request := range requests {
		switch req := request.Content.Content.(type) {
		case *launchtypes.RequestContent_GenesisAccount:
			if req.GenesisAccount.Address == addr {
				return true, nil
			}
		case *launchtypes.RequestContent_GenesisValidator:
			if req.GenesisValidator.Address == addr {
				return true, nil
			}
		}
		request.Content.GetGenesisAccount()
	}
	return false, nil
}

func (b *Blockchain) Join(launchID uint64, valAddress, peer string, gentx, consPubKey []byte, selfDelegation sdk.Coin) (string, error) {
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

// fetchRequests fetches the chain requests from Starport Network from a launch id
func (b *Builder) fetchRequests(ctx context.Context, launchID uint64) ([]launchtypes.Request, error) {
	res, err := launchtypes.NewQueryClient(b.cosmos.Context).RequestAll(ctx, &launchtypes.QueryAllRequestRequest{
		LaunchID: launchID,
	})
	if err != nil {
		return nil, err
	}
	return res.Request, err
}
