package starportcmd

import (
	"github.com/spf13/cobra"
)

// NewNetworkChainStart creates a network chain start commmand.
// this command runs target chain's start command by default. any flags passed after
// a double dash passed to the underlying start command.
func NewNetworkChainStart() *cobra.Command {
	c := &cobra.Command{
		Use:   "start [chain-id] [-- <flags>...]",
		Short: "Start network",
		RunE:  networkChainStartHandler,
		Args:  cobra.MinimumNArgs(1),
	}
	return c
}

func networkChainStartHandler(cmd *cobra.Command, args []string) error {
	nb, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	var startFlags []string
	chainID := args[0]
	if len(args) > 1 { // first arg is always `chain-id`.
		startFlags = args[1:]
	}

	return nb.StartChain(cmd.Context(), chainID, startFlags)
}
