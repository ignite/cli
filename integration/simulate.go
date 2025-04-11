package envtest

import (
	"context"

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

func (a *App) assertJSONList(data []byte, msgName string, fields field.Fields) {
	value, _, _, err := jsonparser.Get(data, msgName, "[0]")
	require.NoError(a.env.t, err)

	for _, f := range fields {
		dt, _ := datatype.IsSupportedType(f.DatatypeName)
		value, _, _, err := jsonparser.Get(value, f.Name.Snake)
		require.NoError(a.env.T(), err)
		require.EqualValues(a.env.T(), string(value), dt.DefaultTestValue)
	}
}

func (a *App) createTx(
	servers Hosts,
	module string,
	name multiformatname.Name,
	args ...string,
) TxResponse {
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

	return tx
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
			a.SendListTxsAndQueryFirst(ctx, servers, module, name, s.fields)
		case "map":
		case "single":
			a.SendSingleTxsAndQuery(ctx, servers, module, name, s.fields)
		case "type":
		case "params":
		case "configs":
		case "message":
		case "query":
		case "packet":
		}
	}
}

// SendSingleTxsAndQuery send a chain transaction to a single store and query.
func (a *App) SendSingleTxsAndQuery(
	ctx context.Context,
	servers Hosts,
	module string,
	name multiformatname.Name,
	fields field.Fields,
) {
	args := txArgs(fields)
	_ = a.createTx(servers, module, name, args...)

	queryResponse := a.CLIQuery(
		servers.RPC,
		module,
		"get-"+name.Kebab,
	)
	a.assertJSONData(queryResponse, name.Snake, fields)

	apiResponse := a.APIQuery(
		ctx,
		servers.API,
		a.namespace,
		module,
		name.Snake,
	)
	a.assertJSONData(apiResponse, name.Snake, fields)

	require.JSONEq(a.env.t, string(queryResponse), string(apiResponse))
}

// SendListTxsAndQueryFirst send a chain transaction and query the first element
func (a *App) SendListTxsAndQueryFirst(
	ctx context.Context,
	servers Hosts,
	module string,
	name multiformatname.Name,
	fields field.Fields,
) {
	args := txArgs(fields)
	_ = a.createTx(servers, module, name, args...)

	queryResponse := a.CLIQuery(
		servers.RPC,
		module,
		"get-"+name.Kebab,
		"0",
	)
	a.assertJSONData(queryResponse, name.Snake, fields)

	apiResponse := a.APIQuery(
		ctx,
		servers.API,
		a.namespace,
		module,
		name.Snake,
		"0",
	)
	a.assertJSONData(apiResponse, name.Snake, fields)

	queryListResponse := a.CLIQuery(
		servers.RPC,
		module,
		"list-"+name.Kebab,
	)
	a.assertJSONList(queryListResponse, name.Snake, fields)

	apiListResponse := a.APIQuery(
		ctx,
		servers.API,
		a.namespace,
		module,
		name.Snake,
	)
	a.assertJSONList(apiListResponse, name.Snake, fields)
}
