package v0_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	v0 "github.com/ignite/cli/ignite/chainconfig/v0"
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
	c2 := c.Clone()

	// Assert
	require.Equal(t, c, c2)
}
