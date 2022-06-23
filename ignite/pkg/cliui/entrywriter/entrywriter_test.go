package entrywriter_test

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cliui/entrywriter"
)

type WriterWithError struct{}

func (WriterWithError) Write(_ []byte) (n int, err error) {
	return 0, errors.New("writer with error")
}

func TestWrite(t *testing.T) {
	header := []string{"foobar", "bar", "foo"}

	entries := [][]string{
		{"foo", "bar", "foobar"},
		{"bar", "foobar", "foo"},
		{"foobar", "foo", "bar"},
	}

	require.NoError(t, entrywriter.Write(io.Discard, header, entries...))
	require.NoError(t, entrywriter.Write(io.Discard, header), "should allow no entry")

	err := entrywriter.Write(io.Discard, []string{})
	require.ErrorIs(t, err, entrywriter.ErrInvalidFormat, "should prevent no header")

	entries[0] = []string{"foo", "bar"}
	err = entrywriter.Write(io.Discard, header, entries...)
	require.ErrorIs(t, err, entrywriter.ErrInvalidFormat, "should prevent entry length mismatch")

	var wErr WriterWithError
	require.Error(t, entrywriter.Write(wErr, header, entries...), "should catch writer errors")
}
