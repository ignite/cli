// Package chaincmdrunner provides a high level access to a blockchain's commands.
package chaincmdrunner

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/lineprefixer"
)

// Runner provides a high level access to a blockchain's commands.
type Runner struct {
	cc                            chaincmd.ChainCmd
	stdout, stderr                io.Writer
	daemonLogPrefix, cliLogPrefix string
}

// Option configures Runner.
type Option func(r *Runner)

// Stdout sets stdout for executed commands.
func Stdout(w io.Writer) Option {
	return func(r *Runner) {
		r.stdout = w
	}
}

// DaemonLogPrefix is a prefix added to app's daemon logs.
func DaemonLogPrefix(prefix string) Option {
	return func(r *Runner) {
		r.daemonLogPrefix = prefix
	}
}

// CLILogPrefix is a prefix added to app's cli logs.
func CLILogPrefix(prefix string) Option {
	return func(r *Runner) {
		r.cliLogPrefix = prefix
	}
}

// Stderr sets stderr for executed commands.
func Stderr(w io.Writer) Option {
	return func(r *Runner) {
		r.stderr = w
	}
}

// New creates a new Runner with cc and options.
func New(ctx context.Context, cc chaincmd.ChainCmd, options ...Option) (Runner, error) {
	r := Runner{
		cc:     cc,
		stdout: ioutil.Discard,
		stderr: ioutil.Discard,
	}

	applyOptions(&r, options)

	// auto detect the chain id and get it applied to chaincmd if auto
	// detection is enabled.
	if cc.IsAutoChainIDDetectionEnabled() {
		status, err := r.Status(ctx)
		if err != nil {
			return Runner{}, err
		}

		r.cc = r.cc.Copy(chaincmd.WithChainID(status.ChainID))
	}

	return r, nil
}

func applyOptions(r *Runner, options []Option) {
	for _, apply := range options {
		apply(r)
	}
}

// Copy makes a copy of runner by overwriting its options with given options.
func (r Runner) Copy(options ...Option) Runner {
	applyOptions(&r, options)

	return r
}

type runOptions struct {
	// wrapStdErr wraps the error logs of the app into the returned error
	wrapStdErr bool

	// wrappedStdErrMaxLen determines the maximum length of the wrapped error logs
	// 0 can be used for no maximum length
	wrappedStdErrMaxLen int

	// stdout and stderr used to collect a copy of command's outputs.
	stdout, stderr io.Writer
}

// run executes a command.
func (r Runner) run(ctx context.Context, roptions runOptions, soptions ...step.Option) error {
	var (
		// we use a truncated buffer to prevent memory leak
		// this is because Stargate app currently send logs to StdErr
		// therefore if the app successfully starts, the written logs can become extensive
		errb = &truncatedBuffer{
			buf: &bytes.Buffer{},
			cap: roptions.wrappedStdErrMaxLen,
		}

		// add optional prefixes to output streams.
		stdout io.Writer = lineprefixer.NewWriter(r.stdout, func() string { return r.daemonLogPrefix })
		stderr io.Writer = lineprefixer.NewWriter(r.stderr, func() string { return r.cliLogPrefix })
	)

	if roptions.stdout != nil {
		stdout = io.MultiWriter(stdout, roptions.stdout)
	}
	if roptions.stderr != nil {
		stderr = io.MultiWriter(stderr, roptions.stderr)
	}

	if roptions.wrapStdErr {
		stderr = io.MultiWriter(stderr, errb)
	}

	rnoptions := []cmdrunner.Option{
		cmdrunner.DefaultStdout(stdout),
		cmdrunner.DefaultStderr(stderr),
	}

	err := cmdrunner.
		New(rnoptions...).
		Run(ctx, step.New(soptions...))

	if roptions.wrapStdErr {
		return errors.Wrap(err, errb.buf.String())
	}
	return err
}

type truncatedBuffer struct {
	buf *bytes.Buffer
	cap int
}

func (b *truncatedBuffer) Write(p []byte) (n int, err error) {
	n, err = b.buf.Write(p)
	if err != nil {
		return n, err
	}
	if b.cap > 0 && b.buf.Len() > b.cap {
		b.buf.Truncate(b.cap)
	}

	return n, nil
}
