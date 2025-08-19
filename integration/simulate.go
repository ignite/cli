package envtest

import (
	"context"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/require"

	sdktypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/templates/field"
	"github.com/ignite/cli/v29/ignite/templates/field/datatype"
)

// testValue determines the default test value for a given datatype.
func testValue(name datatype.Name) string {
	dt, _ := datatype.IsSupportedType(name)
	return dt.DefaultTestValue
}

// txArgs generates transaction arguments as strings from a given set of fields.
func txArgs(fields field.Fields) []string {
	args := make([]string, len(fields))
	for i, f := range fields {
		args[i] = testValue(f.DatatypeName)
	}
	return args
}

// assertJSONData verifies that the JSON data contains expected values for the given fields.
func (a *App) assertJSONData(data []byte, msgName string, fields field.Fields) {
	for _, f := range fields {
		dt := testValue(f.DatatypeName)
		value, _, _, err := jsonparser.Get(data, msgName, f.Name.Snake)
		require.NoError(a.env.T(), err)
		if dt == "{}" {
			continue
		}
		v := string(value)
		switch {
		case f.DatatypeName == datatype.Coin:

			c, err := sdktypes.ParseCoinNormalized(dt)
			require.NoError(a.env.T(), err)
			amount, err := jsonparser.GetString(value, "amount")
			require.NoError(a.env.T(), err)
			require.EqualValues(a.env.T(), amount, c.Amount.String())
			denom, err := jsonparser.GetString(value, "denom")
			require.NoError(a.env.T(), err)
			require.EqualValues(a.env.T(), denom, c.Denom)

		case f.DatatypeName == datatype.Coins || f.DatatypeName == datatype.CoinSliceAlias:

			c, err := sdktypes.ParseCoinsNormalized(dt)
			require.NoError(a.env.T(), err)
			cJSON, err := c.MarshalJSON()
			require.NoError(a.env.T(), err)
			dt = string(cJSON)
			require.JSONEq(a.env.T(), dt, v)

		case f.DatatypeName == datatype.DecCoin || f.DatatypeName == datatype.DecCoins || f.DatatypeName == datatype.DecCoinSliceAlias:

			c, err := sdktypes.ParseCoinNormalized(dt)
			require.NoError(a.env.T(), err)
			// TODO find a better way to compare DecCoins as they have a different result pattern from CLI and Query
			require.Contains(a.env.T(), v, c.Denom)
			require.Contains(a.env.T(), v, c.Amount.String())

		case f.IsSlice():

			var slice []string
			_, err = jsonparser.ArrayEach(value, func(value []byte, _ jsonparser.ValueType, _ int, _ error) {
				slice = append(slice, string(value))
			})
			require.NoError(a.env.T(), err)
			v = strings.Join(slice, ",")
			require.EqualValues(a.env.T(), dt, v)

		default:
			require.EqualValues(a.env.T(), dt, v)
		}
	}
}

// assertJSONList verifies that a JSON array contains expected values for the given fields.
func (a *App) assertJSONList(data []byte, msgName string, fields field.Fields) {
	value, _, _, err := jsonparser.Get(data, msgName)
	require.NoError(a.env.t, err)

	a.assertJSONData(value, "[0]", fields)
}

// createTx sends a transaction to create a resource and verifies the response from the chain.
func (a *App) createTx(
	servers Hosts,
	module string,
	name multiformatname.Name,
	args ...string,
) {
	// Submit the transaction and verify it was accepted
	txResponse := a.CLITx(
		servers.RPC,
		module,
		"create-"+name.Kebab,
		args...,
	)
	require.Equal(a.env.T(), 0, txResponse.Code,
		"tx failed code=%d log=%s", txResponse.Code, txResponse.RawLog)

	// Query the transaction using its hash
	tx := a.CLIQueryTx(
		servers.RPC,
		txResponse.TxHash,
	)
	require.Equal(a.env.T(), 0, tx.Code,
		"tx failed code=%d log=%s", txResponse.Code, txResponse.RawLog)
}

