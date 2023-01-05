package step

import (
	"io"
)

type Step struct {
	Exec      Execution
	PreExec   func() error
	InExec    func() error
	PostExecs []func(error) error
	Stdout    io.Writer
	Stderr    io.Writer
	Stdin     io.Reader
	Workdir   string
	Env       []string
	WriteData []byte
}

type Option func(*Step)

type Options []Option

func NewOptions() Options {
	return Options{}
}

func (o Options) Add(options ...Option) Options {
	return append(o, options...)
}

func New(options ...Option) *Step {
	s := &Step{
		PreExec:   func() error { return nil },
		InExec:    func() error { return nil },
		PostExecs: make([]func(error) error, 0),
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
		s.PostExecs = append(s.PostExecs, hook)
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

func Stdin(r io.Reader) Option {
	return func(s *Step) {
		s.Stdin = r
	}
}

func Workdir(path string) Option {
	return func(s *Step) {
		s.Workdir = path
	}
}

func Env(e ...string) Option {
	return func(s *Step) {
		s.Env = e
	}
}

func Write(data []byte) Option {
	return func(s *Step) {
		s.WriteData = data
	}
}

type Steps []*Step

func NewSteps(steps ...*Step) Steps {
	return steps
}

func (s *Steps) Add(steps ...*Step) Steps {
	*s = append(*s, steps...)
	return *s
}
