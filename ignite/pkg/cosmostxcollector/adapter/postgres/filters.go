package postgres

import (
	"fmt"
	"strconv"

	"github.com/lib/pq"
)

const (
	FieldEventAttrName  = "attribute.name"
	FieldEventAttrValue = "attribute.value"
	FieldEventTXHash    = "event.tx_hash"
	FieldEventType      = "event.type"
)

const (
	filterPlaceholder = "?"
)

// Modifier defines a function that can be used to modify a field name or value.
type Modifier func(field string) string

// CastJSONToText modifier casts a JSON/JSONB field to text.
func CastJSONToText(f string) string {
	return fmt.Sprintf("%s::text", f)
}

// CastJSONToNumeric modifier casts a JSON/JSONB field to numeric.
func CastJSONToNumeric(f string) string {
	return fmt.Sprintf("%s::numeric", f)
}

// FilterOption defines an option for filters.
type FilterOption func(*Filter)

// WithModifiers assigns one or more field modifier functions to the filter.
// Field modifiers can be used to change the behavior of a filtered field.
func WithModifiers(m ...Modifier) FilterOption {
	return func(f *Filter) {
		f.modifiers = m
	}
}

// NewFilter creates a new generic equality filter.
func NewFilter(field string, value any, options ...FilterOption) Filter {
	f := Filter{
		field: field,
		value: value,
	}

	for _, o := range options {
		o(&f)
	}

	return f
}

// Filter defines a generic equality filter.
type Filter struct {
	field     string
	value     any
	modifiers []Modifier
}

func (f Filter) String() string {
	return fmt.Sprintf("%s = %s", f.applyModifiers(f.field), filterPlaceholder)
}

func (f Filter) Field() string {
	return f.field
}

func (f Filter) Value() any {
	return f.value
}

func (f Filter) applyModifiers(field string) string {
	// Apply all the field modifiers in order
	for _, m := range f.modifiers {
		field = m(field)
	}

	return field
}

// NewStringSliceFilter creates a new string slice equality filter.
func NewStringSliceFilter(field string, values []string, options ...FilterOption) SliceFilter {
	return SliceFilter{
		Filter: NewFilter(field, pq.Array(values)),
	}
}

// NewIntSliceFilter creates a new int64 slice equality filter.
func NewIntSliceFilter(field string, values []int64, options ...FilterOption) SliceFilter {
	return SliceFilter{
		Filter: NewFilter(field, pq.Array(values)),
	}
}

// SliceFilter defines a generic slice/array equality filter.
type SliceFilter struct {
	Filter
}

func (f SliceFilter) String() string {
	return fmt.Sprintf("%s = ANY(%s)", f.applyModifiers(f.field), filterPlaceholder)
}

func (f SliceFilter) Value() any {
	return f.Filter.Value()
}

// FilterByEventType creates a new filter to match events by type.
func FilterByEventType(eventType string) Filter {
	return NewFilter(FieldEventType, eventType)
}

// FilterByEventTXs creates a new filter to match events by TX hashes.
func FilterByEventTXs(hashes ...string) SliceFilter {
	return NewStringSliceFilter(FieldEventTXHash, hashes)
}

// FilterByEventAttrName creates a new filter to match events by attribute name.
func FilterByEventAttrName(name string) Filter {
	return NewFilter(FieldEventAttrName, name)
}

// FilterByEventAttrValue creates a new filter to match events by attribute value.
func FilterByEventAttrValue(v string) Filter {
	// The string value must be quoted to match with the JSONB text
	v = strconv.Quote(v)

	// Use a field modifier to cast the event attribute value JSONB field to text
	return NewFilter(FieldEventAttrValue, v, WithModifiers(CastJSONToText))
}

// FilterByEventAttrValueInt creates a new filter to match events by attribute value.
func FilterByEventAttrValueInt(v int64) Filter {
	// Use a field modifier to cast the event attribute value JSONB field to numeric
	return NewFilter(FieldEventAttrValue, v, WithModifiers(CastJSONToNumeric))
}
