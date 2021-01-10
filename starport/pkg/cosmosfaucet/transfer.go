package cosmosfaucet

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
)

func (f Faucet) TotalTransferredAmount(ctx context.Context, toAccountAddress string) (amount uint64, err error) {
	events, err := f.runner.QueryTxEvents(ctx,
		chaincmdrunner.NewEventSelector("message", "sender", f.accountAddress),
		chaincmdrunner.NewEventSelector("transfer", "recipient", toAccountAddress))
	if err != nil {
		return 0, err
	}

	for _, event := range events {
		if event.Type == "transfer" {
			for _, attr := range event.Attributes {
				if attr.Key == "amount" {
					amountStr := strings.TrimRight(attr.Value, f.denom)
					if a, err := strconv.ParseUint(amountStr, 10, 64); err == nil {
						amount += a
					}
				}
			}
		}
	}

	return amount, nil
}

func (f Faucet) Transfer(ctx context.Context, toAccountAddress, amount string) error {
	totalSent, err := f.TotalTransferredAmount(ctx, toAccountAddress)
	if err != nil {
		return err
	}

	if totalSent >= f.maxCredit {
		return fmt.Errorf("account has reached maximum credit allowed per account (%d)", f.maxCredit)
	}

	return f.runner.BankSend(ctx, f.accountAddress, toAccountAddress, amount)
}
