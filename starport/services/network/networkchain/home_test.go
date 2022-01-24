package networkchain_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/services/network/networkchain"
	"github.com/tendermint/starport/starport/services/network/networktypes"
)

func TestChainHome(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	chainHome := networkchain.ChainHome(0)
	require.Equal(t, filepath.Join(home, networktypes.SPN, "0"), chainHome)

	chainHome = networkchain.ChainHome(10)
	require.Equal(t, filepath.Join(home, networktypes.SPN, "10"), chainHome)
}
