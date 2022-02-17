package starportcmd

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/tendermint/starport/starport/pkg/entrywriter"
)

var (
	campaignMainnetsAccSummaryHeader = []string{"Mainnet Account", "Shares"}
	campaignVestingAccSummaryHeader  = []string{"Vesting Account", "Total Shares", "Vesting", "EndTime"}
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
	c.PersistentFlags().AddFlagSet(flagNetworkFrom())
	c.PersistentFlags().AddFlagSet(flagSetKeyringBackend())

	return c
}

func newNetworkCampaignAccountList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list [campaign-id]",
		Short: "Show all mainnet and mainnet vesting of the campaign",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nb, campaignID, err := networkChainLaunch(cmd, args)
			if err != nil {
				return err
			}
			defer nb.Cleanup()
			n, err := nb.Network()
			if err != nil {
				return err
			}

			accountSummary := bytes.NewBufferString("")

			// get all campaign mainnet accounts
			mainnetAccs, err := n.MainnetAccounts(cmd.Context(), campaignID)
			if err != nil {
				return err
			}
			mainnetAccEntries := make([][]string, 0)
			for _, acc := range mainnetAccs {
				mainnetAccEntries = append(mainnetAccEntries, []string{
					acc.Address,
					acc.Shares,
				})
			}
			if len(mainnetAccEntries) > 0 {
				if err = entrywriter.MustWrite(
					accountSummary,
					campaignMainnetsAccSummaryHeader,
					mainnetAccEntries...,
				); err != nil {
					return err
				}
			}

			// get all campaign vesting accounts
			vestingAccs, err := n.MainnetVestingAccounts(cmd.Context(), campaignID)
			if err != nil {
				return err
			}
			mainnetVestingAccEntries := make([][]string, 0)
			for _, acc := range vestingAccs {
				mainnetVestingAccEntries = append(mainnetVestingAccEntries, []string{
					acc.Address,
					acc.TotalShares,
					acc.Vesting,
					strconv.FormatInt(acc.EndTime, 10),
				})
			}
			if len(mainnetVestingAccEntries) > 0 {
				if err = entrywriter.MustWrite(
					accountSummary,
					campaignVestingAccSummaryHeader,
					mainnetVestingAccEntries...,
				); err != nil {
					return err
				}
			}
			nb.Spinner.Stop()
			fmt.Print(accountSummary.String())
			return nil
		},
	}
	return c
}
