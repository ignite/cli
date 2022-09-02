package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/entrywriter"
	"github.com/ignite/cli/ignite/services/network"
	"github.com/ignite/cli/ignite/services/network/networktypes"
)

var LaunchSummaryHeader = []string{"launch ID", "chain ID", "source", "campaign ID", "network", "reward"}

// NewNetworkChainList returns a new command to list all published chains on Ignite
func NewNetworkChainList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List published chains",
		Args:  cobra.NoArgs,
		RunE:  networkChainListHandler,
	}
	return c
}

func networkChainListHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New()
	defer session.Cleanup()

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}
	n, err := nb.Network(network.CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}
	chainLaunches, err := n.ChainLaunchesWithReward(cmd.Context())
	if err != nil {
		return err
	}

	session.StopSpinner()

	return renderLaunchSummaries(chainLaunches, session)
}

// renderLaunchSummaries writes into the provided out, the list of summarized launches
func renderLaunchSummaries(chainLaunches []networktypes.ChainLaunch, session cliui.Session) error {
	var launchEntries [][]string

	for _, c := range chainLaunches {
		campaign := "no campaign"
		if c.CampaignID > 0 {
			campaign = fmt.Sprintf("%d", c.CampaignID)
		}

		reward := entrywriter.None
		if len(c.Reward) > 0 {
			reward = c.Reward
		}

		launchEntries = append(launchEntries, []string{
			fmt.Sprintf("%d", c.ID),
			c.ChainID,
			c.SourceURL,
			campaign,
			c.Network.String(),
			reward,
		})
	}

	return session.PrintTable(LaunchSummaryHeader, launchEntries...)
}
