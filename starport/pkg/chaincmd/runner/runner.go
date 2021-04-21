// Package chaincmdrunner provides a high level access to a blockchain's commands.
package chaincmdrunner

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/lineprefixer"
	"github.com/tendermint/starport/starport/pkg/truncatedbuffer"
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

// Cmd returns underlying chain cmd.
func (r Runner) Cmd() chaincmd.ChainCmd {
	return r.cc
}

type runOptions struct {
	// wrappedStdErrMaxLen determines the maximum length of the wrapped error logs
	// this option is used for long running command to prevent the buffer containing stderr getting too big
	// 0 can be used for no maximum length
	wrappedStdErrMaxLen int

	// stdout and stderr used to collect a copy of command's outputs.
	stdout, stderr io.Writer

	// stdin defines input for the command
	stdin io.Reader
}

// run executes a command.
func (r Runner) run(ctx context.Context, roptions runOptions, soptions ...step.Option) error {
	var (
		// we use a truncated buffer to prevent memory leak
		// this is because Stargate app currently send logs to StdErr
		// therefore if the app successfully starts, the written logs can become extensive
		errb = truncatedbuffer.NewTruncatedBuffer(roptions.wrappedStdErrMaxLen)

		// add optional prefixes to output streams.
		stdout io.Writer = lineprefixer.NewWriter(r.stdout, func() string { return r.daemonLogPrefix })
		stderr io.Writer = lineprefixer.NewWriter(r.stderr, func() string { return r.cliLogPrefix })
	)

	// Set standard outputs
	if roptions.stdout != nil {
		stdout = io.MultiWriter(stdout, roptions.stdout)
	}
	if roptions.stderr != nil {
		stderr = io.MultiWriter(stderr, roptions.stderr)
	}

	stderr = io.MultiWriter(stderr, errb)

	rnoptions := []cmdrunner.Option{
		cmdrunner.DefaultStdout(stdout),
		cmdrunner.DefaultStderr(stderr),
	}

	// Set standard input if defined
	if roptions.stdin != nil {
		rnoptions = append(rnoptions, cmdrunner.DefaultStdin(roptions.stdin))
	}

	err := cmdrunner.
		New(rnoptions...).
		Run(ctx, step.New(soptions...))

	return errors.Wrap(err, errb.GetBuffer().String())
}
