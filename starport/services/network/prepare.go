package network

import (
	"context"
	launchtypes "github.com/tendermint/spn/x/launch/types"
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
	// TODO: detect the correct prefix
	_ := "cosmos"

	return nil
}