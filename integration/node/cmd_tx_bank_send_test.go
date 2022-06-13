package node_test

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
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

func TestNodeTxBankSend(t *testing.T) {
	var (
		name  = "blog"
		alice = "alice"
		bob   = "bob"

		env     = envtest.New(t)
		app     = env.Scaffold(name, "--address-prefix", testPrefix)
		home    = env.AppHome(name)
		servers = app.RandomizeServerPorts()

		accKeyringDir = t.TempDir()
	)

	node, err := xurl.HTTP(servers.RPC)
	require.NoError(t, err)

	defer env.RequireExpectations()

	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringDirPath(filepath.Join(home, keyringTestDirName)),
		cosmosaccount.WithKeyringBackend(cosmosaccount.KeyringTest),
	)
	require.NoError(t, err)

	aliceAccount, aliceMnemonic, err := ca.Create(alice)
	require.NoError(t, err)

	bobAccount, bobMnemonic, err := ca.Create(bob)
	require.NoError(t, err)

	app.EditConfig(func(conf *chainconfig.Config) {
		conf.Accounts = []chainconfig.Account{
			{
				Name:     alice,
				Mnemonic: aliceMnemonic,
				Coins:    []string{"500a", "600b"},
			},
			{
				Name:     bob,
				Mnemonic: bobMnemonic,
				Coins:    []string{"2500a", "2600b"},
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
				"alice",
				"--keyring-dir",
				accKeyringDir,
				"--non-interactive",
				"--secret", aliceMnemonic,
			),
		)),
	))

	env.Must(env.Exec("import bob",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"account",
				"import",
				"bob",
				"--keyring-dir",
				accKeyringDir,
				"--non-interactive",
				"--secret", bobMnemonic,
			),
		)),
	))

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)

	go func() {
		defer cancel()

		if isBackendAliveErr = env.IsAppServed(ctx, servers); isBackendAliveErr != nil {
			return
		}

		// error "account doesn't have any balances" occurs if a sleep is not included
		// TODO find another way without sleep, with retry+ctx routine.
		time.Sleep(time.Second * 1)

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
					"--from",
					"alice",
					"--node",
					node,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
				),
			)),
		)

		if env.HasFailed() {
			return
		}

		env.Exec("send 2stake from bob to alice using addresses",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"tx",
					"bank",
					"send",
					bobAccount.Address(testPrefix),
					aliceAccount.Address(testPrefix),
					"2stake",
					"--from",
					"bob",
					"--node",
					node,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
				),
			)),
		)

		if env.HasFailed() {
			return
		}

		env.Exec("send 5token from alice to bob using a combination of address and account",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"tx",
					"bank",
					"send",
					"alice",
					bobAccount.Address(testPrefix),
					"5token",
					"--from",
					"alice",
					"--node",
					node,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
				),
			)),
		)

		if env.HasFailed() {
			return
		}

		// TODO find another way without sleep, with retry+ctx routine.
		time.Sleep(time.Second * 1)

		b := &bytes.Buffer{}

		env.Exec("query bank balances for alice",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					"alice",
					"--rpc",
					"http://"+servers.RPC,
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
100000002 	stake 	
19895 		token`,
		))

		if env.HasFailed() {
			return
		}

		b.Reset()

		env.Exec("query bank balances for bob",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					"bob",
					"--rpc",
					"http://"+servers.RPC,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
				),
			)),
			envtest.ExecStdout(b),
		)

		assert.True(t, strings.Contains(b.String(), `
Amount 		Denom 	
99999998 	stake 	
10105 		token`,
		))

	}()

	env.Must(app.Serve("should serve with Stargate version", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestNodeTxBankSendGenerateOnly(t *testing.T) {
	var (
		name  = "blog"
		alice = "alice"
		bob   = "bob"

		env     = envtest.New(t)
		app     = env.Scaffold(name, "--address-prefix", testPrefix)
		home    = env.AppHome(name)
		servers = app.RandomizeServerPorts()

		accKeyringDir = t.TempDir()
	)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)

	go func() {
		defer cancel()

		if isBackendAliveErr = env.IsAppServed(ctx, servers); isBackendAliveErr != nil {
			return
		}

		// error "account doesn't exist" occurs if a sleep is not included
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

		var generateOutput = &bytes.Buffer{}
		env.Must(env.Exec("generate unsigned tx",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"tx",
					"bank",
					"send",
					"alice",
					"bob",
					"5token",
					"--rpc",
					"http://"+servers.RPC,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
					"--generate-only",
				),
				step.Workdir(rndWorkdir),
			)),
			envtest.ExecStdout(generateOutput),
		))

		require.True(t, strings.Contains(generateOutput.String(), fmt.Sprintf(`"body":{"messages":[{"@type":"/cosmos.bank.v1beta1.MsgSend","from_address":"%s","to_address":"%s","amount":[{"denom":"token","amount":"5"}]}]`, aliceAddress, bobAddress)))
		require.True(t, strings.Contains(generateOutput.String(), `"signatures":[]`))
	}()

	env.Must(env.Serve("should serve with Stargate version", path, "", "", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestNodeTxBankSendWithGas(t *testing.T) {
	var (
		name  = "blog"
		alice = "alice"
		bob   = "bob"

		env     = envtest.New(t)
		app     = env.Scaffold(name, "--address-prefix", testPrefix)
		home    = env.AppHome(name)
		servers = app.RandomizeServerPorts()

		accKeyringDir = t.TempDir()
	)

	node, err := xurl.HTTP(servers.RPC)
	require.NoError(t, err)

	defer env.RequireExpectations()

	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringDirPath(filepath.Join(home, keyringTestDirName)),
		cosmosaccount.WithKeyringBackend(cosmosaccount.KeyringTest),
	)
	require.NoError(t, err)

	aliceAccount, aliceMnemonic, err := ca.Create(alice)
	require.NoError(t, err)

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

		env.Must(env.Exec("send 100token from alice to bob with gas flags",
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
					"--from",
					"alice",
					"--rpc",
					"http://"+servers.RPC,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
					"--gas",
					"200000",
					"--gas-prices",
					"1stake",
				),
				step.Workdir(rndWorkdir),
			)),
		))

		env.Must(env.Exec("send 100token from alice to bob with too little gas",
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
					"--from",
					"alice",
					"--rpc",
					"http://"+servers.RPC,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
					"--gas",
					"2",
					"--gas-prices",
					"1stake",
				),
				step.Workdir(rndWorkdir),
			)),
			envtest.ExecShouldError(),
		))

		var generateOutput = &bytes.Buffer{}
		env.Must(env.Exec("generate bank send tx with gas flags",
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
					"--from",
					"alice",
					"--rpc",
					"http://"+servers.RPC,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
					"--gas",
					"2000034",
					"--gas-prices",
					"0.089stake",
					"--generate-only",
				),
				step.Workdir(rndWorkdir),
			)),
			envtest.ExecStdout(generateOutput),
		))
		require.True(t, strings.Contains(generateOutput.String(), `"fee":{"amount":[{"denom":"stake","amount":"178004"}],"gas_limit":"2000034"`))

	}()

	env.Must(env.Serve("should serve with Stargate version", path, "", "", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}
