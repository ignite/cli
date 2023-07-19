package cosmosbuf

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindSDKPath(t *testing.T) {
	testCases := []struct {
		name     string
		protoDir string
		want     string
		err      error
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
			name:     "no version",
			protoDir: "test/myproto/test/proto",
			err:      errors.New("invalid sdk mod dir: test/myproto/test/proto"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findSDKProtoPath(tt.protoDir)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
