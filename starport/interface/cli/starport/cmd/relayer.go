package starportcmd

import "github.com/spf13/cobra"

func NewRelayer() *cobra.Command {
	c := &cobra.Command{
		Use:   "relayer",
		Short: "Relay connects blockchains via IBC protocol",
	}
	c.AddCommand(NewRelayerConnect())
	c.AddCommand(NewRelayerStart())
	return c
}
