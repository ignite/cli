package field_test

import (
	"errors"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/field"
	"testing"
)

var (
	noCheck = func(string) error{return nil}
	alwaysInvalid = func(string) error{return errors.New("invalid name")}
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
			"foobar:bool",
			"barfoo:int",
			"foofoo:uint",
		},
		expected: []field.Field{
			{
				Name: "foo",
				Datatype: "string",
				DatatypeName: "string",
			},
			{
				Name: "bar",
				Datatype: "string",
				DatatypeName: "string",
			},
			{
				Name: "foobar",
				Datatype: "bool",
				DatatypeName: "bool",
			},
			{
				Name: "barfoo",
				Datatype: "int32",
				DatatypeName: "int",
			},
			{
				Name: "foofoo",
				Datatype: "uint64",
				DatatypeName: "uint",
			},
		},
	}

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

	// invalid format
	_, err = field.ParseFields([]string{"foo:int:int"}, alwaysInvalid)
	require.Error(t, err)
}