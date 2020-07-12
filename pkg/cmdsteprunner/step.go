package cmdsteprunner

import (
	"io"
)

type Step struct {
	o stepOptions
}

type StepOption func(*stepOptions)

type stepOptions struct {
	afterExec func(error) error
	stdout    io.Writer
	stderr    io.Writer
}

var defaultAfterExec = func(err error) error { return err }

func StepStdout(w io.Writer) StepOption {
	return func(o *stepOptions) {
		o.stdout = w
	}
}

func StepStderr(w io.Writer) StepOption {
	return func(o *stepOptions) {
		o.stderr = w
	}
}

func StepStdouterr(stdout, stderr io.Writer) StepOption {
	return func(o *stepOptions) {
		StepStdout(stdout)(o)
		StepStderr(stderr)(o)
	}
}

// StepPreExec hook is executed just before executing the step.
// returning a non-nil error from the hook will make the other steps not run if
// continue on failure has not been enabled on Runner.
func StepPreExec(hook func() error) StepOption {
	return func(o *stepOptions) {
	}
}

// StepAfterExec hook is executed after command is complated.
// err in the hook is filled if execution is complated with a non-zero code.
//
// returning a non-nil error from the hook will make the other steps not run if
// continue on failure has not been enabled on Runner.
func StepAfterExec(hook func(err error) error) StepOption {
	return func(o *stepOptions) {
	}
}

type Command []string

func StepCommand(name string, arg ...string) Command {
	return Command(append([]string{name}, arg...))
}

func NewStep(cmd Command, options ...StepOption) Step {
	opts := stepOptions{
		afterExec: defaultAfterExec,
	}
	for _, o := range options {
		o(&opts)
	}
	s := Step{
		o: opts,
	}
	return s

}
