package cliui

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/manifoldco/promptui"

	"github.com/ignite/cli/ignite/pkg/cliui/cliquiz"
	"github.com/ignite/cli/ignite/pkg/cliui/clispinner"
	"github.com/ignite/cli/ignite/pkg/cliui/entrywriter"
	uilog "github.com/ignite/cli/ignite/pkg/cliui/log"
	"github.com/ignite/cli/ignite/pkg/events"
)

type sessionOptions struct {
	stdin  io.ReadCloser
	stdout io.WriteCloser
	stderr io.WriteCloser

	spinnerStart bool
	spinnerText  string

	ignoreEvents bool
	verbosity    uilog.Verbosity
}

// Session controls command line interaction with users.
type Session struct {
	options sessionOptions
	ev      events.Bus
	spinner *clispinner.Spinner
	out     uilog.Output
	wg      *sync.WaitGroup
	ended   bool
}

// Option configures session options.
type Option func(s *Session)

// WithStdout sets the starndard output for the session.
func WithStdout(stdout io.WriteCloser) Option {
	return func(s *Session) {
		s.options.stdout = stdout
	}
}

// WithStderr sets base stderr for a Session.
func WithStderr(stderr io.WriteCloser) Option {
	return func(s *Session) {
		s.options.stderr = stderr
	}
}

// WithStdin sets the starndard input for the session.
func WithStdin(stdin io.ReadCloser) Option {
	return func(s *Session) {
		s.options.stdin = stdin
	}
}

// WithVerbosity sets a verbosity level for the Session.
func WithVerbosity(v uilog.Verbosity) Option {
	return func(s *Session) {
		s.options.verbosity = v
	}
}

// IgnoreEvents configures the session to avoid displaying events.
// This is a compatibility option to be able to use the session and
// the events bus when models are used to manage CLI UI. The session
// won't handle the events when this option is present.
func IgnoreEvents() Option {
	return func(s *Session) {
		s.options.ignoreEvents = true
	}
}

// StartSpinner forces spinner to be spinning right after creation.
func StartSpinner() Option {
	return func(s *Session) {
		s.options.spinnerStart = true
	}
}

// StartSpinnerWithText forces spinner to be spinning right after creation
// with a custom status text.
func StartSpinnerWithText(text string) Option {
	return func(s *Session) {
		s.options.spinnerStart = true
		s.options.spinnerText = text
	}
}

// New creates a new Session.
func New(options ...Option) *Session {
	session := Session{
		ev: events.NewBus(),
		wg: &sync.WaitGroup{},
		options: sessionOptions{
			stdin:       os.Stdin,
			stdout:      os.Stdout,
			stderr:      os.Stderr,
			spinnerText: clispinner.DefaultText,
		},
	}

	for _, apply := range options {
		apply(&session)
	}

	logOptions := []uilog.Option{
		uilog.WithStdout(session.options.stdout),
		uilog.WithStderr(session.options.stderr),
	}

	if session.options.verbosity == uilog.VerbosityVerbose {
		logOptions = append(logOptions, uilog.Verbose())
	}

	session.out = uilog.NewOutput(logOptions...)

	if session.options.spinnerStart {
		session.StartSpinner(session.options.spinnerText)
	}

	// The main loop that prints the events uses a wait group to block
	// the session end until all the events are printed.
	if !session.options.ignoreEvents {
		session.wg.Add(1)
		go session.handleEvents()
	}

	return &session
}

// EventBus returns the event bus of the session.
func (s Session) EventBus() events.Bus {
	return s.ev
}

// Verbosity returns the verbosity level for the session output.
func (s Session) Verbosity() uilog.Verbosity {
	return s.options.verbosity
}

