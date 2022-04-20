package protopath

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatPackageName(t *testing.T) {
	cases := []struct {
		name string
		path []string
		want string
	}{
		{
			name: "full",
			path: []string{"a", "b", "c"},
			want: "a.b.c",
		},
		{
			name: "short",
			path: []string{"a", "b"},
			want: "a.b",
		},
		{
			name: "single",
			path: []string{"a"},
			want: "a",
		},
		{
			name: "triplicated",
			path: []string{"a", "a", "a"},
			want: "a",
		},
		{
			name: "duplicated prefix",
			path: []string{"a", "a", "b"},
			want: "a.b",
		},
		{
			name: "duplicated suffix",
			path: []string{"a", "b", "b"},
			want: "a.b",
		},
		{
			name: "empty",
			path: []string{},
			want: "",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, FormatPackageName(tt.path...))
		})
	}
}
