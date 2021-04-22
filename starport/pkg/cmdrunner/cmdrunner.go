package cmdrunner

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/goenv"
	"golang.org/x/sync/errgroup"
)

// Runner is an object to run commands
type Runner struct {
	endSignal   os.Signal
	stdout      io.Writer
	stderr      io.Writer
	stdin       io.Reader
	workdir     string
	runParallel bool
}

// Option defines option to run commands
type Option func(*Runner)

// DefaultStdout provides the default stdout for the commands to run
func DefaultStdout(writer io.Writer) Option {
	return func(r *Runner) {
		r.stdout = writer
	}
}

// DefaultStderr provides the default stderr for the commands to run
func DefaultStderr(writer io.Writer) Option {
	return func(r *Runner) {
		r.stderr = writer
	}
}

// DefaultStdin provides the default stdin for the commands to run
func DefaultStdin(reader io.Reader) Option {
	return func(r *Runner) {
		r.stdin = reader
	}
}

// DefaultWorkdir provides the default working directory for the commands to run
func DefaultWorkdir(path string) Option {
	return func(r *Runner) {
		r.workdir = path
	}
}

// RunParallel allows the commands to run concurrently
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

// New returns a new commands runner
func New(options ...Option) *Runner {
	r := &Runner{
		endSignal: os.Interrupt,
	}
	for _, o := range options {
		o(r)
	}
	return r
}

// Run blocks until all steps have completed their executions.
func (r *Runner) Run(ctx context.Context, steps ...*step.Step) error {
	if len(steps) == 0 {
		return nil
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
		} else if err := runPostExecs(c.Wait()); err != nil {
			return err
		}
	}
	return g.Wait()
}

// Executor is a command to execute
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

type cmdSignalNoWriter struct {
	*exec.Cmd
}

func (c *cmdSignalNoWriter) Signal(s os.Signal) { c.Cmd.Process.Signal(s) }

func (c *cmdSignalNoWriter) Write(data []byte) (n int, err error) { return 0, nil }

// newCommand returns a new command to execute
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
	c.Env = append(c.Env, os.ExpandEnv(fmt.Sprintf("PATH=$PATH:%s", goenv.GetGOBIN())))

	if r.stdin != nil {
		c.Stdin = r.stdin
		return &cmdSignalNoWriter{c}
	}

	w, err := c.StdinPipe()
	if err != nil {
		// TODO do not panic
		panic(err)
	}
	return &cmdSignal{c, w}
}
