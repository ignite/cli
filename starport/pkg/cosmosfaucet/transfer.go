package cosmosfaucet

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
	"github.com/tendermint/starport/starport/pkg/cosmoscoin"
)

// TotalTransferredAmount returns the total transferred amount from faucet account to toAccountAddress.
func (f Faucet) TotalTransferredAmount(ctx context.Context, toAccountAddress, denom string) (amount uint64, err error) {
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
					if !strings.HasSuffix(attr.Value, denom) {
						continue
					}

					if time.Since(event.Time) < f.limitRefreshWindow {
						amountStr := strings.TrimRight(attr.Value, denom)
						if a, err := strconv.ParseUint(amountStr, 10, 64); err == nil {
							amount += a
						}
					}
				}
			}
		}
	}

	return amount, nil
}

// Transfer transfer amount of tokens from the faucet account to toAccountAddress.
func (f *Faucet) Transfer(ctx context.Context, toAccountAddress string, coins []cosmoscoin.Coin) error {
	f.transferMutex.Lock()
	defer f.transferMutex.Unlock()

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
	return f.runner.BankSend(ctx, fromAccount.Address, toAccountAddress, strings.Join(coinsStr, ","))
}
