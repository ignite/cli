package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosaccount"
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
	var (
		name    = args[0]
		session = cliui.New(cliui.StartSpinnerWithText(statusDeleting))
	)
	defer session.End()

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

	return session.Printf("Account %s deleted.\n", name)
}
