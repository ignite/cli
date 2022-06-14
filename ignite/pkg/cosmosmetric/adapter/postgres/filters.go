package postgres

import (
	"fmt"
	"strconv"

	"github.com/ignite-hq/cli/ignite/pkg/cosmosmetric/adapter"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosmetric/query"
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

// NewFilter creates a new generic equality filter for a query field.
// The result is a generic equality filter that maps the query field
// to the database field name.
func NewFilterFromField(f query.Field, value any, options ...FilterOption) Filter {
	return NewFilter(fieldMap[f], value, options...)
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

func (f Filter) GetField() string {
	return f.field
}

func (f Filter) GetValue() any {
	return f.value
}

func (f Filter) applyModifiers(field string) string {
	// Apply all the field modifiers in order
	for _, m := range f.modifiers {
		field = m(field)
	}

	return field
}

// FilterByEventType creates a new filter to match events by type.
func FilterByEventType(eventType string) Filter {
	return NewFilterFromField(adapter.FieldEventType, eventType)
}

// FilterByEventAttrName creates a new filter to match events by attribute name.
func FilterByEventAttrName(name string) Filter {
	return NewFilterFromField(adapter.FieldEventAttrName, name)
}

// FilterByEventAttrValue creates a new filter to match events by attribute value.
func FilterByEventAttrValue(v string) Filter {
	// The string value must be quoted to match with the JSONB text
	v = strconv.Quote(v)

	// Use a field modifier to cast the event attribute value JSONB field to text
	return NewFilterFromField(adapter.FieldEventAttrValue, v, WithModifiers(CastJSONToText))
}
