package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/pkg/cliui"
	"github.com/ignite-hq/cli/ignite/pkg/yaml"
	"github.com/ignite-hq/cli/ignite/services/network"
)

func newNetworkChainShowInfo() *cobra.Command {
	c := &cobra.Command{
		Use:   "info [launch-id]",
		Short: "Show info details of the chain",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainShowInfoHandler,
	}
	return c
}

func networkChainShowInfoHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	nb, launchID, err := networkChainLaunch(cmd, args, session)
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

	reward, err := n.ChainReward(cmd.Context(), launchID)
	if err != nil && err != network.ErrObjectNotFound {
		return err
	}
	chainLaunch.Reward = reward.RemainingCoins.String()

	info, err := yaml.Marshal(cmd.Context(), chainLaunch)
	if err != nil {
		return err
	}

	session.StopSpinner()
	return session.Print(info)
}