// RunChainAndSimulateTxs starts the blockchain network and runs transaction simulations.
func (a *App) RunChainAndSimulateTxs(servers Hosts) {
	ctx, cancel := context.WithCancel(a.env.ctx)
	defer cancel()

	// Start serving the blockchain in a separate goroutine
	go func() {
		a.MustServe(ctx)
	}()

	// Wait until the chain is up and running
	a.WaitChainUp(ctx, servers.API)

	// Run the transaction simulations
	a.RunSimulationTxs(ctx, servers)
}

// RunSimulationTxs runs different types of transactions for modules and queries the chain.
func (a *App) RunSimulationTxs(ctx context.Context, servers Hosts) {
	for _, s := range a.scaffolded {
		module := s.module
		if module == "" {
			module = a.name
		}
		name, err := multiformatname.NewName(s.name)
		require.NoError(a.env.t, err)

		// Handle different types of scaffolds
		switch s.typeName {
		case "module":
			// No transactions for "module" type
		case "list":
			a.SendListTxsAndQueryFirst(ctx, servers, module, name, s.fields)
		case "map":
			a.SendMapTxsAndQuery(ctx, servers, module, name, s.fields, s.index)
		case "single":
			a.SendSingleTxsAndQuery(ctx, servers, module, name, s.fields)
		case "params":
		case "message":
		case "query":
		case "configs":
		case "type":
		case "packet":
		}
	}
}

// SendSingleTxsAndQuery submits a single transaction and queries the result from both CLI and API.
func (a *App) SendSingleTxsAndQuery(
	ctx context.Context,
	servers Hosts,
	module string,
	name multiformatname.Name,
	fields field.Fields,
) {
	// Generate transaction arguments and submit the transaction
	args := txArgs(fields)
	a.createTx(servers, module, name, args...)

	// Query the state via CLI
	queryResponse := a.CLIQuery(
		servers.RPC,
		module,
		"get-"+name.Kebab,
	)
	a.assertJSONData(queryResponse, name.Snake, fields)

	// Query the state via API
	apiResponse := a.APIQuery(
		ctx,
		servers.API,
		a.namespace,
		module,
		name.Snake,
	)
	a.assertJSONData(apiResponse, name.Snake, fields)

	// Ensure CLI and API responses match
	require.JSONEq(a.env.t, string(queryResponse), string(apiResponse))
}

// SendListTxsAndQueryFirst sends a list transaction and queries the first element using both CLI and API.
func (a *App) SendListTxsAndQueryFirst(
	ctx context.Context,
	servers Hosts,
	module string,
	name multiformatname.Name,
	fields field.Fields,
) {
	a.SendTxsAndQuery(ctx, servers, module, name, fields, "0")
}

// SendMapTxsAndQuery sends a map transaction and queries the element using both CLI and API.
func (a *App) SendMapTxsAndQuery(
	ctx context.Context,
	servers Hosts,
	module string,
	name multiformatname.Name,
	fields field.Fields,
	index field.Field,
) {
	a.SendTxsAndQuery(
		ctx,
		servers,
		module,
		name,
		append(field.Fields{index}, fields...),
		testValue(index.DatatypeName),
	)
}

// SendTxsAndQuery sends a transaction and queries the element using both CLI and API.
func (a *App) SendTxsAndQuery(
	ctx context.Context,
	servers Hosts,
	module string,
	name multiformatname.Name,
	fields field.Fields,
	index string,
) {
	// Generate transaction arguments and submit the transaction
	args := txArgs(fields)
	a.createTx(servers, module, name, args...)

	// Query the chain for the first element via CLI
	queryResponse := a.CLIQuery(
		servers.RPC,
		module,
		"get-"+name.Kebab,
		index,
	)
	a.assertJSONData(queryResponse, name.Snake, fields)

	// Query the chain for the first element via API
	apiResponse := a.APIQuery(
		ctx,
		servers.API,
		a.namespace,
		module,
		name.Snake,
		index,
	)
	a.assertJSONData(apiResponse, name.Snake, fields)

	// Query the full list via CLI
	queryListResponse := a.CLIQuery(
		servers.RPC,
		module,
		"list-"+name.Kebab,
	)
	a.assertJSONList(queryListResponse, name.Snake, fields)

	// Query the full list via API
	apiListResponse := a.APIQuery(
		ctx,
		servers.API,
		a.namespace,
		module,
		name.Snake,
	)
	a.assertJSONList(apiListResponse, name.Snake, fields)
}
