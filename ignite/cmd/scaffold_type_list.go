package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/templates/field/datatype"
)

// NewScaffoldTypeList returns a new command to list all scaffold types.
func NewScaffoldTypeList() *cobra.Command {
	c := &cobra.Command{
		Use:   "type-list",
		Short: "List scaffold types",
		Long:  "List all available scaffold types",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			session := cliui.New(
				cliui.StartSpinnerWithText("printing..."),
				cliui.WithoutUserInteraction(getYes(cmd)),
			)
			defer session.End()
			session.StopSpinner()
			return datatype.PrintScaffoldTypeList(cmd.OutOrStdout())
		},
	}

	return c
}
