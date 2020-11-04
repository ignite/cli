package starportcmd

import (
	"github.com/spf13/cobra"
)

var spnAddress string

func NewNetwork() *cobra.Command {
	c := &cobra.Command{
		Use:   "network",
		Short: "Create and start Blochains collaboratively",
		Args:  cobra.ExactArgs(1),
	}

	// configure flags.
	c.Flags().StringVarP(&spnAddress, "spn-address", "s", "localhost:26657", "An SPN node address")

	// add sub commands.
	c.AddCommand(NewNetworkChain())
	c.AddCommand(NewNetworkAccount())
	return c
}
