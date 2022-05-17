// Package events provides functionalities for packages to log their states as events
// for others to consume and display to end users in meaningful ways.
package events

import (
	"fmt"
	"sync"
)

type (
	ProgressIndication uint8

	// Event represents a state.
	Event struct {
		ProgressIndication ProgressIndication
		Icon               string
		Content            interface{}
		Verbose            bool
		Tag                string
		kvstore            map[string]string
	}
	// Option event options
	Option func(*Event)
)

const (
	IndicationNone ProgressIndication = iota
	IndicationStart
	IndicationFinish
)

func ProgressStarted() Option {
	return func(e *Event) {
		e.ProgressIndication = IndicationStart
	}
}

func ProgressFinished() Option {
	return func(e *Event) {
		e.ProgressIndication = IndicationFinish
	}
}

func Verbose() Option {
	return func(e *Event) {
		e.Verbose = true
	}
}

func WithKV(key, value string) Option {
	return func(e *Event) {
		e.kvstore[key] = value
	}
}

// Icon sets the text icon prefix.
func Icon(icon string) Option {
	return func(e *Event) {
		e.Icon = icon
	}
}

func WithTag(tag string) Option {
	return func(e *Event) {
		e.Tag = tag
	}
}

// New creates a new event with given config.
func New(content interface{}, options ...Option) Event {
	ev := Event{
		Content: content,
		kvstore: make(map[string]string),
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

// Text returns the text state of event.
func (e Event) String() string {
	return fmt.Sprintf("%s", e.Content)
}

// GetValue returns a value from underlying event kvstore
func (e Event) GetValue(key string) (string, bool) {
	value, ok := e.kvstore[key]
	return value, ok
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
func (b Bus) Send(content any, options ...Option) {
	if b.evchan == nil {
		return
	}
	if b.buswg != nil {
		b.buswg.Add(1)
	}
	b.evchan <- New(content, options...)
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
