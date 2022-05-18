package cliui

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/cliquiz"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/clispinner"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/entrywriter"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/icons"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/lineprefixer"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/prefixgen"
	"github.com/ignite-hq/cli/ignite/pkg/events"
	"github.com/manifoldco/promptui"
)

const (
	defaultLogStreamLabel = "ignite"
	defaultLogStreamColor = 201
)

type LogStream struct {
	stdout io.WriteCloser
	stderr io.WriteCloser
}

func (ls LogStream) Stdout() io.WriteCloser {
	return ls.stdout
}

func (ls LogStream) Stderr() io.WriteCloser {
	return ls.stderr
}

// Session controls command line interaction with users.
type (
	Verbosity uint8

	Session struct {
		ev       events.Bus
		eventsWg *sync.WaitGroup

		spinner                 *clispinner.Spinner
		startSpinnerImmediately bool

		in     io.ReadCloser
		stdout io.WriteCloser
		stderr io.WriteCloser

		verbosity                     Verbosity
		defaultLogStream              LogStream
		isDefaultLogStreamInitialised bool

		printLoopWg *sync.WaitGroup
	}

	LogStreamer interface {
		NewLogStream(label string, color uint8) (logStream LogStream)
	}

	Option func(s *Session)
)

const (
	VerbosityRegular = iota
	VerbositySilent
	VerbosityVerbose
)

// WithStdout sets base stdout for a session.
func WithStdout(stdout io.WriteCloser) Option {
	return func(s *Session) {
		s.stdout = stdout
	}
}

// WithStderr sets base stderr for a session
func WithStderr(stderr io.WriteCloser) Option {
	return func(s *Session) {
		s.stderr = stderr

	}
}

// WithInput sets input stream for a session.
func WithInput(input io.ReadCloser) Option {
	return func(s *Session) {
		s.in = input
	}
}

func WithVerbosity(v Verbosity) Option {
	return func(s *Session) {
		s.verbosity = v
	}
}

func StartSpinner() Option {
	return func(s *Session) {
		s.startSpinnerImmediately = true
	}
}

// New creates new Session.
func New(options ...Option) Session {
	wg := &sync.WaitGroup{}
	session := Session{
		ev:          events.NewBus(events.WithWaitGroup(wg)),
		in:          os.Stdin,
		stdout:      os.Stdout,
		stderr:      os.Stderr,
		eventsWg:    wg,
		printLoopWg: &sync.WaitGroup{},
	}
	for _, apply := range options {
		apply(&session)
	}

	session.defaultLogStream = session.NewLogStream(defaultLogStreamLabel, defaultLogStreamColor)
	if session.verbosity != VerbosityVerbose {
		session.verbosity = VerbositySilent
	}

	var spinnerOptions = []clispinner.Option{
		clispinner.WithWriter(session.defaultLogStream.Stdout()),
	}
	if session.startSpinnerImmediately {
		spinnerOptions = append(spinnerOptions, clispinner.StartImmediately())
	}
	session.spinner = clispinner.New(spinnerOptions...)

	session.printLoopWg.Add(1)
	go session.printLoop()
	return session
}

func (s Session) NewLogStream(label string, color uint8) (logStream LogStream) {
	prefixed := func(w io.Writer) *lineprefixer.Writer {
		options := prefixgen.Common(prefixgen.Color(color))
		prefixStr := prefixgen.New(label, options...).Gen()
		return lineprefixer.NewWriter(w, func() string { return prefixStr })
	}

	verbosity := s.verbosity
	if s.isDefaultLogStreamInitialised && verbosity != VerbosityVerbose {
		verbosity = VerbositySilent
	}
	s.isDefaultLogStreamInitialised = true

	switch verbosity {
	case VerbositySilent:
		logStream.stdout = ioutils.NopWriteCloser(io.Discard)
		logStream.stderr = ioutils.NopWriteCloser(io.Discard)
	case VerbosityVerbose:
		logStream.stdout = prefixed(s.stdout)
		logStream.stderr = prefixed(s.stderr)
	default:
		logStream.stdout = s.stdout
		logStream.stderr = s.stderr
	}

	return
}

// StopSpinner returns session's event bus.
func (s Session) EventBus() events.Bus {
	return s.ev
}

// StartSpinner starts spinner.
func (s Session) StartSpinner(text string) {
	s.spinner.SetText(text).Start()
}

// StopSpinner stops spinner.
func (s Session) StopSpinner() {
	s.spinner.Stop()
}

// PauseSpinner pauses spinner, returns resume function to start paused spinner again.
func (s Session) PauseSpinner() (mightResume func()) {
	isActive := s.spinner.IsActive()
	f := func() {
		if isActive {
			s.spinner.Start()
		}
	}
	s.spinner.Stop()
	return f
}

// Printf prints formatted arbitrary message.
func (s Session) Printf(format string, a ...interface{}) error {
	s.Wait()
	defer s.PauseSpinner()()
	_, err := fmt.Fprintf(s.defaultLogStream.Stdout(), format, a...)
	return err
}

// Println prints arbitrary message with line break.
func (s Session) Println(messages ...interface{}) error {
	s.Wait()
	defer s.PauseSpinner()()
	_, err := fmt.Fprintln(s.defaultLogStream.Stdout(), messages...)
	return err
}

// PrintSaidNo prints message informing negative was given in a confirmation prompt
func (s Session) PrintSaidNo() error {
	return s.Println("said no")
}

// Println prints arbitrary message
func (s Session) Print(messages ...interface{}) error {
	s.Wait()
	defer s.PauseSpinner()()
	_, err := fmt.Fprint(s.defaultLogStream.Stdout(), messages...)
	return err
}

// Ask asks questions in the terminal and collect answers.
func (s Session) Ask(questions ...cliquiz.Question) error {
	s.Wait()
	defer s.PauseSpinner()()
	// TODO provide writer from the session
	return cliquiz.Ask(questions...)
}

// AskConfirm asks yes/no question in the terminal.
func (s Session) AskConfirm(message string) error {
	s.Wait()
	defer s.PauseSpinner()()
	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
		Stdout:    s.defaultLogStream.Stdout(),
		Stdin:     s.in,
	}
	_, err := prompt.Run()
	return err
}

// PrintTable prints table data.
func (s Session) PrintTable(header []string, entries ...[]string) error {
	s.Wait()
	defer s.PauseSpinner()()
	return entrywriter.MustWrite(s.defaultLogStream.Stdout(), header, entries...)
}

// Wait blocks until all queued events are handled.
func (s Session) Wait() {
	s.eventsWg.Wait()
}

// Cleanup ensure spinner is stopped and printLoop exited correctly.
func (s Session) Cleanup() {
	s.StopSpinner()
	s.ev.Shutdown()
	s.printLoopWg.Wait()
}

// printLoop handles events.
func (s Session) printLoop() {
	for event := range s.ev.Events() {
		switch event.ProgressIndication {
		case events.IndicationStart:
			s.StartSpinner(event.String())

		case events.IndicationFinish:
			if event.Icon == "" {
				event.Icon = icons.OK
			}
			s.StopSpinner()
			fmt.Fprintf(s.defaultLogStream.Stdout(), "%s %s\n", event.Icon, event.String())

		case events.IndicationNone:
			resume := s.PauseSpinner()
			fmt.Fprintln(s.defaultLogStream.Stdout(), event.String())

			resume()
		}

		s.eventsWg.Done()
	}
	s.printLoopWg.Done()
}
