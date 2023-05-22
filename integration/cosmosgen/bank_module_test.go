package cosmosgen_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	chainconfig "github.com/ignite/cli/ignite/config/chain"
	"github.com/ignite/cli/ignite/config/chain/base"
	"github.com/ignite/cli/ignite/pkg/xurl"
	envtest "github.com/ignite/cli/integration"
)

func TestBankModule(t *testing.T) {
	t.Skip()

	var (
		env     = envtest.New(t)
		app     = env.Scaffold("chain", "--no-module")
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
				"toe mail light plug pact length excess predict real artwork laundry when steel ",
				"online adapt clutch debate vehicle dash alter rifle virtual season almost",
			),
			Coins: []string{"10000token", "10000stake"},
		},
		{
			Name:    "account2",
			Address: "cosmos19yy9sf00k00cjcwh532haeq8s63uhdy7qs5m2n",
			Mnemonic: fmt.Sprint(
				"someone major rule wrestle forget want job record coil table enter gold bracket ",
				"zone tent music grow shiver width index radio matter asset when",
			),
			Coins: []string{"100token", "100stake"},
		},
		{
			Name:    "account3",
			Address: "cosmos10957ee377t2xpwyt4mlpedjldp592h0ylt8uz7",
			Mnemonic: fmt.Sprint(
				"edit effort own cat chuckle rookie mechanic side tool sausage other fade math ",
				"joy midnight cabin act plastic spawn loud chest invest budget rebel",
			),
			Coins: []string{"100token", "100stake"},
		},
	}

	app.EditConfig(func(cfg *chainconfig.Config) {
		cfg.Accounts = append(cfg.Accounts, accounts...)
	})

	env.Must(app.GenerateTSClient())

	ctx, cancel := context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
	defer cancel()

	go func() {
		app.Serve("should serve app", envtest.ExecCtx(ctx))
	}()

	// Wait for the server to be up before running the client tests
	err = env.IsAppServed(ctx, servers.API)
	require.NoError(t, err)

	testAccounts, err := json.Marshal(accounts)
	require.NoError(t, err)

	env.Must(app.RunClientTests(
		envtest.ClientTestFile("bank_module_test.ts"),
		envtest.ClientEnv(map[string]string{
			"TEST_QUERY_API": queryAPI,
			"TEST_TX_API":    txAPI,
			"TEST_ACCOUNTS":  string(testAccounts),
		}),
	))
}
