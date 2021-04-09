package names_test

import (
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/names"
	"testing"
)

func TestNoDash(t *testing.T) {
	require.Equal(t, "foo", names.NoDash("foo"))
	require.Equal(t, "foo", names.NoDash("-f-o-o---"))
}

func TestNoNumberPrefix(t *testing.T) {
	require.Equal(t, "foo", names.NoNumberPrefix("foo"))
	require.Equal(t, "_0foo", names.NoNumberPrefix("0foo"))
	require.Equal(t, "_999foo", names.NoNumberPrefix("999foo"))
}
