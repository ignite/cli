package events

import (
	"fmt"
	"sync"
)

type (
	// Bus is a send/receive event bus.
	Bus struct {
		evchan chan Event
		buswg  *sync.WaitGroup
	}

	// BusOption is used to specify Bus parameters
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
func (b Bus) Send(content Content, options ...Option) {
	if b.evchan == nil {
		return
	}
	if b.buswg != nil {
		b.buswg.Add(1)
	}
	b.evchan <- New(content, options...)
}

// SendString sends an Event with a string content to the bus
func (b Bus) SendString(content string, options ...Option) {
	b.Send(StringContent(content), options...)
}

// Sendf sends formatted Event to the bus
func (b Bus) Sendf(options []Option, format string, args ...interface{}) {
	content := fmt.Sprintf(format, args...)
	b.Send(StringContent(content), options...)
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
