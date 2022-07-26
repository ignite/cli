package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
)

func NewAccountList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "Show a list of all accounts",
		RunE:  accountListHandler,
	}

	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	c.Flags().AddFlagSet(flagSetAccountPrefixes())

	return c
}

func accountListHandler(cmd *cobra.Command, args []string) error {
	var (
		keyringBackend = getKeyringBackend(cmd)
		keyringDir     = getKeyringDir(cmd)
	)

	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringBackend(keyringBackend),
		cosmosaccount.WithHome(keyringDir),
	)
	if err != nil {
		return err
	}

	accounts, err := ca.List()
	if err != nil {
		return err
	}

	return printAccounts(cmd, accounts...)
}
