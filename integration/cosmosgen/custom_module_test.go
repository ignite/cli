package cosmosgen_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/config/chain/base"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/pkg/xurl"
	envtest "github.com/ignite/cli/v29/integration"
)

func TestCustomModule(t *testing.T) {
	t.Skip()

	var (
		env     = envtest.New(t)
		app     = env.ScaffoldApp("chain", "--no-module")
		servers = app.RandomizeServerPorts()
	)

	queryAPI, err := xurl.HTTP(servers.API)
	require.NoError(t, err)

	txAPI, err := xurl.HTTP(servers.RPC)
	require.NoError(t, err)

	// Accounts to be included in the genesis
	accounts := []base.Account{
		{
			Name:    "account1",
			Address: "cosmos1j8hw8283hj80hhq8urxaj40syrzqp77dt8qwhm",
			Mnemonic: fmt.Sprint(
				"toe mail light plug pact length excess predict real artwork laundry when ",
				"steel online adapt clutch debate vehicle dash alter rifle virtual season almost",
			),
			Coins: []string{"10000token", "10000stake"},
		},
	}

	app.EditConfig(func(cfg *chainconfig.Config) {
		cfg.Accounts = append(cfg.Accounts, accounts...)
	})

	path := app.SourcePath()

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "disco", "--require-registration", "--yes"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a list type",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "entry", "name", "--module", "disco", "--yes"),
			step.Workdir(path),
		)),
	))

	env.Must(app.GenerateTSClient())

	ctx, cancel := context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
	defer cancel()

	go func() {
		app.Serve("serve app", envtest.ExecCtx(ctx))
	}()

	// Wait for the server to be up before running the client tests
	app.WaitChainUp(ctx, servers.API)

	testAccounts, err := json.Marshal(accounts)
	require.NoError(t, err)

	env.Must(app.RunClientTests(
		envtest.ClientTestFile("custom_module_test.ts"),
		envtest.ClientEnv(map[string]string{
			"TEST_QUERY_API": queryAPI,
			"TEST_TX_API":    txAPI,
			"TEST_ACCOUNTS":  string(testAccounts),
		}),
	))
}
