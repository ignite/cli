package network_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/services/network"
)

func TestChainHome(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	chainHome, err := network.ChainHome(0, false)
	require.NoError(t, err)
	require.Equal(t, filepath.Join(home, network.ChainHomeRoot, network.ChainHomeInitDir, "0"), chainHome)

	chainHome, err = network.ChainHome(10, false)
	require.NoError(t, err)
	require.Equal(t, filepath.Join(home, network.ChainHomeRoot, network.ChainHomeInitDir, "10"), chainHome)

	chainHome, err = network.ChainHome(0, true)
	require.NoError(t, err)
	require.Equal(t, filepath.Join(home, network.ChainHomeRoot, network.ChainHomeLaunchDir, "0"), chainHome)

	chainHome, err = network.ChainHome(10, true)
	require.NoError(t, err)
	require.Equal(t, filepath.Join(home, network.ChainHomeRoot, network.ChainHomeLaunchDir, "10"), chainHome)
}
