//go:build !relayer

package chain_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/v29/integration"
)

func TestModulesList(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("github.com/test/mars")
	)

	var buffer bytes.Buffer

	env.Must(env.Exec("list modules",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"c",
				"modules",
				"list",
			),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecStdout(&buffer),
	))

	output := buffer.String()

	// check for module header
	require.Contains(t, output, "Module")
	require.Contains(t, output, "Version")

	// check for specific modules
	require.Contains(t, output, "client/grpc/cmtservice")
	require.Contains(t, output, "client/grpc/node")
	require.Contains(t, output, "cosmossdk.io/x/circuit")
	require.Contains(t, output, "cosmossdk.io/x/evidence")
	require.Contains(t, output, "cosmossdk.io/x/feegrant/module")
	require.Contains(t, output, "cosmossdk.io/x/nft/module")
	require.Contains(t, output, "cosmossdk.io/x/upgrade")
	require.Contains(t, output, "github.com/test/mars/x/mars")
	require.Contains(t, output, "github.com/test/mars/x/mars/module")
	require.Contains(t, output, "modules/apps/27-interchain-accounts")
	require.Contains(t, output, "modules/apps/transfer")
	require.Contains(t, output, "modules/core")
	require.Contains(t, output, "x/auth")
	require.Contains(t, output, "x/auth/tx")
	require.Contains(t, output, "x/auth/tx/config")
	require.Contains(t, output, "x/auth/vesting")
	require.Contains(t, output, "x/authz")
	require.Contains(t, output, "x/authz/module")
	require.Contains(t, output, "x/bank")
	require.Contains(t, output, "x/consensus")
	require.Contains(t, output, "x/distribution")
	require.Contains(t, output, "x/epochs")
	require.Contains(t, output, "x/gov")
	require.Contains(t, output, "x/group/module")
	require.Contains(t, output, "x/mint")
	require.Contains(t, output, "x/params")
	require.Contains(t, output, "x/slashing")
	require.Contains(t, output, "x/staking")

	app.EnsureSteady()
}
