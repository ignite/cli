package ignitecmd

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/pkg/cliui"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/icons"
)

var (
	chainGenesisAccSummaryHeader = []string{"Genesis Account", "Coins"}
	chainVestingAccSummaryHeader = []string{"Vesting Account", "Total Balance", "Vesting", "EndTime"}
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
	session := cliui.New(cliui.StartSpinner())
	defer session.Cleanup()

	nb, launchID, err := networkChainLaunch(cmd, args, session)
	if err != nil {
		return err
	}
	n, err := nb.Network()
	if err != nil {
		return err
	}

	genesisAccs, err := n.GenesisAccounts(cmd.Context(), launchID)
	if err != nil {
		return err
	}
	vestingAccs, err := n.VestingAccounts(cmd.Context(), launchID)
	if err != nil {
		return err
	}
	if len(genesisAccs)+len(vestingAccs) == 0 {
		session.StopSpinner()
		return session.Printf("%s %s\n", icons.Info, "empty chain account list")
	}

	genesisAccEntries := make([][]string, 0)
	for _, acc := range genesisAccs {
		genesisAccEntries = append(genesisAccEntries, []string{acc.Address, acc.Coins})
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

	session.StopSpinner()
	if len(genesisAccEntries) > 0 {
		if err = session.PrintTable(chainGenesisAccSummaryHeader, genesisAccEntries...); err != nil {
			return err
		}
	}
	if len(genesisVestingAccEntries) > 0 {
		if err = session.PrintTable(chainVestingAccSummaryHeader, genesisVestingAccEntries...); err != nil {
			return err
		}
	}

	return nil
}
