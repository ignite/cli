package networkchain

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network/networktypes"
)

// Prepare prepares the chain to be launched from genesis information
func (c Chain) Prepare(ctx context.Context, gi networktypes.GenesisInformation) error {
	// chain initialization
	chainHome, err := c.chain.Home()
	if err != nil {
		return err
	}

	_, err = os.Stat(chainHome)

	switch {
	case os.IsNotExist(err):
		// if no config exists, perform a full initialization of the chain with a new validator key
		if err := c.Init(ctx); err != nil {
			return err
		}
	case err != nil:
		return err
	default:
		// if config and validator key already exists, build the chain and initialize the genesis
		c.ev.Send(events.New(events.StatusOngoing, "Building the blockchain"))
		if _, err := c.chain.Build(ctx, ""); err != nil {
			return err
		}
		c.ev.Send(events.New(events.StatusDone, "Blockchain built"))

		c.ev.Send(events.New(events.StatusOngoing, "Initializing the genesis"))
		if err := c.initGenesis(ctx); err != nil {
			return err
		}
		c.ev.Send(events.New(events.StatusDone, "Genesis initialized"))
	}

	return c.buildGenesis(ctx, gi)
}

// buildGenesis builds the genesis for the chain from the launch approved requests
func (c Chain) buildGenesis(ctx context.Context, gi networktypes.GenesisInformation) error {
	c.ev.Send(events.New(events.StatusOngoing, "Building the genesis"))

	addressPrefix, err := c.detectPrefix(ctx)
	if err != nil {
		return errors.Wrap(err, "error detecting chain prefix")
	}

	// apply genesis information to the genesis
	if err := c.applyGenesisAccounts(ctx, gi.GenesisAccounts, addressPrefix); err != nil {
		return errors.Wrap(err, "error applying genesis accounts to genesis")
	}
	if err := c.applyVestingAccounts(ctx, gi.VestingAccounts, addressPrefix); err != nil {
		return errors.Wrap(err, "error applying vesting accounts to genesis")
	}
	if err := c.applyGenesisValidators(ctx, gi.GenesisValidators); err != nil {
		return errors.Wrap(err, "error applying genesis validators to genesis")
	}

	// set the genesis time for the chain
	genesisPath, err := c.chain.GenesisPath()
	if err != nil {
		return errors.Wrap(err, "genesis of the blockchain can't be read")
	}
	if err := cosmosutil.SetGenesisTime(genesisPath, c.launchTime); err != nil {
		return errors.Wrap(err, "genesis time can't be set")
	}

	c.ev.Send(events.New(events.StatusDone, "Genesis built"))

	return nil
}

// applyGenesisAccounts adds the genesis account into the genesis using the chain CLI
func (c Chain) applyGenesisAccounts(
	ctx context.Context,
	genesisAccs []networktypes.GenesisAccount,
	addressPrefix string,
) error {
	var err error

	cmd, err := c.chain.Commands(ctx)
	if err != nil {
		return err
	}

	for _, acc := range genesisAccs {
		// change the address prefix to the target chain prefix
		acc.Address, err = cosmosutil.ChangeAddressPrefix(acc.Address, addressPrefix)
		if err != nil {
			return err
		}

		// call the add genesis account CLI command
		err = cmd.AddGenesisAccount(ctx, acc.Address, acc.Coins)
		if err != nil {
			return err
		}
	}

	return nil
}

// applyVestingAccounts adds the genesis vesting account into the genesis using the chain CLI
func (c Chain) applyVestingAccounts(
	ctx context.Context,
	vestingAccs []networktypes.VestingAccount,
	addressPrefix string,
) error {
	cmd, err := c.chain.Commands(ctx)
	if err != nil {
		return err
	}

	for _, acc := range vestingAccs {
		acc.Address, err = cosmosutil.ChangeAddressPrefix(acc.Address, addressPrefix)
		if err != nil {
			return err
		}

		// call the add genesis account CLI command with delayed vesting option
		err = cmd.AddVestingAccount(
			ctx,
			acc.Address,
			acc.StartingBalance,
			acc.Vesting,
			acc.EndTime,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// applyGenesisValidators gathers the validator gentxs into the genesis and adds peers in config
func (c Chain) applyGenesisValidators(ctx context.Context, genesisVals []networktypes.GenesisValidator) error {
	// no validator
	if len(genesisVals) == 0 {
		return nil
	}

	// reset the gentx directory
	gentxDir, err := c.chain.GentxsPath()
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
		if err = ioutil.WriteFile(gentxPath, val.Gentx, 0666); err != nil {
			return err
		}
	}

	// gather gentxs
	cmd, err := c.chain.Commands(ctx)
	if err != nil {
		return err
	}
	if err := cmd.CollectGentxs(ctx); err != nil {
		return err
	}

	return c.updateConfigFromGenesisValidators(genesisVals)
}

// updateConfigFromGenesisValidators adds the peer addresses into the config.toml of the chain
func (c Chain) updateConfigFromGenesisValidators(genesisVals []networktypes.GenesisValidator) error {
	var p2pAddresses []string
	for _, val := range genesisVals {
		p2pAddresses = append(p2pAddresses, val.Peer)
	}

	// set persistent peers
	configPath, err := c.chain.ConfigTOMLPath()
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
