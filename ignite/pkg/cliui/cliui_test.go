package cliui

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cliui/clispinner"
	uilog "github.com/ignite/cli/v29/ignite/pkg/cliui/log"
	"github.com/ignite/cli/v29/ignite/pkg/xio"
)

type fakeSpinner struct {
	active bool
}

func (f *fakeSpinner) SetText(string) clispinner.Spinner      { return f }
func (f *fakeSpinner) SetPrefix(string) clispinner.Spinner    { return f }
func (f *fakeSpinner) SetCharset([]string) clispinner.Spinner { return f }
func (f *fakeSpinner) SetColor(string) clispinner.Spinner     { return f }
func (f *fakeSpinner) Start() clispinner.Spinner {
	f.active = true
	return f
}

func (f *fakeSpinner) Stop() clispinner.Spinner {
	f.active = false
	return f
}
func (f *fakeSpinner) IsActive() bool { return f.active }
func (f *fakeSpinner) Writer() io.Writer {
	return io.Discard
}

func TestNewWithOptions(t *testing.T) {
	var outBuf, errBuf bytes.Buffer
	session := New(
		WithStdout(xio.NopWriteCloser(&outBuf)),
		WithStderr(xio.NopWriteCloser(&errBuf)),
		WithStdin(io.NopCloser(strings.NewReader(""))),
		IgnoreEvents(),
		WithVerbosity(uilog.VerbosityVerbose),
	)
	t.Cleanup(session.End)

	require.EqualValues(t, uilog.VerbosityVerbose, session.Verbosity())
}

func TestAskAndAskConfirmSkipUI(t *testing.T) {
	session := New(IgnoreEvents(), WithoutUserInteraction(true))
	t.Cleanup(session.End)

	require.NoError(t, session.Ask())
	require.NoError(t, session.AskConfirm("continue?"))
}

func TestPauseSpinner(t *testing.T) {
	session := New(IgnoreEvents())
	t.Cleanup(session.End)

	sp := &fakeSpinner{active: true}
	session.spinner = sp

	restart := session.PauseSpinner()
	require.False(t, sp.IsActive())

	restart()
	require.True(t, sp.IsActive())
}

func TestStartSpinnerVerboseWritesText(t *testing.T) {
	var outBuf bytes.Buffer
	session := New(
		WithStdout(xio.NopWriteCloser(&outBuf)),
		WithVerbosity(uilog.VerbosityVerbose),
	)
	t.Cleanup(session.End)

	session.StartSpinner("working")
	require.Contains(t, outBuf.String(), "working")
}

func TestEndIsIdempotent(t *testing.T) {
	session := New(IgnoreEvents())
	session.End()
	require.True(t, session.ended)
	session.End()
	require.True(t, session.ended)
}
