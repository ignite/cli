package v0_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite-hq/cli/ignite/chainconfig/config"
	v0 "github.com/ignite-hq/cli/ignite/chainconfig/v0"
)

func TestClone(t *testing.T) {
	config := &v0.Config{
		Validator: v0.Validator{
			Name:   "alice",
			Staked: "100000000stake",
		},
		Init: config.Init{
			App:    nil,
			Client: nil,
			Config: nil,
		},
	}
	clone := config.Clone()
	require.Equal(t, config, clone)

	clone.(*v0.Config).Validator = v0.Validator{
		Name:   "test",
		Staked: "stakedvalue",
	}
	require.NotEqual(t, config, clone)
	require.Equal(t, v0.Validator{
		Name:   "test",
		Staked: "stakedvalue",
	}, clone.(*v0.Config).Validator)
}
