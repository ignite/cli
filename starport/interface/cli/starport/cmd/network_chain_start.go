package starportcmd

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/networkbuilder"
)

// NewNetworkChainStart creates a network chain start command.
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
	// Check if custom home is provided
	var initOptions []networkbuilder.InitOption
	home, cliHome, err := getHomeFlags(cmd)
	if err != nil {
		return err
	}
	if home != "" {
		initOptions = append(initOptions, networkbuilder.InitializationHomePath(home))
	}
	if cliHome != "" {
		initOptions = append(initOptions, networkbuilder.InitializationCLIHomePath(cliHome))
	}

	nb, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	var startFlags []string
	chainID := args[0]
	if len(args) > 1 { // first arg is always `chain-id`.
		startFlags = args[1:]
	}

	return nb.StartChain(cmd.Context(), chainID, startFlags, initOptions...)
}
