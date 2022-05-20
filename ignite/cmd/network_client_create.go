package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/pkg/cliui"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/icons"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
	"github.com/ignite-hq/cli/ignite/services/network"
	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
)

// NewNetworkClientCreate creates a client id in monitoring consumer modules of SPN
func NewNetworkClientCreate() *cobra.Command {
	c := &cobra.Command{
		Use:   "create [launch-id] [chain-rpc]",
		Short: "Connect the monitoring modules of launched chains with SPN",
		Args:  cobra.ExactArgs(2),
		RunE:  networkClientCreateHandler,
	}

	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().String(flagSPNChainID, networktypes.SPNChainID, "Chain ID of SPN")

	return c
}

func networkClientCreateHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}
	chainRPC := args[1]

	nodeID, chainClientID, spnClientID, err := clientCreate(cmd, launchID, chainRPC)
	if err != nil {
		return err
	}

	session.StopSpinner()
	session.Printf("%s Network client created: %s\n", icons.Info, spnClientID)
	session.Printf("%s Target chain %s client: %s\n", icons.Info, nodeID, chainClientID)
	return nil
}

func clientCreate(cmd *cobra.Command, launchID uint64, nodeAPI string) (string, string, string, error) {
	spnChainID, _ := cmd.Flags().GetString(flagSPNChainID)

	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return "", "", "", err
	}

	nodeClient, err := cosmosclient.New(cmd.Context(), cosmosclient.WithNodeAddress(nodeAPI))
	if err != nil {
		return "", "", "", err
	}
	node, err := network.NewNodeClient(nodeClient)
	if err != nil {
		return "", "", "", err
	}

	nodeClientID, err := node.FindClientID(cmd.Context(), spnChainID)
	if err != nil {
		return "", "", "", err
	}

	rewardsInfo, nodeID, unboundingTime, err := node.RewardsInfo(cmd.Context())
	if err != nil {
		return "", "", "", err
	}

	n, err := nb.Network()
	if err != nil {
		return "", "", "", err
	}

	spnClientID, err := n.CreateClient(launchID, unboundingTime, rewardsInfo)
	if err != nil {
		return "", "", "", err
	}

	return nodeID, nodeClientID, spnClientID, err
}
