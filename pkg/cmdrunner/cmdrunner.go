package cmdrunner

import (
	"context"
	"io"
	"os/exec"

	"github.com/tendermint/starport/pkg/cmdrunner/step"
)

type Runner struct {
	stdout  io.Writer
	stderr  io.Writer
	workdir string
}

type Option func(*Runner)

func DefaultStdout(w io.Writer) Option {
	return func(r *Runner) {
		r.stdout = w
	}
}

func DefaultStderr(w io.Writer) Option {
	return func(r *Runner) {
		r.stderr = w
	}
}

func DefaultWorkdir(path string) Option {
	return func(r *Runner) {
		r.workdir = path
	}
}

func New(options ...Option) *Runner {
	r := &Runner{}
	for _, o := range options {
		o(r)
	}
	return r
}

// Run blocks untill all steps are complated their executions.
func (r *Runner) Run(ctx context.Context, steps ...*step.Step) error {
	if len(steps) == 0 {
		// this is a programmer error so better to panic instead of
		// returning an err.
		panic("no steps to run")
	}
	for _, s := range steps {
		if err := ctx.Err(); err != nil {
			return err
		}
		if s.PreExec != nil {
			if err := s.PreExec(); err != nil {
				return err
			}
		}
		execErr := r.runStep(ctx, s)
		if s.PostExec != nil {
			if err := s.PostExec(execErr); err != nil {
				return err
			}
		}
		if execErr != nil {
			return execErr
		}
	}
	return nil
}

func (r *Runner) runStep(ctx context.Context, s *step.Step) error {
	if s.Exec.Command == "" {
		// this is a programmer error so better to panic instead of
		// returning an err.
		panic("empty command")
	}
	c := exec.CommandContext(ctx, s.Exec.Command, s.Exec.Args...)
	var (
		stdout = s.Stdout
		stderr = s.Stderr
		dir    = s.Workdir
	)
	if stdout == nil {
		stdout = r.stdout
	}
	if stderr == nil {
		stderr = r.stderr
	}
	if dir == "" {
		dir = r.workdir
	}
	c.Stdout = stdout
	c.Stderr = stderr
	c.Dir = dir
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}
