package ignitecmd

import (
	"github.com/spf13/cobra"
)

// NewNetworkClient creates a new client command that holds some other
// sub commands related to connect client to the network.
func NewNetworkClient() *cobra.Command {
	c := &cobra.Command{
		Use:   "client",
		Short: "Connect your network with SPN",
	}

	c.AddCommand(
		NewNetworkClientCreate(),
	)

	return c
}
