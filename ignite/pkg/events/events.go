// Package events provides functionalities for packages to log their states as events
// for others to consume and display to end users in meaningful ways.
package events

import (
	"fmt"

	"github.com/muesli/reflow/indent"
)

// ProgressIndication enumerates possible states of progress indication for an Event.
type ProgressIndication uint8

const (
	GroupError = "error"
)

const (
	IndicationNone ProgressIndication = iota
	IndicationStart
	IndicationUpdate
	IndicationFinish
)

type (
	// Event represents a state.
	Event struct {
		ProgressIndication ProgressIndication
		Icon               string
		Indent             uint
		Message            string
		Verbose            bool
		Group              string
	}

	// Option event options.
	Option func(*Event)
)

// ProgressStart indicates that a status event starts the progress indicator.
func ProgressStart() Option {
	return func(e *Event) {
		e.ProgressIndication = IndicationStart
	}
}

// ProgressUpdate indicates that a status event updated the current progress.
func ProgressUpdate() Option {
	return func(e *Event) {
		e.ProgressIndication = IndicationUpdate
	}
}

// ProgressFinish indicates that a status event finished the ongoing task.
func ProgressFinish() Option {
	return func(e *Event) {
		e.ProgressIndication = IndicationFinish
	}
}

// Verbose sets high verbosity for the Event.
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

// Indent sets the text indentation.
func Indent(indent uint) Option {
	return func(e *Event) {
		e.Indent = indent
	}
}

// Group sets a group name for the event.
func Group(name string) Option {
	return func(e *Event) {
		e.Group = name
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
	s := e.Message
	if e.Icon != "" {
		s = fmt.Sprintf("%s %s", e.Icon, s)
	}

	if e.Indent > 0 {
		s = indent.String(s, e.Indent)
	}

	return s
}

// InProgress returns true when the event is in progress.
func (e Event) InProgress() bool {
	return e.ProgressIndication == IndicationStart || e.ProgressIndication == IndicationUpdate
}
