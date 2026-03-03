package clispinner

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewUsesSimpleSpinnerForNonTerminalWriter(t *testing.T) {
	s := New(WithWriter(&bytes.Buffer{}))
	_, ok := s.(*SimpleSpinner)
	require.True(t, ok)
}

func TestIsRunningInTerminalFalseForNonFileWriter(t *testing.T) {
	require.False(t, isRunningInTerminal(&bytes.Buffer{}))
}

func TestIsRunningInTerminalFalseForRegularFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "spinner-writer-*")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, f.Close())
	})

	require.False(t, isRunningInTerminal(f))
}

func TestNewSimpleSpinnerDefaults(t *testing.T) {
	s := newSimpleSpinner(Options{})

	require.Equal(t, DefaultText, s.text)
	require.Equal(t, simpleCharset, s.charset)
	require.Equal(t, os.Stdout, s.writer)
}

func TestSimpleSpinnerSetters(t *testing.T) {
	s := newSimpleSpinner(Options{})

	require.Same(t, s, s.SetText("text"))
	require.Same(t, s, s.SetPrefix("prefix"))
	require.Same(t, s, s.SetCharset([]string{"1", "2"}))
	require.Same(t, s, s.SetColor("red"))

	require.Equal(t, "text", s.text)
	require.Equal(t, "prefix", s.prefix)
	require.Equal(t, []string{"1", "2"}, s.charset)
	require.Equal(t, "red", s.color)
}

func TestSimpleSpinnerStartAndStop(t *testing.T) {
	oldRefreshRate := simpleRefreshRate
	oldColor := simpleColor
	simpleRefreshRate = time.Millisecond
	simpleColor = func(i ...interface{}) string { return fmt.Sprint(i...) }
	t.Cleanup(func() {
		simpleRefreshRate = oldRefreshRate
		simpleColor = oldColor
	})

	out := &safeBuffer{}
	s := newSimpleSpinner(Options{
		writer:  out,
		text:    "working",
		charset: []string{"."},
	})

	require.False(t, s.IsActive())
	s.Start()
	require.True(t, s.IsActive())

	require.Eventually(t, func() bool {
		return out.Len() > 0
	}, 200*time.Millisecond, 5*time.Millisecond)

	s.Stop()
	require.False(t, s.IsActive())
}

func TestNewTermSpinnerDefaultsAndSetters(t *testing.T) {
	var out bytes.Buffer
	s := newTermSpinner(Options{
		writer:  &out,
		text:    "booting",
		charset: []string{"A", "B"},
	})

	require.Equal(t, []string{"A", "B"}, s.charset)
	require.Equal(t, &out, s.Writer())
	require.Equal(t, " booting", s.sp.Suffix)

	require.Same(t, s, s.SetText("running"))
	require.Same(t, s, s.SetPrefix("ignite"))
	require.Same(t, s, s.SetCharset([]string{"X"}))
	require.Same(t, s, s.SetColor("green"))

	require.Equal(t, "ignite ", s.sp.Prefix)
	require.Equal(t, " running", s.sp.Suffix)
}

type safeBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *safeBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *safeBuffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Len()
}
