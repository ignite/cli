package cosmosfaucet

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
	"sync"
	"time"

	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
	"github.com/tendermint/starport/starport/pkg/cosmoscoin"
)

// transferMutex is a mutex used for transfer request
var transferMutex = &sync.Mutex{}

// TotalTransferredAmount returns the total transferred amount from faucet account to toAccountAddress.
func (f Faucet) TotalTransferredAmount(ctx context.Context, toAccountAddress, denom string) (totalAmount uint64, err error) {
	fromAccount, err := f.runner.ShowAccount(ctx, f.accountName)
	if err != nil {
		return 0, err
	}

	events, err := f.runner.QueryTxEvents(ctx,
		chaincmdrunner.NewEventSelector("message", "sender", fromAccount.Address),
		chaincmdrunner.NewEventSelector("transfer", "recipient", toAccountAddress))
	if err != nil {
		return 0, err
	}

	for _, event := range events {
		if event.Type == "transfer" {
			for _, attr := range event.Attributes {
				if attr.Key == "amount" {
					coins, err := sdk.ParseCoinsNormalized(attr.Value)
					if err != nil {
						return 0, err
					}

					amount := coins.AmountOf(denom).Uint64()

					if amount > 0 && time.Since(event.Time) < f.limitRefreshWindow {
						totalAmount += amount
					}
				}
			}
		}
	}

	return totalAmount, nil
}

// Transfer transfer amount of tokens from the faucet account to toAccountAddress.
func (f* Faucet) Transfer(ctx context.Context, toAccountAddress string, coins []cosmoscoin.Coin) error {
	transferMutex.Lock()
	defer transferMutex.Unlock()

	var coinsStr []string

	// check for each coin, the max transferred amount hasn't been reached
	for _, c := range coins {
		totalSent, err := f.TotalTransferredAmount(ctx, toAccountAddress, c.Denom)
		if err != nil {
			return err
		}

		if f.coinsMax[c.Denom] != 0 {
			if totalSent >= f.coinsMax[c.Denom] {
				return fmt.Errorf("account has reached maximum credit allowed per account (%d)", f.coinsMax[c.Denom])
			}

			if (totalSent + c.Amount) >= f.coinsMax[c.Denom] {
				return fmt.Errorf("account is about to reach maximum credit allowed per account. it can only receive up to (%d) in total", f.coinsMax[c.Denom])
			}
		}

		coinsStr = append(coinsStr, c.String())
	}

	// perform transfer for all coins
	fromAccount, err := f.runner.ShowAccount(ctx, f.accountName)
	if err != nil {
		return err
	}
	txHash, err := f.runner.BankSend(ctx, fromAccount.Address, toAccountAddress, strings.Join(coinsStr, ","))
	if err != nil {
		return err
	}




	return nil
}
