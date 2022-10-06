package xnet_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/xnet"
)

func TestLocalhostIPv4Address(t *testing.T) {
	require.Equal(t, "localhost:42", xnet.LocalhostIPv4Address(42))
}

func TestAnyIPv4Address(t *testing.T) {
	require.Equal(t, "0.0.0.0:42", xnet.AnyIPv4Address(42))
}

func TestIncreasePort(t *testing.T) {
	addr, err := xnet.IncreasePort("localhost:41")

	require.NoError(t, err)
	require.Equal(t, "localhost:42", addr)
}

func TestIncreasePortWithInvalidAddress(t *testing.T) {
	_, err := xnet.IncreasePort("localhost:x:41")

	require.Error(t, err)
}

func TestIncreasePortWithInvalidPort(t *testing.T) {
	_, err := xnet.IncreasePort("localhost:x")

	require.Error(t, err)
}

func TestIncreasePortBy(t *testing.T) {
	addr, err := xnet.IncreasePortBy("localhost:32", 10)

	require.NoError(t, err)
	require.Equal(t, "localhost:42", addr)
}

func TestIncreasePortByWithInvalidAddress(t *testing.T) {
	_, err := xnet.IncreasePortBy("localhost:x:32", 10)

	require.Error(t, err)
}

func TestIncreasePortByWithInvalidPort(t *testing.T) {
	_, err := xnet.IncreasePortBy("localhost:x", 10)

	require.Error(t, err)
}
