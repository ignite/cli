package starportcmd

import (
	"fmt"

	"github.com/ignite-hq/cli/starport/internal/version"
	"github.com/spf13/cobra"
)

// NewVersion creates a new version command to show starport's version.
func NewVersion() *cobra.Command {
	c := &cobra.Command{
		Use:   "version",
		Short: "Print the current build information",
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Println(version.Long(cmd.Context()))
		},
	}
	return c
}
