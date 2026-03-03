package cmdrunner

import (
	"bytes"
	"context"
	stdErrors "errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/pkg/goenv"
)

func TestNewAppliesOptions(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	stdin := strings.NewReader("stdin")

	r := New(
		DefaultStdout(stdout),
		DefaultStderr(stderr),
		DefaultStdin(stdin),
		DefaultWorkdir("/tmp/work"),
		RunParallel(),
		EndSignal(os.Kill),
		EnableDebug(),
	)

	require.Equal(t, stdout, r.stdout)
	require.Equal(t, stderr, r.stderr)
	require.Equal(t, stdin, r.stdin)
	require.Equal(t, "/tmp/work", r.workdir)
	require.True(t, r.runParallel)
	require.Equal(t, os.Kill, r.endSignal)
	require.True(t, r.debug)
}

func TestEnv(t *testing.T) {
	require.Equal(t, "KEY=value", Env("KEY", "value"))
}

func TestNewCommandReturnsDummyExecutorForEmptyCommand(t *testing.T) {
	executor := New().newCommand(step.New())
	_, ok := executor.(*dummyExecutor)
	require.True(t, ok)
}

func TestNewCommandUsesDefaultsWhenStepDoesNotProvideIO(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	r := New(
		DefaultStdout(stdout),
		DefaultStderr(stderr),
		DefaultWorkdir("/tmp/work"),
	)

	executor := r.newCommand(step.New(
		step.Exec("echo", "hello"),
		step.Env("A=B"),
	))

	cmdExecutor, ok := executor.(*cmdSignalWithWriter)
	require.True(t, ok)
	require.Equal(t, stdout, cmdExecutor.Stdout)
	require.Equal(t, stderr, cmdExecutor.Stderr)
	require.Equal(t, "/tmp/work", cmdExecutor.Dir)
	require.Contains(t, cmdExecutor.Env, "A=B")
	require.Contains(t, cmdExecutor.Env, Env("PATH", goenv.Path()))
}

func TestNewCommandWithCustomStdinReturnsCmdSignal(t *testing.T) {
	stdin := strings.NewReader("input")

	executor := New().newCommand(step.New(
		step.Exec("echo"),
		step.Stdin(stdin),
	))

	cmdExecutor, ok := executor.(*cmdSignal)
	require.True(t, ok)
	require.Equal(t, stdin, cmdExecutor.Stdin)
}

func TestRunWithoutStepsReturnsNil(t *testing.T) {
	err := New().Run(context.Background())
	require.NoError(t, err)
}

func TestRunReturnsPreExecError(t *testing.T) {
	expectedErr := stdErrors.New("pre exec error")

	err := New().Run(context.Background(), step.New(
		step.PreExec(func() error { return expectedErr }),
	))

	require.ErrorIs(t, err, expectedErr)
}

func TestRunReturnsStartErrorWithoutPostExec(t *testing.T) {
	err := New().Run(context.Background(), step.New(
		step.Exec("this-command-does-not-exist-cmdrunner-test"),
	))

	require.Error(t, err)
}

func TestRunCanHandleStartErrorInPostExec(t *testing.T) {
	var receivedErr error

	err := New().Run(context.Background(), step.New(
		step.Exec("this-command-does-not-exist-cmdrunner-test"),
		step.PostExec(func(err error) error {
			receivedErr = err
			return nil
		}),
	))

	require.NoError(t, err)
	require.Error(t, receivedErr)
}
