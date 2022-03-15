package starportcmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/tendermint/starport/starport/pkg/entrywriter"
	"github.com/tendermint/starport/starport/services/network/networktypes"
)

var LaunchSummaryHeader = []string{"launch ID", "chain ID", "source", "campaign ID"}

// NewNetworkChainList returns a new command to list all published chains on Starport Network
func NewNetworkChainList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List published chains",
		Args:  cobra.NoArgs,
		RunE:  networkChainListHandler,
	}
	return c
}

func networkChainListHandler(cmd *cobra.Command, args []string) error {
	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}

	nb.Spinner.Stop()

	n, err := nb.Network()
	if err != nil {
		return err
	}
	chainLaunches, err := n.ChainLaunches(cmd.Context())
	if err != nil {
		return err
	}

	nb.Cleanup()
	return renderLaunchSummaries(chainLaunches, os.Stdout)
}

// renderLaunchSummaries writes into the provided out, the list of summarized launches
func renderLaunchSummaries(chainLaunches []networktypes.ChainLaunch, out io.Writer) error {
	var launchEntries [][]string

	for _, c := range chainLaunches {
		campaign := "no campaign"
		if c.CampaignID > 0 {
			campaign = fmt.Sprintf("%d", c.CampaignID)
		}

		launchEntries = append(launchEntries, []string{
			fmt.Sprintf("%d", c.ID),
			c.ChainID,
			c.SourceURL,
			campaign,
		})
	}

	return entrywriter.MustWrite(out, LaunchSummaryHeader, launchEntries...)
}
