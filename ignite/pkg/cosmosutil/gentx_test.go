package cosmosutil_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosutil"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/starport/starport/pkg/cosmosutil"
)

func TestParseGentx(t *testing.T) {
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
				PubKey:           []byte("aeQLCJOjXUyB7evOodI4mbrshIt3vhHGlycJDbUkaMs="),
				SelfDelegation: sdk.Coin{
					Denom:  "stake",
					Amount: sdk.NewInt(95000000),
				},
				Memo: "9b1f4adbfb0c0b513040d914bfb717303c0eaa71@192.168.0.148:26656",
			},
		}, {
			name:      "parse gentx file 2",
			gentxPath: "testdata/gentx2.json",
			wantInfo: cosmosutil.GentxInfo{
				DelegatorAddress: "cosmos1mmlqwyqk7neqegffp99q86eckpm4pjah3ytlpa",
				PubKey:           []byte("OL+EIoo7DwyaBFDbPbgAhwS5rvgIqoUa0x8qWqzfQVQ="),
				SelfDelegation: sdk.Coin{
					Denom:  "stake",
					Amount: sdk.NewInt(95000000),
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

func TestPubKey_Equal(t *testing.T) {
	tests := []struct {
		name   string
		pb     []byte
		cmpKey []byte
		want   bool
	}{
		{
			name:   "equal public keys",
			pb:     []byte("aeQLCJOjXUyB7evOodI4mbrshIt3vhHGlycJDbUkaMs="),
			cmpKey: []byte("aeQLCJOjXUyB7evOodI4mbrshIt3vhHGlycJDbUkaMs="),
			want:   true,
		},
		{
			name:   "not equal public keys",
			pb:     []byte("aeQLCJOjXUyB7evOodI4mbrshIt3vhHGlycJDbUkaMs="),
			cmpKey: []byte("EIoo7DwyaBFDbPbgAhwS5rvgIqoUa0x8qWqzfQVQ="),
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := cosmosutil.PubKey(tt.pb)
			got := pb.Equal(tt.cmpKey)
			require.Equal(t, tt.want, got)
		})
	}
}
