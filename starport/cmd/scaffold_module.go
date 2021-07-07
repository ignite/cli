package starportcmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/validation"
	"github.com/tendermint/starport/starport/services/scaffolder"
	modulecreate "github.com/tendermint/starport/starport/templates/module/create"
)

const (
	flagDep                 = "dep"
	flagIBC                 = "ibc"
	flagOracle              = "oracle"
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
	c.Flags().StringSlice(flagDep, []string{}, "module dependencies (e.g. --dep account,bank)")
	c.Flags().Bool(flagIBC, false, "scaffold an IBC module")
	c.Flags().Bool(flagOracle, false, "scaffold a bandchain oracle module")
	c.Flags().String(flagIBCOrdering, "none", "channel ordering of the IBC module [none|ordered|unordered]")
	c.Flags().Bool(flagRequireRegistration, false, "if true command will fail if module can't be registered")
	return c
}

func scaffoldModuleHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	var options []scaffolder.ModuleCreationOption

	name := args[0]

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

	// Check if the module must be an IBC module
	if ibcModule {
		options = append(options, scaffolder.WithIBCChannelOrdering(ibcOrdering), scaffolder.WithIBC())
	}

	oracle, err := cmd.Flags().GetBool(flagOracle)
	if err != nil {
		return err
	}
	if oracle {
		options = append(options, scaffolder.WithOracleIntegration())
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

	sc, err := scaffolder.New(appPath)
	if err != nil {
		return err
	}
	var msg bytes.Buffer
	fmt.Fprintf(&msg, "\n🎉 Module created %s.\n\n", name)
	sm, err := sc.CreateModule(placeholder.New(), name, options...)
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
		fmt.Println(sourceModificationToString(sm))
	}

	if len(dependencies) > 0 {
		dependencyWarning(dependencies)
	}

	io.Copy(cmd.OutOrStdout(), &msg)
	return nil
}

// in previously scaffolded apps gov keeper is defined below the scaffolded module keeper definition
// therefore we must warn the user to manually move the definition if it's the case
// https://github.com/tendermint/starport/issues/818#issuecomment-865736052
const govWarning = `⚠️ If your app has been scaffolded with Starport 0.16.x or below
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
