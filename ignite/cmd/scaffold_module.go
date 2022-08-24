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
		Long: `Scaffold a new Cosmos SDK module.

Cosmos SDK is a modular framework and each independent piece of functionality is
implemented in a separate module. By default your blockchain imports a set of
standard Cosmos SDK modules. To implement custom functionality of your
blockchain, scaffold a module and implement the logic of your application.

This command does the following:

* Creates a directory with module's protocol buffer files in "proto/"
* Creates a directory with module's boilerplate Go code in "x/"
* Imports the newly created module by modifying "app/app.go"
* Creates a file in "testutil/keeper/" that contains logic to create a keeper
  for testing purposes

This command will proceed with module scaffolding even if "app/app.go" doesn't
have the required default placeholders. If the placeholders are missing, you
will need to modify "app/app.go" manually to import the module. If you want the
command to fail if it can't import the module, use the "--require-registration"
flag.

To scaffold an IBC-enabled module use the "--ibc" flag. An IBC-enabled module is
like a regular module with the addition of IBC-specific logic and placeholders
to scaffold IBC packets with "ignite scaffold packet".

A module can depend on one or more other modules and import their keeper
methods. To scaffold a module with a dependency use the "--dep" flag

For example, your new custom module "foo" might have functionality that requires
sending tokens between accounts. The method for sending tokens is defined in the
bank's module keeper. You can scaffold a "foo" module with the dependency on
"bank" with the following command:

  ignite scaffold module foo --dep bank

You can then define which methods you want to import from the "bank" keeper in
"expected_keepers.go".

Provided with "--dep bank" Ignite looks up "BankKeeper" in "app.go". If
"BankKeeper" doesn't exist in app.go, module scaffolding will fail. Sometimes
you might want to specify the keeper name explicitly after a colon, for example:

  ignite scaffold module foo --dep feegrant:FeeGrantKeeper

You can also scaffold a module with a list of dependencies that can include both
standard and custom modules (provided they exist):

  ignite scaffold module bar --dep foo,mint,account

Note: the "--dep" flag doesn't install third-party modules into your
application, it just generates extra code that specifies which existing modules
your new custom module depends on.

A Cosmos SDK module can have parameters (or "params"). Params are values that
can be set at the genesis of the blockchain and can be modified while the
blockchain is running. An example of a param is "Inflation rate change" of the
"mint" module. A module can be scaffolded with params using the "--params" flag
that accepts a list of param names. By default params are of type "string", but
you can specify a type for each param. For example:

  ignite scaffold module foo --params baz:uint,bar:bool

Refer to Cosmos SDK documentation to learn more about modules, dependencies and
params.
`,
		Args: cobra.ExactArgs(1),
		RunE: scaffoldModuleHandler,
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
