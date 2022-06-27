package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
)

var (
	campaignMainnetsAccSummaryHeader = []string{"Mainnet Account", "Shares"}
)

// NewNetworkCampaignAccount creates a new campaign account command that holds some other
// sub commands related to account for a campaign.
func NewNetworkCampaignAccount() *cobra.Command {
	c := &cobra.Command{
		Use:   "account",
		Short: "Handle campaign accounts",
	}
	c.AddCommand(
		newNetworkCampaignAccountList(),
	)
	return c
}

func newNetworkCampaignAccountList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list [campaign-id]",
		Short: "Show all mainnet and mainnet vesting of the campaign",
		Args:  cobra.ExactArgs(1),
		RunE:  newNetworkCampaignAccountListHandler,
	}
	return c
}

func newNetworkCampaignAccountListHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	nb, campaignID, err := networkChainLaunch(cmd, args, session)
	if err != nil {
		return err
	}
	n, err := nb.Network()
	if err != nil {
		return err
	}

	// get all campaign accounts
	mainnetAccs, err := n.MainnetAccounts(cmd.Context(), campaignID)
	if err != nil {
		return err
	}

	if len(mainnetAccs) == 0 {
		session.StopSpinner()
		return session.Printf("%s %s\n", icons.Info, "no campaign account found")
	}

	mainnetAccEntries := make([][]string, 0)
	for _, acc := range mainnetAccs {
		mainnetAccEntries = append(mainnetAccEntries, []string{acc.Address, acc.Shares.String()})
	}

	session.StopSpinner()
	if len(mainnetAccEntries) > 0 {
		if err = session.PrintTable(campaignMainnetsAccSummaryHeader, mainnetAccEntries...); err != nil {
			return err
		}
	}

	return nil
}
