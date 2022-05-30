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

	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	c.Flags().AddFlagSet(flagSetAccountPrefixes())

	return c
}

func accountShowHandler(cmd *cobra.Command, args []string) error {
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

	acc, err := ca.GetByName(name)
	if err != nil {
		return err
	}

	return printAccounts(cmd, acc)
}
