package envtest

import (
	"context"
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/templates/field"
	"github.com/ignite/cli/v29/ignite/templates/field/datatype"
)

func txArgs(fields field.Fields) []string {
	args := make([]string, len(fields))
	for i, f := range fields {
		dt, _ := datatype.IsSupportedType(f.DatatypeName)
		args[i] = dt.DefaultTestValue
	}
	return args
}

func (a *App) assertJSONData(data []byte, msgName string, fields field.Fields) {
	for _, f := range fields {
		dt, _ := datatype.IsSupportedType(f.DatatypeName)
		value, _, _, err := jsonparser.Get(data, msgName, f.Name.Snake)
		require.NoError(a.env.T(), err)
		require.EqualValues(a.env.T(), string(value), dt.DefaultTestValue)
	}
}

func (a *App) RunChainAndSimulateTxs(servers Hosts) {
	ctx, cancel := context.WithCancel(a.env.ctx)
	defer cancel()

	go func() {
		a.MustServe(ctx)
	}()

	a.WaitChainUp(ctx, servers.API)

	a.RunSimulationTxs(ctx, servers)
}

func (a *App) RunSimulationTxs(ctx context.Context, servers Hosts) {
	for _, s := range a.scaffolded {
		module := s.module
		if module == "" {
			module = a.name
		}
		name, err := multiformatname.NewName(s.name)
		require.NoError(a.env.t, err)

		switch s.typeName {
		case "module":
		case "list":
		case "map":
		case "single":
			a.SendSingleTxAndQuery(ctx, servers, module, name, s.fields)
		case "type":
		case "params":
		case "configs":
		case "message":
		case "query":
		case "packet":
		}
	}
}

// SendSingleTxAndQuery send a chain transaction to a single store and query.
func (a *App) SendSingleTxAndQuery(
	ctx context.Context,
	servers Hosts,
	module string,
	name multiformatname.Name,
	fields field.Fields,
) {
	args := txArgs(fields)
	txResponse := a.CLITx(
		servers.RPC,
		module,
		"create-"+name.Kebab,
		args...,
	)
	require.Equal(a.env.T(), 0, txResponse.Code,
		"tx failed code=%d log=%s", txResponse.Code, txResponse.RawLog)

	tx := a.CLIQueryTx(
		servers.RPC,
		txResponse.TxHash,
	)
	require.Equal(a.env.T(), 0, tx.Code,
		"tx failed code=%d log=%s", txResponse.Code, txResponse.RawLog)

	queryReponse := a.CLIQuery(
		servers.RPC,
		module,
		"get-"+name.Kebab,
	)
	a.assertJSONData(queryReponse, name.Snake, fields)

	apiReponse := a.APIQuery(
		ctx,
		servers.API,
		a.namespace,
		module,
		name.Snake,
	)
	a.assertJSONData(apiReponse, name.Snake, fields)

	require.JSONEq(a.env.t, string(queryReponse), string(apiReponse))
}

// SendTxAndQueryFirst send a chain transaction and query the first element
func (a *App) SendTxAndQueryFirst(ctx context.Context, servers Hosts, namespace, module, msgName string, args ...string) {
	txResponse := a.CLITx(
		servers.RPC,
		module,
		"create-"+msgName,
		args...,
	)
	fmt.Println(txResponse)

	txResponse = a.CLIQueryTx(
		servers.RPC,
		txResponse.TxHash,
	)
	fmt.Println(txResponse)

	queryReponse := a.CLIQuery(
		servers.RPC,
		module,
		"list-"+msgName,
	)
	fmt.Println(queryReponse)

	queryReponse = a.CLIQuery(
		servers.RPC,
		module,
		"get-"+msgName,
		"0",
	)
	fmt.Println(queryReponse)

	apiReponse := a.APIQuery(
		ctx,
		servers.API,
		namespace,
		module,
		msgName,
	)
	fmt.Println(apiReponse)

	apiReponse = a.APIQuery(
		ctx,
		servers.API,
		namespace,
		module,
		msgName,
		"0",
	)
	fmt.Println(apiReponse)
}
