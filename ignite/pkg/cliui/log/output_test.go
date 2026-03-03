package uilog

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/xio"
)

func TestNewOutputDefault(t *testing.T) {
	var outBuf, errBuf bytes.Buffer
	out := NewOutput(
		WithStdout(xio.NopWriteCloser(&outBuf)),
		WithStderr(xio.NopWriteCloser(&errBuf)),
	)

	_, err := out.Stdout().Write([]byte("stdout"))
	require.NoError(t, err)
	_, err = out.Stderr().Write([]byte("stderr"))
	require.NoError(t, err)

	require.EqualValues(t, VerbosityDefault, out.Verbosity())
	require.Equal(t, "stdout", outBuf.String())
	require.Equal(t, "stderr", errBuf.String())
}

func TestNewOutputSilent(t *testing.T) {
	var outBuf, errBuf bytes.Buffer
	out := NewOutput(
		WithStdout(xio.NopWriteCloser(&outBuf)),
		WithStderr(xio.NopWriteCloser(&errBuf)),
		Silent(),
	)

	_, err := out.Stdout().Write([]byte("stdout"))
	require.NoError(t, err)
	_, err = out.Stderr().Write([]byte("stderr"))
	require.NoError(t, err)

	require.EqualValues(t, VerbositySilent, out.Verbosity())
	require.Empty(t, outBuf.String())
	require.Empty(t, errBuf.String())
}

func TestNewOutputVerbose(t *testing.T) {
	var outBuf bytes.Buffer
	out := NewOutput(
		WithStdout(xio.NopWriteCloser(&outBuf)),
		CustomVerbose("ignite", "red"),
	)

	_, err := out.Stdout().Write([]byte("hello\n"))
	require.NoError(t, err)

	require.EqualValues(t, VerbosityVerbose, out.Verbosity())
	require.Contains(t, outBuf.String(), "hello")
	require.Contains(t, outBuf.String(), "IGNITE")
}
