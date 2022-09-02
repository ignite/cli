// Package events provides functionalities for packages to log their states as events
// for others to consume and display to end users in meaningful ways.
package events

import (
	"fmt"
	"sync"

	"github.com/gookit/color"
)

type (
	// Event represents a state.
	Event struct {
		// Description of the state.
		Description string

		// Status shows the current status of event.
		Status Status

		// TextColor of the text.
		TextColor color.Color

		// Icon of the text.
		Icon string
	}

	// Status shows if state is ongoing or completed.
	Status int

	// Option event options
	Option func(*Event)
)

const (
	StatusOngoing Status = iota
	StatusDone
	StatusNeutral
)

// TextColor sets the text color
func TextColor(c color.Color) Option {
	return func(e *Event) {
		e.TextColor = c
	}
}

// Icon sets the text icon prefix.
func Icon(icon string) Option {
	return func(e *Event) {
		e.Icon = icon
	}
}

// New creates a new event with given config.
func New(status Status, description string, options ...Option) Event {
	ev := Event{Status: status, Description: description}
	for _, applyOption := range options {
		applyOption(&ev)
	}
	return ev
}

// NewOngoing creates a new StatusOngoing event.
func NewOngoing(description string) Event {
	return New(StatusOngoing, description)
}

// NewNeutral creates a new StatusNeutral event.
func NewNeutral(description string) Event {
	return New(StatusNeutral, description)
}

// NewDone creates a new StatusDone event.
func NewDone(description, icon string) Event {
	return New(StatusDone, description, Icon(icon))
}

// IsOngoing checks if state change that triggered this event is still ongoing.
func (e Event) IsOngoing() bool {
	return e.Status == StatusOngoing
}

// Text returns the text state of event.
func (e Event) Text() string {
	text := e.Description
	if e.IsOngoing() {
		text = fmt.Sprintf("%s...", e.Description)
	}
	return e.TextColor.Render(text)
}

// Bus is a send/receive event bus.
type (
	Bus struct {
		evchan chan Event
		buswg  *sync.WaitGroup
	}

	BusOption func(*Bus)
)

// WithWaitGroup sets wait group which is blocked if events bus is not empty.
func WithWaitGroup(wg *sync.WaitGroup) BusOption {
	return func(bus *Bus) {
		bus.buswg = wg
	}
}

// WithCustomBufferSize configures buffer size of underlying bus channel
func WithCustomBufferSize(size int) BusOption {
	return func(bus *Bus) {
		bus.evchan = make(chan Event, size)
	}
}

// NewBus creates a new event bus to send/receive events.
func NewBus(options ...BusOption) Bus {
	bus := Bus{
		evchan: make(chan Event),
	}

	for _, apply := range options {
		apply(&bus)
	}

	return bus
}

// Send sends a new event to bus.
func (b Bus) Send(e Event) {
	if b.evchan == nil {
		return
	}
	if b.buswg != nil {
		b.buswg.Add(1)
	}
	b.evchan <- e
}

// Events returns go channel with Event accessible only for read.
func (b *Bus) Events() <-chan Event {
	return b.evchan
}

// Shutdown shutdowns event bus.
func (b Bus) Shutdown() {
	if b.evchan == nil {
		return
	}
	close(b.evchan)
}
