package starportcmd

import (
	"fmt"

	"github.com/ignite-hq/cli/ignite/pkg/cosmosaccount"
	"github.com/spf13/cobra"
)

func NewAccountCreate() *cobra.Command {
	c := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new account",
		Args:  cobra.ExactArgs(1),
		RunE:  accountCreateHandler,
	}

	c.Flags().AddFlagSet(flagSetKeyringBackend())

	return c
}

func accountCreateHandler(cmd *cobra.Command, args []string) error {
	name := args[0]

	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringBackend(getKeyringBackend(cmd)),
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
