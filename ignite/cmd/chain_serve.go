package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/services/chain"
)

const (
	flagConfig          = "config"
	flagForceReset      = "force-reset"
	flagGenerateClients = "generate-clients"
	flagQuitOnFail      = "quit-on-fail"
	flagResetOnce       = "reset-once"
)

// NewChainServe creates a new serve command to serve a blockchain.
func NewChainServe() *cobra.Command {
	c := &cobra.Command{
		Use:   "serve",
		Short: "Start a blockchain node in development",
		Long: `The serve command compiles and installs the binary (like "ignite chain build"),
uses that binary to initialize the blockchain's data directory for one validator
(like "ignite chain init"), and starts the node locally for development purposes
with automatic code reloading.

Automatic code reloading means Ignite starts watching the project directory.
Whenever a file change is detected, Ignite automatically rebuilds, reinitializes
and restarts the node.

Whenever possible Ignite will try to keep the current state of the chain by
exporting and importing the genesis file.

To force Ignite to start from a clean slate even if a genesis file exists, use
the following flag:

	ignite chain serve --reset-once

To force Ignite to reset the state every time the source code is modified, use
the following flag:

	ignite chain serve --force-reset

With Ignite it's possible to start more than one blockchain from the same source
code using different config files. This is handy if you're building
inter-blockchain functionality and, for example, want to try sending packets
from one blockchain to another. To start a node using a specific config file:

	ignite chain serve --config mars.yml

The serve command is meant to be used ONLY FOR DEVELOPMENT PURPOSES. Under the
hood, it runs "appd start", where "appd" is the name of your chain's binary. For
production, you may want to run "appd start" manually.
`,
		Args: cobra.NoArgs,
		RunE: chainServeHandler,
	}

	flagSetPath(c)
	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetCheckDependencies())
	c.Flags().AddFlagSet(flagSetSkipProto())
	c.Flags().BoolP("verbose", "v", false, "Verbose output")
	c.Flags().BoolP(flagForceReset, "f", false, "Force reset of the app state on start and every source change")
	c.Flags().BoolP(flagResetOnce, "r", false, "Reset of the app state on first start")
	c.Flags().Bool(flagGenerateClients, false, "Generate code for the configured clients on reset or source code change")
	c.Flags().Bool(flagQuitOnFail, false, "Quit program if the app fails to start")

	return c
}

func chainServeHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(
		cliui.WithVerbosity(getVerbosity(cmd)),
		cliui.StartSpinner(),
	)
	defer session.End()

	chainOption := []chain.Option{
		chain.WithOutputer(session),
		chain.CollectEvents(session.EventBus()),
	}

	if flagGetCheckDependencies(cmd) {
		chainOption = append(chainOption, chain.CheckDependencies())
	}

	// check if custom config is defined
	config, err := cmd.Flags().GetString(flagConfig)
	if err != nil {
		return err
	}
	if config != "" {
		chainOption = append(chainOption, chain.ConfigFile(config))
	}

	// create the chain
	c, err := NewChainWithHomeFlags(cmd, chainOption...)
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	var serveOptions []chain.ServeOption

	forceUpdate, err := cmd.Flags().GetBool(flagForceReset)
	if err != nil {
		return err
	}

	if forceUpdate {
		serveOptions = append(serveOptions, chain.ServeForceReset())
	}

	resetOnce, err := cmd.Flags().GetBool(flagResetOnce)
	if err != nil {
		return err
	}

	if resetOnce {
		serveOptions = append(serveOptions, chain.ServeResetOnce())
	}

	quitOnFail, err := cmd.Flags().GetBool(flagQuitOnFail)
	if err != nil {
		return err
	}

	if quitOnFail {
		serveOptions = append(serveOptions, chain.QuitOnFail())
	}

	generateClients, err := cmd.Flags().GetBool(flagGenerateClients)
	if err != nil {
		return err
	}

	if generateClients {
		serveOptions = append(serveOptions, chain.GenerateClients())
	}

	if flagGetSkipProto(cmd) {
		serveOptions = append(serveOptions, chain.ServeSkipProto())
	}

	return c.Serve(cmd.Context(), cacheStorage, serveOptions...)
}
