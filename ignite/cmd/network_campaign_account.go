package starportcmd

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/ignite-hq/cli/ignite/pkg/clispinner"
	"github.com/ignite-hq/cli/ignite/pkg/entrywriter"
	"github.com/ignite-hq/cli/ignite/services/network"
	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
)

var (
	campaignMainnetsAccSummaryHeader = []string{"Mainnet Account", "Shares"}
	campaignVestingAccSummaryHeader  = []string{"Vesting Account", "Total Shares", "Vesting", "End Time"}
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
	nb, campaignID, err := networkChainLaunch(cmd, args)
	if err != nil {
		return err
	}
	defer nb.Cleanup()

	n, err := nb.Network()
	if err != nil {
		return err
	}

	accountSummary := &bytes.Buffer{}

	// get all campaign accounts
	mainnetAccs, vestingAccs, err := getAccounts(cmd.Context(), n, campaignID)
	if err != nil {
		return err
	}

	mainnetAccEntries := make([][]string, 0)
	for _, acc := range mainnetAccs {
		mainnetAccEntries = append(mainnetAccEntries, []string{
			acc.Address,
			acc.Shares.String(),
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

	mainnetVestingAccEntries := make([][]string, 0)
	for _, acc := range vestingAccs {
		mainnetVestingAccEntries = append(mainnetVestingAccEntries, []string{
			acc.Address,
			acc.TotalShares.String(),
			acc.Vesting.String(),
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
	if accountSummary.Len() > 0 {
		fmt.Print(accountSummary.String())
	} else {
		fmt.Printf("%s %s\n", clispinner.Info, "no campaign account found")
	}
	return nil
}

// getAccounts get all campaign mainnet and vesting accounts.
func getAccounts(
	ctx context.Context,
	n network.Network,
	campaignID uint64,
) (
	[]networktypes.MainnetAccount,
	[]networktypes.MainnetVestingAccount,
	error,
) {
	// start serving components.
	g, ctx := errgroup.WithContext(ctx)
	var (
		mainnetAccs []networktypes.MainnetAccount
		vestingAccs []networktypes.MainnetVestingAccount
		err         error
	)
	// get all campaign mainnet accounts
	g.Go(func() error {
		mainnetAccs, err = n.MainnetAccounts(ctx, campaignID)
		return err
	})

	// get all campaign vesting accounts
	g.Go(func() error {
		vestingAccs, err = n.MainnetVestingAccounts(ctx, campaignID)
		return err
	})
	return mainnetAccs, vestingAccs, g.Wait()
}
