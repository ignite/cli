package node_test

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/chainconfig"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/ignite/pkg/randstr"
	"github.com/ignite/cli/ignite/pkg/xurl"
	envtest "github.com/ignite/cli/integration"
)

func TestNodeTxBankSend(t *testing.T) {
	var (
		appname = randstr.Runes(10)
		alice   = "alice"
		bob     = "bob"

		env     = envtest.New(t)
		app     = env.Scaffold(appname, "--address-prefix", testPrefix)
		home    = env.AppHome(appname)
		servers = app.RandomizeServerPorts()

		accKeyringDir = t.TempDir()
	)

	node, err := xurl.HTTP(servers.RPC)
	require.NoError(t, err)

	ca, err := cosmosaccount.New(
		cosmosaccount.WithHome(filepath.Join(home, keyringTestDirName)),
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
				Coins:    []string{"2000token", "100000000stake"},
			},
			{
				Name:     bob,
				Mnemonic: bobMnemonic,
				Coins:    []string{"10000token", "100000000stake"},
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
				"--keyring-dir", accKeyringDir,
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
				"--keyring-dir", accKeyringDir,
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

		app.WaitNBlocks(1, "--node", node)

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
					"--keyring-dir", accKeyringDir,
					"--address-prefix", testPrefix,
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
					"--node", node,
					"--keyring-dir", accKeyringDir,
					"--address-prefix", testPrefix,
				),
			)),
		)

		if env.HasFailed() {
			return
		}
		app.WaitNBlocks(1, "--node", node)

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
					"--node", node,
					"--keyring-dir", accKeyringDir,
					"--address-prefix", testPrefix,
				),
			)),
		)

		if env.HasFailed() {
			return
		}
		app.WaitNBlocks(1, "--node", node)

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
					"--node", node,
					"--keyring-dir", accKeyringDir,
					"--address-prefix", testPrefix,
				),
			)),
			envtest.ExecStdout(b),
		)

		assertBankBalanceOutput(t, b.String(), "2stake,1895token")

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
					"--node", node,
					"--keyring-dir", accKeyringDir,
					"--address-prefix", testPrefix,
				),
			)),
			envtest.ExecStdout(b),
		)

		assertBankBalanceOutput(t, b.String(), "99999998stake,10105token")

		// check generated tx
		b.Reset()
		env.Exec("generate unsigned tx",
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
					"--node", node,
					"--keyring-dir", accKeyringDir,
					"--address-prefix", testPrefix,
					"--generate-only",
				),
			)),
			envtest.ExecStdout(b),
		)

		require.Contains(t, b.String(),
			fmt.Sprintf(`"body":{"messages":[{"@type":"/cosmos.bank.v1beta1.MsgSend","from_address":"%s","to_address":"%s","amount":[{"denom":"token","amount":"5"}]}]`,
				aliceAccount.Address(testPrefix), bobAccount.Address(testPrefix)),
		)
		require.Contains(t, b.String(), `"signatures":[]`)

		// test with gas
		env.Exec("send 100token from bob to alice with gas flags",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"tx",
					"bank",
					"send",
					"bob",
					"alice",
					"100token",
					"--node", node,
					"--keyring-dir", accKeyringDir,
					"--address-prefix", testPrefix,
					"--gas", "200000",
					"--gas-prices", "1stake",
				),
			)),
		)

		// not enough minerals
		env.Exec("send 100token from alice to bob with too little gas",
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
					"--keyring-dir", accKeyringDir,
					"--address-prefix", testPrefix,
					"--gas", "2",
					"--gas-prices", "1stake",
				),
			)),
			envtest.ExecShouldError(),
		)

		b.Reset()
		env.Exec("generate bank send tx with gas flags",
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
					"--keyring-dir", accKeyringDir,
					"--address-prefix", testPrefix,
					"--gas", "2000034",
					"--gas-prices", "0.089stake",
					"--generate-only",
				),
			)),
			envtest.ExecStdout(b),
		)
		require.Contains(t, b.String(), `"fee":{"amount":[{"denom":"stake","amount":"178004"}],"gas_limit":"2000034"`)
	}()

	env.Must(app.Serve("should serve with Stargate version", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}
