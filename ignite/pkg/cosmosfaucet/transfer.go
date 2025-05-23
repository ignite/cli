package cosmosfaucet

import (
	"context"
	"sync"
	"time"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ignite/cli/v29/ignite/pkg/chaincmd"
	chaincmdrunner "github.com/ignite/cli/v29/ignite/pkg/chaincmd/runner"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// transferMutex is a mutex used for keeping transfer requests in a queue so checking account balance and sending tokens is atomic.
var transferMutex = &sync.Mutex{}

// TotalTransferredAmount returns the total transferred amount from faucet account to toAccountAddress.
func (f Faucet) TotalTransferredAmount(ctx context.Context, toAccountAddress, denom string) (totalAmount sdkmath.Int, err error) {
	fromAccount, err := f.runner.ShowAccount(ctx, f.accountName)
	if err != nil {
		return sdkmath.NewInt(0), err
	}

	opts := []chaincmdrunner.EventSelector{
		chaincmdrunner.NewEventSelector("message", "sender", fromAccount.Address),
		chaincmdrunner.NewEventSelector("transfer", "recipient", toAccountAddress),
	}

	var events []chaincmdrunner.Event
	if f.version.GTE(cosmosver.StargateFiftyVersion) {
		events, err = f.runner.QueryTxByQuery(ctx, opts...)
		if err != nil {
			return sdkmath.NewInt(0), err
		}
	} else {
		events, err = f.runner.QueryTxByEvents(ctx, opts...)
		if err != nil {
			return sdkmath.NewInt(0), err
		}
	}

	totalAmount = sdkmath.NewInt(0)
	for _, event := range events {
		if event.Type == "transfer" {
			for _, attr := range event.Attributes {
				if attr.Key == "amount" {
					coins, err := sdk.ParseCoinsNormalized(attr.Value)
					if err != nil {
						return sdkmath.NewInt(0), err
					}

					amount := coins.AmountOf(denom)
					if amount.GT(sdkmath.NewInt(0)) && time.Since(event.Time) < f.limitRefreshWindow {
						totalAmount = totalAmount.Add(amount)
					}
				}
			}
		}
	}

	return totalAmount, nil
}

// Transfer transfers amount of tokens from the faucet account to toAccountAddress.
func (f *Faucet) Transfer(ctx context.Context, toAccountAddress string, coins sdk.Coins) (string, error) {
	transferMutex.Lock()
	defer transferMutex.Unlock()

	transfer := sdk.NewCoins()
	// check for each coin, the max transferred amount hasn't been reached
	for _, c := range coins {
		if f.indexerDisabled { // we cannot check the transfer history if indexer is disabled
			transfer = transfer.Add(c)
			continue
		}

		totalSent, err := f.TotalTransferredAmount(ctx, toAccountAddress, c.Denom)
		if err != nil {
			return "", err
		}
		coinMax, found := f.coinsMax[c.Denom]
		if found && !coinMax.IsNil() && !coinMax.Equal(sdkmath.NewInt(0)) {
			if totalSent.GTE(coinMax) {
				return "", errors.Errorf(
					"account has reached to the max. allowed amount (%d) for %q denom",
					coinMax,
					c.Denom,
				)
			}

			if (totalSent.Add(c.Amount)).GT(coinMax) {
				return "", errors.Errorf(
					`ask less amount for %q denom. account is reaching to the limit (%d) that faucet can tolerate`,
					c.Denom,
					coinMax,
				)
			}
		}

		transfer = transfer.Add(c)
	}

	// perform transfer for all coins
	fromAccount, err := f.runner.ShowAccount(ctx, f.accountName)
	if err != nil {
		return "", err
	}
	txHash, err := f.runner.BankSend(ctx, fromAccount.Address, toAccountAddress, transfer.String(), chaincmd.BankSendWithFees(f.feeAmount))
	if err != nil {
		return "", err
	}

	if f.indexerDisabled {
		return txHash, nil // we cannot check the tx status if indexer is disabled
	}

	// wait for send tx to be confirmed
	return txHash, f.runner.WaitTx(ctx, txHash, time.Second, 30)
}
