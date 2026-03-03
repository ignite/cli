package step

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDefaults(t *testing.T) {
	s := New()
	require.NoError(t, s.PreExec())
	require.NoError(t, s.InExec())
	require.Empty(t, s.PostExecs)
}

func TestNewAppliesOptions(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	stdin := strings.NewReader("stdin")
	postErr := errors.New("post")

	s := New(
		Exec("go", "version"),
		Stdout(stdout),
		Stderr(stderr),
		Stdin(stdin),
		Workdir("/tmp/work"),
		Env("A=B", "C=D"),
		Write([]byte("payload")),
		PreExec(func() error { return postErr }),
		InExec(func() error { return postErr }),
		PostExec(func(err error) error { return err }),
	)

	require.Equal(t, Execution{Command: "go", Args: []string{"version"}}, s.Exec)
	require.Equal(t, stdout, s.Stdout)
	require.Equal(t, stderr, s.Stderr)
	require.Equal(t, stdin, s.Stdin)
	require.Equal(t, "/tmp/work", s.Workdir)
	require.Equal(t, []string{"A=B", "C=D"}, s.Env)
	require.Equal(t, []byte("payload"), s.WriteData)
	require.ErrorIs(t, s.PreExec(), postErr)
	require.ErrorIs(t, s.InExec(), postErr)
	require.Len(t, s.PostExecs, 1)
	require.ErrorIs(t, s.PostExecs[0](postErr), postErr)
}

func TestOptionsAdd(t *testing.T) {
	options := NewOptions().Add(Exec("go"), Workdir("/tmp"))
	require.Len(t, options, 2)
}

func TestStepsAdd(t *testing.T) {
	s1 := New(Exec("one"))
	s2 := New(Exec("two"))
	steps := NewSteps(s1)
	got := (&steps).Add(s2)
	require.Len(t, got, 2)
	require.Equal(t, "one", got[0].Exec.Command)
	require.Equal(t, "two", got[1].Exec.Command)
}
