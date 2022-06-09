package postgres

import (
	"fmt"
	"strconv"

	"github.com/ignite-hq/cli/ignite/pkg/cosmosmetric/query"
)

const (
	filterPlaceholder = "?"
)

func NewTextFilter(field query.Field, value string) TextFilter {
	return TextFilter{field, value}
}

type TextFilter struct {
	field query.Field
	value string
}

func (f TextFilter) String() string {
	return fmt.Sprintf("%s = %s", fieldMap[f.field], filterPlaceholder)
}

func (f TextFilter) GetField() query.Field {
	return f.field
}

func (f TextFilter) GetValue() any {
	return f.value
}

func NewTextEventAttrValueFilter(value string) TextEventAttrValueFilter {
	return TextEventAttrValueFilter{
		// The string value must be quoted to match with the JSONB text
		NewTextFilter(query.FieldEventAttrValue, strconv.Quote(value)),
	}
}

type TextEventAttrValueFilter struct {
	TextFilter
}

func (f TextEventAttrValueFilter) String() string {
	return fmt.Sprintf("%s::text = %s", fieldMap[f.field], filterPlaceholder)
}

func FilterByEventType(eventType string) TextFilter {
	return NewTextFilter(query.FieldEventType, eventType)
}

func FilterByEventAttrName(name string) TextFilter {
	return NewTextFilter(query.FieldEventAttrName, name)
}

func FilterByEventAttrValue(v string) TextEventAttrValueFilter {
	return NewTextEventAttrValueFilter(v)
}
