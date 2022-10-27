// Package events provides functionalities for packages to log their states as events
// for others to consume and display to end users in meaningful ways.
package events

import "fmt"

// ProgressIndication enumerates possible states of progress indication for an Event
type ProgressIndication uint8

const (
	IndicationNone ProgressIndication = iota
	IndicationStart
	IndicationFinish
)

type (
	// Event represents a state.
	Event struct {
		ProgressIndication ProgressIndication
		Icon               string
		Message            string
		Verbose            bool
	}

	// Option event options
	Option func(*Event)
)

// ProgressStarted sets ProgressIndication as started
func ProgressStarted() Option {
	return func(e *Event) {
		e.ProgressIndication = IndicationStart
	}
}

// ProgressFinished sets ProgressIndication as finished
func ProgressFinished() Option {
	return func(e *Event) {
		e.ProgressIndication = IndicationFinish
	}
}

// Verbose sets high verbosity for the Event
func Verbose() Option {
	return func(e *Event) {
		e.Verbose = true
	}
}

// Icon sets the text icon prefix.
func Icon(icon string) Option {
	return func(e *Event) {
		e.Icon = icon
	}
}

// New creates a new event with given config.
func New(message string, options ...Option) Event {
	ev := Event{Message: message}

	for _, applyOption := range options {
		applyOption(&ev)
	}

	return ev
}

func (e Event) String() string {
	if e.Icon != "" {
		return fmt.Sprintf("%s %s", e.Icon, e.Message)
	}

	return e.Message
}

// InProgress returns true when the event is in progress.
func (e Event) InProgress() bool {
	return e.ProgressIndication == IndicationStart
}
