package cosmosbuf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindSDKPath(t *testing.T) {
	testCases := []struct {
		name     string
		protoDir string
		want     string
	}{
		{
			name:     "full path",
			protoDir: "/mod/github.com/cosmos/cosmos-sdk@v0.47.2/test/path/proto",
			want:     "/mod/github.com/cosmos/cosmos-sdk@v0.47.2/proto",
		},
		{
			name:     "simple path",
			protoDir: "myproto@v1/test/proto/animo/sdk",
			want:     "myproto@v1/proto",
		},
		{
			name:     "only version",
			protoDir: "test/myproto@v1",
			want:     "test/myproto@v1/proto",
		},
		{
			name:     "semantic version",
			protoDir: "test/myproto@v0.3.1/test/proto",
			want:     "test/myproto@v0.3.1/proto",
		},
		{
			name:     "no version (local)",
			protoDir: "test/myproto/test/proto",
			want:     "test/myproto/test/proto",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := findSDKProtoPath(tt.protoDir)
			require.Equal(t, tt.want, got)
		})
	}
}
