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
	expected field.Fields
}

func TestParseFields(t *testing.T) {
	names := []string{
		"foo",
		"bar",
		"fooBar",
		"bar-foo",
		"foo_foo",
	}
	cases := testCases{
		provided: []string{
			names[0],
			names[1] + ":string",
			names[2] + ":bool",
			names[3] + ":int",
			names[4] + ":uint",
		},
		expected: field.Fields{
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
	for i, name := range names {
		cases.expected[i].Name, _ = multiformatname.NewName(name)
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

	// invalid field name
	_, err = field.ParseFields([]string{"foo@bar:int"}, alwaysInvalid)
	require.Error(t, err)

	// invalid format
	_, err = field.ParseFields([]string{"foo:int:int"}, alwaysInvalid)
	require.Error(t, err)
}
