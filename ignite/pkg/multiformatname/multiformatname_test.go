package multiformatname_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
)

func TestNewName(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want multiformatname.Name
		err  error
	}{
		{
			name: "simple lowercase name",
			arg:  "foo",
			want: multiformatname.Name{
				Original:   "foo",
				LowerCamel: "foo",
				UpperCamel: "Foo",
				PascalCase: "Foo",
				LowerCase:  "foo",
				UpperCase:  "FOO",
				Kebab:      "foo",
				Snake:      "foo",
			},
		},
		{
			name: "camelCase name",
			arg:  "fooBar",
			want: multiformatname.Name{
				Original:   "fooBar",
				LowerCamel: "fooBar",
				UpperCamel: "FooBar",
				PascalCase: "FooBar",
				LowerCase:  "foobar",
				UpperCase:  "FOOBAR",
				Kebab:      "foo-bar",
				Snake:      "foo_bar",
			},
		},
		{
			name: "kebab-case name",
			arg:  "foo-bar",
			want: multiformatname.Name{
				Original:   "foo-bar",
				LowerCamel: "fooBar",
				UpperCamel: "FooBar",
				PascalCase: "FooBar",
				LowerCase:  "foobar",
				UpperCase:  "FOOBAR",
				Kebab:      "foo-bar",
				Snake:      "foo_bar",
			},
		},
		{
			name: "snake_case name",
			arg:  "foo_bar",
			want: multiformatname.Name{
				Original:   "foo_bar",
				LowerCamel: "fooBar",
				UpperCamel: "FooBar",
				PascalCase: "FooBar",
				LowerCase:  "foobar",
				UpperCase:  "FOOBAR",
				Kebab:      "foo-bar",
				Snake:      "foo_bar",
			},
		},
		{
			name: "mixed snake_case and camelCase name",
			arg:  "foo_barFoobar",
			want: multiformatname.Name{
				Original:   "foo_barFoobar",
				LowerCamel: "fooBarFoobar",
				UpperCamel: "FooBarFoobar",
				PascalCase: "FooBarFoobar",
				LowerCase:  "foobarfoobar",
				UpperCase:  "FOOBARFOOBAR",
				Kebab:      "foo-bar-foobar",
				Snake:      "foo_bar_foobar",
			},
		},
		{
			name: "mixed underscores and dashes",
			arg:  "foo_-_bar",
			want: multiformatname.Name{
				Original:   "foo_-_bar",
				LowerCamel: "fooBar",
				UpperCamel: "Foo__Bar",
				PascalCase: "FooBar",
				LowerCase:  "foobar",
				UpperCase:  "FOOBAR",
				Kebab:      "foo---bar",
				Snake:      "foo___bar",
			},
		},
		{
			name: "mixed underscores, dashes, and numbers",
			arg:  "foo_-_Bar1",
			want: multiformatname.Name{
				Original:   "foo_-_Bar1",
				LowerCamel: "fooBar1",
				UpperCamel: "Foo__Bar_1",
				PascalCase: "FooBar1",
				LowerCase:  "foobar1",
				UpperCase:  "FOOBAR1",
				Kebab:      "foo---bar-1",
				Snake:      "foo___bar_1",
			},
		},
		{
			name: "uppercase variant in simple name",
			arg:  "fooBAR",
			want: multiformatname.Name{
				Original:   "fooBAR",
				LowerCamel: "fooBar",
				UpperCamel: "FooBar",
				PascalCase: "FooBar",
				LowerCase:  "foobar",
				UpperCase:  "FOOBAR",
				Kebab:      "foo-bar",
				Snake:      "foo_bar",
			},
		},
		{
			name: "uppercase variant with starting capital",
			arg:  "FooBAR",
			want: multiformatname.Name{
				Original:   "FooBAR",
				LowerCamel: "fooBar",
				UpperCamel: "FooBar",
				PascalCase: "FooBar",
				LowerCase:  "foobar",
				UpperCase:  "FOOBAR",
				Kebab:      "foo-bar",
				Snake:      "foo_bar",
			},
		},
		{
			name: "camelCase name with numbers",
			arg:  "fooBar123",
			want: multiformatname.Name{
				Original:   "fooBar123",
				LowerCamel: "fooBar123",
				UpperCamel: "FooBar_123",
				PascalCase: "FooBar123",
				LowerCase:  "foobar123",
				UpperCase:  "FOOBAR123",
				Kebab:      "foo-bar-123",
				Snake:      "foo_bar_123",
			},
		},
		{
			name: "multiple numbers in name",
			arg:  "para_2_m_s_43_tr_1",
			want: multiformatname.Name{
				Original:   "para_2_m_s_43_tr_1",
				LowerCamel: "para2MS43Tr1",
				UpperCamel: "Para_2MS_43Tr_1",
				PascalCase: "Para2MS43Tr1",
				LowerCase:  "para2ms43tr1",
				UpperCase:  "PARA2MS43TR1",
				Kebab:      "para-2-m-s-43-tr-1",
				Snake:      "para_2_m_s_43_tr_1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := multiformatname.NewName(tt.arg)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.want, got)
		})
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
