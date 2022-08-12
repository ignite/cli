package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
)

func NewAccountDelete() *cobra.Command {
	c := &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete an account by name",
		Args:  cobra.ExactArgs(1),
		RunE:  accountDeleteHandler,
	}

	return c
}

func accountDeleteHandler(cmd *cobra.Command, args []string) error {
	name := args[0]

	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringBackend(getKeyringBackend(cmd)),
		cosmosaccount.WithHome(getKeyringDir(cmd)),
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
