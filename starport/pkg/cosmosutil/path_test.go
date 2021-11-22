package cosmosutil_test

import (
	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChainHome(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	chainHome, err := cosmosutil.ChainHome(0)
	require.NoError(t, err)
	require.Equal(t, filepath.Join(home, cosmosutil.ChainHomeRoot, "0"), chainHome)

	chainHome, err = cosmosutil.ChainHome(10)
	require.NoError(t, err)
	require.Equal(t, filepath.Join(home, cosmosutil.ChainHomeRoot, "10"), chainHome)
}
