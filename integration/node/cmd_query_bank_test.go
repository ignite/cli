package node_test

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chainconfig "github.com/ignite/cli/ignite/config/chain"
	"github.com/ignite/cli/ignite/config/chain/base"
	"github.com/ignite/cli/ignite/pkg/cliui/entrywriter"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/randstr"
	"github.com/ignite/cli/ignite/pkg/xurl"
	xyaml "github.com/ignite/cli/ignite/pkg/yaml"
	envtest "github.com/ignite/cli/integration"
)

const (
	keyringTestDirName = "keyring-test"
	testPrefix         = "testpref"
)

func assertBankBalanceOutput(t *testing.T, output string, balances string) {
	var table [][]string
	coins, err := sdktypes.ParseCoinsNormalized(balances)
	require.NoError(t, err, "wrong balances %s", balances)
	for _, c := range coins {
		table = append(table, []string{c.Amount.String(), c.Denom})
	}
	var expectedBalances strings.Builder
	entrywriter.MustWrite(&expectedBalances, []string{"Amount", "Denom"}, table...)
	assert.Contains(t, output, expectedBalances.String())
}

func TestNodeQueryBankBalances(t *testing.T) {
	var (
		appname = randstr.Runes(10)
		alice   = "alice"

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

	app.EditConfig(func(c *chainconfig.Config) {
		c.Accounts = []base.Account{
			{
				Name:     alice,
				Mnemonic: aliceMnemonic,
				Coins:    []string{"5600atoken", "1200btoken", "100000000stake"},
			},
		}
		c.Faucet = base.Faucet{}
		c.Validators[0].Client = xyaml.Map{
			"keyring-backend": keyring.BackendTest,
		}
	})

	env.Must(env.Exec("import alice",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"account",
				"import",
				alice,
				"--keyring-dir", accKeyringDir,
				"--non-interactive",
				"--secret", aliceMnemonic,
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

		env.Exec("query bank balances by account name",
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

		if env.HasFailed() {
			return
		}

		assertBankBalanceOutput(t, b.String(), "5600atoken,1200btoken")

		b.Reset()

		aliceAddr, err := aliceAccount.Address(testPrefix)
		require.NoError(t, err)

		env.Exec("query bank balances by address",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					aliceAddr,
					"--node", node,
					"--keyring-dir", accKeyringDir,
					"--address-prefix", testPrefix,
				),
			)),
			envtest.ExecStdout(b),
		)

		if env.HasFailed() {
			return
		}

		assertBankBalanceOutput(t, b.String(), "5600atoken,1200btoken")

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
					"--node", node,
					"--keyring-dir", accKeyringDir,
					"--address-prefix", testPrefix,
					"--limit", "1",
					"--page", "1",
				),
			)),
			envtest.ExecStdout(b),
		)

		if env.HasFailed() {
			return
		}

		assertBankBalanceOutput(t, b.String(), "5600atoken")
		assert.NotContains(t, b.String(), "btoken")

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
					"--node", node,
					"--keyring-dir", accKeyringDir,
					"--address-prefix", testPrefix,
					"--limit", "1",
					"--page", "2",
				),
			)),
			envtest.ExecStdout(b),
		)

		if env.HasFailed() {
			return
		}

		assertBankBalanceOutput(t, b.String(), "1200btoken")
		assert.NotContains(t, b.String(), "atoken")

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
					"--node", node,
					"--keyring-dir", accKeyringDir,
					"--address-prefix", testPrefix,
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
					"--node", node,
					"--keyring-dir", accKeyringDir,
					"--address-prefix", testPrefix,
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
					"--node", node,
					"--keyring-dir", accKeyringDir,
					// the default prefix will fail this test, which is on purpose.
				),
			)),
			envtest.ExecShouldError(),
		)
	}()

	env.Must(app.Serve("should serve", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}
