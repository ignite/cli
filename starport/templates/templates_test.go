package templates_test

import (
	"github.com/tendermint/starport/starport/templates"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNoNumberPrefix(t *testing.T) {
	require.Equal(t, "foo", templates.NoNumberPrefix("foo"))
	require.Equal(t, "_0foo", templates.NoNumberPrefix("0foo"))
	require.Equal(t, "_999foo", templates.NoNumberPrefix("999foo"))
}