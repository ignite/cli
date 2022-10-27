package ignitecmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
)

const (
	flagForceReset = "force-reset"
	flagResetOnce  = "reset-once"
	flagConfig     = "config"
	flagQuitOnFail = "quit-on-fail"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New(cliui.WithVerbosity(getVerbosity(cmd)))
			defer session.End()

			// TODO: How to deal with verbose mode?
			// TODO: Should all commands use bubbletea pattern?
			m := initialChainServeModel(cmd, session)
			if err := tea.NewProgram(m).Start(); err != nil {
				return err
			}

			return nil
		},
	}

	flagSetPath(c)
	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetProto3rdParty(""))
	c.Flags().AddFlagSet(flagSetCheckDependencies())
	c.Flags().AddFlagSet(flagSetSkipProto())
	c.Flags().BoolP("verbose", "v", false, "Verbose output")
	c.Flags().BoolP(flagForceReset, "f", false, "Force reset of the app state on start and every source change")
	c.Flags().BoolP(flagResetOnce, "r", false, "Reset of the app state on first start")
	c.Flags().Bool(flagQuitOnFail, false, "Quit program if the app fails to start")

	return c
}
