//go:build !relayer

package app_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	envtest "github.com/ignite/cli/v28/integration"
)

// TestGenerateAppCheckBufPulsarPath tests scaffolding a new chain and checks if the buf.gen.pulsar.yaml file is correct
func TestGenerateAppCheckBufPulsarPath(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/blog")
	)

	bufGenPulsarPath := filepath.Join(app.SourcePath(), "proto", "buf.gen.pulsar.yaml")
	_, statErr := os.Stat(bufGenPulsarPath)
	require.False(t, os.IsNotExist(statErr), "buf.gen.pulsar.yaml should be scaffolded")

	result, err := os.ReadFile(bufGenPulsarPath)
	require.NoError(t, err)

	require.True(t, strings.Contains(string(result), "default: github.com/test/blog/api"), "buf.gen.pulsar.yaml should contain the correct api override")

	app.EnsureSteady()
}
