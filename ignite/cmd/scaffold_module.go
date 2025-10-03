package ignitecmd

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/services/scaffolder"
	modulecreate "github.com/ignite/cli/v29/ignite/templates/module/create"
)

// moduleNameKeeperAlias is a map of well known module names that have a different keeper name than the usual <module-name>Keeper.
var moduleNameKeeperAlias = map[string]string{}

const (
	flagDep                 = "dep"
	flagIBC                 = "ibc"
	flagParams              = "params"
	flagModuleConfigs       = "module-configs"
	flagIBCOrdering         = "ordering"
	flagRequireRegistration = "require-registration"
)

// NewScaffoldModule returns the command to scaffold a Cosmos SDK module.
func NewScaffoldModule() *cobra.Command {
	c := &cobra.Command{
		Use:   "module [name]",
		Short: "Custom Cosmos SDK module",
		Long: `Scaffold a new Cosmos SDK module.

Cosmos SDK is a modular framework and each independent piece of functionality is
implemented in a separate module. By default your blockchain imports a set of
standard Cosmos SDK modules. To implement custom functionality of your
blockchain, scaffold a module and implement the logic of your application.

This command does the following:

* Creates a directory with module's protocol buffer files in "proto/"
* Creates a directory with module's boilerplate Go code in "x/"
* Imports the newly created module by modifying "app/app.go"

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
sending tokens between accounts. The method for sending tokens is a defined in
the "bank"'s module keeper. You can scaffold a "foo" module with the dependency
on "bank" with the following command:

	ignite scaffold module foo --dep bank

You can then define which methods you want to import from the "bank" keeper in
"expected_keepers.go".

You can also scaffold a module with a list of dependencies that can include both
standard and custom modules (provided they exist):

	ignite scaffold module bar --dep foo,mint,account,FeeGrant

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
		Args:    cobra.ExactArgs(1),
		PreRunE: migrationPreRunHandler,
		RunE:    scaffoldModuleHandler,
	}

	flagSetPath(c)
	flagSetClearCache(c)

	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().StringSlice(flagDep, []string{}, "add a dependency on another module")
	c.Flags().Bool(flagIBC, false, "add IBC functionality")
	c.Flags().String(flagIBCOrdering, "none", "channel ordering of the IBC module [none|ordered|unordered]")
	c.Flags().Bool(flagRequireRegistration, false, "fail if module can't be registered")
	c.Flags().StringSlice(flagParams, []string{}, "add module parameters")
	c.Flags().StringSlice(flagModuleConfigs, []string{}, "add module configs")

	return c
}

func scaffoldModuleHandler(cmd *cobra.Command, args []string) error {
	var (
		name    = args[0]
		appPath = flagGetPath(cmd)
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

	ibcModule, _ := cmd.Flags().GetBool(flagIBC)
	ibcOrdering, _ := cmd.Flags().GetString(flagIBCOrdering)
	requireRegistration, _ := cmd.Flags().GetBool(flagRequireRegistration)
	params, _ := cmd.Flags().GetStringSlice(flagParams)

	moduleConfigs, err := cmd.Flags().GetStringSlice(flagModuleConfigs)
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	options := []scaffolder.ModuleCreationOption{
		scaffolder.WithParams(params),
		scaffolder.WithModuleConfigs(moduleConfigs),
	}

	// Check if the module must be an IBC module
	if ibcModule {
		options = append(options, scaffolder.WithIBCChannelOrdering(ibcOrdering), scaffolder.WithIBC())
	}

	// Get module dependencies
	dependencies, _ := cmd.Flags().GetStringSlice(flagDep)
	if len(dependencies) > 0 {
		var deps []modulecreate.Dependency

		isValid := regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

		for _, name := range dependencies {
			if !isValid(name) {
				return errors.Errorf("invalid module dependency name format '%s'", name)
			}

			if alias, ok := moduleNameKeeperAlias[strings.ToLower(name)]; ok {
				name = alias
			}

			deps = append(deps, modulecreate.NewDependency(name))
		}

		options = append(options, scaffolder.WithDependencies(deps))
	}

	var msg bytes.Buffer
	fmt.Fprintf(&msg, "\n🎉 Module created %s.\n\n", name)

	sc, err := scaffolder.New(cmd.Context(), appPath, cfg.Build.Proto.Path)
	if err != nil {
		return err
	}

	if err := sc.CreateModule(name, options...); err != nil {
		var validationErr errors.ValidationError
		if !requireRegistration && errors.As(err, &validationErr) {
			fmt.Fprintf(&msg, "Can't register module '%s'.\n", name)
			fmt.Fprintln(&msg, validationErr.ValidationInfo())
		} else {
			return err
		}
	}

	sm, err := sc.ApplyModifications(xgenny.ApplyPreRun(scaffolder.AskOverwriteFiles(session)))
	if err != nil {
		return err
	}

	if err := sc.PostScaffold(cmd.Context(), cacheStorage, false); err != nil {
		return err
	}

	modificationsStr, err := sm.String()
	if err != nil {
		return err
	}

	session.Println(modificationsStr)

	return session.Print(msg.String())
}
