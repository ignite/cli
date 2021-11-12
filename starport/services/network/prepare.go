package network

import (
	"context"
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

	// Get the genesis validators, gather gentxs and modify config to include the peers
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

	for i := range genesisAccs {
		// change the address prefix to the target chain prefix
		genesisAccs[i].Address, err = cosmosaddress.ChangePrefix(genesisAccs[i].Address, prefix)
		if err != nil {
			return err
		}
		acc := genesisAccs[i]

		// call add genesis account cli command
		err = cmd.AddGenesisAccount(ctx, acc.Address, acc.Coins.String())
		if err != nil {
			return err
		}
	}

	return nil
}

