package ignitecmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/services/scaffolder"
)

// NewScaffoldParams returns the command to scaffold a Cosmos SDK parameters into a module.
func NewScaffoldParams() *cobra.Command {
	c := &cobra.Command{
		Use:   "params [param]...",
		Short: "Parameters for a custom Cosmos SDK module",
		Long: `Scaffold a new parameter for a Cosmos SDK module.

A Cosmos SDK module can have parameters (or "params"). Params are values that
can be set at the genesis of the blockchain and can be modified while the
blockchain is running. An example of a param is "Inflation rate change" of the
"mint" module. A params can be scaffolded into a module using the "--params" into
the scaffold module command or using the scaffold params command. By default 
params are of type "string", but you can specify a type for each param. For example:

	ignite scaffold params foo baz:uint bar:bool

Refer to Cosmos SDK documentation to learn more about modules, dependencies and
params.
`,
		Args:    cobra.MinimumNArgs(1),
		PreRunE: migrationPreRunHandler,
		RunE:    scaffoldParamsHandler,
	}

	flagSetPath(c)
	flagSetClearCache(c)

	c.Flags().AddFlagSet(flagSetYes())

	c.Flags().String(flagModule, "", "module to add the query into. Default: app's main module")

	return c
}

func scaffoldParamsHandler(cmd *cobra.Command, args []string) error {
	var (
		params     = args[0:]
		appPath    = flagGetPath(cmd)
		moduleName = flagGetModule(cmd)
	)

	session := cliui.New(cliui.StartSpinnerWithText(statusScaffolding))
	defer session.End()

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	sc, err := scaffolder.New(appPath)
	if err != nil {
		return err
	}

	sm, err := sc.CreateParams(cmd.Context(), cacheStorage, placeholder.New(), moduleName, params...)
	if err != nil {
		return err
	}

	modificationsStr, err := sourceModificationToString(sm)
	if err != nil {
		return err
	}

	session.Println(modificationsStr)
	session.Printf("\n🎉 New parameters added to the module:\n- %s.\n\n", moduleName, strings.Join(params, "\n- "))

	return nil
}
