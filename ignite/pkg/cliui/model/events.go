package model

import (
	"container/list"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/cliui/style"
	"github.com/ignite/cli/ignite/pkg/events"
)

// EventMsg defines a message for events.
type EventMsg struct {
	events.Event

	start    time.Time
	duration time.Duration
}

// NewStatusEvents returns a new events model.
func NewStatusEvents(bus events.Bus, maxHistory int) StatusEvents {
	// TODO: Using latest github.com/charmbracelet/bubbles is not possible because
	//       of https://github.com/charmbracelet/glow/issues/268, we have dependency
	//       conflicts with markdownviewer module.
	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.ForegroundColor = ColorSpinner

	return StatusEvents{
		events:     list.New(),
		spinner:    s,
		maxHistory: maxHistory,
		bus:        bus,
	}
}

// StatusEvents defines a model for status events.
// The model renders a view that can be divided in three sections.
// The first one displays the "static" events which are the ones
// that are not status events. The second section displays a spinner
// with the status event that is in progress, and the third one
// displays a list with the past status events.
type StatusEvents struct {
	static     []events.Event
	events     *list.List
	spinner    spinner.Model
	maxHistory int
	bus        events.Bus
}

func (m StatusEvents) Wait() tea.Cmd {
	return tea.Batch(spinner.Tick, m.WaitEvent)
}

func (m StatusEvents) WaitEvent() tea.Msg {
	e := <-m.bus.Events()

	return EventMsg{
		start: time.Now(),
		Event: e,
	}
}

func (m StatusEvents) Update(msg tea.Msg) (StatusEvents, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case EventMsg:
		if msg.InProgress() {
			// Save the duration of the current ongoing event before setting a new one
			if e := m.events.Front(); e != nil {
				evt := e.Value.(EventMsg)
				evt.duration = time.Since(evt.start)
				e.Value = evt
			}

			// Add the event to the queue
			m.events.PushFront(msg)

			// Only show a reduced history of events
			if m.events.Len() > m.maxHistory {
				m.events.Remove(m.events.Back())
			}
		} else {
			// Events that have no progress status are considered static
			// so they will be printed without the spinner and won't be
			// removed from the output until the view is removed.
			m.static = append(m.static, msg.Event)
		}

		// Return a command to wait for the next event
		cmd = m.Wait()
	default:
		// Update the spinner state and get a new tick command
		m.spinner, cmd = m.spinner.Update(msg)
	}

	return m, cmd
}

func (m StatusEvents) View() string {
	var view strings.Builder

	// Show static events
	for _, evt := range m.static {
		view.WriteString(evt.String())

		if !strings.HasSuffix(evt.Message, "\n") {
			view.WriteRune(EOL)
		}
	}

	if m.static != nil {
		view.WriteRune(EOL)
	}

	// Show status events
	if m.events.Len() > 0 {
		for e := m.events.Front(); e != nil; e = e.Next() {
			evt := e.Value.(EventMsg)

			// The first event is displayed using a spinner
			if e.Prev() == nil {
				fmt.Fprintf(&view, "%s%s\n", m.spinner.View(), evt)

				if e.Next() != nil {
					view.WriteRune(EOL)
				}

				continue
			}

			// Display finished status event
			d := evt.duration.Round(time.Second)
			s := strings.TrimSuffix(evt.String(), "...")

			fmt.Fprintf(&view, "%s %s %s\n", icons.OK, s, style.Faint.Render(d.String()))
		}

		view.WriteRune(EOL)
	}

	return view.String()
}
