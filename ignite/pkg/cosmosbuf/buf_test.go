package cosmosbuf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_join(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		want   string
	}{
		{
			name:   "empty",
			values: []string{},
		},
		{
			name:   "single",
			values: []string{"foo"},
			want:   `"foo"`,
		},
		{
			name:   "two",
			values: []string{"foo", "bar"},
			want:   `"foo","bar"`,
		},
		{
			name:   "three",
			values: []string{"foo", "bar", "baz"},
			want:   `"foo","bar","baz"`,
		},
		{
			name:   "repeated",
			values: []string{"foo", "bar", "baz", "foo", "bar", "baz"},
			want:   `"foo","bar","baz","foo","bar","baz"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := join(tt.values...)
			require.Equal(t, tt.want, got)
		})
	}
}
