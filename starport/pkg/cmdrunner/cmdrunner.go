package cmdrunner

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"golang.org/x/sync/errgroup"
)

type Runner struct {
	endSignal   os.Signal
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

// EndSignal configures s to be signaled to the processes to end them.
func EndSignal(s os.Signal) Option {
	return func(r *Runner) {
		r.endSignal = s
	}
}

func New(options ...Option) *Runner {
	r := &Runner{
		endSignal: os.Interrupt,
	}
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
		runPostExec := func(processErr error) error {
			// if context is canceled, then we can ignore exit error of the
			// process because it should be exited because of the cancellation.
			var err error
			ctxErr := ctx.Err()
			if ctxErr != nil {
				err = ctxErr
			} else {
				err = processErr
			}
			return s.PostExec(err)
		}
		c := r.newCommand(ctx, s)
		startErr := c.Start()
		if startErr != nil {
			if err := runPostExec(startErr); err != nil {
				return err
			}
			continue
		}
		if err := s.InExec(); err != nil {
			return err
		}
		if r.runParallel {
			g.Go(func() error {
				return runPostExec(c.Wait())
			})
		} else {
			if err := runPostExec(c.Wait()); err != nil {
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
	c := exec.Command(s.Exec.Command, s.Exec.Args...)
	c.Stdout = stdout
	c.Stderr = stderr
	c.Dir = dir
	go func() {
		<-ctx.Done()
		if c.Process != nil {
			c.Process.Signal(r.endSignal)
		}
	}()
	return c
}
