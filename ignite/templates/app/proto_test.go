package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/config/chain/defaults"
)

func TestBufFiles(t *testing.T) {
	want := []string{"buf.work.yaml"}
	protoDir, err := os.ReadDir(filepath.Join("files", defaults.ProtoDir))
	require.NoError(t, err)
	for _, e := range protoDir {
		want = append(want, filepath.Join(defaults.ProtoDir, e.Name()))
	}

	got, err := BufFiles()
	require.NoError(t, err)
	require.ElementsMatch(t, want, got)
}
