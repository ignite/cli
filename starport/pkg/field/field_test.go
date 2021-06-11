package field_test

import (
	"errors"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/field"
	"github.com/tendermint/starport/starport/pkg/multiformatname"
	"testing"
)

var (
	noCheck       = func(string) error { return nil }
	alwaysInvalid = func(string) error { return errors.New("invalid name") }
)

type testCases struct {
	provided []string
	expected []field.Field
}

func TestParseFields(t *testing.T) {
	cases := testCases{
		provided: []string{
			"foo",
			"bar:string",
			"fooBar:bool",
			"bar-foo:int",
			"foo_foo:uint",
		},
		expected: []field.Field{
			{
				Datatype:     "string",
				DatatypeName: "string",
			},
			{
				Datatype:     "string",
				DatatypeName: "string",
			},
			{
				Datatype:     "bool",
				DatatypeName: "bool",
			},
			{
				Datatype:     "int32",
				DatatypeName: "int",
			},
			{
				Datatype:     "uint64",
				DatatypeName: "uint",
			},
		},
	}
	cases.expected[0].Name, _ = multiformatname.NewMultiFormatName("foo")
	cases.expected[1].Name, _ = multiformatname.NewMultiFormatName("bar")
	cases.expected[2].Name, _ = multiformatname.NewMultiFormatName("fooBar")
	cases.expected[3].Name, _ = multiformatname.NewMultiFormatName("bar-foo")
	cases.expected[4].Name, _ = multiformatname.NewMultiFormatName("foo_foo")

	actual, err := field.ParseFields(cases.provided, noCheck)
	require.NoError(t, err)
	require.Equal(t, cases.expected, actual)

	// No field provided
	actual, err = field.ParseFields([]string{}, noCheck)
	require.NoError(t, err)
	require.Empty(t, actual)
}

func TestParseFields2(t *testing.T) {
	// test failing cases

	// check doesn't pass
	_, err := field.ParseFields([]string{"foo"}, alwaysInvalid)
	require.Error(t, err)

	// duplicated field
	_, err = field.ParseFields([]string{"foo", "foo:int"}, noCheck)
	require.Error(t, err)

	// invalid type
	_, err = field.ParseFields([]string{"foo:invalid"}, alwaysInvalid)
	require.Error(t, err)

	// invalid field name
	_, err = field.ParseFields([]string{"foo@bar:int"}, alwaysInvalid)
	require.Error(t, err)

	// invalid format
	_, err = field.ParseFields([]string{"foo:int:int"}, alwaysInvalid)
	require.Error(t, err)
}
