package starportcmd

import (
	"github.com/ignite-hq/cli/starport/pkg/cosmosaccount"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NewRelayer returns a new relayer command.
func NewRelayer() *cobra.Command {
	c := &cobra.Command{
		Use:     "relayer",
		Aliases: []string{"r"},
		Short:   "Connect blockchains by using IBC protocol",
	}

	c.AddCommand(NewRelayerConfigure())
	c.AddCommand(NewRelayerConnect())

	return c
}

func handleRelayerAccountErr(err error) error {
	var accountErr *cosmosaccount.AccountDoesNotExistError
	if !errors.As(err, &accountErr) {
		return err
	}

	return errors.Wrap(accountErr, `make sure to create or import your account through "starport account" commands`)
}
