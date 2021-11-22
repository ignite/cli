package starportcmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/services/network"
	"strconv"
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
	nb, s, shutdown, err := initializeNetwork(cmd)
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

	chainHome, exist, err := network.IsChainHomeExist(launchID, true)
	if err != nil {
		return err
	}
	if !getYes(cmd) && exist {
		prompt := promptui.Prompt{
			Label: fmt.Sprintf("The chain launch has already been prepared under: %s. Would you like to overwrite the home directory",
				chainHome,
			),
			IsConfirm: true,
		}
		s.Stop()
		if _, err := prompt.Run(); err != nil {
			fmt.Println("said no")
			return nil
		}
		s.Start()
	}

	// initialize the blockchain from the launch ID
	initOptions := initOptionWithHomeFlag(cmd, network.InitializationPrepareLaunch())
	sourceOption := network.SourceLaunchID(launchID)
	blockchain, err := nb.Blockchain(cmd.Context(), sourceOption, initOptions...)
	if err != nil {
		return err
	}

	if err := blockchain.Init(cmd.Context()); err != nil {
		return err
	}

	// prepare for launch
	if err := blockchain.Prepare(cmd.Context()); err != nil {
		return err
	}

	return nil
}
