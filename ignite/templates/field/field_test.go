package field

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/templates/field/datatype"
)

// TestField_IsSlice tests the IsSlice method of Field struct.
func TestField_IsSlice(t *testing.T) {
	testCases := []struct {
		name     string
		field    Field
		expected bool
	}{
		{
			name: "array type should be slice",
			field: Field{
				Name:         multiformatname.Name{},
				DatatypeName: datatype.IntSlice,
				Datatype:     string(datatype.IntSlice),
			},
			expected: true,
		},
		{
			name: "array type should be slice",
			field: Field{
				Name:         multiformatname.Name{},
				DatatypeName: datatype.Bytes,
				Datatype:     string(datatype.Bytes),
			},
			expected: true,
		},
		{
			name: "array type should be slice",
			field: Field{
				Name:         multiformatname.Name{},
				DatatypeName: datatype.CoinSliceAlias,
				Datatype:     string(datatype.CoinSliceAlias),
			},
			expected: true,
		},
		{
			name: "coin type should not be slice",
			field: Field{
				Name:         multiformatname.Name{},
				DatatypeName: datatype.Coin,
				Datatype:     string(datatype.Coin),
			},
			expected: false,
		},
		{
			name: "string type should not be slice",
			field: Field{
				Name:         multiformatname.Name{},
				DatatypeName: datatype.String,
				Datatype:     "",
			},
			expected: false,
		},
		{
			name: "int type should not be slice",
			field: Field{
				Name:         multiformatname.Name{},
				DatatypeName: datatype.Int,
				Datatype:     "",
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, tc.field.IsSlice())
		})
	}
}
