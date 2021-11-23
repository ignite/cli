package network

import (
	"context"
	"fmt"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/cosmosaddress"
	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Prepare queries launch information and prepare the chain to be launched from these information
func (b Blockchain) Prepare(ctx context.Context) error {
	if !b.isInitialized {
		return errors.New("the blockchain must be initialized to prepare for launch")
	}

	// get the genesis accounts and apply them to the genesis
	genesisAccounts, err := b.builder.GenesisAccounts(ctx, b.launchID)
	if err != nil {
		return errors.Wrap(err, "error querying genesis accounts")
	}
	if err := b.applyGenesisAccounts(ctx, genesisAccounts); err != nil {
		return errors.Wrap(err, "error applying genesis accounts to genesis")
	}

	// get the genesis vesting accounts and apply them to the genesis
	vestingAccounts, err := b.builder.VestingAccounts(ctx, b.launchID)
	if err != nil {
		return errors.Wrap(err, "error querying vesting accounts")
	}
	if err := b.applyVestingAccounts(ctx, vestingAccounts); err != nil {
		return errors.Wrap(err, "error applying vesting accounts to genesis")
	}

	// get the genesis validators, gather gentxs and modify config to include the peers
	genesisValidators, err := b.builder.GenesisValidators(ctx, b.launchID)
	if err != nil {
		return errors.Wrap(err, "error querying genesis validators")
	}
	if err := b.applyGenesisValidators(ctx, genesisValidators); err != nil {
		return errors.Wrap(err, "error applying genesis validators to genesis")
	}

	// set the genesis time for the chain
	genesisPath, err := b.chain.GenesisPath()
	if err != nil {
		return errors.Wrap(err, "genesis of the blockchain can't be read")
	}
	return cosmosutil.SetGenesisTime(genesisPath, b.launchTime)
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

		// only delayed vesting option is supported for now
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

// applyGenesisValidators gathers the validator gentxs into the genesis and add peers in config
func (b Blockchain) applyGenesisValidators(ctx context.Context, genesisVals []launchtypes.GenesisValidator) error {
	// no validator
	if len(genesisVals) == 0 {
		return nil
	}

	// reset the gentx directory
	gentxDir, err := b.chain.GentxPath()
	if err != nil {
		return err
	}
	if err := os.RemoveAll(gentxDir); err != nil {
		return err
	}
	if err := os.MkdirAll(gentxDir, 0700); err != nil {
		return err
	}

	// write gentxs
	for i, val := range genesisVals {
		gentxPath := filepath.Join(gentxDir, fmt.Sprintf("gentx%d.json", i))
		if err = ioutil.WriteFile(gentxPath, val.GenTx, 0666); err != nil {
			return err
		}
	}

	// gather gentxs
	cmd, err := b.chain.Commands(ctx)
	if err != nil {
		return err
	}
	if err := cmd.CollectGentxs(ctx); err != nil {
		return err
	}

	return b.updateConfigFromGenesisValidators(genesisVals)
}

// updateConfigFromGenesisValidators adds the peer addresses into the config.toml of the chain
func (b Blockchain) updateConfigFromGenesisValidators(genesisVals []launchtypes.GenesisValidator) error {
	var p2pAddresses []string
	for _, val := range genesisVals {
		p2pAddresses = append(p2pAddresses, val.Peer)
	}

	// set persistent peers
	configPath, err := b.chain.ConfigTOMLPath()
	if err != nil {
		return err
	}
	configToml, err := toml.LoadFile(configPath)
	if err != nil {
		return err
	}
	configToml.Set("p2p.persistent_peers", strings.Join(p2pAddresses, ","))
	if err != nil {
		return err
	}

	// save config.toml file
	configTomlFile, err := os.OpenFile(configPath, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer configTomlFile.Close()
	_, err = configToml.WriteTo(configTomlFile)
	return err
}
