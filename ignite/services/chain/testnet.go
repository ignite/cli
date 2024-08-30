package chain

import (
	"context"
	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
	"os"
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
