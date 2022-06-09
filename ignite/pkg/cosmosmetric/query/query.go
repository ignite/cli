package query

import "fmt"

const (
	DefaultPageSize = 30

	SortOrderAsc  = "asc"
	SortOrderDesc = "desc"
)

const (
	EntityTX Entity = iota
	EntityEventAttr
)

const (
	FieldTXHash Field = iota
	FieldTXIndex
	FieldTXBlockHeight
	FieldTXBlockTime
	FieldTXCreateTime

	FieldEventTXHash
	FieldEventType
	FieldEventIndex
	FieldEventAttrName
	FieldEventAttrValue
	FieldEventCreateTime
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

// Call defines a data backend function or view call to get the query results.
type Call struct {
	Name string
	Args []any
}

// Filter describes a filter to apply to a query.
type Filter interface {
	fmt.Stringer

	// GetField returns the filtered field.
	GetField() Field

	// GetValue returns the value to use for filtering.
	GetValue() any
}

// New creates a new data backend query.
func New(e Entity) Query {
	return Query{
		entity:   e,
		pageSize: DefaultPageSize,
		atPage:   1,
	}
}

// NewCall creates a new data backend query that selects a function or view.
func NewCall(name string, args ...any) Query {
	return Query{
		call: Call{
			Name: name,
			Args: args,
		},
		pageSize: DefaultPageSize,
		atPage:   1,
	}
}

// Query describes a data backend query.
type Query struct {
	entity   Entity
	fields   []Field
	sortBy   []SortBy
	pageSize uint64
	atPage   uint64
	call     Call
	filters  []Filter
}

// GetEntity returns the name of the data entity to select.
func (q Query) GetEntity() Entity {
	return q.entity
}

// GetFields returns list of fields to select.
func (q Query) GetFields() []Field {
	return q.fields
}

// GetSortBy returns the sort info for the query.
func (q Query) GetSortBy() []SortBy {
	return q.sortBy
}

// GetPageSize returns the size for each query result set.
func (q Query) GetPageSize() uint64 {
	return q.pageSize
}

// GetAtPage returns the result set page to query.
func (q Query) GetAtPage() uint64 {
	return q.atPage
}

// GetCall returns the function or view to query.
func (q Query) GetCall() Call {
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
	return q.call.Name != ""
}

// AtPage assigns a page to select.
func (q Query) AtPage(page uint64) Query {
	q.atPage = page

	return q
}

// WithPageSize assigns the number of results to select per page.
func (q Query) WithPageSize(size uint64) Query {
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

// AppendFields appends one or more fields to select.
func (q Query) AppendFields(f ...Field) Query {
	q.fields = append(q.fields, f...)

	return q
}

// AppendFilters appends one or more filters to apply to the query.
func (q Query) AppendFilters(f ...Filter) Query {
	q.filters = append(q.filters, f...)

	return q
}
