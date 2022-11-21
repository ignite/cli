package node_test

import (
	"bytes"
	"context"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/randstr"
	"github.com/ignite/cli/ignite/pkg/xurl"
	envtest "github.com/ignite/cli/integration"
)

func TestNodeQueryTx(t *testing.T) {
	var (
		appname = randstr.Runes(10)
		// alice   = "alice"
		// bob     = "bob"

		env     = envtest.New(t)
		app     = env.Scaffold(appname)
		home    = env.AppHome(appname)
		servers = app.RandomizeServerPorts()

		// accKeyringDir = t.TempDir()
	)

	node, err := xurl.HTTP(servers.RPC)
	require.NoError(t, err)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)

	go func() {
		defer cancel()

		if isBackendAliveErr = env.IsAppServed(ctx, servers.API); isBackendAliveErr != nil {
			return
		}
		client, err := cosmosclient.New(context.Background(),
			cosmosclient.WithAddressPrefix(testPrefix),
			cosmosclient.WithNodeAddress(node),
		)
		require.NoError(t, err)
		waitForNextBlock(env, client)

		b := &bytes.Buffer{}
		env.Exec("send 100token from alice to bob",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"tx",
					"bank",
					"send",
					"alice",
					"bob",
					"100token",
					"--node", node,
					"--keyring-dir", home,
				),
				step.Stdout(b),
			)),
		)
		require.False(t, env.HasFailed(), b.String())

		// Parse tx hash from output
		res := regexp.MustCompile(`\(hash = (\w+)\)`).FindAllStringSubmatch(b.String(), -1)
		require.Len(t, res[0], 2, "can't extract hash from command output")
		hash := res[0][1]
		waitForNextBlock(env, client)

		env.Must(env.Exec("query tx",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"tx",
					hash,
					"--node", node,
				),
			)),
		))
	}()

	env.Must(app.Serve("should serve", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}
