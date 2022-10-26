package ignitecmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/entrywriter"
	"github.com/ignite/cli/ignite/services/network"
	"github.com/ignite/cli/ignite/services/network/networktypes"
)

var LaunchSummaryHeader = []string{
	"launch ID",
	"chain ID",
	"source",
	"phase",
}

var LaunchSummaryAdvancedHeader = []string{
	"campaign ID",
	"network",
	"reward",
}

// NewNetworkChainList returns a new command to list all published chains on Ignite
func NewNetworkChainList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List published chains",
		Args:  cobra.NoArgs,
		RunE:  networkChainListHandler,
	}
	c.Flags().Bool(flagAdvanced, false, "Show advanced information about the chains")
	c.Flags().Uint64(flagLimit, 100, "Limit of results per page")
	c.Flags().Uint64(flagPage, 1, "Page for chain list result")

	return c
}

func networkChainListHandler(cmd *cobra.Command, _ []string) error {
	var (
		advanced, _ = cmd.Flags().GetBool(flagAdvanced)
		limit, _    = cmd.Flags().GetUint64(flagLimit)
		page, _     = cmd.Flags().GetUint64(flagPage)
	)

	session := cliui.New(cliui.StartSpinner())

	defer session.End()

	if page == 0 {
		return errors.New("invalid page value")
	}

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}
	n, err := nb.Network(network.CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}
	chainLaunches, err := n.ChainLaunchesWithReward(cmd.Context(), &query.PageRequest{
		Offset: limit * (page - 1),
		Limit:  limit,
	})
	if err != nil {
		return err
	}

	return renderLaunchSummaries(chainLaunches, session, advanced)
}

// renderLaunchSummaries writes into the provided out, the list of summarized launches
func renderLaunchSummaries(chainLaunches []networktypes.ChainLaunch, session *cliui.Session, advanced bool) error {
	header := LaunchSummaryHeader
	if advanced {
		// advanced information show the campaign ID, type of network and rewards for incentivized testnet
		header = append(header, LaunchSummaryAdvancedHeader...)
	}

	var launchEntries [][]string

	// iterate and fetch summary for chains
	for _, c := range chainLaunches {

		// get the current phase of the chain
		var phase string
		switch {
		case !c.LaunchTriggered:
			phase = "coordinating"
		case time.Now().Before(c.LaunchTime):
			phase = "launching"
		default:
			phase = "launched"
		}

		entry := []string{
			fmt.Sprintf("%d", c.ID),
			c.ChainID,
			c.SourceURL,
			phase,
		}

		// add advanced information
		if advanced {
			campaign := "no campaign"
			if c.CampaignID > 0 {
				campaign = fmt.Sprintf("%d", c.CampaignID)
			}

			reward := entrywriter.None
			if len(c.Reward) > 0 {
				reward = c.Reward
			}

			entry = append(entry,
				campaign,
				c.Network.String(),
				reward)
		}

		launchEntries = append(launchEntries, entry)
	}

	return session.PrintTable(header, launchEntries...)
}
