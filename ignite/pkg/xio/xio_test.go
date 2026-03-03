package xio

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNopWriteCloser(t *testing.T) {
	var b bytes.Buffer
	w := NopWriteCloser(&b)

	n, err := w.Write([]byte("ignite"))
	require.NoError(t, err)
	require.Equal(t, 6, n)
	require.Equal(t, "ignite", b.String())
	require.NoError(t, w.Close())
}
