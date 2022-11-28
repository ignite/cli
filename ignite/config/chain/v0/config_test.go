package v0_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	v0 "github.com/ignite/cli/ignite/config/chain/v0"
)

func TestClone(t *testing.T) {
	// Arrange
	c := &v0.Config{
		Validator: v0.Validator{
			Name:   "alice",
			Staked: "100000000stake",
		},
	}

	// Act
	c2, err := c.Clone()

	// Assert
	require.NoError(t, err)
	require.Equal(t, c, c2)
}
