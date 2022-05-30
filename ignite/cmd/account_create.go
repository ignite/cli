package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/pkg/cosmosaccount"
)

func NewAccountCreate() *cobra.Command {
	c := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new account",
		Args:  cobra.ExactArgs(1),
		RunE:  accountCreateHandler,
	}

	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())

	return c
}

func accountCreateHandler(cmd *cobra.Command, args []string) error {
	var (
		name           = args[0]
		keyringBackend = getKeyringBackend(cmd)
		keyringDir     = getKeyringDir(cmd)
	)

	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringBackend(keyringBackend),
		cosmosaccount.WithKeyringDir(keyringDir),
	)
	if err != nil {
		return err
	}

	_, mnemonic, err := ca.Create(name)
	if err != nil {
		return err
	}

	fmt.Printf("Account %q created, keep your mnemonic in a secret place:\n\n%s\n", name, mnemonic)
	return nil
}
