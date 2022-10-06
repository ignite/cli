package query

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID         int64
	TXHash     string
	Index      uint64
	Type       string
	Attributes []Attribute
	CreatedAt  time.Time
}

func NewAttribute(name string, value []byte) Attribute {
	return Attribute{
		value: value,
		Name:  name,
	}
}

type Attribute struct {
	value []byte

	Name string
}

func (a Attribute) Value() (v any, err error) {
	if a.value == nil {
		return
	}

	err = json.Unmarshal(a.value, &v)
	return
}

func NewEventQuery() EventQuery {
	return EventQuery{
		pageSize: DefaultPageSize,
		atPage:   1,
	}
}

type EventQuery struct {
	filters  []Filter
	pageSize uint32
	atPage   uint32
}

// GetFilters returns the list of filters to apply to the query.
func (q EventQuery) GetFilters() []Filter {
	return q.filters
}

// GetPageSize returns the size for each query result set.
func (q EventQuery) GetPageSize() uint32 {
	return q.pageSize
}

// GetAtPage returns the result set page to query.
func (q EventQuery) GetAtPage() uint32 {
	return q.atPage
}

// IsPagingEnabled checks if the query results should be paginated.
func (q EventQuery) IsPagingEnabled() bool {
	return q.pageSize > 0
}

// AtPage assigns a page to select.
// Pages start from page one, so assigning page zero selects the first page.
func (q EventQuery) AtPage(page uint32) EventQuery {
	if page == 0 {
		q.atPage = 1
	} else {
		q.atPage = page
	}

	return q
}

// WithPageSize assigns the number of results to select per page.
// The default page size is used when size zero is assigned.
func (q EventQuery) WithPageSize(size uint32) EventQuery {
	if size == 0 {
		q.pageSize = DefaultPageSize
	} else {
		q.pageSize = size
	}

	return q
}

// WithoutPaging disables the paging of results.
// All results are selected when paging is disabled.
func (q EventQuery) WithoutPaging() EventQuery {
	q.pageSize = 0

	return q
}

func (q EventQuery) AppendFilters(f ...Filter) EventQuery {
	q.filters = append(q.filters, f...)

	return q
}
