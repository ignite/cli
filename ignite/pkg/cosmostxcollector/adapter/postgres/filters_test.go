package postgres_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
)

func TestFilter(t *testing.T) {
	// Arrange
	name := "string_field"
	value := "test"
	repr := fmt.Sprintf("%s = ?", name)

	// Act
	filter := postgres.NewFilter(name, value)

	// Assert
	require.Equal(t, repr, filter.String())
	require.Equal(t, name, filter.Field())
	require.Equal(t, value, filter.Value())
}

func TestFilterModifiers(t *testing.T) {
	cases := []struct {
		name     string
		modifier postgres.Modifier
		want     string
	}{
		{
			name:     "CastJSONToText",
			modifier: postgres.CastJSONToText,
			want:     "field::text = ?",
		},
		{
			name:     "CastJSONToNumeric",
			modifier: postgres.CastJSONToNumeric,
			want:     "field::numeric = ?",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			field := postgres.NewFilter("field", nil, postgres.WithModifiers(tt.modifier))

			require.EqualValues(t, tt.want, field.String())
		})
	}
}
