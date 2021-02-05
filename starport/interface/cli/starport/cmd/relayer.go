package starportcmd

import "github.com/spf13/cobra"

func NewRelayer() *cobra.Command {
	c := &cobra.Command{
		Use:   "relayer",
		Short: "Relay connects blockchains via IBC protocol",
	}
	c.AddCommand(NewRelayerConfigure())
	c.AddCommand(NewRelayerConnect())
	return c
}
