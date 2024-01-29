package ignitecmd

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	cmdmodel "github.com/ignite/cli/v28/ignite/cmd/model"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	uilog "github.com/ignite/cli/v28/ignite/pkg/cliui/log"
	cliuimodel "github.com/ignite/cli/v28/ignite/pkg/cliui/model"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/events"
	"github.com/ignite/cli/v28/ignite/services/chain"
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
	c.Flags().BoolP("verbose", "v", false, "verbose output")
	c.Flags().BoolP(flagForceReset, "f", false, "force reset of the app state on start and every source change")
	c.Flags().BoolP(flagResetOnce, "r", false, "reset the app state once on init")
	c.Flags().Bool(flagGenerateClients, false, "generate code for the configured clients on reset or source code change")
	c.Flags().Bool(flagQuitOnFail, false, "quit program if the app fails to start")
	c.Flags().StringSlice(flagBuildTags, []string{}, "parameters to build the chain binary")

	return c
}

func chainServeHandler(cmd *cobra.Command, _ []string) error {
	var options []cliui.Option

	// Session must not handle events when the verbosity is the default
	// to allow render of the UI and events using bubbletea. The custom
	// UI is not used for other verbosity levels in which the session
	// must handle the events to use custom output prefixes.
	verbosity := getVerbosity(cmd)
	if verbosity == uilog.VerbosityDefault {
		options = append(options, cliui.IgnoreEvents())
	} else {
		options = append(options, cliui.WithVerbosity(verbosity))
	}

	session := cliui.New(options...)
	defer session.End()

	// Depending on the verbosity execute the serve command within
	// a bubbletea context to display the custom UI.
	if verbosity == uilog.VerbosityDefault {
		bus := session.EventBus()
		bus.Send("Initializing...", events.ProgressStart())

		// Render UI
		m := cmdmodel.NewChainServe(cmd, bus, chainServeCmd(cmd, session))
		_, err := tea.NewProgram(m).Run()
		return err
	}

	// Otherwise run the serve command directly
	return chainServe(cmd, session)
}

func chainServeCmd(cmd *cobra.Command, session *cliui.Session) tea.Cmd {
	return func() tea.Msg {
		if err := chainServe(cmd, session); err != nil && !errors.Is(err, context.Canceled) {
			return cliuimodel.ErrorMsg{Error: err}
		}
		return cliuimodel.QuitMsg{}
	}
}

func chainServe(cmd *cobra.Command, session *cliui.Session) error {
	chainOption := []chain.Option{
		chain.WithOutputer(session),
		chain.CollectEvents(session.EventBus()),
		chain.CheckCosmosSDKVersion(),
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
	c, err := newChainWithHomeFlags(cmd, chainOption...)
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	// serve the chain
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

	buildTags, err := cmd.Flags().GetStringSlice(flagBuildTags)
	if err != nil {
		return err
	}

	if len(buildTags) > 0 {
		serveOptions = append(serveOptions, chain.BuildTags(buildTags...))
	}

	if flagGetSkipProto(cmd) {
		serveOptions = append(serveOptions, chain.ServeSkipProto())
	}

	if quitOnFail {
		serveOptions = append(serveOptions, chain.QuitOnFail())
	}

	return c.Serve(cmd.Context(), cacheStorage, serveOptions...)
}
