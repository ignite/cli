package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/services/network"
)

// NewNetworkClientCreate connects the monitoring modules of launched chains with SPN
func NewNetworkClientCreate() *cobra.Command {
	c := &cobra.Command{
		Use:   "create [launch-id] [node-api-url]",
		Short: "Connect the monitoring modules of launched chains with SPN",
		Args:  cobra.ExactArgs(2),
		RunE:  networkClientCreateHandler,
	}
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	return c
}

func networkClientCreateHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}
	nodeAPI := args[1]

	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}

	nodeClient, err := cosmosclient.New(cmd.Context(), cosmosclient.WithNodeAddress(nodeAPI))
	if err != nil {
		return err
	}
	node, err := network.NewNodeClient(nodeClient)
	if err != nil {
		return err
	}

	rewardsInfo, unboundingTime, err := node.RewardsInfo(cmd.Context())
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	clientID, err := n.CreateClient(launchID, unboundingTime, rewardsInfo)
	if err != nil {
		return err
	}

	session.StopSpinner()
	session.Printf("%s Client created: %s\n", icons.Info, clientID)
	return nil
}
