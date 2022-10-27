package ctxreader

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadAndCancel(t *testing.T) {
	// create the ctx.
	ctx, cancel := context.WithCancel(context.Background())

	// create a buffer and write some initial data.
	buf := &bytes.Buffer{}
	buf.Write([]byte{1, 2, 3})

	// initialize cancelableReader with buf.
	r := New(ctx, buf)

	// make sure that cancelableReader will read the first 2 bytes
	// of previously written data and.
	data := make([]byte, 2)

	n, err := r.Read(data)
	require.NoError(t, err)
	require.Equal(t, 2, n)
	require.Equal(t, []byte{1, 2}, data)

	// cancel ctx and try to read again, this time Read will unblock and return
	// with a context.Canceled instead of reading the 3rd byte from initial data.
	cancel()

	n, err = r.Read(data)
	require.Equal(t, context.Canceled, err)
	require.Equal(t, 0, n)
}
