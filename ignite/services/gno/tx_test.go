package gno

import (
	"testing"

	core_types "github.com/gnolang/gno/tm2/pkg/bft/rpc/core/types"
	"gotest.tools/v3/assert"
)

func TestCallBuildsMsgCall(t *testing.T) {
	withTestGnoHome(t)

	var got txPlan
	restore := stubBroadcast(func(plan txPlan) (*core_types.ResultBroadcastTxCommit, error) {
		got = plan
		return &core_types.ResultBroadcastTxCommit{}, nil
	})
	defer restore()

	assert.NilError(t, Call(CallOptions{
		PkgPath: "gno.land/r/counter",
		Func:    "Set",
		Args:    []string{"42"},
	}))
	assert.Equal(t, 1, len(got.tx.Msgs))
}

func TestCallValidation(t *testing.T) {
	withTestGnoHome(t)

	err := Call(CallOptions{Func: "NoPkg"})
	assert.ErrorContains(t, err, "package path is required")

	err = Call(CallOptions{PkgPath: "gno.land/r/x"})
	assert.ErrorContains(t, err, "function name is required")

	err = Call(CallOptions{PkgPath: "gno.land/r/x", Func: "F", Send: "bogus"})
	assert.ErrorContains(t, err, "parsing send coins")
}

func TestSendValidation(t *testing.T) {
	withTestGnoHome(t)

	err := Send(SendOptions{})
	assert.ErrorContains(t, err, "beneficiary is required")

	err = Send(SendOptions{To: "g1xxx"})
	assert.ErrorContains(t, err, "amount is required")

	err = Send(SendOptions{To: "g1xxx", Amount: "notcoins"})
	assert.ErrorContains(t, err, "parsing amount")
}

func TestSendResolvesKeyNames(t *testing.T) {
	withTestGnoHome(t)

	info, _, err := CreateKey("alice", "", "", 0, 0)
	assert.NilError(t, err)

	var got txPlan
	restore := stubBroadcast(func(plan txPlan) (*core_types.ResultBroadcastTxCommit, error) {
		got = plan
		return &core_types.ResultBroadcastTxCommit{}, nil
	})
	defer restore()

	assert.NilError(t, Send(SendOptions{
		TxBaseOptions: TxBaseOptions{From: "test1"},
		To:            "alice",
		Amount:        "10ugnot",
	}))
	assert.Equal(t, 1, len(got.tx.Msgs))
	_ = info
}

func TestQueryValidation(t *testing.T) {
	_, err := Query("127.0.0.1:26657", "")
	assert.ErrorContains(t, err, "expression is required")
}

func TestRenderQueryResult(t *testing.T) {
	tt := []struct {
		name string
		data string
		want string
	}{
		{
			name: "primitive int",
			// real amino-json shape returned by vm/qeval_json for []TypedValue{int 2}
			data: `{"results":[{"T":{"@type":"/gno.PrimitiveType","value":"32"},"N":"AgAAAAAAAAA="}]}`,
			want: "2",
		},
		{
			name: "no results",
			data: `{"results":[]}`,
			want: "",
		},
		{
			name: "raw fallback",
			data: "not-json",
			want: "not-json",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := renderQueryResult([]byte(tc.data))
			assert.NilError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestRenderQueryResultError(t *testing.T) {
	_, err := renderQueryResult([]byte(`{"results":[],"@error":"undefined: Nope"}`))
	assert.ErrorContains(t, err, "undefined: Nope")
}
