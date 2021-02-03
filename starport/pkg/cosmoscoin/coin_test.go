package cosmoscoin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	amount, denom, err := Parse("100token")
	require.NoError(t, err)
	require.Equal(t, uint64(100), amount)
	require.Equal(t, "token", denom)
}

func TestParseInvalid(t *testing.T) {
	_, _, err := Parse("!100token")
	require.Equal(t, errInvalidCoin, err)
}
