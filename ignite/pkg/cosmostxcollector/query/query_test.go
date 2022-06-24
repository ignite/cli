package query_test

import (
	"testing"

	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/query"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/query/call"
	"github.com/stretchr/testify/require"
)

var (
	entity query.Entity = 1
	field  query.Field  = 1
)

func TestQuery(t *testing.T) {
	qry := query.New(entity, field)

	// Assert entity query
	require.Equal(t, entity, qry.GetEntity())
	require.Equal(t, []query.Field{field}, qry.GetFields())
	require.False(t, qry.IsCall())

	// Assert defaults
	require.Nil(t, qry.GetSortBy())
	require.Nil(t, qry.GetFilters())
	require.True(t, qry.IsPagingEnabled())
	require.EqualValues(t, query.DefaultPageSize, qry.GetPageSize())
	require.EqualValues(t, 1, qry.GetAtPage())
}

func TestCallQuery(t *testing.T) {
	c := call.New("test")
	qry := query.NewCall(c)

	// Assert call query
	require.True(t, qry.IsCall())
	require.True(t, qry.GetCall().Name() == c.Name())

	// Assert defaults
	require.Nil(t, qry.GetSortBy())
	require.Nil(t, qry.GetFilters())
	require.True(t, qry.IsPagingEnabled())
	require.EqualValues(t, query.DefaultPageSize, qry.GetPageSize())
	require.EqualValues(t, 1, qry.GetAtPage())
}

func TestPaging(t *testing.T) {
	var (
		page     uint32 = 42
		pageSize uint32 = 100
	)

	qry := query.
		New(entity, field).
		WithPageSize(pageSize).
		AtPage(page)

	// Assert
	require.True(t, qry.IsPagingEnabled())
	require.EqualValues(t, pageSize, qry.GetPageSize())
	require.EqualValues(t, page, qry.GetAtPage())
}

func TestDisablePaging(t *testing.T) {
	qry := query.
		New(entity, field).
		WithoutPaging()

	// Assert
	require.False(t, qry.IsPagingEnabled())
	require.EqualValues(t, 0, qry.GetPageSize())
}

func TestAtPageZero(t *testing.T) {
	qry := query.
		New(entity, field).
		AtPage(0)

	// Assert
	require.True(t, qry.IsPagingEnabled())
	require.EqualValues(t, 1, qry.GetAtPage())
}
