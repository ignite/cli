package gentx

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParseGentx(t *testing.T) {
	tests := []struct {
		name      string
		gentxPath string
		wantInfo  Info
		wantErr   bool
	}{
		{
			name:      "parse gentx file 1",
			gentxPath: "test/gentx1.json",
			wantInfo: Info{
				DelegatorAddress: "cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8ygndsj",
				PubKey:           []byte("aeQLCJOjXUyB7evOodI4mbrshIt3vhHGlycJDbUkaMs="),
				SelfDelegation: sdk.Coin{
					Denom:  "stake",
					Amount: sdk.NewInt(95000000),
				},
			},
		}, {
			name:      "parse gentx file 2",
			gentxPath: "test/gentx2.json",
			wantInfo: Info{
				DelegatorAddress: "cosmos1mmlqwyqk7neqegffp99q86eckpm4pjah3ytlpa",
				PubKey:           []byte("OL+EIoo7DwyaBFDbPbgAhwS5rvgIqoUa0x8qWqzfQVQ="),
				SelfDelegation: sdk.Coin{
					Denom:  "stake",
					Amount: sdk.NewInt(95000000),
				},
			},
		}, {
			name:      "parse invalid file",
			gentxPath: "test/gentx_invalid.json",
			wantErr:   true,
		}, {
			name:      "not found file",
			gentxPath: "test/gentx_not_found.json",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfo, _, err := FromPath(tt.gentxPath)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantInfo, gotInfo)
		})
	}
}
