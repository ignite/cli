package cosmosfaucet

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	chaincmdrunner "github.com/ignite/cli/ignite/pkg/chaincmd/runner"
)

// transferMutex is a mutex used for keeping transfer requests in a queue so checking account balance and sending tokens is atomic.
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

// Transfer transfers amount of tokens from the faucet account to toAccountAddress.
func (f *Faucet) Transfer(ctx context.Context, toAccountAddress string, coins sdk.Coins) error {
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
				return fmt.Errorf(
					"account has reached to the max. allowed amount (%d) for %q denom",
					f.coinsMax[c.Denom],
					c.Denom,
				)
			}

			if (totalSent + c.Amount.Uint64()) > f.coinsMax[c.Denom] {
				return fmt.Errorf(
					`ask less amount for %q denom. account is reaching to the limit (%d) that faucet can tolerate`,
					c.Denom,
					f.coinsMax[c.Denom],
				)
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

	// wait for send tx to be confirmed
	return f.runner.WaitTx(ctx, txHash, time.Second, 30)
}
