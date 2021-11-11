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
	validatorMsg,
	accountMsg sdk.Msg,
) (string, error) {
	msgs := make([]sdk.Msg, 0)
	if accountMsg != nil {
		msgs = append(msgs, accountMsg)
	}
	if validatorMsg != nil {
		msgs = append(msgs, validatorMsg)
	}
	b.ev.Send(events.New(events.StatusOngoing, "Broadcasting transactions"))
	response, err := b.cosmos.BroadcastTx(b.account.Name, msgs...)
	if err != nil {
		return "", err
	}

	out, err := b.cosmos.Context.Codec.MarshalJSON(response)
	if err != nil {
		return "", err
	}
	b.ev.Send(events.New(events.StatusDone, "Transactions broadcasted"))

	return string(out), err
}

// CreateValidatorRequestMsg creates an add AddValidator request message
func (b *Builder) CreateValidatorRequestMsg(
	ctx context.Context,
	launchID uint64,
	peer,
	valAddress string,
	gentx,
	consPubKey []byte,
	selfDelegation sdk.Coin,
) (sdk.Msg, error) {
	// Check if the validator request already exist
	exist, launchID, err := b.checkRequestValidator(ctx, launchID, valAddress)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.New("validator already exist: " + valAddress)
	}

	return launchtypes.NewMsgRequestAddValidator(
		valAddress,
		launchID,
		gentx,
		consPubKey,
		selfDelegation,
		peer,
	), nil
}

// CreateAccountRequestMsg creates an add AddAccount request message
func (b *Builder) CreateAccountRequestMsg(
	ctx context.Context,
	chainHome string,
	launchID uint64,
	amount sdk.Coin,
) (msg sdk.Msg, err error) {
	addr := b.account.Address(SPNAddressPrefix)
	b.ev.Send(events.New(events.StatusOngoing, "Verifying account already exists "+addr))

	shouldCreateAcc := false
	if !amount.IsZero() {
		exist, err := CheckGenesisAddress(chainHome, addr)
		if err != nil {
			return msg, err
		}
		if !exist {
			exist, err = b.CheckRequestAccount(ctx, launchID, addr)
			if err != nil {
				return msg, err
			}
		}
		shouldCreateAcc = !exist
	}
	if shouldCreateAcc {
		b.ev.Send(events.New(events.StatusDone, "Account message created"))
		msg = launchtypes.NewMsgRequestAddAccount(
			addr,
			launchID,
			sdk.NewCoins(amount),
		)
	} else {
		b.ev.Send(events.New(events.StatusDone, "Account message not created"))
	}
	return msg, err

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

// CheckRequestAccount check if the add account request already exist
func (b *Builder) CheckRequestAccount(ctx context.Context, launchID uint64, addr string) (bool, error) {
	requests, err := b.fetchRequests(ctx, launchID)
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
func (b *Builder) checkRequestValidator(ctx context.Context, launchID uint64, addr string) (bool, uint64, error) {
	requests, err := b.fetchRequests(ctx, launchID)
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
