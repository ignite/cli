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
)

// Runner provides a high level access to a blockchain's commands.
type Runner struct {
	cc             chaincmd.ChainCmd
	stdout, stderr io.Writer
	workdir        string
}

// Option configures Runner.
type Option func(r *Runner)

// Stdout sets stdout for executed commands.
func Stdout(w io.Writer) Option {
	return func(r *Runner) {
		r.stdout = w
	}
}

// Stderr sets stderr for executed commands.
func Stderr(w io.Writer) Option {
	return func(r *Runner) {
		r.stderr = w
	}
}

// Workdir sets current working directory.
func Workdir(path string) Option {
	return func(r *Runner) {
		r.workdir = path
	}
}

// New creates a new Runner with cc and options.
func New(cc chaincmd.ChainCmd, options ...Option) Runner {
	r := Runner{
		cc:     cc,
		stdout: ioutil.Discard,
		stderr: ioutil.Discard,
	}

	// apply options.
	for _, apply := range options {
		apply(&r)
	}

	return r
}

type runOptions struct {
	// longRunning indicates that command expected to run for a long period of time.
	longRunning bool

	// stdout and stderr used to collect a copy of command's outputs.
	stdout, stderr io.Writer
}

// run executes a command.
func (r Runner) run(ctx context.Context, roptions runOptions, soptions ...step.Option) error {
	errb := &bytes.Buffer{}
	stdout := r.stdout
	stderr := r.stderr

	if roptions.stdout != nil {
		stdout = io.MultiWriter(stdout, roptions.stdout)
	}
	if roptions.stderr != nil {
		stderr = io.MultiWriter(stderr, roptions.stderr)
	}

	if !roptions.longRunning {
		stderr = io.MultiWriter(r.stderr, errb)
	}

	rnoptions := []cmdrunner.Option{
		cmdrunner.DefaultStdout(stdout),
		cmdrunner.DefaultStderr(stderr),
	}

	err := cmdrunner.
		New(rnoptions...).
		Run(ctx, step.New(soptions...))

	return errors.Wrap(err, errb.String())
}
