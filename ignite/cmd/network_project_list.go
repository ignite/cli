package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/entrywriter"
	"github.com/ignite/cli/ignite/services/network/networktypes"
)

var ProjectSummaryHeader = []string{
	"id",
	"name",
	"coordinator id",
	"mainnet id",
}

// NewNetworkProjectList returns a new command to list all published Projects on Ignite.
func NewNetworkProjectList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List published projects",
		Args:  cobra.NoArgs,
		RunE:  networkProjectListHandler,
	}
	return c
}

func networkProjectListHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(cliui.StartSpinner())

	defer session.End()

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}
	projects, err := n.Projects(cmd.Context())
	if err != nil {
		return err
	}

	return renderProjectSummaries(projects, session)
}

// renderProjectSummaries writes into the provided out, the list of summarized projects
func renderProjectSummaries(projects []networktypes.Project, session *cliui.Session) error {
	var projectEntries [][]string

	for _, c := range projects {
		mainnetID := entrywriter.None
		if c.MainnetInitialized {
			mainnetID = fmt.Sprintf("%d", c.MainnetID)
		}

		projectEntries = append(projectEntries, []string{
			fmt.Sprintf("%d", c.ID),
			c.Name,
			fmt.Sprintf("%d", c.CoordinatorID),
			mainnetID,
		})
	}

	return session.PrintTable(ProjectSummaryHeader, projectEntries...)
}
