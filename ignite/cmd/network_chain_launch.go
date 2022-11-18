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
		Short: "Trigger the launch of a chain",
		Long: `The launch command communicates to the world that the chain is ready to be
launched.

Only the coordinator of the chain can execute the launch command.

	ignite network chain launch 42

After the launch command is executed no changes to the genesis are accepted. For
example, validators will no longer be able to successfully execute the "ignite
network chain join" command to apply as a validator.

The launch command sets the date and time after which the chain will start. By
default, the current time is set. To give validators more time to prepare for
the launch, set the time with the "--launch-time" flag:

	ignite network chain launch 42 --launch-time 2023-01-01T00:00:00Z

After the launch command is executed, validators can generate the finalized
genesis and prepare their nodes for the launch. For example, validators can run
"ignite network chain prepare" to generate the genesis and populate the peer
list.

If you want to change the launch time or open up the genesis file for changes
you can use "ignite network chain revert-launch" to make it possible, for
example, to accept new validators and add accounts.
`,
		Args: cobra.ExactArgs(1),
		RunE: networkChainLaunchHandler,
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
