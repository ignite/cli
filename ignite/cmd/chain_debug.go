package ignitecmd

import (
	"context"
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	cmdmodel "github.com/ignite/cli/ignite/cmd/model"
	"github.com/ignite/cli/ignite/pkg/chaincmd"
	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	cliuimodel "github.com/ignite/cli/ignite/pkg/cliui/model"
	"github.com/ignite/cli/ignite/pkg/debugger"
	"github.com/ignite/cli/ignite/pkg/events"
	"github.com/ignite/cli/ignite/pkg/xurl"
	"github.com/ignite/cli/ignite/services/chain"
)

const (
	flagServer        = "server"
	flagServerAddress = "server-address"
)

// NewChainDebug returns a new debug command to debug a blockchain app.
func NewChainDebug() *cobra.Command {
	c := &cobra.Command{
		Use:   "debug",
		Short: "Launch a debugger for a blockchain app",
		Long: `The debug command starts a debug server and launches a debugger.

Ignite uses the Delve debugger by default. Delve enables you to interact with
your program by controlling the execution of the process, evaluating variables,
and providing information of thread / goroutine state, CPU register state and
more.

A debug server can optionally be started in cases where default terminal client
is not desirable. When the server starts it first runs the blockchain app,
attaches to it and finally waits for a client connection. It accepts both
JSON-RPC or DAP client connections.

To start a debug server use the following flag:

	ignite chain debug --server

To start a debug server with a custom address use the following flags:

	ignite chain debug --server --server-address 127.0.0.1:30500

The debug server stops automatically when the client connection is closed.
`,
		Args: cobra.NoArgs,
		RunE: chainDebugHandler,
	}

	flagSetPath(c)
	c.Flags().Bool(flagServer, false, "start a debug server")
	c.Flags().String(flagServerAddress, debugger.DefaultAddress, "debug server address")

	return c
}

func chainDebugHandler(cmd *cobra.Command, _ []string) error {
	// Prepare session options.
	// Events are ignored by the session when the debug server UI is used.
	options := []cliui.Option{cliui.StartSpinnerWithText("Initializing...")}
	server, _ := cmd.Flags().GetBool(flagServer)
	if server {
		options = append(options, cliui.IgnoreEvents())
	}

	session := cliui.New(options...)
	defer session.End()

	// Start debug server
	if server {
		bus := session.EventBus()
		m := cmdmodel.NewChainDebug(cmd, bus, chainDebugCmd(cmd, session))
		_, err := tea.NewProgram(m).Run()
		return err
	}

	return chainDebug(cmd, session)
}

func chainDebugCmd(cmd *cobra.Command, session *cliui.Session) tea.Cmd {
	return func() tea.Msg {
		if err := chainDebug(cmd, session); err != nil && !errors.Is(err, context.Canceled) {
			return cliuimodel.ErrorMsg{Error: err}
		}
		return cliuimodel.QuitMsg{}
	}
}

func chainDebug(cmd *cobra.Command, session *cliui.Session) error {
	chainOptions := []chain.Option{
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
	}

	config, err := cmd.Flags().GetString(flagConfig)
	if err != nil {
		return err
	}
	if config != "" {
		chainOptions = append(chainOptions, chain.ConfigFile(config))
	}

	c, err := newChainWithHomeFlags(cmd, chainOptions...)
	if err != nil {
		return err
	}

	cfg, err := c.Config()
	if err != nil {
		return err
	}

	// TODO: Replace by config.FirstValidator when PR #3199 is merged
	validator := cfg.Validators[0]
	servers, err := validator.GetServers()
	if err != nil {
		return err
	}

	home, err := c.Home()
	if err != nil {
		return err
	}

	binPath, err := c.AbsBinaryPath()
	if err != nil {
		return err
	}

	// Common debugger options
	debugOptions := []debugger.Option{
		debugger.WorkingDir(flagGetPath(cmd)),
		debugger.BinaryArgs(
			"start",
			"--pruning", "nothing",
			"--grpc.address", servers.GRPC.Address,
			"--home", home,
		),
	}

	// Start debug server
	ctx := cmd.Context()
	bus := session.EventBus()
	if server, _ := cmd.Flags().GetBool(flagServer); server {
		addr, _ := cmd.Flags().GetString(flagServerAddress)
		tcpAddr, err := xurl.TCP(addr)
		if err != nil {
			return err
		}

		debugOptions = append(debugOptions,
			debugger.Address(addr),
			debugger.ServerStartHook(func() {
				bus.Send(
					fmt.Sprintf("Debug server: %s", tcpAddr),
					events.Icon(icons.Earth),
					events.ProgressFinish(),
				)
			}),
		)

		bus.Send("Launching debug server", events.ProgressUpdate())
		return debugger.Start(ctx, binPath, debugOptions...)
	}

	// Launch a debugger client
	debugOptions = append(debugOptions,
		debugger.ClientRunHook(func() {
			// End session to allow debugger to gain control of stdout
			session.End()
		}),
	)

	bus.Send("Launching debugger", events.ProgressUpdate())
	return debugger.Run(ctx, binPath, debugOptions...)
}
