package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
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
		Short: "Make requests to a live blockchain node",
		Args:  cobra.ExactArgs(1),
	}

	c.PersistentFlags().String(flagNode, cosmosRPCAddress, "<host>:<port> to tendermint rpc interface for this chain")

	c.AddCommand(
		NewNodeQuery(),
		NewNodeTx(),
	)

	return c
}

func newNodeCosmosClient(cmd *cobra.Command) (cosmosclient.Client, error) {
	var (
		home           = getHome(cmd)
		prefix         = getAddressPrefix(cmd)
		node           = getNode(cmd)
		keyringBackend = getKeyringBackend(cmd)
		keyringDir     = getKeyringDir(cmd)
		gas            = getGas(cmd)
		gasPrices      = getGasPrices(cmd)
		gasAdjustment  = getGasAdjustment(cmd)
		fees           = getFees(cmd)
		generateOnly   = getGenerateOnly(cmd)
	)
	if keyringBackend == "" {
		// Makes cosmosclient usable for commands that doesn't expose the keyring
		// backend flag (cosmosclient.New returns an error if it's empty).
		keyringBackend = cosmosaccount.KeyringTest
	}

	options := []cosmosclient.Option{
		cosmosclient.WithAddressPrefix(prefix),
		cosmosclient.WithHome(home),
		cosmosclient.WithKeyringBackend(keyringBackend),
		cosmosclient.WithKeyringDir(keyringDir),
		cosmosclient.WithNodeAddress(xurl.HTTPEnsurePort(node)),
		cosmosclient.WithGenerateOnly(generateOnly),
	}

	if gas != "" {
		options = append(options, cosmosclient.WithGas(gas))
	}
	if gasPrices != "" {
		options = append(options, cosmosclient.WithGasPrices(gasPrices))
	}
	if gasAdjustment != 0 {
		options = append(options, cosmosclient.WithGasAdjustment(gasAdjustment))
	}
	if fees != "" {
		options = append(options, cosmosclient.WithFees(fees))
	}

	return cosmosclient.New(cmd.Context(), options...)
}

func getNode(cmd *cobra.Command) (node string) {
	node, _ = cmd.Flags().GetString(flagNode)
	return
}
