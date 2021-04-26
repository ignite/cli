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
	return func(runner *Runner) {
		runner.stdout = writer
	}
}

// DefaultStderr provides the default stderr for the commands to run
func DefaultStderr(writer io.Writer) Option {
	return func(runner *Runner) {
		runner.stderr = writer
	}
}

// DefaultStdin provides the default stdin for the commands to run
func DefaultStdin(reader io.Reader) Option {
	return func(runner *Runner) {
		runner.stdin = reader
	}
}

// DefaultWorkdir provides the default working directory for the commands to run
func DefaultWorkdir(path string) Option {
	return func(runner *Runner) {
		runner.workdir = path
	}
}

// RunParallel allows the commands to run concurrently
func RunParallel() Option {
	return func(runner *Runner) {
		runner.runParallel = true
	}
}

// EndSignal configures s to be signaled to the processes to end them.
func EndSignal(s os.Signal) Option {
	return func(runner *Runner) {
		runner.endSignal = s
	}
}

// New returns a new commands runner
func New(options ...Option) *Runner {
	runner := &Runner{
		endSignal: os.Interrupt,
	}
	for _, option := range options {
		option(runner)
	}
	return runner
}

// Run blocks until all steps have completed their executions.
func (r *Runner) Run(ctx context.Context, steps ...*step.Step) error {
	if len(steps) == 0 {
		return nil
	}
	g, ctx := errgroup.WithContext(ctx)
	for _, step := range steps {
		// copy s to a new variable to allocate a new address
		// so we can safely use it inside goroutines spawned in this loop.
		step := step
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := step.PreExec(); err != nil {
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
			for _, exec := range step.PostExecs {
				if err := exec(err); err != nil {
					return err
				}
			}
			if len(step.PostExecs) > 0 {
				return nil
			}
			return err
		}
		command := r.newCommand(step)
		startErr := command.Start()
		if startErr != nil {
			if err := runPostExecs(startErr); err != nil {
				return err
			}
			continue
		}
		go func() {
			<-ctx.Done()
			command.Signal(r.endSignal)
		}()
		if err := step.InExec(); err != nil {
			return err
		}
		if len(step.WriteData) > 0 {
			if _, err := command.Write(step.WriteData); err != nil {
				return err
			}
		}
		if r.runParallel {
			g.Go(func() error {
				return runPostExecs(command.Wait())
			})
		} else if err := runPostExecs(command.Wait()); err != nil {
			return err
		}
	}
	return g.Wait()
}

// Executor represents a command to execute
type Executor interface {
	Wait() error
	Start() error
	Signal(os.Signal)
	Write(data []byte) (n int, err error)
}

// dummyExecutor is an executor that does nothing
type dummyExecutor struct{}

func (exec *dummyExecutor) Start() error { return nil }

func (exec *dummyExecutor) Wait() error { return nil }

func (exec *dummyExecutor) Signal(os.Signal) {}

func (exec *dummyExecutor) Write([]byte) (int, error) { return 0, nil }

// cmdSignal is an executor with signal processing
type cmdSignal struct {
	*exec.Cmd
}

func (exec *cmdSignal) Signal(s os.Signal) { exec.Cmd.Process.Signal(s) }

func (exec *cmdSignal) Write(data []byte) (n int, err error) { return 0, nil }

// cmdSignalWithWriter is an executor with signal processing and that can write into stdin
type cmdSignalWithWriter struct {
	*exec.Cmd
	w io.WriteCloser
}

func (exec *cmdSignalWithWriter) Signal(s os.Signal) { exec.Cmd.Process.Signal(s) }

func (exec *cmdSignalWithWriter) Write(data []byte) (n int, err error) {
	defer exec.w.Close()
	return exec.w.Write(data)
}

// newCommand returns a new command to execute
func (r *Runner) newCommand(step *step.Step) Executor {
	// Return a dummy executor in case of an empty command
	if step.Exec.Command == "" {
		return &dummyExecutor{}
	}
	var (
		stdout = step.Stdout
		stderr = step.Stderr
		stdin  = step.Stdin
		dir    = step.Workdir
	)

	// Define standard input and outputs
	if stdout == nil {
		stdout = r.stdout
	}
	if stderr == nil {
		stderr = r.stderr
	}
	if stdin == nil {
		stdin = r.stdin
	}
	if dir == "" {
		dir = r.workdir
	}

	// Initialize command
	command := exec.Command(step.Exec.Command, step.Exec.Args...)
	command.Stdout = stdout
	command.Stderr = stderr
	command.Dir = dir
	command.Env = append(os.Environ(), step.Env...)
	command.Env = append(command.Env, os.ExpandEnv(
		fmt.Sprintf("PATH=$PATH:%s", goenv.GetGOBIN()),
	))

	// If a custom stdin is provided it will be as the stdin for the command
	if stdin != nil {
		command.Stdin = stdin
		return &cmdSignal{command}
	}

	// If no custom stdin, the executor can write into the stdin of the program
	writer, err := command.StdinPipe()
	if err != nil {
		// TODO do not panic
		panic(err)
	}
	return &cmdSignalWithWriter{command, writer}
}
