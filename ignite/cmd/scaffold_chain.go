package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/config/chain/defaults"
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/env"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xfilepath"
	"github.com/ignite/cli/v29/ignite/pkg/xgit"
	"github.com/ignite/cli/v29/ignite/services/scaffolder"
)

const defaultIgniteDenom = "stake"

const (
	flagMinimal         = "minimal"
	flagNoDefaultModule = "no-module"
	flagSkipGit         = "skip-git"
	flagDefaultDenom    = "default-denom"

	tplScaffoldChainSuccess = `
⭐️ Successfully created a new blockchain '%[1]v'.
👉 Get started with the following commands:

 %% cd %[1]v
 %% ignite chain serve

Documentation: https://docs.ignite.com
`
)

// NewScaffoldChain creates new command to scaffold a Comos-SDK based blockchain.
func NewScaffoldChain() *cobra.Command {
	c := &cobra.Command{
		Use:   "chain [name]",
		Short: "New Cosmos SDK blockchain",
		Long: `Create a new application-specific Cosmos SDK blockchain.

For example, the following command will create a blockchain called "hello" in
the "hello/" directory:

	ignite scaffold chain hello

A project name can be a simple name or a URL. The name will be used as the Go
module path for the project. Examples of project names:

	ignite scaffold chain foo
	ignite scaffold chain foo/bar
	ignite scaffold chain example.org/foo
	ignite scaffold chain github.com/username/foo
		
A new directory with source code files will be created in the current directory.
To use a different path use the "--path" flag.

Most of the logic of your blockchain is written in custom modules. Each module
effectively encapsulates an independent piece of functionality. Following the
Cosmos SDK convention, custom modules are stored inside the "x/" directory. By
default, Ignite creates a module with a name that matches the name of the
project. To create a blockchain without a default module use the "--no-module"
flag. Additional modules can be added after a project is created with "ignite
scaffold module" command.

Account addresses on Cosmos SDK-based blockchains have string prefixes. For
example, the Cosmos Hub blockchain uses the default "cosmos" prefix, so that
addresses look like this: "cosmos12fjzdtqfrrve7zyg9sv8j25azw2ua6tvu07ypf". To
use a custom address prefix use the "--address-prefix" flag. For example:

	ignite scaffold chain foo --address-prefix bar

By default when compiling a blockchain's source code Ignite creates a cache to
speed up the build process. To clear the cache when building a blockchain use
the "--clear-cache" flag. It is very unlikely you will ever need to use this
flag.

The blockchain is using the Cosmos SDK modular blockchain framework. Learn more
about Cosmos SDK on https://docs.cosmos.network
`,
		Args: cobra.ExactArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			if verbose := flagGetVerbose(cmd); verbose {
				env.SetDebug()
			}
		},
		RunE: scaffoldChainHandler,
	}

	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagSetAccountPrefixes())
	c.Flags().AddFlagSet(flagSetCoinType())
	c.Flags().String(flagDefaultDenom, defaultIgniteDenom, "default staking denom")
	c.Flags().StringP(flagPath, "p", "", "create a project in a specific path")
	c.Flags().Bool(flagNoDefaultModule, false, "create a project without a default module")
	c.Flags().StringSlice(flagParams, []string{}, "add default module parameters")
	c.Flags().StringSlice(flagModuleConfigs, []string{}, "add module configs")
	c.Flags().Bool(flagSkipGit, false, "skip Git repository initialization")
	c.Flags().Bool(flagSkipProto, false, "skip proto generation")
	c.Flags().Bool(flagMinimal, false, "create a minimal blockchain (with the minimum required Cosmos SDK modules)")
	c.Flags().String(flagProtoDir, defaults.ProtoDir, "chain proto directory")

	// consumer scaffolding have been migrated to an ignite app
	_ = c.Flags().MarkDeprecated("consumer", "use 'ignite consumer' app instead")

	return c
}

func scaffoldChainHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(
		cliui.StartSpinnerWithText(statusScaffolding),
		cliui.WithoutUserInteraction(getYes(cmd)),
	)
	defer session.End()

	var (
		name          = args[0]
		addressPrefix = getAddressPrefix(cmd)
		coinType      = getCoinType(cmd)
		appPath       = flagGetPath(cmd)

		noDefaultModule, _ = cmd.Flags().GetBool(flagNoDefaultModule)
		skipGit, _         = cmd.Flags().GetBool(flagSkipGit)
		minimal, _         = cmd.Flags().GetBool(flagMinimal)
		params, _          = cmd.Flags().GetStringSlice(flagParams)
		moduleConfigs, _   = cmd.Flags().GetStringSlice(flagModuleConfigs)
		skipProto, _       = cmd.Flags().GetBool(flagSkipProto)
		protoDir, _        = cmd.Flags().GetString(flagProtoDir)
		defaultDenom, _    = cmd.Flags().GetString(flagDefaultDenom)
	)

	if cmd.Flags().Changed(flagDefaultDenom) && len(defaultDenom) <= 2 {
		return errors.New("default denom must be at least 3 characters and maximum 128 characters")
	}

	if noDefaultModule {
		if len(params) > 0 {
			return errors.New("params flag is only supported if the default module is enabled")
		} else if len(moduleConfigs) > 0 {
			return errors.New("module configs flag is only supported if the default module is enabled")
		}
		skipProto = true
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	appDir, goModule, err := scaffolder.Init(
		cmd.Context(),
		appPath,
		name,
		addressPrefix,
		coinType,
		defaultDenom,
		protoDir,
		noDefaultModule,
		minimal,
		params,
		moduleConfigs,
	)
	if err != nil {
		return err
	}

	path, err := xfilepath.RelativePath(appDir)
	if err != nil {
		return err
	}

	if err := scaffolder.PostScaffold(cmd.Context(), cacheStorage, appDir, protoDir, goModule, skipProto); err != nil {
		return err
	}

	if !skipGit {
		// Initialize git repository and perform the first commit
		if err := xgit.InitAndCommit(path); err != nil {
			return err
		}
	}

	return session.Printf(tplScaffoldChainSuccess, path)
}
