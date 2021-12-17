package networkchain

import (
	"context"
	"fmt"
	"os"

	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"github.com/tendermint/starport/starport/pkg/events"
)

// Init initializes blockchain by building the binaries and running the init command and
// create the initial genesis of the chain, and set up a validator key
func (c *Chain) Init(ctx context.Context) error {
	chainHome, err := c.chain.Home()
	if err != nil {
		return err
	}

	// cleanup home dir of app if exists.
	if err := os.RemoveAll(chainHome); err != nil {
		return err
	}

	// build the chain and initialize it with a new validator key
	c.ev.Send(events.New(events.StatusOngoing, "Building the blockchain"))
	if _, err := c.chain.Build(ctx, ""); err != nil {
		return err
	}

	c.ev.Send(events.New(events.StatusDone, "Blockchain built"))
	c.ev.Send(events.New(events.StatusOngoing, "Initializing the blockchain"))

	if err := c.chain.Init(ctx, false); err != nil {
		return err
	}

	c.ev.Send(events.New(events.StatusDone, "Blockchain initialized"))

	// initialize and verify the genesis
	if err := c.initGenesis(ctx); err != nil {
		return err
	}

	c.isInitialized = true

	return nil
}

// initGenesis creates the initial genesis of the genesis depending on the initial genesis type (default, url, ...)
func (c *Chain) initGenesis(ctx context.Context) error {
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
	if c.genesisURL != "" {
		genesis, hash, err := cosmosutil.GenesisAndHashFromURL(ctx, c.genesisURL)
		if err != nil {
			return err
		}

		// if the blockchain has been initialized with no genesis hash, we assign the fetched hash to it
		// otherwise we check the genesis integrity with the existing hash
		if c.genesisHash == "" {
			c.genesisHash = hash
		} else if hash != c.genesisHash {
			return fmt.Errorf("genesis from URL %s is invalid. Expected hash %s, actual hash %s", c.genesisURL, c.genesisHash, hash)
		}

		// replace the default genesis with the fetched genesis
		if err := os.WriteFile(genesisPath, genesis, 0644); err != nil {
			return err
		}
	} else {
		// default genesis is used, init CLI command is used to generate it
		cmd, err := c.chain.Commands(ctx)
		if err != nil {
			return err
		}

		// TODO: use validator moniker https://github.com/tendermint/starport/issues/1834
		if err := cmd.Init(ctx, "moniker"); err != nil {
			return err
		}

	}

	// check the genesis is valid
	return c.checkGenesis(ctx)
}

// checkGenesis checks the stored genesis is valid
func (c *Chain) checkGenesis(ctx context.Context) error {
	// perform static analysis of the chain with the validate-genesis command.
	chainCmd, err := c.chain.Commands(ctx)
	if err != nil {
		return err
	}

	return chainCmd.ValidateGenesis(ctx)

	// TODO: static analysis of the genesis with validate-genesis doesn't check the full validity of the genesis
	// example: gentxs formats are not checked
	// to perform a full validity check of the genesis we must try to start the chain with sample accounts
}
