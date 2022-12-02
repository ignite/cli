package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/chaincmd"
	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/debugger"
	"github.com/ignite/cli/ignite/pkg/events"
	"github.com/ignite/cli/ignite/pkg/xexec"
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
		Short: "Debug a blockchain app",
		Args:  cobra.NoArgs,
		RunE:  chainDebugHandler,
	}

	// TODO: Add --reset-once support
	flagSetPath(c)
	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagSetCheckDependencies())
	c.Flags().AddFlagSet(flagSetSkipProto())
	c.Flags().Bool(flagServer, false, "start a debug server")
	c.Flags().String(flagServerAddress, debugger.DefaultAddress, "debug server address")

	return c
}

func chainDebugHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(cliui.StartSpinnerWithText("Initializing..."))
	defer session.End()

	ev := session.EventBus()
	chainOptions := []chain.Option{
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
		chain.WithOutputer(session),
		chain.CollectEvents(ev),
	}

	if flagGetCheckDependencies(cmd) {
		chainOptions = append(chainOptions, chain.CheckDependencies())
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

	cache, err := newCache(cmd)
	if err != nil {
		return err
	}

	ctx := cmd.Context()
	binaryName, err := c.Build(ctx, cache, "", flagGetSkipProto(cmd), true)
	if err != nil {
		return err
	}

	binaryPath, err := xexec.ResolveAbsPath(binaryName)
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
	if serve, _ := cmd.Flags().GetBool(flagServer); serve {
		addr, _ := cmd.Flags().GetString(flagServerAddress)
		tcpAddr, err := xurl.TCP(addr)
		if err != nil {
			return err
		}

		debugOptions = append(debugOptions,
			debugger.Address(addr),
			debugger.ServerStartHook(func() {
				ev.Send(
					fmt.Sprintf("Debug server: %s", tcpAddr),
					events.Icon(icons.Earth),
					events.ProgressFinish(),
				)
			}),
		)

		// TODO: Use bubbletea for the debug server UI
		ev.Send("Launching debug server", events.ProgressUpdate())
		return debugger.Start(ctx, binaryPath, debugOptions...)
	}

	// Launch a debugger client
	debugOptions = append(debugOptions,
		debugger.ClientRunHook(func() {
			// End session to allow debugger to gain control of stdout
			session.End()
		}),
	)

	ev.Send("Launching debugger", events.ProgressUpdate())
	return debugger.Run(ctx, binaryPath, debugOptions...)
}
