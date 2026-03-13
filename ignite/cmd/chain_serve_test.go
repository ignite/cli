package ignitecmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShouldRunServeInDaemonMode(t *testing.T) {
	cmd := NewChainServe()

	require.NoError(t, cmd.Flags().Set(flagOutputFile, "/tmp/serve.log"))
	require.True(t, shouldRunServeInDaemonMode(cmd))
}

func TestShouldRunServeInDaemonModeNoTerminal(t *testing.T) {
	cmd := NewChainServe()
	cmd.SetIn(bytes.NewBufferString(""))
	cmd.SetOut(bytes.NewBuffer(nil))

	require.True(t, shouldRunServeInDaemonMode(cmd))
}

func TestShouldRunServeInDaemonModeInteractiveTerminal(t *testing.T) {
	cmd := NewChainServe()
	cmd.SetIn(os.Stdin)
	cmd.SetOut(os.Stdout)

	originalIsTerminal := isTerminal
	t.Cleanup(func() {
		isTerminal = originalIsTerminal
	})
	isTerminal = func(_ int) bool { return true }

	require.False(t, shouldRunServeInDaemonMode(cmd))
}

func TestShouldRunServeInDaemonModeMixedTerminal(t *testing.T) {
	cmd := NewChainServe()
	cmd.SetIn(os.Stdin)
	cmd.SetOut(os.Stdout)

	originalIsTerminal := isTerminal
	t.Cleanup(func() {
		isTerminal = originalIsTerminal
	})
	calls := 0
	isTerminal = func(_ int) bool {
		calls++
		return calls == 1
	}

	require.True(t, shouldRunServeInDaemonMode(cmd))
}
