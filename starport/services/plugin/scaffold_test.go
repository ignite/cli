package plugin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScaffold(t *testing.T) {
	tmp, err := os.MkdirTemp("", "plugin")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)

	err = Scaffold(tmp, "github.com/foo/bar")

	require.NoError(t, err)
}
