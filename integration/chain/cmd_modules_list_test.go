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

	require.Equal(t, buffer.String(), "\r\r\x1b[K\r\r\x1b[KModule \t\t\t\t\tVersion \nclient/grpc/cmtservice \t\t\tv0.53.2 \nclient/grpc/node \t\t\tv0.53.2 \ncosmossdk.io/x/circuit \t\t\tv0.1.1 \t\ncosmossdk.io/x/evidence \t\tv0.1.1 \t\ncosmossdk.io/x/feegrant/module \t\tv0.1.1 \t\ncosmossdk.io/x/nft/module \t\tv0.1.0 \t\ncosmossdk.io/x/upgrade \t\t\tv0.2.0 \t\ngithub.com/test/mars/x/mars \t\tmain \t\ngithub.com/test/mars/x/mars/module \tmain \t\nmodules/apps/27-interchain-accounts \tv10.2.0 \nmodules/apps/transfer \t\t\tv10.2.0 \nmodules/core \t\t\t\tv10.2.0 \nx/auth \t\t\t\t\tv0.53.2 \nx/auth/tx \t\t\t\tv0.53.2 \nx/auth/tx/config \t\t\tv0.53.2 \nx/auth/vesting \t\t\t\tv0.53.2 \nx/authz \t\t\t\tv0.53.2 \nx/authz/module \t\t\t\tv0.53.2 \nx/bank \t\t\t\t\tv0.53.2 \nx/consensus \t\t\t\tv0.53.2 \nx/distribution \t\t\t\tv0.53.2 \nx/epochs \t\t\t\tv0.53.2 \nx/gov \t\t\t\t\tv0.53.2 \nx/group/module \t\t\t\tv0.53.2 \nx/mint \t\t\t\t\tv0.53.2 \nx/params \t\t\t\tv0.53.2 \nx/slashing \t\t\t\tv0.53.2 \nx/staking \t\t\t\tv0.53.2 \n\n")

	app.EnsureSteady()
}
