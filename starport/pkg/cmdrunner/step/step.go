package step

import "io"

type Step struct {
	Exec     Execution
	PreExec  func() error
	InExec   func() error
	PostExec func(error) error
	Stdout   io.Writer
	Stderr   io.Writer
	Workdir  string
}

type Option func(*Step)

func New(options ...Option) *Step {
	s := &Step{
		PreExec:  func() error { return nil },
		InExec:   func() error { return nil },
		PostExec: func(exitErr error) error { return exitErr },
	}
	for _, o := range options {
		o(s)
	}
	return s
}

type Execution struct {
	Command string
	Args    []string
}

func Exec(command string, args ...string) Option {
	return func(s *Step) {
		s.Exec = Execution{command, args}
	}
}

func PreExec(hook func() error) Option {
	return func(s *Step) {
		s.PreExec = hook
	}
}

func InExec(hook func() error) Option {
	return func(s *Step) {
		s.InExec = hook
	}
}

func PostExec(hook func(exitErr error) error) Option { // *os.ExitError
	return func(s *Step) {
		s.PostExec = hook
	}
}

func Stdout(w io.Writer) Option {
	return func(s *Step) {
		s.Stdout = w
	}
}

func Stderr(w io.Writer) Option {
	return func(s *Step) {
		s.Stderr = w
	}
}

func Workdir(path string) Option {
	return func(s *Step) {
		s.Workdir = path
	}
}

type Steps []*Step

func (s *Steps) Add(step *Step) {
	*s = append(*s, step)
}
