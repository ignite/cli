package multiformatname_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/multiformatname"
	"testing"
)

func TestNewMultiFormatName(t *testing.T) {
	// [valueToTest, lowerCamel, upperCamel, kebabCase]
	cases := [][4]string{
		{"foo", "foo", "Foo", "foo"},
		{"fooBar", "fooBar", "FooBar", "foo-bar"},
		{"foo-bar", "fooBar", "FooBar", "foo-bar"},
		{"foo_bar", "fooBar", "FooBar", "foo-bar"},
		{"foo_barFoobar", "fooBarFoobar", "FooBarFoobar", "foo-bar-foobar"},
		{"foo_-_bar", "fooBar", "FooBar", "foo---bar"},
		{"foo_-_Bar", "fooBar", "FooBar", "foo---bar"},
		{"fooBAR", "fooBAR", "FooBAR", "foo-bar"},
	}

	// test cases
	for _, testCase := range cases {
		name, err := multiformatname.NewMultiFormatName(testCase[0])
		require.NoError(t, err)
		require.Equal(
			t,
			testCase[0],
			name.Original,
		)
		require.Equal(
			t,
			testCase[1],
			name.LowerCamel,
			fmt.Sprintf("%s should be converted to the correct lower camel format", testCase[0]),
		)
		require.Equal(
			t,
			testCase[2],
			name.UpperCamel,
			fmt.Sprintf("%s should be converted the correct upper camel format", testCase[0]),
		)
		require.Equal(
			t,
			testCase[3],
			name.Kebab,
			fmt.Sprintf("%s should be converted the correct kebab format", testCase[0]),
		)
	}
}

func TestNewMultiFormatName2(t *testing.T) {
	// Test forbidden names
	cases := []string{
		"",
		"foo bar",
		"foo123",
		"@foo",
	}

	// test cases
	for _, testCase := range cases {
		_, err := multiformatname.NewMultiFormatName(testCase)
		require.Error(t, err)
	}
}
