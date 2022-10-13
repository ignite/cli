package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/query"
)

func TestQuery(t *testing.T) {
	// Arrange
	name := "entity"
	field := "foo"

	// Act
	qry := query.New(name, query.Fields(field))

	// Assert
	require.Equal(t, name, qry.Name())
	require.Equal(t, []string{field}, qry.Fields())
	require.Nil(t, qry.SortBy())
	require.Nil(t, qry.Filters())
	require.Nil(t, qry.Args())
	require.True(t, qry.IsPagingEnabled())
	require.EqualValues(t, query.DefaultPageSize, qry.PageSize())
	require.EqualValues(t, 1, qry.AtPage())
}

func TestPaging(t *testing.T) {
	// Arrange
	var (
		page     uint32 = 42
		pageSize uint32 = 100
	)

	// Act
	qry := query.New(
		"name",
		query.WithPageSize(pageSize),
		query.AtPage(page),
	)

	// Assert
	require.True(t, qry.IsPagingEnabled())
	require.EqualValues(t, pageSize, qry.PageSize())
	require.EqualValues(t, page, qry.AtPage())
}

func TestDisablePaging(t *testing.T) {
	// Act
	qry := query.New("name", query.WithoutPaging())

	// Assert
	require.False(t, qry.IsPagingEnabled())
	require.EqualValues(t, 0, qry.PageSize())
}

func TestAtPageZero(t *testing.T) {
	// Act
	qry := query.New("name", query.AtPage(0))

	// Assert
	require.True(t, qry.IsPagingEnabled())
	require.EqualValues(t, 1, qry.AtPage())
}
