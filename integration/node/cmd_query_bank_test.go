package node_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/stretchr/testify/require"

	"github.com/ignite-hq/cli/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite-hq/cli/integration"
)

const testPrefix = "testpref"
const aliceMnemonic = "trade physical mention claw forum fork night rate distance steak monster among soldier custom cave cloud addict runway melody current witness destroy version forward"
const aliceAddress = "testpref148akaazpnhce4gjcxy8l59969dtaxxrceju4m6"
const bobMnemonic = "alcohol alert unknown tissue clap basic slide air treat liquid proof toward outdoor loyal depart toddler cabbage glimpse warm outer switch output theme try"
const bobAddress = "testpref1nrzh528qngagy6vzgt2yc8p9quv8adjxn7rk65"

func TestNodeQueryBankBalances(t *testing.T) {
	var (
		env           = envtest.New(t)
		path          = env.Scaffold("github.com/test/blog", "--address-prefix", testPrefix)
		servers       = env.RandomizeServerPorts(path, "")
		rndWorkdir    = t.TempDir() // To make sure we can run these commands from anywhere
		accKeyringDir = t.TempDir()
	)

	env.SetKeyringBackend(path, "", keyring.BackendTest)
	env.SetConfigMnemonic(path, "", "alice", aliceMnemonic)
	env.SetConfigMnemonic(path, "", "bob", bobMnemonic)

	var (
		ctx, cancel = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
	)

	go func() {
		defer cancel()
		isBackendAliveErr := env.IsAppServed(ctx, servers)
		require.NoError(t, isBackendAliveErr, "app cannot get online in time")

		// error "account doesn't have any balances" occurs if a sleep is not included
		time.Sleep(time.Second * 1)

		env.Must(env.Exec("import alice",
			step.NewSteps(step.New(
				step.Exec(envtest.IgniteApp, "account", "import", "alice", "--keyring-dir", accKeyringDir, "--non-interactive", "--secret", aliceMnemonic),
			)),
		))
		env.Must(env.Exec("import bob",
			step.NewSteps(step.New(
				step.Exec(envtest.IgniteApp, "account", "import", "bob", "--keyring-dir", accKeyringDir, "--non-interactive", "--secret", bobMnemonic),
			)),
		))

		var accountOutputBuffer = &bytes.Buffer{}
		env.Must(env.Exec("query bank balances",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					"alice",
					"--node",
					"http://"+servers.RPC,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
				),
				step.Workdir(rndWorkdir),
			)),
			envtest.ExecStdout(accountOutputBuffer),
		))
		require.True(t, strings.Contains(accountOutputBuffer.String(), `Amount 		Denom 	
100000000 	stake 	
20000 		token`))

		var addressOutputBuffer = &bytes.Buffer{}
		env.Must(env.Exec("query bank balances",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					aliceAddress,
					"--node",
					"http://"+servers.RPC,
					"--keyring-dir",
					accKeyringDir,
				),
				step.Workdir(rndWorkdir),
			)),
			envtest.ExecStdout(addressOutputBuffer),
		))
		require.True(t, strings.Contains(addressOutputBuffer.String(), `Amount 		Denom 	
100000000 	stake 	
20000 		token`))

		var paginationFirstPageOutput = &bytes.Buffer{}
		env.Must(env.Exec("query bank balances with pagination",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					"alice",
					"--node",
					"http://"+servers.RPC,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
					"--limit",
					"1",
					"--page",
					"1",
				),
				step.Workdir(rndWorkdir),
			)),
			envtest.ExecStdout(paginationFirstPageOutput),
		))
		require.True(t, strings.Contains(paginationFirstPageOutput.String(), `Amount 		Denom 	
100000000 	stake`))
		require.False(t, strings.Contains(paginationFirstPageOutput.String(), `token`))

		var paginationSecondPageOutput = &bytes.Buffer{}
		env.Must(env.Exec("query bank balances with pagination",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					"alice",
					"--node",
					"http://"+servers.RPC,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
					"--limit",
					"1",
					"--page",
					"2",
				),
				step.Workdir(rndWorkdir),
			)),
			envtest.ExecStdout(paginationSecondPageOutput),
		))
		require.True(t, strings.Contains(paginationSecondPageOutput.String(), `Amount 	Denom 	
20000 	token`))
		require.False(t, strings.Contains(paginationSecondPageOutput.String(), `stake`))

		env.Must(env.Exec("query bank balances fail with non-existent account name",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					"nonexistentaccount",
					"--node",
					"http://"+servers.RPC,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
				),
				step.Workdir(rndWorkdir),
			)),
			envtest.ExecShouldError(),
		))

		env.Must(env.Exec("query bank balances fail with non-existent address",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					testPrefix+"1gspvt8qsk8cryrsxnqt452cjczjm5ejdgla24e",
					"--node",
					"http://"+servers.RPC,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
				),
				step.Workdir(rndWorkdir),
			)),
			envtest.ExecShouldError(),
		))

		env.Must(env.Exec("query bank balances fail with wrong prefix",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					"alice",
					"--node",
					"http://"+servers.RPC,
					"--keyring-dir",
					accKeyringDir,
				),
				step.Workdir(rndWorkdir),
			)),
			envtest.ExecShouldError(),
		))
	}()
	env.Must(env.Serve("should serve with Stargate version", path, "", "", envtest.ExecCtx(ctx)))
}
