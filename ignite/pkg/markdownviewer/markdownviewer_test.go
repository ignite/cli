package markdownviewer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigReturnsErrorWhenStdoutIsNotTTY(t *testing.T) {
	tempStdout := setNonTTYStdout(t)

	_, err := config(t.TempDir())
	require.Error(t, err)

	require.NoError(t, tempStdout.Close())
}

func TestViewReturnsConfigErrorWhenStdoutIsNotTTY(t *testing.T) {
	tempStdout := setNonTTYStdout(t)

	err := View(t.TempDir())
	require.Error(t, err)

	require.NoError(t, tempStdout.Close())
}

func setNonTTYStdout(t *testing.T) *os.File {
	t.Helper()

	file, err := os.CreateTemp(t.TempDir(), "stdout-*")
	require.NoError(t, err)

	originalStdout := os.Stdout
	os.Stdout = file
	t.Cleanup(func() {
		os.Stdout = originalStdout
	})

	return file
}
