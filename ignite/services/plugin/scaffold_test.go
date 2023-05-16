package plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScaffold(t *testing.T) {
	tmp := t.TempDir()

	path, err := Scaffold(context.Background(), tmp, "github.com/foo/bar", false)

	require.NoError(t, err)
	require.DirExists(t, path)
}
