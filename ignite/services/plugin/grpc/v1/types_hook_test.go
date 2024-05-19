package v1_test

import (
	"testing"

	v1 "github.com/ignite/cli/v29/ignite/services/plugin/grpc/v1"
	"github.com/stretchr/testify/require"
)

func TestHookCommandPath(t *testing.T) {
	cases := []struct {
		name, wantPath string
		hook           *v1.Hook
	}{
		{
			name: "relative path",
			hook: &v1.Hook{
				PlaceHookOn: "chain",
			},
			wantPath: "ignite chain",
		},
		{
			name: "full path",
			hook: &v1.Hook{
				PlaceHookOn: "ignite chain",
			},
			wantPath: "ignite chain",
		},
		{
			name: "path with spaces",
			hook: &v1.Hook{
				PlaceHookOn: " ignite scaffold  ",
			},
			wantPath: "ignite scaffold",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := tc.hook.CommandPath()
			require.Equal(t, tc.wantPath, path)
		})
	}
}
