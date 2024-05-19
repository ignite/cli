package ignitecmd

import (
	"fmt"

	"github.com/ignite/cli/v29/ignite/services/scaffolder"
	"github.com/spf13/cobra"
)

// NewScaffoldType returns a new command to scaffold a type.
func NewScaffoldType() *cobra.Command {
	c := &cobra.Command{
		Use:     "type NAME [field:type] ...",
		Short:   "Type definition",
		Long:    fmt.Sprintf("Type information\n%s\n", supportFieldTypes),
		Example: "  ignite scaffold type todo-item priority:int desc:string tags:array.string done:bool",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: migrationPreRunHandler,
		RunE:    scaffoldTypeHandler,
	}

	flagSetPath(c)
	flagSetClearCache(c)

	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().AddFlagSet(flagSetScaffoldType())

	return c
}

func scaffoldTypeHandler(cmd *cobra.Command, args []string) error {
	return scaffoldType(cmd, args, scaffolder.DryType())
}
