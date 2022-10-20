package query

import "fmt"

const (
	// DefaultPageSize defines the default number of results to select per page.
	DefaultPageSize = 30

	SortOrderAsc  = "asc"
	SortOrderDesc = "desc"
)

// Pager describes support for paging query results.
type Pager interface {
	// PageSize returns the size for each query result set.
	PageSize() uint32

	// AtPage returns the result set page to query.
	AtPage() uint32

	// IsPagingEnabled checks if the query results should be paginated.
	IsPagingEnabled() bool
}

// SortBy contains info on how to sort query results.
type SortBy struct {
	Field string
	Order string
}

// Filter describes a filter to apply to a query.
type Filter interface {
	fmt.Stringer

	// Field returns the name of the filtered field.
	Field() string

	// Value returns the value to use for filtering.
	Value() any
}

// Filterer describes support for filtering query results.
type Filterer interface {
	// Filters returns the list of filters to apply to the query.
	Filters() []Filter
}

type queryOptions struct {
	pageSize uint32
	atPage   uint32
	args     []any
	fields   []string
	sortBy   []SortBy
	filters  []Filter
}

// Option configures queries.
type Option func(*Query)

// AtPage assigns a page to select.
// Pages start from page one, so assigning page zero selects the first page.
func AtPage(page uint32) Option {
	return func(q *Query) {
		if page == 0 {
			q.options.atPage = 1
		} else {
			q.options.atPage = page
		}
	}
}

// WithPageSize assigns the number of results to select per page.
// The default page size is used when size zero is assigned.
func WithPageSize(size uint32) Option {
	return func(q *Query) {
		if size == 0 {
			q.options.pageSize = DefaultPageSize
		} else {
			q.options.pageSize = size
		}
	}
}

// WithoutPaging disables the paging of results.
// All results are selected when paging is disabled.
func WithoutPaging() Option {
	return func(q *Query) {
		q.options.pageSize = 0
	}
}

// WithFilters adds one or more filters to apply to the query.
func WithFilters(f ...Filter) Option {
	return func(q *Query) {
		q.options.filters = f
	}
}

// WithArgs adds one or more arguments to the query.
func WithArgs(args ...any) Option {
	return func(q *Query) {
		q.options.args = args
	}
}

// Fields assigns the field names to query.
// The default is to select all fields.
func Fields(fields ...string) Option {
	return func(q *Query) {
		q.options.fields = fields
	}
}

// SortByFields orders the query by one or more fields.
// Use `WithSortBy` option when multiple order by directions are needed.
func SortByFields(order string, fields ...string) Option {
	return func(q *Query) {
		for _, f := range fields {
			q.options.sortBy = append(q.options.sortBy, SortBy{
				Field: f,
				Order: order,
			})
		}
	}
}

// WithSortBy orders the query by one or more fields.
func WithSortBy(o ...SortBy) Option {
	return func(q *Query) {
		q.options.sortBy = append(q.options.sortBy, o...)
	}
}

// New creates a new query that selects results from an entity.
// The name is the name of an entity which depending on the data backend have
// different meanings. In a relational database the name should be a table,
// function or view, while in a NoSQL database it should be a collection.
func New(name string, options ...Option) Query {
	q := Query{
		name: name,
		options: queryOptions{
			pageSize: DefaultPageSize,
			atPage:   1,
		},
	}

	for _, apply := range options {
		apply(&q)
	}

	return q
}

// Query describes how to select values from a data backend.
type Query struct {
	name    string
	options queryOptions
}

// Name returns the name of the database table, collection, view or function to select.
func (q Query) Name() string {
	return q.name
}

// Args returns the arguments for query.
// Arguments are used when the query calls a function in the data backend.
func (q Query) Args() []any {
	return q.options.args
}

// Fields returns list of field names to select.
func (q Query) Fields() []string {
	return q.options.fields
}

// SortBy returns the sort info for the query.
func (q Query) SortBy() []SortBy {
	return q.options.sortBy
}

// PageSize returns the size for each query result set.
func (q Query) PageSize() uint32 {
	return q.options.pageSize
}

// AtPage returns the result set page to query.
func (q Query) AtPage() uint32 {
	return q.options.atPage
}

// Filters returns the list of filters to apply to the query.
func (q Query) Filters() []Filter {
	return q.options.filters
}

// IsPagingEnabled checks if the query results should be paginated.
func (q Query) IsPagingEnabled() bool {
	return q.options.pageSize > 0
}
