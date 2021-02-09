package starportcmd

import (
	relayercmd "github.com/cosmos/relayer/cmd"
	"github.com/spf13/cobra"
)

func NewRelayer() *cobra.Command {
	c := &cobra.Command{
		Use:   "relayer",
		Short: "Connects blockchains via IBC protocol",
	}
	c.AddCommand(NewRelayerConfigure())
	c.AddCommand(NewRelayerConnect())
	c.AddCommand(relayercmd.NewRootCmd())
	return c
}
