package gno

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/tm2/pkg/amino"
	"github.com/gnolang/gno/tm2/pkg/bft/rpc/client"
)

// evalQueryPath is the ABCI query path evaluating a read-only gno expression.
const evalQueryPath = "vm/qeval_json"

// queryEnvelope is the response envelope of vm/qeval_json.
type queryEnvelope struct {
	Results json.RawMessage `json:"results"`
	Error   *string         `json:"@error"`
}

// Query evaluates a read-only gno expression against a gno.land chain and
// returns the rendered result, e.g. Query(remote, "gno.land/r/counter.Get()").
func Query(remote, expr string) (string, error) {
	if remote == "" {
		remote = "127.0.0.1:26657"
	}
	if expr == "" {
		return "", fmt.Errorf("expression is required")
	}

	cli, err := client.NewHTTPClient(remote)
	if err != nil {
		return "", fmt.Errorf("connecting to %s: %w", remote, err)
	}

	qres, err := cli.ABCIQueryWithOptions(context.Background(), evalQueryPath, []byte(expr), client.ABCIQueryOptions{})
	if err != nil {
		return "", fmt.Errorf("querying %s: %w", remote, err)
	}
	if qres.Response.Error != nil {
		return "", fmt.Errorf("query failed: %w, log: %s", qres.Response.Error, qres.Response.Log)
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
		return "", fmt.Errorf("query returned an error: %s", *env.Error)
	}

	var tvs []gnolang.TypedValue
	if err := amino.UnmarshalJSON(env.Results, &tvs); err != nil {
		return string(data), nil // fallback to raw
	}

	parts := make([]string, 0, len(tvs))
	for _, tv := range tvs {
		parts = append(parts, tv.Sprint(nil))
	}
	return strings.Join(parts, " "), nil
}
