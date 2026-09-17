package envtest

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/require"

	sdktypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/templates/field"
	"github.com/ignite/cli/v29/ignite/templates/field/datatype"
)

const (
	// chainServeTimeout bounds each chain serve attempt. An attempt that does
	// not bring the chain API up in time is considered wedged and is stopped.
	chainServeTimeout = 8 * time.Minute

	// chainServeAttempts is how many times chain serving is attempted before
	// the test is failed with the serve logs.
	chainServeAttempts = 2
)

// testValue determines the default test value for a given datatype.
func testValue(name datatype.Name) string {
	// Repeated custom message fields in AutoCLI expect one JSON object per argument.
	// Using "{}" keeps simulation compatible with both singular and repeated custom types.
	if name == datatype.CustomSlice {
		return "{}"
	}

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

	var (
		serveCtx    context.Context
		cancelServe context.CancelFunc
		serveDone   chan struct{}
		serveLogs   *bytes.Buffer
	)

	// startServe runs the chain serve command in the background. Its output is
	// captured so that it can be reported when serving fails. Serve output is
	// otherwise invisible because the exec step only prints it on failure.
	startServe := func() {
		sctx, scancel := context.WithCancel(ctx)
		logs := &bytes.Buffer{}
		done := make(chan struct{})
		serveCtx, cancelServe, serveLogs, serveDone = sctx, scancel, logs, done
		go func() {
			defer close(done)
			a.Serve("should serve chain",
				ExecCtx(sctx),
				ExecStdout(logs),
				ExecStderr(logs),
			)
		}()
	}

	// waitChainUp waits until the chain API responds, serving exits or the
	// wait context is done.
	waitChainUp := func(waitCtx context.Context) (servingExited bool, err error) {
		apiErr := make(chan error, 1)
		go func() {
			apiErr <- a.env.IsAppServed(waitCtx, servers.API)
		}()
		select {
		case err = <-apiErr:
		case <-serveDone:
			servingExited = true
		case <-waitCtx.Done():
			err = waitCtx.Err()
		}
		return servingExited, err
	}

	startServe()
	for attempt := 1; ; attempt++ {
		waitCtx, cancelWait := context.WithTimeout(serveCtx, chainServeTimeout)
		servingExited, err := waitChainUp(waitCtx)
		cancelWait()

		if err == nil {
			break // the chain API is up.
		}

		if servingExited {
			cancelServe()
			cancel()
			a.env.t.Fatalf("chain serve exited before the chain API %s was up\n\nServe logs:\n\n%s",
				servers.API, serveLogs.String())
		}

		if attempt >= chainServeAttempts || ctx.Err() != nil {
			// Stop serving and wait for the serve goroutine before reading the
			// logs it wrote.
			cancelServe()
			<-serveDone
			cancel()
			a.env.t.Fatalf("chain API %s did not come up: %v\n\nServe logs:\n\n%s",
				servers.API, err, serveLogs.String())
		}

		// The serve attempt looks wedged: stop it and try again.
		cancelServe()
		<-serveDone
		startServe()
	}

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
