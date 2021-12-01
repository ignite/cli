package starportcmd

import (
	"fmt"
	"github.com/tendermint/starport/starport/services/network"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/pkg/entrywriter"
)

var LaunchSummaryHeader = []string{"launch ID", "chain ID", "source", "campaign ID"}

// LaunchSummary holds summarized information about a chain launch
type LaunchSummary struct {
	LaunchID   string
	ChainID    string
	Source     string
	CampaignID string
}

// NewNetworkChainList returns a new command to list all published chains on Starport Network
func NewNetworkChainList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List published chains",
		Args:  cobra.NoArgs,
		RunE:  networkChainListHandler,
	}
	c.Flags().String(flagFrom, cosmosaccount.DefaultAccount, "Account name to use for sending transactions to SPN")
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetHome())

	return c
}

func networkChainListHandler(cmd *cobra.Command, args []string) error {
	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}
	defer nb.Cleanup()

	n, err := nb.Network()
	if err != nil {
		return err
	}
	launchesInfo, err := n.LaunchesInfo(cmd.Context())
	if err != nil {
		return err
	}
	return renderLaunchSummaries(launchesInfo, os.Stdout)
}

// renderLaunchSummaries writes into the provided out, the list of summarized launches
func renderLaunchSummaries(launchesInfo []network.LaunchInfo, out io.Writer) error {
	var launchEntries [][]string

	for _, info := range launchesInfo {
		campaign := "no campaign"
		if info.CampaignID > 0 {
			campaign = fmt.Sprintf("%d", info.CampaignID)
		}

		launchEntries = append(launchEntries, []string{
			fmt.Sprintf("%d", info.ID),
			info.ChainID,
			info.SourceURL,
			campaign,
		})
	}

	if err := entrywriter.Write(out, LaunchSummaryHeader, launchEntries...); err != nil {
		return errors.Wrap(err, "error printing chain summaries")
	}
	return nil
}
