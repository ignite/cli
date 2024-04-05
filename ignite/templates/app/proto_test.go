package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBufFiles(t *testing.T) {
	want := []string{"buf.work.yaml"}
	protoDir, err := os.ReadDir("files/{{protoDir}}")
	require.NoError(t, err)
	for _, e := range protoDir {
		want = append(want, filepath.Join("{{protoDir}}", e.Name()))
	}

	got, err := BufFiles()
	require.NoError(t, err)
	require.ElementsMatch(t, want, got)
}