// NewOutput returns a new logging output bound to the session.
// The new output will use the session's verbosity, stderr and stdout.
// Label and color arguments are used to prefix the output when the
// session verbosity is verbose.
func (s Session) NewOutput(label, color string) uilog.Output {
	options := []uilog.Option{
		uilog.WithStdout(s.options.stdout),
		uilog.WithStderr(s.options.stderr),
	}

	if s.options.verbosity == uilog.VerbosityVerbose {
		options = append(options, uilog.CustomVerbose(label, color))
	}

	return uilog.NewOutput(options...)
}

// StartSpinner starts the spinner.
func (s *Session) StartSpinner(text string) {
	if s.options.ignoreEvents {
		return
	}

	// Verbose mode must not render the spinner but instead
	// it should just print the text to display next to the
	// app label otherwise the verbose logs would be printed
	// with an invalid format.
	if s.options.verbosity == uilog.VerbosityVerbose {
		fmt.Fprint(s.out.Stdout(), text)
		return
	}

	if s.spinner == nil {
		s.spinner = clispinner.New(clispinner.WithWriter(s.out.Stdout()))
	}

	s.spinner.SetText(text).Start()
}

// StopSpinner stops the spinner.
func (s Session) StopSpinner() {
	if s.spinner == nil {
		return
	}

	s.spinner.Stop()
}

// PauseSpinner pauses spinner and returns a function to restart the spinner.
func (s Session) PauseSpinner() (restart func()) {
	isActive := s.spinner != nil && s.spinner.IsActive()
	if isActive {
		s.spinner.Stop()
	}

	return func() {
		if isActive {
			s.spinner.Start()
		}
	}
}

// Printf prints formatted arbitrary message.
func (s Session) Printf(format string, a ...interface{}) error {
	defer s.PauseSpinner()()
	_, err := fmt.Fprintf(s.out.Stdout(), format, a...)
	return err
}

// Println prints arbitrary message with line break.
func (s Session) Println(messages ...interface{}) error {
	defer s.PauseSpinner()()
	_, err := fmt.Fprintln(s.out.Stdout(), messages...)
	return err
}

// Print prints arbitrary message.
func (s Session) Print(messages ...interface{}) error {
	defer s.PauseSpinner()()
	_, err := fmt.Fprint(s.out.Stdout(), messages...)
	return err
}

// Ask asks questions in the terminal and collect answers.
func (s Session) Ask(questions ...cliquiz.Question) error {
	defer s.PauseSpinner()()
	// TODO provide writer from the session
	return cliquiz.Ask(questions...)
}

// AskConfirm asks yes/no question in the terminal.
func (s Session) AskConfirm(message string) error {
	defer s.PauseSpinner()()
	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
		Stdout:    s.out.Stdout(),
		Stdin:     s.options.stdin,
	}
	_, err := prompt.Run()
	return err
}

// PrintTable prints table data.
func (s Session) PrintTable(header []string, entries ...[]string) error {
	defer s.PauseSpinner()()
	return entrywriter.MustWrite(s.out.Stdout(), header, entries...)
}

// End finishes the session by stopping the spinner and the event bus.
// Once the session is ended it should not be used anymore.
func (s *Session) End() {
	if s.ended {
		return
	}

	s.StopSpinner()
	s.ev.Stop()
	s.wg.Wait()
	s.ended = true
}

func (s *Session) handleEvents() {
	defer s.wg.Done()

	stdout := s.out.Stdout()

	for e := range s.ev.Events() {
		switch e.ProgressIndication {
		case events.IndicationStart:
			s.StartSpinner(e.String())
		case events.IndicationUpdate:
			if s.spinner == nil {
				// When the spinner is not initialized print the event
				fmt.Fprintf(stdout, "%s\n", e)
			} else {
				// Otherwise update the spinner with a new text
				s.spinner.SetText(e.String())
			}
		case events.IndicationFinish:
			s.StopSpinner()
			fmt.Fprintf(stdout, "%s\n", e)
		case events.IndicationNone:
		default:
			// The text printed here won't be removed when the spinner stops
			resume := s.PauseSpinner()
			fmt.Fprintf(stdout, "%s\n", e)
			resume()
		}
	}
}
