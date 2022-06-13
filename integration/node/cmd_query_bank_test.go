package node_test

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/stretchr/testify/require"

	"github.com/ignite-hq/cli/ignite/chainconfig"
	"github.com/ignite-hq/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite-hq/cli/ignite/pkg/xurl"
	envtest "github.com/ignite-hq/cli/integration"
)

const keyringTestDirName = "keyring-test"
const testPrefix = "testpref"

func TestNodeQueryBankBalances(t *testing.T) {
	var (
		name  = "blog"
		alice = "alice"

		env     = envtest.New(t)
		app     = env.Scaffold(name, "--address-prefix", testPrefix)
		home    = env.AppHome(name)
		servers = app.RandomizeServerPorts()

		accKeyringDir = t.TempDir()
	)

	node, err := xurl.HTTP(servers.RPC)
	require.NoError(t, err)

	defer env.RequireExpectations()

	// TODO use INMEM
	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringDirPath(filepath.Join(home, keyringTestDirName)),
		cosmosaccount.WithKeyringBackend(cosmosaccount.KeyringTest),
	)
	require.NoError(t, err)

	aliceAccount, aliceMnemonic, err := ca.Create(alice)
	require.NoError(t, err)

	app.EditConfig(func(conf *chainconfig.Config) {
		conf.Accounts = []chainconfig.Account{
			{
				Name:     alice,
				Mnemonic: aliceMnemonic,
				Coins:    []string{"5600a", "1200b"},
			},
		}
		conf.Init.KeyringBackend = keyring.BackendTest
	})

	env.Must(env.Exec("import alice",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"account",
				"import",
				alice,
				"--keyring-dir",
				accKeyringDir,
				"--non-interactive",
				"--secret",
				aliceMnemonic,
			),
		)),
	))

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)

	// do not fail the test in a goroutine, it has to be done in the main.
	go func() {
		defer cancel()

		if isBackendAliveErr = env.IsAppServed(ctx, servers); isBackendAliveErr != nil {
			return
		}

		// error "account doesn't have any balances" occurs if a sleep is not included
		// TODO find another way without sleep, with retry+ctx routine.
		time.Sleep(time.Second * 1)

		b := &bytes.Buffer{}

		env.Exec("query bank balances by account name",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					"alice",
					"--node",
					node,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
				),
			)),
			envtest.ExecStdout(b),
		)

		assert.True(t, envtest.Contains(b.String(), `
Amount 		Denom 	
5600		a	
1200		b`,
		))

		if env.HasFailed() {
			return
		}

		b.Reset()

		env.Exec("query bank balances by address",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					aliceAccount.Address(testPrefix),
					"--node",
					node,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
				),
			)),
			envtest.ExecStdout(b),
		)

		assert.True(t, envtest.Contains(b.String(), `,
Amount 		Denom 	
5600		a	
1200		b`,
		))

		if env.HasFailed() {
			return
		}

		b.Reset()

		env.Exec("query bank balances with pagination -page 1",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					"alice",
					"--node",
					node,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
					"--limit",
					"1",
					"--page",
					"1",
				),
			)),
			envtest.ExecStdout(b),
		)

		assert.True(t, envtest.Contains(b.String(), `
Amount 		Denom 	
5600		a`,
		))
		assert.False(t, envtest.Contains(b.String(), `token`))

		if env.HasFailed() {
			return
		}

		b.Reset()

		env.Exec("query bank balances with pagination -page 2",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					"alice",
					"--node",
					node,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
					"--limit",
					"1",
					"--page",
					"2",
				),
			)),
			envtest.ExecStdout(b),
		)

		assert.True(t, envtest.Contains(b.String(), `
Amount 	Denom 	
1200	b`,
		))
		assert.False(t, envtest.Contains(b.String(), `stake`))

		if env.HasFailed() {
			return
		}

		b.Reset()

		env.Exec("query bank balances fail with non-existent account name",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					"nonexistentaccount",
					"--node",
					node,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
				),
			)),
			envtest.ExecShouldError(),
		)

		if env.HasFailed() {
			return
		}

		env.Exec("query bank balances fail with non-existent address",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					testPrefix+"1gspvt8qsk8cryrsxnqt452cjczjm5ejdgla24e",
					"--node",
					node,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
				),
			)),
			envtest.ExecShouldError(),
		)

		if env.HasFailed() {
			return
		}

		env.Exec("query bank balances should fail with a wrong prefix",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					"alice",
					"--node",
					node,
					"--keyring-dir",
					accKeyringDir,
					// the default prefix will fail this test, which is on purpose.
				),
			)),
			envtest.ExecShouldError(),
		)
	}()

	env.Must(app.Serve("should serve with Stargate version", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}
