package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/chaincmd"
	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/debugger"
	"github.com/ignite/cli/ignite/pkg/events"
	"github.com/ignite/cli/ignite/pkg/xexec"
	"github.com/ignite/cli/ignite/pkg/xurl"
	"github.com/ignite/cli/ignite/services/chain"
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
	// TODO: Add --server & --server-address flags
	flagSetPath(c)
	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagSetCheckDependencies())
	c.Flags().AddFlagSet(flagSetSkipProto())

	return c
}

func chainDebugHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(cliui.StartSpinner())
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

	rpcAddr, err := xurl.TCP(servers.RPC.Address)
	if err != nil {
		return err
	}

	home, err := c.Home()
	if err != nil {
		return err
	}

	debugOptions := []debugger.Option{
		debugger.WorkingDir(flagGetPath(cmd)),
		debugger.BinaryArgs(
			"start",
			"--pruning", "nothing",
			"--grpc.address", rpcAddr,
			"--home", home,
		),
		debugger.ClientRunHook(func() {
			// End session to allow debugger to gain control of stdout
			session.End()
		}),
	}

	ev.Send("Launching debugger", events.ProgressUpdate())
	return debugger.Run(ctx, binaryPath, debugOptions...)
}
