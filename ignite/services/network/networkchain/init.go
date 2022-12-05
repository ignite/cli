package networkchain

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	chainconfig "github.com/ignite/cli/ignite/config/chain"
	"github.com/ignite/cli/ignite/pkg/cache"
	cosmosgenesis "github.com/ignite/cli/ignite/pkg/cosmosutil/genesis"
	"github.com/ignite/cli/ignite/pkg/events"
)

// Init initializes blockchain by building the binaries and running the init command and
// create the initial genesis of the chain, and set up a validator key
func (c *Chain) Init(ctx context.Context, cacheStorage cache.Storage) error {
	chainHome, err := c.chain.Home()
	if err != nil {
		return err
	}

	// cleanup home dir of app if exists.
	if err = os.RemoveAll(chainHome); err != nil {
		return err
	}

	// build the chain and initialize it with a new validator key
	if _, err := c.Build(ctx, cacheStorage); err != nil {
		return err
	}

	c.ev.Send("Initializing the blockchain", events.ProgressStart())

	if err = c.chain.Init(ctx, false); err != nil {
		return err
	}

	c.ev.Send("Blockchain initialized", events.ProgressFinish())

	// initialize and verify the genesis
	if err = c.initGenesis(ctx); err != nil {
		return err
	}

	c.isInitialized = true

	return nil
}

// initGenesis creates the initial genesis of the genesis depending on the initial genesis type (default, url, ...)
func (c *Chain) initGenesis(ctx context.Context) error {
	c.ev.Send("Computing the Genesis", events.ProgressStart())

	genesisPath, err := c.chain.GenesisPath()
	if err != nil {
		return err
	}

	// remove existing genesis
	if err := os.RemoveAll(genesisPath); err != nil {
		return err
	}

	// if the blockchain has a genesis URL, the initial genesis is fetched from the URL
	// otherwise, the default genesis is used, which requires no action since the default genesis is generated from the init command
	switch {
	case c.genesisURL != "":
		c.ev.Send("Fetching custom Genesis from URL", events.ProgressUpdate())
		genesis, err := cosmosgenesis.FromURL(ctx, c.genesisURL, genesisPath)
		if err != nil {
			return err
		}

		if genesis.TarballPath() != "" {
			c.ev.Send(
				fmt.Sprintf("Extracted custom Genesis from tarball at %s", genesis.TarballPath()),
				events.ProgressFinish(),
			)
		} else {
			c.ev.Send("Custom Genesis JSON from URL fetched", events.ProgressFinish())
		}

		hash, err := genesis.Hash()
		if err != nil {
			return err
		}

		// if the blockchain has been initialized with no genesis hash, we assign the fetched hash to it
		// otherwise we check the genesis integrity with the existing hash
		if c.genesisHash == "" {
			c.genesisHash = hash
		} else if hash != c.genesisHash {
			return fmt.Errorf("genesis from URL %s is invalid. expected hash %s, actual hash %s", c.genesisURL, c.genesisHash, hash)
		}

		genBytes, err := genesis.Bytes()
		if err != nil {
			return err
		}

		// replace the default genesis with the fetched genesis
		if err := os.WriteFile(genesisPath, genBytes, 0o644); err != nil {
			return err
		}
	case c.genesisConfig != "":
		c.ev.Send("Fetching custom genesis from chain config", events.ProgressUpdate())

		// first, initialize with default genesis
		cmd, err := c.chain.Commands(ctx)
		if err != nil {
			return err
		}

		// TODO: use validator moniker https://github.com/ignite/cli/issues/1834
		if err := cmd.Init(ctx, "moniker"); err != nil {
			return err
		}

		// find config in downloaded source
		path := filepath.Join(c.path, c.genesisConfig)
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("the config for genesis doesn't exist: %w", err)
		}

		config, err := chainconfig.ParseNetworkFile(path)
		if err != nil {
			return err
		}

		// make sure that chain id given during chain.New() has the most priority.
		chainID, err := c.ID()
		if err != nil {
			return err
		}
		if config.Genesis != nil {
			config.Genesis["chain_id"] = chainID
		}

		// update genesis file with the genesis values defined in the config
		if err := c.chain.UpdateGenesisFile(config.Genesis); err != nil {
			return err
		}

		if err := c.chain.InitAccounts(ctx, config); err != nil {
			return err
		}

	default:
		// default genesis is used, init CLI command is used to generate it
		cmd, err := c.chain.Commands(ctx)
		if err != nil {
			return err
		}

		// TODO: use validator moniker https://github.com/ignite/cli/issues/1834
		if err := cmd.Init(ctx, "moniker"); err != nil {
			return err
		}
	}

	// check the initial genesis is valid
	if err := c.checkInitialGenesis(ctx); err != nil {
		return err
	}

	c.ev.Send("Genesis initialized", events.ProgressFinish())
	return nil
}

// checkGenesis checks the stored genesis is valid
func (c *Chain) checkInitialGenesis(ctx context.Context) error {
	// perform static analysis of the chain with the validate-genesis command.
	chainCmd, err := c.chain.Commands(ctx)
	if err != nil {
		return err
	}

	// the chain initial genesis should not contain gentx, gentxs should be added through requests
	genesisPath, err := c.chain.GenesisPath()
	if err != nil {
		return err
	}

	chainGenesis, err := cosmosgenesis.FromPath(genesisPath)
	if err != nil {
		return err
	}

	gentxCount, err := chainGenesis.GentxCount()
	if err != nil {
		return err
	}

	if gentxCount > 0 {
		return errors.New("the initial genesis for the chain should not contain gentx")
	}

	return chainCmd.ValidateGenesis(ctx)

	// TODO: static analysis of the genesis with validate-genesis doesn't check the full validity of the genesis
	// example: gentxs formats are not checked
	// to perform a full validity check of the genesis we must try to start the chain with sample accounts
}
