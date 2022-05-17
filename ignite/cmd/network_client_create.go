package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/pkg/cliui"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/icons"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
	"github.com/ignite-hq/cli/ignite/services/network"
)

// NewNetworkClientCreate creates a client id in monitoring consumer modules of SPN
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

	clientID, err := clientCreate(cmd, launchID, nodeAPI)
	if err != nil {
		return err
	}

	session.StopSpinner()
	session.Printf("%s Client created: %s\n", icons.Info, clientID)
	return nil
}

func clientCreate(cmd *cobra.Command, launchID uint64, nodeAPI string) (string, error) {
	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return "", err
	}

	nodeClient, err := cosmosclient.New(cmd.Context(), cosmosclient.WithNodeAddress(nodeAPI))
	if err != nil {
		return "", err
	}
	node, err := network.NewNodeClient(nodeClient)
	if err != nil {
		return "", err
	}

	ibcInfo, err := node.IBCInfo(cmd.Context())
	if err != nil {
		return "", err
	}

	n, err := nb.Network()
	if err != nil {
		return "", err
	}

	return n.CreateClient(launchID, ibcInfo)
}
