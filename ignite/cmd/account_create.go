package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

func NewAccountCreate() *cobra.Command {
	c := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new account",
		Args:  cobra.ExactArgs(1),
		RunE:  accountCreateHandler,
	}

	return c
}

func accountCreateHandler(cmd *cobra.Command, args []string) error {
	name := args[0]

	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringBackend(getKeyringBackend(cmd)),
		cosmosaccount.WithHome(getKeyringDir(cmd)),
	)
	if err != nil {
		return errors.Errorf("unable to create registry: %w", err)
	}

	_, mnemonic, err := ca.Create(name)
	if err != nil {
		return errors.Errorf("unable to create account: %w", err)
	}

	cmd.Printf("Account %q created, keep your mnemonic in a secret place:\n\n%s\n", name, mnemonic)
	return nil
}
