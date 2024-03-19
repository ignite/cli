package ignitecmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/services/scaffolder"
)

// NewScaffoldConfigs returns the command to scaffold a Cosmos SDK configs into a module.
func NewScaffoldConfigs() *cobra.Command {
	c := &cobra.Command{
		Use:   "configs [configs]...",
		Short: "Configs for a custom Cosmos SDK module",
		Long: `Scaffold a new config for a Cosmos SDK module.

A Cosmos SDK module can have configurations. An example of a config is "address prefix" of the
"auth" module. A config can be scaffolded into a module using the "--module-configs" into
the scaffold module command or using the "scaffold configs" command. By default 
configs are of type "string", but you can specify a type for each config. For example:

	ignite scaffold configs foo baz:uint bar:bool

Refer to Cosmos SDK documentation to learn more about modules, dependencies and
configs.
`,
		Args:    cobra.MinimumNArgs(1),
		PreRunE: migrationPreRunHandler,
		RunE:    scaffoldConfigsHandler,
	}

	flagSetPath(c)
	flagSetClearCache(c)

	c.Flags().AddFlagSet(flagSetYes())

	c.Flags().String(flagModule, "", "module to add the query into (default: app's main module)")

	return c
}

func scaffoldConfigsHandler(cmd *cobra.Command, args []string) error {
	var (
		configs    = args[0:]
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

	runner := xgenny.NewRunner(cmd.Context(), appPath)
	err = sc.CreateConfigs(runner, moduleName, configs...)
	if err != nil {
		return err
	}

	modificationsStr, err := runner.ApplyModifications()
	if err != nil {
		return err
	}

	if err := sc.PostScaffold(cmd.Context(), cacheStorage, false); err != nil {
		return err
	}

	session.Println(modificationsStr)
	session.Printf("\n🎉 New configs added to the module:\n\n- %s\n\n", strings.Join(configs, "\n- "))

	return nil
}
