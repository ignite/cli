package ignitecmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/goenv"
	"github.com/ignite/cli/ignite/services/network"
	"github.com/ignite/cli/ignite/services/network/networkchain"
)

const (
	flagForce = "force"
)

// NewNetworkChainPrepare returns a new command to prepare the chain for launch
func NewNetworkChainPrepare() *cobra.Command {
	c := &cobra.Command{
		Use:   "prepare [launch-id]",
		Short: "Prepare the chain for launch",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainPrepareHandler,
	}

	flagSetClearCache(c)
	c.Flags().BoolP(flagForce, "f", false, "Force the prepare command to run even if the chain is not launched")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetHome())

	return c
}

func networkChainPrepareHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	force, _ := cmd.Flags().GetBool(flagForce)

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// parse launch ID
	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	// fetch chain information
	chainLaunch, err := n.ChainLaunch(cmd.Context(), launchID)
	if err != nil {
		return err
	}

	if !force && !chainLaunch.LaunchTriggered {
		return fmt.Errorf("chain %d launch has not been triggered yet. use --force to prepare anyway", launchID)
	}

	c, err := nb.Chain(networkchain.SourceLaunch(chainLaunch))
	if err != nil {
		return err
	}

	// fetch the information to construct genesis
	genesisInformation, err := n.GenesisInformation(cmd.Context(), launchID)
	if err != nil {
		return err
	}

	rewardsInfo, lastBlockHeight, unboundingTime, err := n.RewardsInfo(
		cmd.Context(),
		launchID,
		chainLaunch.ConsumerRevisionHeight,
	)
	if err != nil {
		return err
	}

	spnChainID, err := n.ChainID(cmd.Context())
	if err != nil {
		return err
	}

	if err := c.Prepare(
		cmd.Context(),
		cacheStorage,
		genesisInformation,
		rewardsInfo,
		spnChainID,
		lastBlockHeight,
		unboundingTime,
	); err != nil {
		return err
	}

	chainHome, err := c.Home()
	if err != nil {
		return err
	}
	binaryName, err := c.BinaryName()
	if err != nil {
		return err
	}
	binaryDir := filepath.Dir(filepath.Join(goenv.Bin(), binaryName))

	session.StopSpinner()
	session.Printf("%s Chain is prepared for launch\n", icons.OK)
	session.Println("\nYou can start your node by running the following command:")
	commandStr := fmt.Sprintf("%s start --home %s", binaryName, chainHome)
	session.Printf("\t%s/%s\n", binaryDir, colors.Info(commandStr))

	return nil
}
