package ignitecmd

import (
	"github.com/ignite/cli/v29/ignite/version"
	"github.com/spf13/cobra"
)

// NewVersion creates a new version command to show the Ignite CLI version.
func NewVersion() *cobra.Command {
	c := &cobra.Command{
		Use:   "version",
		Short: "Print the current build information",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Println(version.Long(cmd.Context()))
		},
	}
	return c
}
