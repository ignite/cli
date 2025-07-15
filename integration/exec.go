package envtest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

type execOptions struct {
	ctx                    context.Context
	shouldErr, shouldRetry bool
	stdout, stderr         io.Writer
	stdin                  io.Reader
	tty                    bool
}

type ExecOption func(*execOptions)

// ExecShouldError sets the expectations of a command's execution to end with a failure.
func ExecShouldError() ExecOption {
	return func(o *execOptions) {
		o.shouldErr = true
	}
}

// ExecCtx sets cancelation context for the execution.
func ExecCtx(ctx context.Context) ExecOption {
	return func(o *execOptions) {
		o.ctx = ctx
	}
}

// ExecStdout captures stdout of an execution.
func ExecStdout(w io.Writer) ExecOption {
	return func(o *execOptions) {
		o.stdout = w
	}
}

// ExecStderr captures stderr of an execution.
func ExecStderr(w io.Writer) ExecOption {
	return func(o *execOptions) {
		o.stderr = w
	}
}

// ExecStdin captures stdin of an execution.
func ExecStdin(r io.Reader) ExecOption {
	return func(o *execOptions) {
		o.stdin = r
	}
}

// ExecRetry retries command until it is successful before context is canceled.
func ExecRetry() ExecOption {
	return func(o *execOptions) {
		o.shouldRetry = true
	}
}

// TTY simulates a TTY device.
func TTY() ExecOption {
	return func(o *execOptions) {
		o.tty = true
	}
}

// Exec executes a command step with options where msg describes the expectation from the test.
// unless calling with Must(), Exec() will not exit test runtime on failure.
func (e Env) Exec(msg string, steps step.Steps, options ...ExecOption) (ok bool) {
	opts := &execOptions{
		ctx:    e.ctx,
		stdout: io.Discard,
		stderr: io.Discard,
	}
	for _, o := range options {
		o(opts)
	}
	var (
		stdout = &bytes.Buffer{}
		stderr = &bytes.Buffer{}
	)
	copts := []cmdrunner.Option{
		cmdrunner.DefaultStdout(io.MultiWriter(stdout, opts.stdout)),
		cmdrunner.DefaultStderr(io.MultiWriter(stderr, opts.stderr)),
	}
	if opts.stdin != nil {
		copts = append(copts, cmdrunner.DefaultStdin(opts.stdin))
	}
	if opts.tty {
		copts = append(copts, cmdrunner.TTY())
	}

	if HasTestVerboseFlag() {
		fmt.Printf("Executing %d step(s) for %q\n", len(steps), msg)
		copts = append(copts, cmdrunner.EnableDebug())
	}
	if IsCI {
		copts = append(copts, cmdrunner.EndSignal(os.Kill))
	}
	err := cmdrunner.
		New(copts...).
		Run(opts.ctx, steps...)
	if errors.Is(err, context.Canceled) {
		err = nil
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if opts.shouldRetry && opts.ctx.Err() == nil {
			time.Sleep(time.Second)
			return e.Exec(msg, steps, options...)
		}

		msg = fmt.Sprintf("%s\n\nLogs:\n\n%s\n\nError Logs:\n\n%s\n",
			msg,
			stdout.String(),
			stderr.String())
	}

	if opts.shouldErr {
		return assert.Error(e.t, err, msg)
	}
	return assert.NoError(e.t, err, msg)
}
