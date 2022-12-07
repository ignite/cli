package ignitecmd

import (
	"github.com/spf13/cobra"
)

// NewNetworkProject creates a new project command that holds other
// subcommands related to launching a network for a project.
func NewNetworkProject() *cobra.Command {
	c := &cobra.Command{
		Use:    "project",
		Short:  "Handle projects",
		Hidden: true,
	}
	c.AddCommand(
		NewNetworkProjectPublish(),
		NewNetworkProjectList(),
		NewNetworkProjectShow(),
		NewNetworkProjectUpdate(),
		NewNetworkProjectAccount(),
	)
	return c
}
