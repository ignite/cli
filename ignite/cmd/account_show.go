package starportcmd

import (
	"github.com/ignite-hq/cli/ignite/pkg/cosmosaccount"
	"github.com/spf13/cobra"

	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
)

func NewAccountShow() *cobra.Command {
	c := &cobra.Command{
		Use:   "show [name]",
		Short: "Show detailed information about a particular account",
		Args:  cobra.ExactArgs(1),
		RunE:  accountShowHandler,
	}

	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetAccountPrefixes())

	return c
}

func accountShowHandler(cmd *cobra.Command, args []string) error {
	name := args[0]

	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringBackend(getKeyringBackend(cmd)),
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
