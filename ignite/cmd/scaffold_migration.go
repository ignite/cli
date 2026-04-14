package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/services/scaffolder"
)

// NewScaffoldMigration returns the command to scaffold a module migration.
func NewScaffoldMigration() *cobra.Command {
	c := &cobra.Command{
		Use:   "migration [module]",
		Short: "Module migration boilerplate",
		Long: `Scaffold no-op migration boilerplate for an existing Cosmos SDK module.

This command creates a new migration file in "x/<module>/migrations/vN/",
increments the module consensus version, and registers the new migration handler
inside "x/<module>/module/module.go".`,
		Args:    cobra.ExactArgs(1),
		PreRunE: migrationPreRunHandler,
		RunE:    scaffoldMigrationHandler,
	}

	flagSetPath(c)
	c.Flags().AddFlagSet(flagSetYes())

	return c
}

func scaffoldMigrationHandler(cmd *cobra.Command, args []string) error {
	var (
		moduleName = args[0]
		appPath    = flagGetPath(cmd)
	)

	session := cliui.New(
		cliui.StartSpinnerWithText(statusScaffolding),
		cliui.WithoutUserInteraction(getYes(cmd)),
	)
	defer session.End()

	cfg, _, err := getChainConfig(cmd)
	if err != nil {
		return err
	}

	sc, err := scaffolder.New(cmd.Context(), appPath, cfg.Build.Proto.Path)
	if err != nil {
		return err
	}

	if err := sc.CreateModuleMigration(moduleName); err != nil {
		return err
	}

	sm, err := sc.ApplyModifications(xgenny.ApplyPreRun(scaffolder.AskOverwriteFiles(session)))
	if err != nil {
		return err
	}

	if err := sc.PostScaffold(cmd.Context(), cache.Storage{}, true); err != nil {
		return err
	}

	modificationsStr, err := sm.String()
	if err != nil {
		return err
	}

	session.Println(modificationsStr)
	session.Printf("\n🎉 Migration added to module %s.\n\n", moduleName)

	return nil
}
