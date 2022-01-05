package cosmoscoin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	c, err := Parse("100token")
	require.NoError(t, err)
	require.Equal(t, uint64(100), c.Amount)
	require.Equal(t, "token", c.Denom)
}

func TestParseInvalid(t *testing.T) {
	_, err := Parse("!100token")
	require.Equal(t, errInvalidCoin, err)
}

func TestCoin_String(t *testing.T) {
	require.Equal(t, "1000foo", Coin{Amount: 1000, Denom: "foo"}.String())
	require.Equal(t, "2000bar", Coin{Amount: 2000, Denom: "bar"}.String())
}
