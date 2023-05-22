package datatype_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/templates/field/datatype"
)

func TestIsSupportedType(t *testing.T) {
	tests := []struct {
		name     string
		typename datatype.Name
		ok       bool
	}{
		{
			name:     "string",
			typename: datatype.String,
			ok:       true,
		},
		{
			name:     "string slice",
			typename: datatype.StringSlice,
			ok:       true,
		},
		{
			name:     "bool",
			typename: datatype.Bool,
			ok:       true,
		},
		{
			name:     "int",
			typename: datatype.Int,
			ok:       true,
		},
		{
			name:     "int slice",
			typename: datatype.IntSlice,
			ok:       true,
		},
		{
			name:     "uint",
			typename: datatype.Uint,
			ok:       true,
		},
		{
			name:     "uint slice",
			typename: datatype.UintSlice,
			ok:       true,
		},
		{
			name:     "coin",
			typename: datatype.Coin,
			ok:       true,
		},
		{
			name:     "coin slice",
			typename: datatype.Coins,
			ok:       true,
		},
		{
			name:     "custom",
			typename: datatype.Custom,
			ok:       true,
		},
		{
			name:     "string slice alias",
			typename: datatype.StringSliceAlias,
			ok:       true,
		},
		{
			name:     "int slice alias",
			typename: datatype.IntSliceAlias,
			ok:       true,
		},
		{
			name:     "uint slice alias",
			typename: datatype.UintSliceAlias,
			ok:       true,
		},
		{
			name:     "coin slice alias",
			typename: datatype.CoinSliceAlias,
			ok:       true,
		},
		{
			name:     "invalid type name",
			typename: datatype.Name("invalid"),
			ok:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, ok := datatype.IsSupportedType(tc.typename)
			require.Equal(t, tc.ok, ok)
		})
	}
}
