package gno

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/tm2/pkg/amino"
	"github.com/gnolang/gno/tm2/pkg/bft/rpc/client"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// evalQueryPath is the ABCI query path evaluating a read-only gno expression.
const evalQueryPath = "vm/qeval_json"

// queryRefHint is appended to results containing references to on-chain
// objects (functions or non-primitive state), which have no printable value
// client-side.
const queryRefHint = "This is a reference to on-chain data (a function or a non-primitive state variable), not a printable value.\nTo read state, call a realm function, e.g.: ignite chain query \"gno.land/r/helloworld.Get()\""

// queryEnvelope is the response envelope of vm/qeval_json.
type queryEnvelope struct {
	Results json.RawMessage `json:"results"`
	Error   *string         `json:"@error"`
}

// Query evaluates a read-only gno expression against a gno.land chain and
// returns the rendered result. An expression that names a function without
// parentheses is evaluated as a call, e.g. "gno.land/r/helloworld.Get" is
// queried as "gno.land/r/helloworld.Get()".
func Query(remote, expr string) (string, error) {
	if remote == "" {
		remote = defaultRemote
	}
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return "", errors.Errorf("expression is required")
	}
	if !strings.HasSuffix(expr, ")") {
		// a bare function reference evaluates to an unhelpful ref(...) value:
		// try it as a call first, then fall back to the raw expression.
		if res, err := evalQuery(remote, expr+"()"); err == nil {
			return res, nil
		}
	}
	return evalQuery(remote, expr)
}

// evalQuery evaluates expr against the chain at remote.
func evalQuery(remote, expr string) (string, error) {
	cli, err := client.NewHTTPClient(remote)
	if err != nil {
		return "", errors.Errorf("connecting to %s: %w", remote, err)
	}

	qres, err := cli.ABCIQueryWithOptions(context.Background(), evalQueryPath, []byte(expr), client.ABCIQueryOptions{})
	if err != nil {
		return "", errors.Errorf("querying %s: %w", remote, err)
	}
	if qres.Response.Error != nil {
		return "", errors.Errorf("query failed: %w, log: %s", qres.Response.Error, qres.Response.Log)
	}

	return renderQueryResult(qres.Response.Data)
}

// renderQueryResult turns a vm/qeval_json response into a human friendly
// string: decoded TypedValues when possible, raw JSON otherwise.
func renderQueryResult(data []byte) (string, error) {
	var env queryEnvelope
	if err := json.Unmarshal(data, &env); err != nil {
		return string(data), nil // not the expected envelope, show raw
	}
	if env.Error != nil && *env.Error != "" {
		return "", errors.Errorf("query returned an error: %s", *env.Error)
	}

	var tvs []gnolang.TypedValue
	if err := amino.UnmarshalJSON(env.Results, &tvs); err != nil {
		return string(data), nil // fallback to raw
	}

	parts := make([]string, 0, len(tvs))
	refs := 0
	for _, tv := range tvs {
		switch tv.V.(type) {
		case gnolang.RefValue, *gnolang.RefValue:
			refs++
		}
		parts = append(parts, tv.Sprint(nil))
	}
	res := strings.Join(parts, " ")
	if refs > 0 {
		res += "\n" + queryRefHint
	}
	return res, nil
}
