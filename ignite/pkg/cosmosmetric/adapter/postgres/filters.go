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

func NewFilter(field, value string) Filter {
	return Filter{field, value}
}

type Filter struct {
	field, value string
}

func (f Filter) String() string {
	return fmt.Sprintf("%s = %s", f.field, filterPlaceholder)
}

func (f Filter) GetField() string {
	return f.field
}

func (f Filter) GetValue() any {
	return f.value
}

func NewTextFieldFilter(f query.Field, value string) Filter {
	field := fieldMap[f]

	return Filter{field, value}
}

func NewTextEventAttrValueFilter(value string) TextEventAttrValueFilter {
	return TextEventAttrValueFilter{
		// The string value must be quoted to match with the JSONB text
		NewTextFieldFilter(adapter.FieldEventAttrValue, strconv.Quote(value)),
	}
}

type TextEventAttrValueFilter struct {
	Filter
}

func (f TextEventAttrValueFilter) String() string {
	return fmt.Sprintf("%s::text = %s", f.GetField(), filterPlaceholder)
}

func FilterByEventType(eventType string) Filter {
	return NewTextFieldFilter(adapter.FieldEventType, eventType)
}

func FilterByEventAttrName(name string) Filter {
	return NewTextFieldFilter(adapter.FieldEventAttrName, name)
}

func FilterByEventAttrValue(v string) TextEventAttrValueFilter {
	return NewTextEventAttrValueFilter(v)
}
