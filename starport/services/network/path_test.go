package network_test

import (
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/services/network"
	"os"
	"path/filepath"
	"testing"
)

func TestChainHome(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	chainHome, err := network.ChainHome(0)
	require.NoError(t, err)
	require.Equal(t, filepath.Join(home, "spn", "0"), chainHome)

	chainHome, err = network.ChainHome(10)
	require.NoError(t, err)
	require.Equal(t, filepath.Join(home, "spn", "10"), chainHome)
}
