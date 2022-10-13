package query

import (
	"encoding/json"
	"time"
)

// Event defines a transaction event.
type Event struct {
	ID         int64
	TXHash     string
	Index      uint64
	Type       string
	Attributes []Attribute
	CreatedAt  time.Time
}

// NewAttribute creates a new transaction event attribute.
func NewAttribute(name string, value []byte) Attribute {
	return Attribute{
		value: value,
		Name:  name,
	}
}

// Attribute defines a transaction event attribute.
type Attribute struct {
	value []byte

	Name string
}

// Value returns the attribute value.
// Event attribute values are originally encoded as JSON.
// This method decodes the event value into its Go representation.
func (a Attribute) Value() (v any, err error) {
	if a.value == nil {
		return
	}

	err = json.Unmarshal(a.value, &v)
	return
}

// NewEventQuery creates a new query that selects events.
func NewEventQuery(options ...Option) EventQuery {
	return New("event", options...)
}

// EventQuery describes how to select event values from a data backend.
type EventQuery interface {
	Pager
	Filterer
}
