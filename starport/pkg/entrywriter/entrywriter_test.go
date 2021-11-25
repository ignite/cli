package entrywriter_test

import (
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/entrywriter"
	"io"
	"testing"
)

func TestWrite(t *testing.T) {
	header := []string{"foobar", "bar", "foo"}

	entries := [][]string{
		{"foo", "bar", "foobar"},
		{"bar", "foobar", "foo"},
		{"foobar", "foo", "bar"},
	}

	require.NoError(t, entrywriter.Write(io.Discard, header, entries...))
	require.NoError(t, entrywriter.Write(io.Discard, header), "should allow no entry")
	require.Error(t, entrywriter.Write(io.Discard, []string{}), "should prevent no header")

	entries[0] = []string{"foo", "bar"}
	require.Error(t, entrywriter.Write(io.Discard, header, entries...), "should prevent entry length mismatch")
}
