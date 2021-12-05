package starportcmd

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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

	c.Flags().String(flagRemainingTime, "", "The remaining time for validator preparation before the chain is effectively launched")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())

	return c
}

func networkChainLaunchHandler(cmd *cobra.Command, args []string) error {
	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}
	defer nb.Cleanup()

	// parse launch ID.
	launchID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return errors.Wrap(err, "error parsing launchID")
	}
	if launchID == 0 {
		return errors.New("launch ID must be greater than 0")
	}

	remainingTime, _ := cmd.Flags().GetUint64(flagRemainingTime)

	n, err := nb.Network()
	if err != nil {
		return err
	}

	return n.TriggerLaunch(cmd.Context(), launchID, remainingTime)
}
