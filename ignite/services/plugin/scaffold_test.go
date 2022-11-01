package plugin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScaffold(t *testing.T) {
	tmp := t.TempDir()

	path, err := Scaffold(tmp, "github.com/foo/bar")

	require.NoError(t, err)
	require.DirExists(t, path)
}
