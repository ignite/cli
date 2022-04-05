package cosmosutil_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ignite-hq/cli/starport/pkg/cosmosutil"
	"github.com/stretchr/testify/require"
)

const (
	genesisSample = `
{
	"foo": "bar",
	"genesis_time": "foobar"
}
`
	unixTime = 1600000000
	rfcTime  = "2020-09-13T12:26:40Z"
)

func TestSetGenesisTime(t *testing.T) {
	tmp, err := os.MkdirTemp("", "")
	t.Cleanup(func() { os.RemoveAll(tmp) })
	tmpGenesis := filepath.Join(tmp, "genesis.json")

	// fails with no file
	require.NoError(t, err)
	require.Error(t, cosmosutil.SetGenesisTime(tmpGenesis, 0))

	require.NoError(t, os.WriteFile(tmpGenesis, []byte(genesisSample), 0644))
	require.NoError(t, cosmosutil.SetGenesisTime(tmpGenesis, unixTime))

	// check genesis modified value
	var actual struct {
		Foo         string `json:"foo"`
		GenesisTime string `json:"genesis_time"`
	}
	actualBytes, err := os.ReadFile(tmpGenesis)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(actualBytes, &actual))
	require.Equal(t, "bar", actual.Foo)
	require.Equal(t, rfcTime, actual.GenesisTime)
}
