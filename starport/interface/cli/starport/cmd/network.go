package starportcmd

import (
	"github.com/spf13/cobra"
)

func NewNetwork() *cobra.Command {
	c := &cobra.Command{
		Use:   "network",
		Short: "Create and start Blochains collaboratively",
		Args:  cobra.ExactArgs(1),
	}
	c.AddCommand(NewNetworkChain())
	return c
}
