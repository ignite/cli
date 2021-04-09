package templateutils_test

import (
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/templateutils"
	"testing"
)

func TestNoDash(t *testing.T) {
	require.Equal(t, "foo", templateutils.NoDash("foo"))
	require.Equal(t, "foo", templateutils.NoDash("-f-o-o---"))
}

func TestNoNumberPrefix(t *testing.T) {
	require.Equal(t, "foo", templateutils.NoNumberPrefix("foo"))
	require.Equal(t, "_0foo", templateutils.NoNumberPrefix("0foo"))
	require.Equal(t, "_999foo", templateutils.NoNumberPrefix("999foo"))
}
