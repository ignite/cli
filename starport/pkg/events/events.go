// Package events provides functionalities for packages to log their states as events
// for others to consume and display to end users in meaningful ways.
package events

import "fmt"

// Event represents a state.
type Event struct {
	// Status shows the current status of event.
	status Status

	// Description of the state.
	Description string
}

// Status shows if state is ongoing or completed.
type Status int

const (
	StatusOngoing Status = iota
	StatusDone
)

// New creates a new event with given config.
func New(status Status, description string) Event {
	return Event{status, description}
}

// IsOngoing checks if state change that triggered this event is still ongoing.
func (e Event) IsOngoing() bool {
	return e.status == StatusOngoing
}

// Text returns the text state of event.
func (e Event) Text() string {
	if e.IsOngoing() {
		return fmt.Sprintf("%s...", e.Description)
	}
	return e.Description
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
