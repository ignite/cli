package network

import (
	"context"
	"fmt"

	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/cosmosaddress"
)

// Prepare queries launch information and prepare the chain to be launched from these information
func (b Blockchain) Prepare(ctx context.Context) error {
	// Get the genesis accounts and apply them to the genesis
	// TODO: include launchID in the init process
	genesisAccounts, err := b.builder.GenesisAccounts(ctx, b.launchID)
	if err != nil {
		return err
	}
	if err := b.applyGenesisAccounts(ctx, genesisAccounts); err != nil {
		return err
	}

	// Get the genesis vesting accounts and apply them to the genesis
	vestingAccounts, err := b.builder.VestingAccounts(ctx, b.launchID)
	if err != nil {
		return err
	}
	if err := b.applyVestingAccounts(ctx, vestingAccounts); err != nil {
		return err
	}

	// Get the genesis validators, gather gentxs and modify config to include the peers

	return nil
}

// applyGenesisAccounts adds the genesis account into the genesis using the chain CLI
func (b Blockchain) applyGenesisAccounts(ctx context.Context, genesisAccs []launchtypes.GenesisAccount) error {
	var err error
	// TODO: detect the correct prefix
	prefix := "cosmos"

	cmd, err := b.chain.Commands(ctx)
	if err != nil {
		return err
	}

	for _, acc := range genesisAccs {
		// change the address prefix to the target chain prefix
		acc.Address, err = cosmosaddress.ChangePrefix(acc.Address, prefix)
		if err != nil {
			return err
		}

		// call add genesis account cli command
		err = cmd.AddGenesisAccount(ctx, acc.Address, acc.Coins.String())
		if err != nil {
			return err
		}
	}

	return nil
}

// applyVestingAccounts adds the genesis vesting account into the genesis using the chain CLI
func (b Blockchain) applyVestingAccounts(ctx context.Context, vestingAccs []launchtypes.VestingAccount) error {
	var err error
	prefix := "cosmos"

	cmd, err := b.chain.Commands(ctx)
	if err != nil {
		return err
	}

	for _, acc := range vestingAccs {
		acc.Address, err = cosmosaddress.ChangePrefix(acc.Address, prefix)
		if err != nil {
			return err
		}

		// Only delayed vesting option is supported for now
		delayedVesting := acc.VestingOptions.GetDelayedVesting()
		if delayedVesting == nil {
			return fmt.Errorf("invalid vesting option for account %s", acc.Address)
		}

		// call add genesis account cli command with delayed vesting option
		err = cmd.AddVestingAccount(
			ctx,
			acc.Address,
			acc.StartingBalance.String(),
			delayedVesting.Vesting.String(),
			delayedVesting.EndTime,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
