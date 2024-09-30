package chain

import (
	"context"
	"os"

	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
)

type InPlaceArgs struct {
	NewChainID         string
	NewOperatorAddress string
	PrvKeyValidator    string
	AccountsToFund     string
}

func (c Chain) TestnetInPlace(ctx context.Context, args InPlaceArgs) error {
	commands, err := c.Commands(ctx)
	if err != nil {
		return err
	}

	// make sure that config.yml exists
	if c.options.ConfigFile != "" {
		if _, err := os.Stat(c.options.ConfigFile); err != nil {
			return err
		}
	} else if _, err := chainconfig.LocateDefault(c.app.Path); err != nil {
		return err
	}

	err = c.InPlace(ctx, commands, args)
	if err != nil {
		return err
	}
	return nil
}

type MultiNodeArgs struct {
	ChainID               string
	OutputDir             string
	NumValidator          string
	ValidatorsStakeAmount string
	NodeDirPrefix         string
}

// If the app state still exists, TestnetMultiNode will reuse it.
// Otherwise, it will automatically re-initialize the app state from the beginning.
func (c Chain) TestnetMultiNode(ctx context.Context, args MultiNodeArgs) error {
	commands, err := c.Commands(ctx)
	if err != nil {
		return err
	}

	// make sure that config.yml exists
	if c.options.ConfigFile != "" {
		if _, err := os.Stat(c.options.ConfigFile); err != nil {
			return err
		}
	} else if _, err := chainconfig.LocateDefault(c.app.Path); err != nil {
		return err
	}

	err = c.MultiNode(ctx, commands, args)
	if err != nil {
		return err
	}
	return nil
}
