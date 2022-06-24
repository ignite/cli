package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/services/network"
)

const (
	flagRemainingTime = "remaining-time"
)

// NewNetworkChainLaunch creates a new chain launch command to launch
// the network as a coordinator.
func NewNetworkChainLaunch() *cobra.Command {
	c := &cobra.Command{
		Use:   "launch [launch-id]",
		Short: "Launch a network as a coordinator",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainLaunchHandler,
	}

	c.Flags().Duration(flagRemainingTime, 0, "Duration of time in seconds before the chain is effectively launched")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())

	return c
}

func networkChainLaunchHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// parse launch ID
	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	remainingTime, _ := cmd.Flags().GetDuration(flagRemainingTime)

	n, err := nb.Network()
	if err != nil {
		return err
	}

	return n.TriggerLaunch(cmd.Context(), launchID, remainingTime)
}
