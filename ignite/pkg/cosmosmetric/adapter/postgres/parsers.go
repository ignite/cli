package postgres

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ignite-hq/cli/ignite/pkg/cosmosmetric/query"
)

const (
	entityTX        = "tx"
	entityEventAttr = "attribute"

	fieldCreatedAt      = "created_at"
	fieldEventAttrName  = "name"
	fieldEventAttrValue = "value"
	fieldEventIndex     = "event_index"
	fieldEventTXHash    = "tx_hash"
	fieldEventType      = "event_type"
	fieldTXBlockHeight  = "height"
	fieldTXBlockTime    = "block_time"
	fieldTXHash         = "hash"
	fieldTXIndex        = "index"

	sqlSelectAll = "SELECT *"
)

var (
	ErrUnknownEntity    = errors.New("unknown query entity")
	ErrInvalidSortOrder = errors.New("invalid query sort order")
)

var (
	entityMap = map[query.Entity]string{
		query.EntityTX:        entityTX,
		query.EntityEventAttr: entityEventAttr,
	}

	fieldMap = map[query.Field]string{
		query.FieldTXHash:          fieldTXHash,
		query.FieldTXIndex:         fieldTXIndex,
		query.FieldTXBlockHeight:   fieldTXBlockHeight,
		query.FieldTXBlockTime:     fieldTXBlockTime,
		query.FieldTXCreateTime:    fieldCreatedAt,
		query.FieldEventTXHash:     fieldEventTXHash,
		query.FieldEventType:       fieldEventType,
		query.FieldEventIndex:      fieldEventIndex,
		query.FieldEventAttrName:   fieldEventAttrName,
		query.FieldEventAttrValue:  fieldEventAttrValue,
		query.FieldEventCreateTime: fieldCreatedAt,
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
	sections := []string{
		// Add SELECT
		parseFields(q.GetFields()),
		// Add FROM
		parseCall(q.GetCall()),
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

func parseFields(fields []query.Field) string {
	if len(fields) == 0 {
		// By default select all fields
		return sqlSelectAll
	}

	var names []string

	for _, f := range fields {
		if n := fieldMap[f]; n != "" {
			names = append(names, n)
		}
	}

	return fmt.Sprintf("SELECT %s", strings.Join(names, ", "))
}

func parseCall(c query.Call) string {
	var params []string

	// Init the function call placeholders for the arguments
	for i := range c.Args {
		params = append(params, fmt.Sprintf("$%d", i+1))
	}

	// When there are arguments it means it is a postgres function
	// call otherwise the call is treated as a view
	s := fmt.Sprintf("FROM %s", c.Name)
	if len(params) > 0 {
		s = fmt.Sprintf("%s(%s)", s, strings.Join(params, ", "))
	}

	return s
}

func parseEntity(e query.Entity) (string, error) {
	if name := entityMap[e]; name != "" {
		return fmt.Sprintf("FROM %s", name), nil
	}

	return "", ErrUnknownEntity
}

func parseFilters(filters []query.Filter) string {
	if len(filters) == 0 {
		return ""
	}

	var (
		items []string
		pos   uint
	)

	for _, f := range filters {
		// Render the filter so it can be applied to the query
		expr := f.String()

		// When the filter has a value replace the "?" by a positional
		// postgres placeholder like "$1", "$2", and so on
		if v := f.GetValue(); v != nil {
			index := strings.LastIndex(expr, filterPlaceholder)
			expr = expr[:index] + fmt.Sprintf("$%d", pos+1) + expr[index+1:]
			pos += 1
		}

		items = append(items, expr)
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
