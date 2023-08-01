package cliuimodel

import (
	"container/list"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/events"
)

// EventMsg defines a message for events.
type EventMsg struct {
	events.Event

	Start    time.Time
	Duration time.Duration
}

// NewStatusEvents returns a new events model.
func NewStatusEvents(bus events.Provider, maxHistory int) StatusEvents {
	return StatusEvents{
		events:     list.New(),
		spinner:    NewSpinner(),
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
	bus        events.Provider
}

func (m *StatusEvents) ClearEvents() {
	m.static = nil
	m.events.Init()
}

func (m StatusEvents) Wait() tea.Cmd {
	return tea.Batch(spinner.Tick, m.WaitEvent)
}

func (m StatusEvents) WaitEvent() tea.Msg {
	e := <-m.bus.Events()

	return EventMsg{
		Start: time.Now(),
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
				evt.Duration = time.Since(evt.Start)
				e.Value = evt
			}

			// Add the event to the queue
			m.events.PushFront(msg)

			// Only show a reduced history of events
			if m.events.Len() > m.maxHistory {
				m.events.Remove(m.events.Back())
			}
		} else {
			// Events that have no progress status are considered static,
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

	// Display static events first
	for _, evt := range m.static {
		view.WriteString(evt.String())

		if !strings.HasSuffix(evt.Message, "\n") {
			view.WriteRune('\n')
		}
	}

	// Make sure there is a line between the static and status events
	if m.static != nil && m.events.Len() > 0 {
		view.WriteRune('\n')
	}

	// Display status events
	if m.events.Len() > 0 {
		for e := m.events.Front(); e != nil; e = e.Next() {
			evt := e.Value.(EventMsg)

			// The first event is displayed using a spinner
			if e.Prev() == nil {
				fmt.Fprintf(&view, "%s%s\n", m.spinner.View(), evt)

				if e.Next() != nil {
					view.WriteRune('\n')
				}

				continue
			}

			// Display finished status event
			d := evt.Duration.Round(time.Second)
			s := strings.TrimSuffix(evt.String(), "...")

			fmt.Fprintf(&view, "%s %s %s\n", icons.OK, s, colors.Faint(d.String()))
		}
	}

	return view.String()
}

// NewEvents returns a new events model.
func NewEvents(bus events.Provider) Events {
	return Events{
		events:  list.New(),
		bus:     bus,
		spinner: NewSpinner(),
	}
}

// Events defines a model for events.
// The model renders a view that prints all received events one after
// the other. Status events are displayed with a spinner and removed
// from the list once they finish.
type Events struct {
	events  *list.List
	bus     events.Provider
	spinner spinner.Model
}

func (m *Events) ClearEvents() {
	m.events.Init()
}

func (m Events) Wait() tea.Cmd {
	// Check if the last added event is a status event
	// and if so make sure that the spinner is updated.
	if e := m.events.Back(); e != nil {
		if evt := e.Value.(events.Event); evt.InProgress() {
			return tea.Batch(spinner.Tick, m.WaitEvent)
		}
	}

	// By default, just wait until the next event is received
	return m.WaitEvent
}

func (m Events) WaitEvent() tea.Msg {
	e := <-m.bus.Events()

	return EventMsg{
		Event: e,
		Start: time.Now(),
	}
}

func (m Events) Update(msg tea.Msg) (Events, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case EventMsg:
		// Remove the last event if is a status one.
		// Status events must always be the last event in the list so the
		// spinner is displayed at the bottom and not in between events.
		// They are removed when another status event is received.
		if e := m.events.Back(); e != nil {
			if evt := e.Value.(events.Event); evt.InProgress() {
				m.events.Remove(e)
			}
		}

		// Append event at the end of the list
		m.events.PushBack(msg.Event)

		// Return a command to wait for the next event
		cmd = m.Wait()
	default:
		// Update the spinner state and get a new tick command
		m.spinner, cmd = m.spinner.Update(msg)
	}

	return m, cmd
}

func (m Events) View() string {
	var (
		view  strings.Builder
		group string
	)

	// Display the list of events
	for e := m.events.Front(); e != nil; e = e.Next() {
		evt := e.Value.(events.Event)

		// Add an empty line when the event group changes but omit it
		// for the first event to avoid adding an initial empty line.
		if group != evt.Group && e.Prev() != nil {
			// Update the group being displayed
			group = evt.Group

			view.WriteRune('\n')
		}

		if e.Next() == nil && evt.InProgress() {
			// When the event is the last one and is a status event display a spinner...
			fmt.Fprintf(&view, "\n%s%s", m.spinner.View(), evt)
		} else {
			// Otherwise display the event without the spinner
			view.WriteString(evt.String())
		}

		// Make sure that events have an EOL, so they are displayed right below each other
		if !strings.HasSuffix(evt.Message, "\n") {
			view.WriteRune('\n')
		}
	}

	return view.String()
}
