// Package events provides functionalities for packages to log their states as events
// for others to consume and display to end users in meaningful ways.
package events

import (
	"fmt"

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
)

// TextColor sets the text color
func TextColor(c color.Color) Option {
	return func(e *Event) {
		e.TextColor = c
	}
}

// Icon sets the text icon prefix
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
type Bus chan Event

// NewBus creates a new event bus to send/receive events.
func NewBus() Bus {
	return make(Bus)
}

// Send sends a new event to bus.
func (b Bus) Send(e Event) {
	if b == nil {
		return
	}
	b <- e
}

// Shutdown shutdowns event bus.
func (b Bus) Shutdown() {
	if b == nil {
		return
	}
	close(b)
}
