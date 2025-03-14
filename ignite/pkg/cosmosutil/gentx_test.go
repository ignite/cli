package cosmosutil_test

import (
	"encoding/base64"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosutil"
)

func TestParseGentx(t *testing.T) {
	pk1, err := base64.StdEncoding.DecodeString("aeQLCJOjXUyB7evOodI4mbrshIt3vhHGlycJDbUkaMs=")
	require.NoError(t, err)
	pk2, err := base64.StdEncoding.DecodeString("OL+EIoo7DwyaBFDbPbgAhwS5rvgIqoUa0x8qWqzfQVQ=")
	require.NoError(t, err)

	tests := []struct {
		name      string
		gentxPath string
		wantInfo  cosmosutil.GentxInfo
		wantErr   bool
	}{
		{
			name:      "parse gentx file 1",
			gentxPath: "testdata/gentx1.json",
			wantInfo: cosmosutil.GentxInfo{
				DelegatorAddress: "cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8ygndsj",
				ValidatorAddress: "cosmosvaloper1dd246yq6z5vzjz9gh8cff46pll75yyl8pu8cup",
				PubKey:           ed25519.PubKey(pk1),
				SelfDelegation: sdk.Coin{
					Denom:  "stake",
					Amount: sdkmath.NewInt(95000000),
				},
				Memo: "9b1f4adbfb0c0b513040d914bfb717303c0eaa71@192.168.0.148:26656",
			},
		}, {
			name:      "parse gentx file 2",
			gentxPath: "testdata/gentx2.json",
			wantInfo: cosmosutil.GentxInfo{
				DelegatorAddress: "",
				ValidatorAddress: "cosmosvaloper1mmlqwyqk7neqegffp99q86eckpm4pjah5sl2dw",
				PubKey:           ed25519.PubKey(pk2),
				SelfDelegation: sdk.Coin{
					Denom:  "stake",
					Amount: sdkmath.NewInt(95000000),
				},
				Memo: "a412c917cb29f73cc3ad0592bbd0152fe0e690bd@192.168.0.148:26656",
			},
		}, {
			name:      "parse invalid file",
			gentxPath: "testdata/gentx_invalid.json",
			wantErr:   true,
		}, {
			name:      "not found file",
			gentxPath: "testdata/gentx_not_found.json",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfo, _, err := cosmosutil.GentxFromPath(tt.gentxPath)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantInfo, gotInfo)
		})
	}
}
