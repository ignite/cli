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
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/events"
)

// Session controls command line interaction with users.
type Session struct {
	ev       events.Bus
	eventsWg *sync.WaitGroup

	spinner *clispinner.Spinner

	in          io.Reader
	out         io.Writer
	printLoopWg *sync.WaitGroup
}

type Option func(s *Session)

// WithOutput sets output stream for a session.
func WithOutput(output io.Writer) Option {
	return func(s *Session) {
		s.out = output
	}
}

// WithInput sets input stream for a session.
func WithInput(input io.Reader) Option {
	return func(s *Session) {
		s.in = input
	}
}

// New creates new Session.
func New(options ...Option) Session {
	wg := &sync.WaitGroup{}
	session := Session{
		ev:          events.NewBus(events.WithWaitGroup(wg)),
		in:          os.Stdin,
		out:         os.Stdout,
		eventsWg:    wg,
		printLoopWg: &sync.WaitGroup{},
	}
	for _, apply := range options {
		apply(&session)
	}
	session.spinner = clispinner.New(clispinner.WithWriter(session.out))
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
	_, err := fmt.Fprintf(s.out, format, a...)
	return err
}

// Println prints arbitrary message with line break.
func (s Session) Println(messages ...interface{}) error {
	s.Wait()
	defer s.PauseSpinner()()
	_, err := fmt.Fprintln(s.out, messages...)
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
	_, err := fmt.Fprint(s.out, messages...)
	return err
}

// Ask asks questions in the terminal and collect answers.
func (s Session) Ask(questions ...cliquiz.Question) error {
	s.Wait()
	defer s.PauseSpinner()()
	return cliquiz.Ask(questions...)
}

// AskConfirm asks yes/no question in the terminal.
func (s Session) AskConfirm(message string) error {
	s.Wait()
	defer s.PauseSpinner()()
	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
	}
	_, err := prompt.Run()
	return err
}

// PrintTable prints table data.
func (s Session) PrintTable(header []string, entries ...[]string) error {
	s.Wait()
	defer s.PauseSpinner()()
	return entrywriter.MustWrite(s.out, header, entries...)
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
		switch event.Status {
		case events.StatusOngoing:
			s.StartSpinner(event.Text())

		case events.StatusDone:
			if event.Icon == "" {
				event.Icon = icons.OK
			}
			s.StopSpinner()
			fmt.Fprintf(s.out, "%s %s\n", event.Icon, event.Text())

		case events.StatusNeutral:
			resume := s.PauseSpinner()
			fmt.Fprintf(s.out, event.Text())
			resume()
		}

		s.eventsWg.Done()
	}
	s.printLoopWg.Done()
}
