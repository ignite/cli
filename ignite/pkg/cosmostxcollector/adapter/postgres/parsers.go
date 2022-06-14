package postgres

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ignite-hq/cli/ignite/pkg/cosmostxcollector/adapter"
	"github.com/ignite-hq/cli/ignite/pkg/cosmostxcollector/query"
	"github.com/ignite-hq/cli/ignite/pkg/cosmostxcollector/query/call"
)

const (
	entityEvent     = "event"
	entityEventAttr = "attribute"
)

const (
	fieldEventAttrName  = "attribute.name"
	fieldEventAttrValue = "attribute.value"
	fieldEventCreatedAt = "event.created_at"
	fieldEventID        = "event.id"
	fieldEventIndex     = "event.index"
	fieldEventTXHash    = "event.tx_hash"
	fieldEventType      = "event.type"
	fieldTXCreatedAt    = "tx.created_at"
	fieldTXBlockHeight  = "tx.height"
	fieldTXBlockTime    = "tx.block_time"
	fieldTXHash         = "tx.hash"
	fieldTXIndex        = "tx.index"
)

const (
	sqlSelectAll = "SELECT *"
	sqlFromTX    = "FROM tx"
)

const (
	tplFromEventSQL = "FROM %[1]v INNER JOIN %[2]v ON %[1]v.id = %[2]v.event_id"
)

var (
	ErrUnknownEntity    = errors.New("unknown query entity")
	ErrInvalidSortOrder = errors.New("invalid query sort order")
)

var (
	fieldMap = map[query.Field]string{
		adapter.FieldTXHash:          fieldTXHash,
		adapter.FieldTXIndex:         fieldTXIndex,
		adapter.FieldTXBlockHeight:   fieldTXBlockHeight,
		adapter.FieldTXBlockTime:     fieldTXBlockTime,
		adapter.FieldTXCreateTime:    fieldTXCreatedAt,
		adapter.FieldEventID:         fieldEventID,
		adapter.FieldEventTXHash:     fieldEventTXHash,
		adapter.FieldEventType:       fieldEventType,
		adapter.FieldEventIndex:      fieldEventIndex,
		adapter.FieldEventAttrName:   fieldEventAttrName,
		adapter.FieldEventAttrValue:  fieldEventAttrValue,
		adapter.FieldEventCreateTime: fieldEventCreatedAt,
	}
)

// TODO: use an SQL builder/parser to build the queries
func parseQuery(q query.Query) (string, error) {
	if q.IsCall() {
		return parseCallQuery(q)
	}

	return parseEntityQuery(q)
}

func parseCallQuery(q query.Query) (string, error) {
	call := q.GetCall()
	sections := []string{
		// Add SELECT
		parseCustomFields(call.Fields()),
		// Add FROM
		parseCall(call),
	}

	// Add WHERE
	if s := parseFilters(q.GetFilters()); s != "" {
		sections = append(sections, s)
	}

	// Add ORDER BY
	sortBy, err := parseSortBy(q.GetSortBy())
	if err != nil {
		return "", err
	}

	if sortBy != "" {
		sections = append(sections, sortBy)
	}

	// Add LIMIT/OFFSET
	if q.IsPagingEnabled() {
		sections = append(sections, parsePaging(q))
	}

	return strings.Join(sections, " "), nil
}

func parseEntityQuery(q query.Query) (string, error) {
	// TODO: entities can be inferred from the fields
	fromEntity, err := parseEntity(q.GetEntity())
	if err != nil {
		return "", err
	}

	sections := []string{
		// Add SELECT
		parseFields(q.GetFields()),
		// Add FROM
		fromEntity,
	}

	// Add WHERE
	if s := parseFilters(q.GetFilters()); s != "" {
		sections = append(sections, s)
	}

	// Add ORDER BY
	sortBy, err := parseSortBy(q.GetSortBy())
	if err != nil {
		return "", err
	}

	if sortBy != "" {
		sections = append(sections, sortBy)
	}

	// Add LIMIT/OFFSET
	if q.IsPagingEnabled() {
		sections = append(sections, parsePaging(q))
	}

	return strings.Join(sections, " "), nil
}

func parseCustomFields(fields []string) string {
	if len(fields) == 0 {
		// By default select all fields
		return sqlSelectAll
	}

	return fmt.Sprintf("SELECT DISTINCT %s", strings.Join(fields, ", "))
}

func parseFields(fields []query.Field) string {
	var names []string

	for _, f := range fields {
		if n := fieldMap[f]; n != "" {
			names = append(names, n)
		}
	}

	return parseCustomFields(names)
}

func parseCall(c call.Call) string {
	args := c.Args()
	params := make([]string, len(args))

	// Init the function call placeholders for the arguments
	for i := range args {
		params[i] = fmt.Sprintf("$%d", i+1)
	}

	// When there are arguments it means it is a postgres function
	// call otherwise the call is treated as a view
	s := fmt.Sprintf("FROM %s", c.Name())
	if len(params) > 0 {
		s = fmt.Sprintf("%s(%s)", s, strings.Join(params, ", "))
	}

	return s
}

func parseEntity(e query.Entity) (string, error) {
	switch e {
	case adapter.EntityTX:
		return sqlFromTX, nil
	case adapter.EntityEvent:
		return fmt.Sprintf(tplFromEventSQL, entityEvent, entityEventAttr), nil
	}

	return "", ErrUnknownEntity
}

func parseFilters(filters []query.Filter) string {
	if len(filters) == 0 {
		return ""
	}

	pos := 0
	items := make([]string, len(filters))

	for i, f := range filters {
		// Render the filter so it can be applied to the query
		expr := f.String()

		// When the filter has a value replace the "?" by a positional
		// postgres placeholder like "$1", "$2", and so on
		if v := f.GetValue(); v != nil {
			index := strings.LastIndex(expr, filterPlaceholder)
			expr = expr[:index] + fmt.Sprintf("$%d", pos+1) + expr[index+1:]
			pos++
		}

		items[i] = expr
	}

	return fmt.Sprintf("WHERE %s", strings.Join(items, " AND "))
}

func parseSortBy(sortInfo []query.SortBy) (string, error) {
	if len(sortInfo) == 0 {
		return "", nil
	}

	var items []string

	for _, s := range sortInfo {
		if s.Order != query.SortOrderAsc && s.Order != query.SortOrderDesc {
			return "", ErrInvalidSortOrder
		}

		if n := fieldMap[s.Field]; n != "" {
			items = append(items, fmt.Sprintf("%s %s", n, s.Order))
		}
	}

	orderBy := fmt.Sprintf("ORDER BY %s", strings.Join(items, ", "))

	return orderBy, nil
}

func parsePaging(q query.Query) string {
	// Get the current page and make sure that the page number is valid
	page := q.GetAtPage()
	if page == 0 {
		page = 1
	}

	limit := q.GetPageSize()
	offset := limit * (page - 1)

	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}
