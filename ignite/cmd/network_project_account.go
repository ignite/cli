package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
)

var projectMainnetsAccSummaryHeader = []string{"Mainnet Account", "Shares"}

// NewNetworkProjectAccount creates a new project account command that holds some other
// sub commands related to account for a project.
func NewNetworkProjectAccount() *cobra.Command {
	c := &cobra.Command{
		Use:   "account",
		Short: "Handle project accounts",
	}
	c.AddCommand(
		newNetworkProjectAccountList(),
	)
	return c
}

func newNetworkProjectAccountList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list [project-id]",
		Short: "Show all mainnet and mainnet vesting of the project",
		Args:  cobra.ExactArgs(1),
		RunE:  newNetworkProjectAccountListHandler,
	}
	return c
}

func newNetworkProjectAccountListHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	nb, projectID, err := networkChainLaunch(cmd, args, session)
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	// get all project accounts
	mainnetAccs, err := n.MainnetAccounts(cmd.Context(), projectID)
	if err != nil {
		return err
	}

	if len(mainnetAccs) == 0 {
		return session.Printf("%s no project account found\n", icons.Info)
	}

	mainnetAccEntries := make([][]string, 0)
	for _, acc := range mainnetAccs {
		mainnetAccEntries = append(mainnetAccEntries, []string{acc.Address, acc.Shares.String()})
	}

	if len(mainnetAccEntries) > 0 {
		if err = session.PrintTable(projectMainnetsAccSummaryHeader, mainnetAccEntries...); err != nil {
			return err
		}
	}

	return nil
}
