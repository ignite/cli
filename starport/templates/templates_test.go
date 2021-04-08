package templates_test

import (
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/templates"
	"testing"
)

func TestNoNumberPrefix(t *testing.T) {
	require.Equal(t, "foo", templates.NoNumberPrefix("foo"))
	require.Equal(t, "_0foo", templates.NoNumberPrefix("0foo"))
	require.Equal(t, "_999foo", templates.NoNumberPrefix("999foo"))
}
