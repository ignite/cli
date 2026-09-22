package module

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSwitchCosmosSDKPackagePath(t *testing.T) {
	cases := []struct {
		name, srcPath, sdkDir, want string
	}{
		{
			name:    "cosmossdk.io module with x/ submodule",
			srcPath: "/home/user/go/pkg/mod/cosmossdk.io/x/evidence@v0.1.1",
			sdkDir:  "/home/user/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.55.0",
			want:    filepath.Join("/home/user/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.55.0", "proto", "cosmos", "evidence"),
		},
		{
			name:    "cosmos-sdk module root",
			srcPath: "/home/user/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.55.0",
			sdkDir:  "/home/user/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.55.0",
			want:    filepath.Join("/home/user/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.55.0", "proto"),
		},
		{
			name:    "cosmos-sdk module root without version",
			srcPath: "/home/user/go/pkg/mod/github.com/cosmos/cosmos-sdk",
			sdkDir:  "/home/user/go/pkg/mod/github.com/cosmos/cosmos-sdk",
			want:    filepath.Join("/home/user/go/pkg/mod/github.com/cosmos/cosmos-sdk", "proto"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, switchCosmosSDKPackagePath(tt.srcPath, tt.sdkDir))
		})
	}
}
