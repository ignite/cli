package ignitecmd

import (
	"time"

	timeparser "github.com/aws/smithy-go/time"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/services/network"
)

const (
	flagLauchTime = "launch-time"
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

	c.Flags().String(
		flagLauchTime,
		"",
		"Timestamp the chain is effectively launched (example \"2022-01-01T00:00:00Z\")",
	)
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())

	return c
}

func networkChainLaunchHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// parse launch ID
	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	// parse launch time
	var launchTime time.Time
	launchTimeStr, _ := cmd.Flags().GetString(flagLauchTime)

	if launchTimeStr != "" {
		launchTime, err = timeparser.ParseDateTime(launchTimeStr)
		if err != nil {
			return err
		}
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	return n.TriggerLaunch(cmd.Context(), launchID, launchTime)
}
