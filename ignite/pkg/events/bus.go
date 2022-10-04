package events

import "fmt"

// DefaultBufferSize defines the default maximum number
// of events that the bus can cache before they are handled.
const DefaultBufferSize = 50

type (
	// Bus is a send/receive event bus.
	Bus struct {
		evChan chan Event
	}

	// BusOption is used to specify Bus parameters
	BusOption func(*Bus)
)

// WithBufferSize assigns the size of the buffer to use for buffering events.
func WithBufferSize(size int) BusOption {
	return func(bus *Bus) {
		bus.evChan = make(chan Event, size)
	}
}

// NewBus creates a new event bus to send/receive events.
func NewBus(options ...BusOption) Bus {
	bus := Bus{
		evChan: make(chan Event, DefaultBufferSize),
	}

	for _, apply := range options {
		apply(&bus)
	}

	return bus
}

// Send sends a new event to bus.
// This method will block if the event bus buffer is full.
// The default buffer size can be changed using the `WithBufferSize` option.
func (b Bus) Send(content string, options ...Option) {
	if b.evChan == nil {
		return
	}

	b.evChan <- New(content, options...)
}

// Sendf sends formatted Event to the bus.
func (b Bus) Sendf(options []Option, format string, args ...interface{}) {
	b.Send(fmt.Sprintf(format, args...), options...)
}

// SendView sends a new event for a view to the bus.
func (b Bus) SendView(s fmt.Stringer, options ...Option) {
	b.Send(s.String(), options...)
}

// SendError sends a new event for an error to the bus.
func (b Bus) SendError(err error, options ...Option) {
	b.Send(err.Error(), options...)
}

// Events returns a read only channel to read the events.
func (b *Bus) Events() <-chan Event {
	return b.evChan
}

// Shutdown shutdowns event bus.
func (b Bus) Shutdown() {
	if b.evChan == nil {
		return
	}

	close(b.evChan)
}
