// Package events provides functionalities for packages to log their states as events
// for others to consume and display to end users in meaningful ways.
package events

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
		Content            Content
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
func New(content Content, options ...Option) Event {
	ev := Event{
		Content: content,
	}
	for _, applyOption := range options {
		applyOption(&ev)
	}
	return ev
}

// IsOngoing checks if state change that triggered this event is still ongoing.
func (e Event) IsOngoing() bool {
	return e.ProgressIndication == IndicationStart
}

// HasIcon checks if event contains an icon
func (e Event) HasIcon() bool {
	return e.Icon != ""
}
