package plugin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScaffold(t *testing.T) {
	tmp := t.TempDir()

	err := Scaffold(tmp, "github.com/foo/bar")

	require.NoError(t, err)
}
