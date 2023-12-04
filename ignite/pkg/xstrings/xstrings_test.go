package xstrings_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v28/ignite/pkg/xstrings"
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

func TestStringBetween(t *testing.T) {
	require.Equal(t, "bar", xstrings.StringBetween("foobarbaz", "foo", "baz"))
	require.Equal(t, "bar", xstrings.StringBetween("0foobarbaz1", "foo", "baz"))
	require.Equal(t, "", xstrings.StringBetween("0foo", "0", ""))
	require.Equal(t, "", xstrings.StringBetween("foo0", "", "0"))
	require.Equal(t, "", xstrings.StringBetween("", "0", "1"))
}
