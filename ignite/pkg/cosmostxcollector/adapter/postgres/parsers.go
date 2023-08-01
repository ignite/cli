package postgres

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/query"
)

const (
	eventAttrPrefix = "attribute."

	sqlSelectAll = "SELECT *"
	sqlWhereTrue = "WHERE true"

	tplSelectEventsSQL = `
		SELECT event.id, event.index, event.tx_hash, event.type, event.created_at
		FROM event INNER JOIN tx ON event.tx_hash = tx.hash
		%s
		ORDER BY tx.height, tx.index, event.index
	`
	tplSelectEventsWithAttrSQL = `
		SELECT DISTINCT events.*
		FROM (
			SELECT event.id, event.index, event.tx_hash, event.type, event.created_at
			FROM event
				INNER JOIN tx ON event.tx_hash = tx.hash
				INNER JOIN attribute ON event.id = attribute.event_id
			%s
			ORDER BY tx.height, tx.index, event.index
		) AS events
	`
)

var (
	ErrUnknownEntity    = errors.New("unknown query entity")
	ErrInvalidSortOrder = errors.New("invalid query sort order")
)

// TODO: Use an SQL builder/parser to build the queries?
func parseQuery(q query.Query) (string, error) {
	sections := []string{
		// Add SELECT
		parseFields(q.Fields()),
		// Add FROM
		parseFrom(q),
	}

	// Add WHERE
	sections = append(sections, parseFilters(q.Filters()))

	// Add ORDER BY
	sortBy, err := parseSortBy(q.SortBy())
	if err != nil {
		return "", err
	}

	if sortBy != "" {
		sections = append(sections, sortBy)
	}

	// Add LIMIT/OFFSET
	if s, ok := parsePaging(q); ok {
		sections = append(sections, s)
	}

	return strings.Join(sections, " "), nil
}

func parseEventQuery(q query.EventQuery) string {
	sql := tplSelectEventsSQL
	filters := q.Filters()

	// Check if any of the filters references an event attribute
	// and if so add the required INNER JOIN to the raw SQL query.
	// The JOIN is not present by default to improve events queries.
	for _, f := range filters {
		if strings.HasPrefix(f.Field(), eventAttrPrefix) {
			sql = tplSelectEventsWithAttrSQL

			break
		}
	}

	// Add SELECT
	sections := []string{
		fmt.Sprintf(sql, parseFilters(q.Filters())),
	}

	// Add LIMIT/OFFSET
	if s, ok := parsePaging(q); ok {
		sections = append(sections, s)
	}

	return strings.Join(sections, " ")
}

func parseFields(fields []string) string {
	if len(fields) == 0 {
		// By default select all fields
		return sqlSelectAll
	}

	return fmt.Sprintf("SELECT DISTINCT %s", strings.Join(fields, ", "))
}

func parseFrom(q query.Query) string {
	// Init the function call placeholders for the arguments
	args := q.Args()
	placeholders := make([]string, len(args))
	for i := range args {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	// When there are arguments it means it is a postgres function
	// call otherwise the call is treated as a table or view.
	s := fmt.Sprintf("FROM %s", q.Name())
	if len(placeholders) > 0 {
		s = fmt.Sprintf("%s(%s)", s, strings.Join(placeholders, ", "))
	}

	return s
}

func parseFilters(filters []query.Filter) string {
	if len(filters) == 0 {
		return sqlWhereTrue
	}

	pos := 0
	items := make([]string, len(filters))

	for i, f := range filters {
		// Render the filter, so it can be applied to the query
		expr := f.String()

		// When the filter has a value replace the "?" by a positional
		// postgres placeholder like "$1", "$2", and so on
		if v := f.Value(); v != nil {
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

		items = append(items, fmt.Sprintf("%s %s", s.Field, s.Order))
	}

	return fmt.Sprintf("ORDER BY %s", strings.Join(items, ", ")), nil
}

func parsePaging(q query.Pager) (string, bool) {
	if !q.IsPagingEnabled() {
		return "", false
	}

	// Get the current page and make sure that the page number is valid
	page := q.AtPage()
	if page == 0 {
		page = 1
	}

	limit := q.PageSize()
	offset := limit * (page - 1)

	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset), true
}
