package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/pkg/cosmosaccount"
)

func NewAccountDelete() *cobra.Command {
	c := &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete an account by name",
		Args:  cobra.ExactArgs(1),
		RunE:  accountDeleteHandler,
	}

	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())

	return c
}

func accountDeleteHandler(cmd *cobra.Command, args []string) error {
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

	if err := ca.DeleteByName(name); err != nil {
		return err
	}

	fmt.Printf("Account %s deleted.\n", name)
	return nil
}
