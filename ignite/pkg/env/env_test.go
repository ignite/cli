package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetDebugAndIsDebug(t *testing.T) {
	t.Setenv(DebugEnvVar, "")
	require.False(t, IsDebug())

	SetDebug()
	require.True(t, IsDebug())
}

func TestConfigDirFromEnv(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "ignite-config")
	t.Setenv(ConfigDirEnvVar, dir)

	got, err := ConfigDir()()
	require.NoError(t, err)
	require.Equal(t, dir, got)
}

func TestConfigDirPanicsWithRelativePath(t *testing.T) {
	t.Setenv(ConfigDirEnvVar, "relative/path")
	require.Panics(t, func() {
		_, _ = ConfigDir()()
	})
}

func TestSetConfigDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "cfg")
	SetConfigDir(dir)
	require.Equal(t, dir, os.Getenv(ConfigDirEnvVar))
}
