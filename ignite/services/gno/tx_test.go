package gno

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/tm2/pkg/amino"
	abci "github.com/gnolang/gno/tm2/pkg/bft/abci/types"
	core_types "github.com/gnolang/gno/tm2/pkg/bft/rpc/core/types"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"gotest.tools/v3/assert"
)

// broadcastStubResult is the broadcast result returned by test stubs, with
// gas and hash set so BroadcastResult propagation can be asserted.
var broadcastStubResult = &core_types.ResultBroadcastTxCommit{
	Hash:      []byte{0xab, 0xcd},
	DeliverTx: abci.ResponseDeliverTx{GasUsed: 42},
}

func TestCallBuildsMsgCall(t *testing.T) {
	withTestGnoHome(t)

	var got txPlan
	restore := stubBroadcast(func(plan txPlan) (*core_types.ResultBroadcastTxCommit, error) {
		got = plan
		return broadcastStubResult, nil
	})
	defer restore()

	res, err := Call(CallOptions{
		PkgPath: "gno.land/r/counter",
		Func:    "Set",
		Args:    []string{"42"},
	})
	assert.NilError(t, err)
	assert.Equal(t, 1, len(got.tx.Msgs))
	assert.Assert(t, res != nil, "Call should return a broadcast result")
	assert.Equal(t, int64(42), res.GasUsed)
	assert.Equal(t, "ABCD", res.TxHash)
}

func TestCallValidation(t *testing.T) {
	withTestGnoHome(t)

	_, err := Call(CallOptions{Func: "NoPkg"})
	assert.ErrorContains(t, err, "package path is required")

	_, err = Call(CallOptions{PkgPath: "gno.land/r/x"})
	assert.ErrorContains(t, err, "function name is required")

	_, err = Call(CallOptions{PkgPath: "gno.land/r/x", Func: "F", Send: "bogus"})
	assert.ErrorContains(t, err, "parsing send coins")
}

func TestSendValidation(t *testing.T) {
	withTestGnoHome(t)

	_, err := Send(SendOptions{})
	assert.ErrorContains(t, err, "beneficiary is required")

	_, err = Send(SendOptions{To: "g1xxx"})
	assert.ErrorContains(t, err, "amount is required")

	_, err = Send(SendOptions{To: "g1xxx", Amount: "notcoins"})
	assert.ErrorContains(t, err, "parsing amount")
}

func TestSendRejectsInvalidAddress(t *testing.T) {
	withTestGnoHome(t)

	// an invalid address must fail with an error, not a panic
	_, err := Send(SendOptions{To: "g1typo", Amount: "1ugnot"})
	assert.ErrorContains(t, err, "neither a bech32 address nor a key")
}

func TestSendResolvesKeyNames(t *testing.T) {
	withTestGnoHome(t)

	info, _, err := CreateKey("alice", "", "", 0, 0)
	assert.NilError(t, err)

	var got txPlan
	restore := stubBroadcast(func(plan txPlan) (*core_types.ResultBroadcastTxCommit, error) {
		got = plan
		return broadcastStubResult, nil
	})
	defer restore()

	_, err = Send(SendOptions{
		TxBaseOptions: TxBaseOptions{From: "test1"},
		To:            "alice",
		Amount:        "10ugnot",
	})
	assert.NilError(t, err)
	assert.Equal(t, 1, len(got.tx.Msgs))
	_ = info
}

func TestSendAcceptsBech32Address(t *testing.T) {
	withTestGnoHome(t)

	restore := stubBroadcast(func(txPlan) (*core_types.ResultBroadcastTxCommit, error) {
		return broadcastStubResult, nil
	})
	defer restore()

	// a valid address must be accepted even when no key matches it in the keybase
	addr := crypto.AddressToBech32(crypto.AddressFromPreimage([]byte("ignite-test-recipient")))
	_, err := Send(SendOptions{
		TxBaseOptions: TxBaseOptions{From: "test1"},
		To:            addr,
		Amount:        "10ugnot",
	})
	assert.NilError(t, err)
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

func TestRenderQueryResultRef(t *testing.T) {
	// amino round-trip of a RefValue TV, mirroring the vm/qeval_json response
	// for expressions that evaluate to a stored object (e.g. a bare function
	// reference)
	tvs := []gnolang.TypedValue{{T: gnolang.BoolType, V: gnolang.RefValue{ObjectID: gnolang.ObjectID{NewTime: 6}}}}
	tvsJSON, err := amino.MarshalJSON(tvs)
	assert.NilError(t, err)
	data := []byte(fmt.Sprintf(`{"results":%s}`, tvsJSON))

	res, err := renderQueryResult(data)
	assert.NilError(t, err)
	assert.Assert(t, strings.Contains(res, "ref("), "ref value should be rendered: %s", res)
	assert.Assert(t, strings.Contains(res, queryRefHint), "ref value should come with a hint: %s", res)
}
