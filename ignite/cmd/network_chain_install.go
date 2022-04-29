package ignitecmd

import (
	"path/filepath"

	"github.com/ignite-hq/cli/ignite/pkg/cliui"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/colors"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/icons"
	"github.com/ignite-hq/cli/ignite/pkg/goenv"
	"github.com/ignite-hq/cli/ignite/services/network"
	"github.com/ignite-hq/cli/ignite/services/network/networkchain"
	"github.com/spf13/cobra"
)

// NewNetworkChainInstall returns a new command to install a chain's binary by the launch id.
func NewNetworkChainInstall() *cobra.Command {
	c := &cobra.Command{
		Use:   "install [launch-id]",
		Short: "Install chain binary for a launch",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainInstallHandler,
	}
	c.Flags().AddFlagSet(flagNetworkFrom())
	return c
}

func networkChainInstallHandler(cmd *cobra.Command, args []string) error {
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

	n, err := nb.Network()
	if err != nil {
		return err
	}

	chainLaunch, err := n.ChainLaunch(cmd.Context(), launchID)
	if err != nil {
		return err
	}

	c, err := nb.Chain(networkchain.SourceLaunch(chainLaunch))
	if err != nil {
		return err
	}

	binaryName, err := c.Build(cmd.Context())
	if err != nil {
		return err
	}
	binaryPath := filepath.Join(goenv.Bin(), binaryName)

	session.StopSpinner()
	session.Printf("%s Binary installed\n", icons.OK)
	session.Printf("%s Binary's name: %s\n", icons.Info, colors.Info(binaryName))
	session.Printf("%s Binary's path: %s\n", icons.Info, colors.Info(binaryPath))

	return nil
}
