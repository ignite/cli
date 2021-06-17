package starportcmd

import (
	"github.com/spf13/cobra"
)

// NewRelayer returns a new relayer command.
func NewRelayer() *cobra.Command {
	c := &cobra.Command{
		Use:   "relayer",
		Short: "Connects blockchains via IBC protocol",
	}

	c.AddCommand(NewRelayerConfigure())
	c.AddCommand(NewRelayerConnect())
	c.AddCommand(NewRelayerLowLevel())

	return c
}
