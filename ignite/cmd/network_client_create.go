package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/pkg/clispinner"
	"github.com/ignite-hq/cli/ignite/services/network"
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
	node, err := network.NewNode(cmd.Context(), nodeAPI)
	if err != nil {
		return err
	}

	ibcInfo, err := node.IBCInfo(cmd.Context())
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	clientID, err := n.CreateClient(
		launchID,
		ibcInfo.ConsensusState,
		ibcInfo.ValidatorSet,
		ibcInfo.UnbondingTime,
		ibcInfo.Height,
	)
	if err != nil {
		return err
	}

	nb.Spinner.Stop()
	fmt.Printf("%s Client created: %s\n", clispinner.Info, clientID)
	return nil
}
