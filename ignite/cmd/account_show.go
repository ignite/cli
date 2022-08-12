package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
)

func NewAccountShow() *cobra.Command {
	c := &cobra.Command{
		Use:   "show [name]",
		Short: "Show detailed information about a particular account",
		Args:  cobra.ExactArgs(1),
		RunE:  accountShowHandler,
	}

	c.Flags().AddFlagSet(flagSetAccountPrefixes())

	return c
}

func accountShowHandler(cmd *cobra.Command, args []string) error {
	name := args[0]

	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringBackend(getKeyringBackend(cmd)),
		cosmosaccount.WithHome(getKeyringDir(cmd)),
	)
	if err != nil {
		return err
	}

	acc, err := ca.GetByName(name)
	if err != nil {
		return err
	}

	return printAccounts(cmd, acc)
}
