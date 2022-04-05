package networkchain_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ignite-hq/cli/starport/services/network/networkchain"
)

func TestChainHome(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	chainHome := networkchain.ChainHome(0)
	require.Equal(t, filepath.Join(home, networkchain.SPN, "0"), chainHome)

	chainHome = networkchain.ChainHome(10)
	require.Equal(t, filepath.Join(home, networkchain.SPN, "10"), chainHome)
}
