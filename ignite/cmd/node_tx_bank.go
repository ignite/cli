package ignitecmd

import "github.com/spf13/cobra"

func NewNodeTxBank() *cobra.Command {
	c := &cobra.Command{
		Use:   "bank",
		Short: "Bank transaction subcommands",
	}

	c.AddCommand(NewNodeTxBankSend())

	return c
}
