package cosmosgen_test

import (
	"context"
	"testing"

	"github.com/ignite-hq/cli/ignite/chainconfig"
	"github.com/ignite-hq/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite-hq/cli/ignite/pkg/xurl"
	envtest "github.com/ignite-hq/cli/integration"
	"github.com/stretchr/testify/require"
)

func TestBankModule(t *testing.T) {
	var (
		env  = envtest.New(t)
		path = env.Scaffold("chain")
		host = env.RandomizeServerPorts(path, "")
	)

	queryAPI, err := xurl.HTTP(host.API)
	require.NoError(t, err)

	txAPI, err := xurl.TCP(host.RPC)
	require.NoError(t, err)

	// Accounts to be included in the genesis
	env.UpdateConfig(path, "", func(cfg *chainconfig.Config) error {
		accounts := []chainconfig.Account{
			{
				Name:    "account1",
				Address: "cosmos1j8hw8283hj80hhq8urxaj40syrzqp77dt8qwhm",
				Coins:   []string{"10000token", "10000stake"},
			},
			{
				Name:    "account2",
				Address: "cosmos19yy9sf00k00cjcwh532haeq8s63uhdy7qs5m2n",
				Coins:   []string{"100token", "100stake"},
			},
			{
				Name:    "account3",
				Address: "cosmos10957ee377t2xpwyt4mlpedjldp592h0ylt8uz7",
				Coins:   []string{"100token", "100stake"},
			},
		}

		cfg.Accounts = append(cfg.Accounts, accounts...)

		return nil
	})

	env.Must(env.Exec("generate vuex store", step.NewSteps(
		step.New(
			step.Exec(envtest.IgniteApp, "g", "vuex", "--proto-all-modules", "--yes"),
			step.Workdir(path),
		),
	)))

	ctx, cancel := context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
	defer cancel()

	go func() {
		env.Serve("should serve app", path, "", "", envtest.ExecCtx(ctx))
	}()

	// Wait for the server to be up before running the client tests
	err = env.IsAppServed(ctx, host)
	require.NoError(t, err)

	env.Must(env.RunClientTests(
		path,
		envtest.ClientTestName("Bank"),
		envtest.ClientEnv(map[string]string{
			"TEST_QUERY_API": queryAPI,
			"TEST_TX_API":    txAPI,
		}),
	))
}
