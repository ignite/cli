package networkchain

import (
	"context"
	"fmt"
	"os"

	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"github.com/tendermint/starport/starport/pkg/events"
)

// Init initializes blockchain by building the binaries and running the init command and
// create the initial genesis of the chain.
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

	// initialize depending on the initial genesis type (default, url, ...) and verify it.
	// if the blockchain has a genesis URL, the initial genesis is fetched from the url
	// otherwise, default genesis is used, which requires no action since the default genesis is generated from the init command
	if c.genesisURL != "" {
		genesis, hash, err := cosmosutil.GenesisAndHashFromURL(ctx, c.genesisURL)
		if err != nil {
			return err
		}

		if hash != c.genesisHash {
			return fmt.Errorf("genesis from URL %s is invalid. expected hash %s, actual hash %s", c.genesisURL, c.genesisHash, hash)
		}

		// replace the default genesis with the fetched genesis
		genesisPath, err := c.chain.GenesisPath()
		if err != nil {
			return err
		}
		if err := os.WriteFile(genesisPath, genesis, 0644); err != nil {
			return err
		}
	}

	if err := c.checkGenesis(ctx); err != nil {
		return err
	}

	c.isInitialized = true

	return nil
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
