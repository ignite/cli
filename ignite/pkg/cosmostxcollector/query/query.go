package query

import (
	"fmt"

	"github.com/ignite-hq/cli/ignite/pkg/cosmostxcollector/query/call"
)

const (
	DefaultPageSize = 30
)

const (
	SortOrderAsc  = "asc"
	SortOrderDesc = "desc"
)

// Entity defines a data backend entity.
type Entity uint

// Field defines an entity field.
type Field uint

// SortBy contains info on how to sort query results.
type SortBy struct {
	Field Field
	Order string
}

// Filter describes a filter to apply to a query.
type Filter interface {
	fmt.Stringer

	// GetField returns the name of the filtered field.
	GetField() string

	// GetValue returns the value to use for filtering.
	GetValue() any
}

// New creates a new query.
func New(e Entity, f ...Field) Query {
	return Query{
		entity:   e,
		fields:   f,
		pageSize: DefaultPageSize,
		atPage:   1,
	}
}

// NewCall creates a new query that selects results from a view or function.
func NewCall(c call.Call) Query {
	return Query{
		call:     c,
		pageSize: DefaultPageSize,
		atPage:   1,
	}
}

// Query describes a how to select values from a data backend.
type Query struct {
	entity   Entity
	fields   []Field
	sortBy   []SortBy
	pageSize uint32
	atPage   uint32
	call     call.Call
	filters  []Filter
}

// GetEntity returns the name of the data entity to select.
func (q Query) GetEntity() Entity {
	return q.entity
}

// GetFields returns list of data entity fields to select.
func (q Query) GetFields() []Field {
	return q.fields
}

// GetSortBy returns the sort info for the query.
func (q Query) GetSortBy() []SortBy {
	return q.sortBy
}

// GetPageSize returns the size for each query result set.
func (q Query) GetPageSize() uint32 {
	return q.pageSize
}

// GetAtPage returns the result set page to query.
func (q Query) GetAtPage() uint32 {
	return q.atPage
}

// GetCall returns the function or view to query.
func (q Query) GetCall() call.Call {
	return q.call
}

// GetFilters returns the list of filters to apply to the query.
func (q Query) GetFilters() []Filter {
	return q.filters
}

// IsPagingEnabled checks if the query results should be paginated.
func (q Query) IsPagingEnabled() bool {
	return q.pageSize > 0
}

// IsCall checks if the query is a call to a function or view.
func (q Query) IsCall() bool {
	return q.call.Name() != ""
}

// AtPage assigns a page to select.
func (q Query) AtPage(page uint32) Query {
	q.atPage = page

	return q
}

// WithPageSize assigns the number of results to select per page.
func (q Query) WithPageSize(size uint32) Query {
	q.pageSize = size

	return q
}

// WithoutPaging disables the paging of results.
// All results are selected when paging is disabled.
func (q Query) WithoutPaging() Query {
	q.pageSize = 0

	return q
}

// AppendSortBy appends ordering information for one or more fields.
func (q Query) AppendSortBy(order string, fields ...Field) Query {
	for _, f := range fields {
		q.sortBy = append(q.sortBy, SortBy{
			Field: f,
			Order: order,
		})
	}

	return q
}

// AppendFilters appends one or more filters to apply to the query.
func (q Query) AppendFilters(f ...Filter) Query {
	q.filters = append(q.filters, f...)

	return q
}
