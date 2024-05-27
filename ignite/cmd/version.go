package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/version"
)

// NewVersion creates a new version command to show the Ignite CLI version.
func NewVersion() *cobra.Command {
	c := &cobra.Command{
		Use:   "version",
		Short: "Print the current build information",
		RunE: func(cmd *cobra.Command, _ []string) error {
			v, err := version.Long(cmd.Context())
			if err != nil {
				return err
			}
			cmd.Println(v)
			return nil
		},
	}
	return c
}
