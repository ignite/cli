package modulecreate

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func readFixture(t *testing.T, relativePath string) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)

	path := filepath.Clean(filepath.Join(filepath.Dir(currentFile), relativePath))

	content, err := os.ReadFile(path)
	require.NoError(t, err)

	return string(content)
}
