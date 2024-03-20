package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// NewRelayer returns a new relayer command.
func NewRelayer() *cobra.Command {
	c := &cobra.Command{
		Use:     "relayer",
		Aliases: []string{"r"},
		Short:   "Connect blockchains with an IBC relayer",
	}

	c.AddCommand(
		NewRelayerConfigure(),
		NewRelayerConnect(),
	)

	return c
}

func handleRelayerAccountErr(err error) error {
	var accountErr *cosmosaccount.AccountDoesNotExistError
	if !errors.As(err, &accountErr) {
		return err
	}

	return errors.Wrap(accountErr, `make sure to create or import your account through "ignite account" commands`)
}
