package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/pkg/clispinner"
	"github.com/ignite-hq/cli/ignite/services/network"
	"github.com/ignite-hq/cli/ignite/services/network/node"
)

// NewNetworkConnect connects the monitoring modules of launched chains with SPN
func NewNetworkConnect() *cobra.Command {
	c := &cobra.Command{
		Use:   "connect [launch-id] [node-api-url]",
		Short: "Connect the monitoring modules of launched chains with SPN",
		Args:  cobra.ExactArgs(2),
		RunE:  networkConnectHandler,
	}
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	return c
}

func networkConnectHandler(cmd *cobra.Command, args []string) error {
	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}
	defer nb.Cleanup()

	// parse launch ID.
	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	nodeAPI := args[1]
	node, err := node.New(cmd.Context(), nodeAPI)
	if err != nil {
		return err
	}

	nodeInfo, err := node.Info(cmd.Context())
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	clientID, err := n.CreateClient(
		launchID,
		nodeInfo.ConsensusState,
		nodeInfo.ValidatorSet,
		nodeInfo.UnbondingTime,
		nodeInfo.Height,
	)
	if err != nil {
		return err
	}

	nb.Spinner.Stop()
	fmt.Printf("%s Client created: %s\n", clispinner.Info, clientID)
	return nil
}
