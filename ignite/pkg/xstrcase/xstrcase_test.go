package xstrcase

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLowercase(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "simple lowercase",
			arg:  "Example-Test",
			want: "exampletest",
		},
		{
			name: "already lowercase",
			arg:  "example_test",
			want: "exampletest",
		},
		{
			name: "mixed case with dash",
			arg:  "Mixed-Case_String",
			want: "mixedcasestring",
		},
		{
			name: "uppercase input",
			arg:  "UPPER-CASE",
			want: "uppercase",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Lowercase(tt.arg)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUpperCamel(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "simple camel case",
			arg:  "example_test",
			want: "ExampleTest",
		},
		{
			name: "mixed case with dash",
			arg:  "Mixed-Case_String",
			want: "MixedCaseString",
		},
		{
			name: "uppercase input",
			arg:  "UPPER_CASE",
			want: "UpperCase",
		},
		{
			name: "lowercase input",
			arg:  "lower_case",
			want: "LowerCase",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UpperCamel(tt.arg)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUppercase(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "simple uppercase",
			arg:  "example-test",
			want: "EXAMPLETEST",
		},
		{
			name: "already uppercase",
			arg:  "EXAMPLE_TEST",
			want: "EXAMPLETEST",
		},
		{
			name: "mixed case input",
			arg:  "Mixed-Case_String",
			want: "MIXEDCASESTRING",
		},
		{
			name: "lowercase input",
			arg:  "lower-case",
			want: "LOWERCASE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Uppercase(tt.arg)
			require.Equal(t, tt.want, got)
		})
	}
}
