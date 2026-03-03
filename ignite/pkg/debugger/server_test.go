package debugger

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplyDebuggerOptionsDefaults(t *testing.T) {
	o := applyDebuggerOptions()
	require.Equal(t, DefaultAddress, o.address)
	require.Equal(t, DefaultWorkingDir, o.workingDir)
	require.NotNil(t, o.disconnectChan)
}

func TestApplyDebuggerOptionsWithOverrides(t *testing.T) {
	c := make(chan struct{})
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, l.Close())
	})

	clientRan := false
	serverStarted := false

	o := applyDebuggerOptions(
		Address("127.0.0.1:9999"),
		DisconnectChannel(c),
		Listener(l),
		WorkingDir("/tmp/work"),
		BinaryArgs("a", "b"),
		ClientRunHook(func() { clientRan = true }),
		ServerStartHook(func() { serverStarted = true }),
	)

	require.Equal(t, "127.0.0.1:9999", o.address)
	require.Equal(t, c, o.disconnectChan)
	require.Equal(t, l, o.listener)
	require.Equal(t, "/tmp/work", o.workingDir)
	require.Equal(t, []string{"a", "b"}, o.binaryArgs)

	o.clientRunHook()
	o.serverStartHook()
	require.True(t, clientRan)
	require.True(t, serverStarted)
}

func TestDisableDelveLogging(t *testing.T) {
	require.NoError(t, disableDelveLogging())
}
