package xstrings_test

import (
	"github.com/ignite-hq/cli/starport/pkg/xstrings"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNoDash(t *testing.T) {
	require.Equal(t, "foo", xstrings.NoDash("foo"))
	require.Equal(t, "foo", xstrings.NoDash("-f-o-o---"))
}

func TestNoNumberPrefix(t *testing.T) {
	require.Equal(t, "foo", xstrings.NoNumberPrefix("foo"))
	require.Equal(t, "_0foo", xstrings.NoNumberPrefix("0foo"))
	require.Equal(t, "_999foo", xstrings.NoNumberPrefix("999foo"))
}
