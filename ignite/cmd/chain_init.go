package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v28/ignite/pkg/chaincmd"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v28/ignite/services/chain"
)

func NewChainInit() *cobra.Command {
	c := &cobra.Command{
		Use:   "init",
		Short: "Initialize your chain",
		Long: `The init command compiles and installs the binary (like "ignite chain build")
and uses that binary to initialize the blockchain's data directory for one
validator. To learn how the build process works, refer to "ignite chain build
--help".

By default, the data directory will be initialized in $HOME/.mychain, where
"mychain" is the name of the project. To set a custom data directory use the
--home flag or set the value in config.yml:

	validators:
	  - name: alice
	    bonded: '100000000stake'
	    home: "~/.customdir"

The data directory contains three files in the "config" directory: app.toml,
config.toml, client.toml. These files let you customize the behavior of your
blockchain node and the client executable. When a chain is re-initialized the
data directory can be reset. To make some values in these files persistent, set
them in config.yml:

	validators:
	  - name: alice
	    bonded: '100000000stake'
	    app:
	      minimum-gas-prices: "0.025stake"
	    config:
	      consensus:
	        timeout_commit: "5s"
	        timeout_propose: "5s"
	    client:
	      output: "json"

The configuration above changes the minimum gas price of the validator (by
default the gas price is set to 0 to allow "free" transactions), sets the block
time to 5s, and changes the output format to JSON. To see what kind of values
this configuration accepts see the generated TOML files in the data directory.

As part of the initialization process Ignite creates on-chain accounts with
token balances. By default, config.yml has two accounts in the top-level
"accounts" property. You can add more accounts and change their token balances.
Refer to config.yml guide to see which values you can set.

One of these accounts is a validator account and the amount of self-delegated
tokens can be set in the top-level "validator" property.

One of the most important components of an initialized chain is the genesis
file, the 0th block of the chain. The genesis file is stored in the data
directory "config" subdirectory and contains the initial state of the chain,
including consensus and module parameters. You can customize the values of the
genesis in config.yml:

	genesis:
	  app_state:
	    staking:
	      params:
	        bond_denom: "foo"

The example above changes the staking token to "foo". If you change the staking
denom, make sure the validator account has the right tokens.

The init command is meant to be used ONLY FOR DEVELOPMENT PURPOSES. Under the
hood it runs commands like "appd init", "appd add-genesis-account", "appd
gentx", and "appd collect-gentx". For production, you may want to run these
commands manually to ensure a production-level node initialization.
`,
		Args: cobra.NoArgs,
		RunE: chainInitHandler,
	}

	flagSetPath(c)
	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetCheckDependencies())
	c.Flags().AddFlagSet(flagSetSkipProto())
	c.Flags().AddFlagSet(flagSetDebug())
	c.Flags().StringSlice(flagBuildTags, []string{}, "parameters to build the chain binary")

	return c
}

func chainInitHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(
		cliui.WithVerbosity(getVerbosity(cmd)),
		cliui.StartSpinner(),
	)
	defer session.End()

	chainOption := []chain.Option{
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
		chain.WithOutputer(session),
		chain.CollectEvents(session.EventBus()),
		chain.CheckCosmosSDKVersion(),
	}

	if flagGetCheckDependencies(cmd) {
		chainOption = append(chainOption, chain.CheckDependencies())
	}

	c, err := newChainWithHomeFlags(cmd, chainOption...)
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	var (
		ctx          = cmd.Context()
		buildTags, _ = cmd.Flags().GetStringSlice(flagBuildTags)
	)
	if _, err = c.Build(ctx, cacheStorage, buildTags, "", flagGetSkipProto(cmd), flagGetDebug(cmd)); err != nil {
		return err
	}

	if err := c.Init(ctx, chain.InitArgsAll); err != nil {
		return err
	}

	home, err := c.Home()
	if err != nil {
		return err
	}

	return session.Printf("🗃  Initialized. Checkout your chain's home (data) directory: %s\n", colors.Info(home))
}
