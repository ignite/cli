package ignitecmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui/clispinner"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/validation"
	"github.com/ignite/cli/ignite/services/scaffolder"
	modulecreate "github.com/ignite/cli/ignite/templates/module/create"
)

const (
	flagDep                 = "dep"
	flagIBC                 = "ibc"
	flagParams              = "params"
	flagIBCOrdering         = "ordering"
	flagRequireRegistration = "require-registration"
)

// NewScaffoldModule returns the command to scaffold a Cosmos SDK module
func NewScaffoldModule() *cobra.Command {
	c := &cobra.Command{
		Use:   "module [name]",
		Short: "Scaffold a Cosmos SDK module",
		Long:  "Scaffold a new Cosmos SDK module in the `x` directory",
		Args:  cobra.MinimumNArgs(1),
		RunE:  scaffoldModuleHandler,
	}

	flagSetPath(c)
	flagSetClearCache(c)
	c.Flags().StringSlice(flagDep, []string{}, "module dependencies (e.g. --dep account,bank)")
	c.Flags().Bool(flagIBC, false, "scaffold an IBC module")
	c.Flags().String(flagIBCOrdering, "none", "channel ordering of the IBC module [none|ordered|unordered]")
	c.Flags().Bool(flagRequireRegistration, false, "if true command will fail if module can't be registered")
	c.Flags().StringSlice(flagParams, []string{}, "scaffold module params")

	return c
}

func scaffoldModuleHandler(cmd *cobra.Command, args []string) error {
	var (
		name    = args[0]
		appPath = flagGetPath(cmd)
	)
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	ibcModule, err := cmd.Flags().GetBool(flagIBC)
	if err != nil {
		return err
	}

	ibcOrdering, err := cmd.Flags().GetString(flagIBCOrdering)
	if err != nil {
		return err
	}
	requireRegistration, err := cmd.Flags().GetBool(flagRequireRegistration)
	if err != nil {
		return err
	}

	params, err := cmd.Flags().GetStringSlice(flagParams)
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	options := []scaffolder.ModuleCreationOption{
		scaffolder.WithParams(params),
	}

	// Check if the module must be an IBC module
	if ibcModule {
		options = append(options, scaffolder.WithIBCChannelOrdering(ibcOrdering), scaffolder.WithIBC())
	}

	// Get module dependencies
	dependencies, err := cmd.Flags().GetStringSlice(flagDep)
	if err != nil {
		return err
	}
	if len(dependencies) > 0 {
		var formattedDependencies []modulecreate.Dependency

		// Parse the provided dependencies
		for _, dependency := range dependencies {
			var formattedDependency modulecreate.Dependency

			splitted := strings.Split(dependency, ":")
			switch len(splitted) {
			case 1:
				formattedDependency = modulecreate.NewDependency(splitted[0], "")
			case 2:
				formattedDependency = modulecreate.NewDependency(splitted[0], splitted[1])
			default:
				return fmt.Errorf("dependency %s is invalid, must have <depName> or <depName>.<depKeeperName>", dependency)
			}
			formattedDependencies = append(formattedDependencies, formattedDependency)
		}
		options = append(options, scaffolder.WithDependencies(formattedDependencies))
	}

	var msg bytes.Buffer
	fmt.Fprintf(&msg, "\n🎉 Module created %s.\n\n", name)

	sc, err := newApp(appPath)
	if err != nil {
		return err
	}

	sm, err := sc.CreateModule(cacheStorage, placeholder.New(), name, options...)
	s.Stop()
	if err != nil {
		var validationErr validation.Error
		if !requireRegistration && errors.As(err, &validationErr) {
			fmt.Fprintf(&msg, "Can't register module '%s'.\n", name)
			fmt.Fprintln(&msg, validationErr.ValidationInfo())
		} else {
			return err
		}
	} else {
		modificationsStr, err := sourceModificationToString(sm)
		if err != nil {
			return err
		}

		fmt.Println(modificationsStr)
	}

	if len(dependencies) > 0 {
		dependencyWarning(dependencies)
	}

	io.Copy(cmd.OutOrStdout(), &msg)
	return nil
}

// in previously scaffolded apps gov keeper is defined below the scaffolded module keeper definition
// therefore we must warn the user to manually move the definition if it's the case
// https://github.com/ignite/cli/issues/818#issuecomment-865736052
const govWarning = `⚠️ If your app has been scaffolded with Ignite CLI 0.16.x or below
Please make sure that your module keeper definition is defined after gov module keeper definition in app/app.go:

app.GovKeeper = ...
...
[your module keeper definition]
`

// dependencyWarning is used to print a warning if gov is provided as a dependency
func dependencyWarning(dependencies []string) {
	for _, dep := range dependencies {
		if dep == "gov" {
			fmt.Print(govWarning)
		}
	}
}
