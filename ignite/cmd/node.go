package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/xurl"
)

const (
	flagNode         = "node"
	cosmosRPCAddress = "https://rpc.cosmos.network:443"
)

func NewNode() *cobra.Command {
	c := &cobra.Command{
		Use:   "node [command]",
		Short: "Make calls to a live blockchain node",
		Args:  cobra.ExactArgs(1),
	}

	c.PersistentFlags().String(flagNode, cosmosRPCAddress, "<host>:<port> to tendermint rpc interface for this chain")

	c.AddCommand(NewNodeQuery())
	c.AddCommand(NewNodeTx())

	return c
}

func newNodeCosmosClient(cmd *cobra.Command) (cosmosclient.Client, error) {
	var (
		home           = getHome(cmd)
		prefix         = getAddressPrefix(cmd)
		node           = getRPC(cmd)
		keyringBackend = getKeyringBackend(cmd)
		keyringDir     = getKeyringDir(cmd)
		gas            = getGas(cmd)
		gasPrices      = getGasPrices(cmd)
		fees           = getFees(cmd)
		broadcastMode  = getBroadcastMode(cmd)
	)

	options := []cosmosclient.Option{
		cosmosclient.WithAddressPrefix(prefix),
		cosmosclient.WithHome(home),
		cosmosclient.WithKeyringBackend(keyringBackend),
		cosmosclient.WithKeyringDir(keyringDir),
		cosmosclient.WithNodeAddress(xurl.HTTPEnsurePort(node)),
		cosmosclient.WithBroadcastMode(broadcastMode),
	}

	if gas != "" {
		options = append(options, cosmosclient.WithGas(gas))
	}
	if gasPrices != "" {
		options = append(options, cosmosclient.WithGasPrices(gasPrices))
	}
	if fees != "" {
		options = append(options, cosmosclient.WithFees(fees))
	}

	return cosmosclient.New(cmd.Context(), options...)
}

func getRPC(cmd *cobra.Command) (rpc string) {
	rpc, _ = cmd.Flags().GetString(flagNode)
	return
}
