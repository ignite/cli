package chainregistry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChainSaveJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "chain.json")
	in := Chain{
		ChainName: "ignite",
		ChainID:   "ignite-1",
	}

	require.NoError(t, in.SaveJSON(path))

	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	var got Chain
	require.NoError(t, json.Unmarshal(raw, &got))
	require.Equal(t, in.ChainName, got.ChainName)
	require.Equal(t, in.ChainID, got.ChainID)
}

func TestAssetListSaveJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "assetlist.json")
	in := AssetList{
		ChainName: "ignite",
		Assets: []Asset{
			{
				Name:   "Ignite",
				Base:   "uignite",
				Symbol: "IGNT",
			},
		},
	}

	require.NoError(t, in.SaveJSON(path))

	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	var got AssetList
	require.NoError(t, json.Unmarshal(raw, &got))
	require.Equal(t, in.ChainName, got.ChainName)
	require.Len(t, got.Assets, 1)
	require.Equal(t, "IGNT", got.Assets[0].Symbol)
}
