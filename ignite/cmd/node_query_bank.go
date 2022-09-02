package ignitecmd

import "github.com/spf13/cobra"

func NewNodeQueryBank() *cobra.Command {
	c := &cobra.Command{
		Use:   "bank",
		Short: "Querying commands for the bank module",
	}

	c.AddCommand(NewNodeQueryBankBalances())

	return c
}
