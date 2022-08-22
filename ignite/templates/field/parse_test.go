package field

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/multiformatname"
	"github.com/ignite/cli/ignite/templates/field/datatype"
)

var (
	noCheck       = func(string) error { return nil }
	alwaysInvalid = func(string) error { return errors.New("invalid Name") }
)

func TestForbiddenParseFields(t *testing.T) {
	// check doesn't pass
	_, err := ParseFields([]string{"foo"}, alwaysInvalid)
	require.Error(t, err)

	// duplicated field
	_, err = ParseFields([]string{"foo", "foo:int"}, noCheck)
	require.Error(t, err)

	// invalid type
	_, err = ParseFields([]string{"foo:invalid"}, alwaysInvalid)
	require.Error(t, err)

	// invalid field Name
	_, err = ParseFields([]string{"foo@bar:int"}, alwaysInvalid)
	require.Error(t, err)

	// invalid format
	_, err = ParseFields([]string{"foo:int:int"}, alwaysInvalid)
	require.Error(t, err)
}

func TestParseFields1(t *testing.T) {
	name1, err := multiformatname.NewName("foo")
	require.NoError(t, err)
	name2, err := multiformatname.NewName("fooBar")
	require.NoError(t, err)
	name3, err := multiformatname.NewName("bar-foo")
	require.NoError(t, err)
	name4, err := multiformatname.NewName("foo_foo")
	require.NoError(t, err)

	tests := []struct {
		name   string
		fields []string
		want   Fields
		err    error
	}{
		{
			name: "test string types",
			fields: []string{
				name1.Original,
				name2.Original + ":string",
			},
			want: Fields{
				{
					Name:         name1,
					DatatypeName: datatype.String,
				},
				{
					Name:         name2,
					DatatypeName: datatype.String,
				},
			},
		},
		{
			name: "test number types",
			fields: []string{
				name1.Original + ":uint",
				name2.Original + ":int",
				name3.Original + ":bool",
			},
			want: Fields{
				{
					Name:         name1,
					DatatypeName: datatype.Uint,
				},
				{
					Name:         name2,
					DatatypeName: datatype.Int,
				},
				{
					Name:         name3,
					DatatypeName: datatype.Bool,
				},
			},
		},
		{
			name: "test list types",
			fields: []string{
				name1.Original + ":array.uint",
				name2.Original + ":array.int",
				name3.Original + ":array.string",
			},
			want: Fields{
				{
					Name:         name1,
					DatatypeName: datatype.UintSlice,
				},
				{
					Name:         name2,
					DatatypeName: datatype.IntSlice,
				},
				{
					Name:         name3,
					DatatypeName: datatype.StringSlice,
				},
			},
		},
		{
			name: "test mixed types",
			fields: []string{
				name1.Original + ":uint",
				name2.Original + ":array.coin",
				name3.Original,
				name4.Original + ":strings",
			},
			want: Fields{
				{
					Name:         name1,
					DatatypeName: datatype.Uint,
				},
				{
					Name:         name2,
					DatatypeName: datatype.Coins,
				},
				{
					Name:         name3,
					DatatypeName: datatype.String,
				},
				{
					Name:         name4,
					DatatypeName: datatype.StringSliceAlias,
				},
			},
		},
		{
			name: "test custom types",
			fields: []string{
				name1.Original + ":Bla",
				name2.Original + ":Test",
				name3.Original,
			},
			want: Fields{
				{
					Name:         name1,
					DatatypeName: datatype.Custom,
					Datatype:     "Bla",
				},
				{
					Name:         name2,
					DatatypeName: datatype.Custom,
					Datatype:     "Test",
				},
				{
					Name:         name3,
					DatatypeName: datatype.String,
				},
			},
		},
		{
			name: "test sdk.Coin types",
			fields: []string{
				name1.Original + ":coin",
				name2.Original + ":array.coin",
			},
			want: Fields{
				{
					Name:         name1,
					DatatypeName: datatype.Coin,
				},
				{
					Name:         name2,
					DatatypeName: datatype.Coins,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFields(tt.fields, noCheck)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.want, got)
		})
	}
}
