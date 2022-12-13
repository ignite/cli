package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/yaml"
	"github.com/ignite/cli/ignite/services/network"
)

// NewNetworkProjectShow returns a new command to show published project on Ignite.
func NewNetworkProjectShow() *cobra.Command {
	c := &cobra.Command{
		Use:   "show [project-id]",
		Short: "Show published project",
		Args:  cobra.ExactArgs(1),
		RunE:  networkProjectShowHandler,
	}
	return c
}

func networkProjectShowHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	// parse project ID
	projectID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	project, err := n.Project(cmd.Context(), projectID)
	if err != nil {
		return err
	}

	info, err := yaml.Marshal(cmd.Context(), project)
	if err != nil {
		return err
	}

	return session.Println(info)
}
