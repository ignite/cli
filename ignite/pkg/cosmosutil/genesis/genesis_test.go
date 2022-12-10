package genesis_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	cosmosgenesis "github.com/ignite/cli/ignite/pkg/cosmosutil/genesis"
)

func TestModuleParamField(t *testing.T) {
	tests := []struct {
		name   string
		module string
		param  string
		want   string
	}{
		{
			name:   "valid 1",
			module: "foo",
			param:  "bar",
			want:   "app_state.foo.params.bar",
		},
		{
			name:   "valid 2",
			module: "bar",
			param:  "foo",
			want:   "app_state.bar.params.foo",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := cosmosgenesis.ModuleParamField(tc.module, tc.param)
			require.Equal(t, tc.want, got)
		})
	}
}
