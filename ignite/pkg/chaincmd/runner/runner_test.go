package chaincmdrunner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewKV(t *testing.T) {
	kv := NewKV("k", "v")
	require.Equal(t, "k", kv.key)
	require.Equal(t, "v", kv.value)
}

func TestFindMatchingCloseBracket(t *testing.T) {
	data := []byte(`{"a":{"b":1}} trailing`)
	idx := findMatchingCloseBracket(data, '{', '}')
	require.Equal(t, 12, idx)

	require.Equal(t, -1, findMatchingCloseBracket([]byte(`{"a":1`), '{', '}'))
}

func TestCleanAndValidateJSON(t *testing.T) {
	raw := []byte("logs...\n{\"code\":0,\"raw_log\":\"ok\",\"txhash\":\"ABC\"}\nmore")
	got, err := cleanAndValidateJSON(raw)
	require.NoError(t, err)
	require.Equal(t, `{"code":0,"raw_log":"ok","txhash":"ABC"}`, string(got))
}

func TestFallbackFormatDetectionConvertsYAML(t *testing.T) {
	raw := []byte("code: 0\nraw_log: ok\ntxhash: ABC\n")
	got, err := fallbackFormatDetection(raw)
	require.NoError(t, err)
	require.JSONEq(t, `{"code":0,"raw_log":"ok","txhash":"ABC"}`, string(got))
}

func TestJSONEnsuredBytes(t *testing.T) {
	b := newBuffer()
	_, err := b.WriteString("noise\n{\"k\":\"v\"}\n")
	require.NoError(t, err)

	got, err := b.JSONEnsuredBytes()
	require.NoError(t, err)
	require.JSONEq(t, `{"k":"v"}`, string(got))
}

func TestDecodeTxResult(t *testing.T) {
	b := newBuffer()
	_, err := b.WriteString("code: 0\nraw_log: ok\ntxhash: HASH\n")
	require.NoError(t, err)

	got, err := decodeTxResult(b)
	require.NoError(t, err)
	require.Equal(t, 0, got.Code)
	require.Equal(t, "ok", got.RawLog)
	require.Equal(t, "HASH", got.TxHash)
}

func TestQueryTxByEventsRequiresSelectors(t *testing.T) {
	r := Runner{}
	events, err := r.QueryTxByEvents(context.Background())
	require.Error(t, err)
	require.Nil(t, events)
}

func TestQueryTxByQueryRequiresSelectors(t *testing.T) {
	r := Runner{}
	events, err := r.QueryTxByQuery(context.Background())
	require.Error(t, err)
	require.Nil(t, events)
}
