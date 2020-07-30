package cmdrunner

import (
	"context"
	"io"
	"os/exec"

	"github.com/tendermint/starport/pkg/cmdrunner/step"
	"golang.org/x/sync/errgroup"
)

type Runner struct {
	stdout      io.Writer
	stderr      io.Writer
	workdir     string
	runParallel bool
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

func RunParallel() Option {
	return func(r *Runner) {
		r.runParallel = true
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
	g, ctx := errgroup.WithContext(ctx)
	for _, s := range steps {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := s.PreExec(); err != nil {
			return err
		}
		c := r.newCommand(ctx, s)
		startErr := c.Start()
		if startErr != nil {
			if err := s.PostExec(startErr); err != nil {
				return err
			}
			continue
		}
		if err := s.InExec(); err != nil {
			return err
		}
		if r.runParallel {
			g.Go(func() error {
				return s.PostExec(c.Wait())
			})
		} else {
			if err := s.PostExec(c.Wait()); err != nil {
				return err
			}
		}
	}
	return g.Wait()
}

func (r *Runner) newCommand(ctx context.Context, s *step.Step) *exec.Cmd {
	if s.Exec.Command == "" {
		// this is a programmer error so better to panic instead of
		// returning an err.
		panic("empty command")
	}
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
	c := exec.CommandContext(ctx, s.Exec.Command, s.Exec.Args...)
	c.Stdout = stdout
	c.Stderr = stderr
	c.Dir = dir
	return c
}
