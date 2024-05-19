package ignitecmd

import (
	"github.com/ignite/cli/v29/ignite/pkg/cosmosaccount"
	"github.com/spf13/cobra"
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

	cmd.Printf("Account %s deleted.\n", name)
	return nil
}
