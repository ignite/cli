package ignitecmd

import (
	"strconv"

	"github.com/ignite-hq/cli/ignite/pkg/cliui"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/icons"
	"github.com/spf13/cobra"
)

func newNetworkChainShowAccounts() *cobra.Command {
	c := &cobra.Command{
		Use:   "accounts [launch-id]",
		Short: "Show all vesting and genesis accounts of the chain",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainShowAccountsHandler,
	}

	return c
}

func networkChainShowAccountsHandler(cmd *cobra.Command, args []string) error {
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

	// get all chain genesis accounts
	genesisAccs, err := n.GenesisAccounts(cmd.Context(), launchID)
	if err != nil {
		return err
	}
	genesisAccEntries := make([][]string, 0)
	for _, acc := range genesisAccs {
		genesisAccEntries = append(genesisAccEntries, []string{
			acc.Address,
			acc.Coins,
		})
	}
	if len(genesisAccEntries) > 0 {
		if err = session.PrintTable(chainGenesisAccSummaryHeader, genesisAccEntries...); err != nil {
			return err
		}
	}

	// get all chain vesting accounts
	vestingAccs, err := n.VestingAccounts(cmd.Context(), launchID)
	if err != nil {
		return err
	}
	genesisVestingAccEntries := make([][]string, 0)
	for _, acc := range vestingAccs {
		genesisVestingAccEntries = append(genesisVestingAccEntries, []string{
			acc.Address,
			acc.TotalBalance,
			acc.Vesting,
			strconv.FormatInt(acc.EndTime, 10),
		})
	}
	if len(genesisVestingAccEntries) > 0 {
		if err = session.PrintTable(chainVestingAccSummaryHeader, genesisVestingAccEntries...); err != nil {
			return err
		}
	}

	if len(genesisVestingAccEntries)+len(genesisAccEntries) == 0 {
		return session.Printf("%s %s\n", icons.Info, "empty chain account list")
	}

	return nil
}
