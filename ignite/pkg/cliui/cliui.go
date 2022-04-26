package cliui

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"io"
	"os"
	"sync"

	"github.com/ignite-hq/cli/ignite/pkg/cliui/cliquiz"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/clispinner"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/entrywriter"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/icons"
	"github.com/ignite-hq/cli/ignite/pkg/events"
)

type Session struct {
	ev       *events.Bus
	eventsWg *sync.WaitGroup

	spinner          *clispinner.Spinner
	spinnerWasPaused bool

	in          io.Reader
	out         io.Writer
	printLoopWg *sync.WaitGroup
}

type Option func(s *Session)

func WithOutput(output io.Writer) Option {
	return func(s *Session) {
		s.out = output
	}
}

func WithInput(input io.Reader) Option {
	return func(s *Session) {
		s.in = input
	}
}

func New(options ...Option) Session {
	wg := &sync.WaitGroup{}
	session := Session{
		spinner:     clispinner.New(),
		ev:          events.NewBus(events.WithWaitGroup(wg)),
		in:          os.Stdin,
		out:         os.Stdout,
		eventsWg:    wg,
		printLoopWg: &sync.WaitGroup{},
	}
	for _, apply := range options {
		apply(&session)
	}
	session.printLoopWg.Add(1)
	go session.printLoop()
	return session
}

func (s Session) EventBus() *events.Bus {
	return s.ev
}

func (s Session) StartSpinner(text string) {
	s.spinner.SetText(text).Start()
}

func (s Session) StopSpinner() {
	s.spinnerWasPaused = false
	s.spinner.Stop()
}

func (s Session) PauseSpinner() {
	s.spinner.Stop()
	s.spinnerWasPaused = true
}

func (s Session) UnpauseSpinner() {
	if s.spinnerWasPaused {
		s.spinner.Start()
		s.spinnerWasPaused = false
	}
}

func (s Session) Printf(format string, a ...interface{}) error {
	s.eventsWg.Wait()
	s.PauseSpinner()
	defer s.UnpauseSpinner()
	_, err := fmt.Fprintf(s.out, format, a...)
	return err
}

func (s Session) Println(messages ...interface{}) error {
	s.eventsWg.Wait()
	s.PauseSpinner()
	defer s.UnpauseSpinner()
	_, err := fmt.Fprintln(s.out, messages...)
	return err
}

func (s Session) Ask(questions ...cliquiz.Question) error {
	s.eventsWg.Wait()
	s.PauseSpinner()
	defer s.UnpauseSpinner()
	if s.in != os.Stdin && s.out != os.Stdout {
		return errors.New("cannot use quiz with customized io")
	}
	return cliquiz.Ask(questions...)
}

func (s Session) AskConfirm(message string) error {
	s.eventsWg.Wait()
	s.PauseSpinner()
	defer s.UnpauseSpinner()
	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
	}
	_, err := prompt.Run()
	return err
}

func (s Session) Table(header []string, entries ...[]string) error {
	s.eventsWg.Wait()
	s.PauseSpinner()
	defer s.UnpauseSpinner()
	return entrywriter.MustWrite(s.out, header, entries...)
}

func (s Session) Wait() {
	s.eventsWg.Wait()
}

func (s Session) Cleanup() {
	s.StopSpinner()
	s.ev.Shutdown()
	s.printLoopWg.Wait()
}

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
			s.PauseSpinner()
			fmt.Fprintf(s.out, event.Text())
			s.UnpauseSpinner()
		}

		s.eventsWg.Done()
	}
	s.printLoopWg.Done()
}
