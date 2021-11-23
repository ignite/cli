package starportcmd

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/services/network"
)

// NewNetworkChainPrepare returns a new command to prepare the chain for launch
func NewNetworkChainPrepare() *cobra.Command {
	c := &cobra.Command{
		Use:   "prepare [launch-id]",
		Short: "Prepare the chain for launch",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainPrepareHandler,
	}

	c.Flags().String(flagFrom, cosmosaccount.DefaultAccount, "Account name to use for sending transactions to SPN")
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetHome())

	return c
}

func networkChainPrepareHandler(cmd *cobra.Command, args []string) error {
	nb, _, shutdown, err := initializeNetwork(cmd)
	if err != nil {
		return err
	}
	defer shutdown()

	// parse launch ID
	launchID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return errors.Wrap(err, "error parsing launchID")
	}
	if launchID == 0 {
		return errors.New("launch ID must be greater than 0")
	}

	// initialize the blockchain from the launch ID
	initOptions := initOptionWithHomeFlag(cmd, network.InitializationPrepareLaunch())
	sourceOption := network.SourceLaunchID(launchID)
	blockchain, err := nb.Blockchain(cmd.Context(), sourceOption, initOptions...)
	if err != nil {
		return err
	}

	// prepare for launch
	if err := blockchain.Prepare(cmd.Context()); err != nil {
		return err
	}

	return nil
}
