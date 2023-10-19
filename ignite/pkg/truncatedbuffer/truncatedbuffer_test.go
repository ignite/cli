package truncatedbuffer

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriter(t *testing.T) {
	ranBytes10 := make([]byte, 10)
	_, err := rand.Read(ranBytes10)
	require.NoError(t, err)

	ranBytes1000 := make([]byte, 1000)
	_, err = rand.Read(ranBytes1000)
	require.NoError(t, err)

	// TruncatedBuffer has a max capacity
	b := NewTruncatedBuffer(100)

	require.Equal(t, 100, b.GetCap())

	n, err := b.Write(ranBytes10)
	require.NoError(t, err)
	require.Equal(t, 10, n)
	require.Equal(t, ranBytes10, b.GetBuffer().Bytes())

	n, err = b.Write(ranBytes1000)
	require.NoError(t, err)
	require.Equal(t, 1000, n)
	require.Equal(t, append(ranBytes10, ranBytes1000[:90]...), b.GetBuffer().Bytes())

	// TruncatedBuffer has no max capacity
	b = NewTruncatedBuffer(0)
	n, err = b.Write(ranBytes1000)
	require.NoError(t, err)
	require.Equal(t, 1000, n)
	require.Equal(t, ranBytes1000, b.GetBuffer().Bytes())
}
