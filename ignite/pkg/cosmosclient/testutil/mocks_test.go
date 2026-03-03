package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTendermintClientMock(t *testing.T) {
	m := NewTendermintClientMock(t)
	require.NotNil(t, m)
	require.NotNil(t, m.OnStatus())
	require.NotNil(t, m.OnBlock())
	require.NotNil(t, m.OnTxSearch())
}

func TestRepeatMockArgs(t *testing.T) {
	args := RepeatMockArgs(3)
	require.Len(t, args, 3)
	for _, arg := range args {
		require.NotNil(t, arg)
	}
}
