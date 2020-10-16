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
		// copy s to a new variable to allocate a new address
		// so we can safely use it inside goroutines spawned in this loop.
		s := s
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := s.PreExec(); err != nil {
			return err
		}
		runPostExecs := func(processErr error) error {
			// if context is canceled, then we can ignore exit error of the
			// process because it should be exited because of the cancellation.
			var err error
			ctxErr := ctx.Err()
			if ctxErr != nil {
				err = ctxErr
			} else {
				err = processErr
			}
			for _, exec := range s.PostExecs {
				if err := exec(err); err != nil {
					return err
				}
			}
			if len(s.PostExecs) > 0 {
				return nil
			}
			return err
		}
		c := r.newCommand(s)
		startErr := c.Start()
		if startErr != nil {
			if err := runPostExecs(startErr); err != nil {
				return err
			}
			continue
		}
		go func() {
			<-ctx.Done()
			c.Signal(r.endSignal)
		}()
		if err := s.InExec(); err != nil {
			return err
		}
		if len(s.WriteData) > 0 {
			if _, err := c.Write(s.WriteData); err != nil {
				return err
			}
		}
		if r.runParallel {
			g.Go(func() error {
				return runPostExecs(c.Wait())
			})
		} else {
			if err := runPostExecs(c.Wait()); err != nil {
				return err
			}
		}
	}
	return g.Wait()
}

type Executor interface {
	Wait() error
	Start() error
	Signal(os.Signal)
	Write(data []byte) (n int, err error)
}

type dummyExecutor struct{}

func (s *dummyExecutor) Start() error { return nil }

func (s *dummyExecutor) Wait() error { return nil }

func (s *dummyExecutor) Signal(os.Signal) {}

func (s *dummyExecutor) Write([]byte) (int, error) { return 0, nil }

type cmdSignal struct {
	*exec.Cmd
	w io.WriteCloser
}

func (c *cmdSignal) Signal(s os.Signal) { c.Cmd.Process.Signal(s) }

func (c *cmdSignal) Write(data []byte) (n int, err error) {
	defer c.w.Close()
	return c.w.Write(data)
}

func (r *Runner) newCommand(s *step.Step) Executor {
	if s.Exec.Command == "" {
		return &dummyExecutor{}
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
	c.Env = append(os.Environ(), s.Env...)
	w, err := c.StdinPipe()
	if err != nil {
		// TODO do not panic
		panic(err)
	}
	return &cmdSignal{c, w}
}
