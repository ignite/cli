package cosmosutil_test

import (
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSetGenesisTime(t *testing.T) {
	require.Error(t, cosmosutil.SetGenesisTime(filepath.Join(os.TempDir(),"no", "genesis.json"), 0))


}