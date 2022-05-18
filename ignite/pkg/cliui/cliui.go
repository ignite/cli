package cliui

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/manifoldco/promptui"

	"github.com/ignite-hq/cli/ignite/pkg/cliui/cliquiz"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/clispinner"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/entrywriter"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/icons"
	"github.com/ignite-hq/cli/ignite/pkg/events"
)

// Verbosity enumerates possible verbosity levels for cli output
type Verbosity uint8

const (
	VerbosityRegular = iota
	VerbositySilent
	VerbosityVerbose
)

// Session controls command line interaction with users.
type Session struct {
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

// Option is used to customize Session during creation
type Option func(s *Session)

// WithStdout sets base stdout for a Session.
func WithStdout(stdout io.WriteCloser) Option {
	return func(s *Session) {
		s.stdout = stdout
	}
}

// WithStderr sets base stderr for a Session
func WithStderr(stderr io.WriteCloser) Option {
	return func(s *Session) {
		s.stderr = stderr

	}
}

// WithInput sets input stream for a Session.
func WithInput(input io.ReadCloser) Option {
	return func(s *Session) {
		s.in = input
	}
}

// WithVerbosity sets verbosity level for Session
func WithVerbosity(v Verbosity) Option {
	return func(s *Session) {
		s.verbosity = v
	}
}

// StartSpinner forces spinner to be spinning right after creation
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

	if session.startSpinnerImmediately {
		session.spinner = clispinner.New(clispinner.WithWriter(session.defaultLogStream.Stdout()))
	}

	session.defaultLogStream = session.NewLogStream(defaultLogStreamLabel, defaultLogStreamColor)
	if session.verbosity != VerbosityVerbose {
		session.verbosity = VerbositySilent
	}

	session.printLoopWg.Add(1)
	go session.printLoop()
	return session
}

// StopSpinner returns session's event bus.
func (s Session) EventBus() events.Bus {
	return s.ev
}

// StartSpinner starts spinner.
func (s Session) StartSpinner(text string) {
	if s.spinner == nil {
		s.spinner = clispinner.New(clispinner.WithWriter(s.defaultLogStream.Stdout()))
	}
	s.spinner.SetText(text).Start()
}

// StopSpinner stops spinner.
func (s Session) StopSpinner() {
	if s.spinner == nil {
		return
	}
	s.spinner.Stop()
}

// PauseSpinner pauses spinner, returns resume function to start paused spinner again.
func (s Session) PauseSpinner() (mightResume func()) {
	isActive := s.spinner != nil && s.spinner.IsActive()
	f := func() {
		if isActive {
			s.spinner.Start()
		}
	}

	if isActive {
		s.spinner.Stop()
	}
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
			s.StartSpinner(event.Content.String())

		case events.IndicationFinish:
			if event.Icon == "" {
				event.Icon = icons.OK
			}
			s.StopSpinner()
			if event.HasIcon() {
				fmt.Fprintf(s.defaultLogStream.Stdout(), "%s %s\n", event.Icon, event.Content.String())
			} else {
				fmt.Fprintf(s.defaultLogStream.Stdout(), "%s\n", event.Content.String())
			}

		case events.IndicationNone:
			resume := s.PauseSpinner()
			if event.HasIcon() {
				fmt.Fprintf(s.defaultLogStream.Stdout(), "%s %s\n", event.Icon, event.Content.String())
			} else {
				fmt.Fprintf(s.defaultLogStream.Stdout(), "%s\n", event.Content.String())
			}
			resume()
		}

		s.eventsWg.Done()
	}
	s.printLoopWg.Done()
}
