package ignitecmd

import (
	"sort"

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
		RunE:  scaffoldTypeListHandler,
	}

	return c
}

func scaffoldTypeListHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(
		cliui.StartSpinnerWithText("printing..."),
		cliui.WithoutUserInteraction(getYes(cmd)),
	)
	defer session.End()

	supported := datatype.SupportedTypes()
	entries := make([][]string, 0, len(supported))
	for name, usage := range supported {
		entries = append(entries, []string{name, usage})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i][0] < entries[j][0]
	})

	return session.PrintTable([]string{"types", "usage"}, entries...)
}
