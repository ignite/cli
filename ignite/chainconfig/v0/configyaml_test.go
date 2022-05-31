package v0

import (
	"testing"

	"github.com/ignite/cli/ignite/chainconfig/common"
	"github.com/stretchr/testify/require"
)

func TestClone(t *testing.T) {
	config := &Config{
		Validator: Validator{
			Name:   "alice",
			Staked: "100000000stake",
		},
		Init: common.Init{
			App:    nil,
			Client: nil,
			Config: nil,
		},
	}
	clone := config.Clone()
	require.Equal(t, config, clone)

	clone.(*Config).Validator = Validator{
		Name:   "test",
		Staked: "stakedvalue",
	}
	require.NotEqual(t, config, clone)
	require.Equal(t, Validator{
		Name:   "test",
		Staked: "stakedvalue",
	}, clone.(*Config).Validator)
}
