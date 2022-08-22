package multiformatname_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/multiformatname"
)

func TestNewMultiFormatName(t *testing.T) {
	// [valueToTest, lowerCamel, upperCamel, kebabCase]
	cases := [][6]string{
		{"foo", "foo", "Foo", "foo", "foo", "foo"},
		{"fooBar", "fooBar", "FooBar", "foo-bar", "foo_bar", "foobar"},
		{"foo-bar", "fooBar", "FooBar", "foo-bar", "foo_bar", "foobar"},
		{"foo_bar", "fooBar", "FooBar", "foo-bar", "foo_bar", "foobar"},
		{"foo_barFoobar", "fooBarFoobar", "FooBarFoobar", "foo-bar-foobar", "foo_bar_foobar", "foobarfoobar"},
		{"foo_-_bar", "fooBar", "FooBar", "foo---bar", "foo___bar", "foobar"},
		{"foo_-_Bar", "fooBar", "FooBar", "foo---bar", "foo___bar", "foobar"},
		{"fooBAR", "fooBAR", "FooBAR", "foo-bar", "foo_bar", "foobar"},
		{"fooBar123", "fooBar123", "FooBar123", "foo-bar-123", "foo_bar_123", "foobar123"},
	}

	// test cases
	for _, testCase := range cases {
		name, err := multiformatname.NewName(testCase[0])
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
		require.Equal(
			t,
			testCase[4],
			name.Snake,
			fmt.Sprintf("%s should be converted the correct snake format", testCase[0]),
		)
		require.Equal(
			t,
			testCase[5],
			name.LowerCase,
			fmt.Sprintf("%s should be converted the correct lowercase format", testCase[0]),
		)
	}
}

func TestNewMultiFormatName2(t *testing.T) {
	// Test basic forbidden names
	cases := []string{
		"",
		"foo bar",
		"1foo",
		"-foo",
		"_foo",
		"@foo",
	}
	for _, testCase := range cases {
		_, err := multiformatname.NewName(testCase)
		require.Error(t, err)
	}

	// Test custom check
	alwaysWrong := func(string) error { return errors.New("always wrong") }
	_, err := multiformatname.NewName("foo", alwaysWrong)
	require.Error(t, err)

	alwaysGood := func(string) error { return nil }
	_, err = multiformatname.NewName("foo", alwaysGood)
	require.NoError(t, err)
}

func TestNoNumber(t *testing.T) {
	require.NoError(t, multiformatname.NoNumber("foo"))
	require.Error(t, multiformatname.NoNumber("foo1"))
}
