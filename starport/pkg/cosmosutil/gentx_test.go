package cosmosutil

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestChainGenesis_HasAccount(t *testing.T) {
	tests := []struct {
		name     string
		accounts []acc
		address  string
		want     bool
	}{
		{
			name:    "found account",
			address: "cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8ygndsj",
			accounts: []acc{
				{Address: "cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8ygndsj"},
				{Address: "cosmos1mmlqwyqk7neqegffp99q86eckpm4pjah3ytlpa"},
			},
			want: true,
		}, {
			name:    "not found account",
			address: "cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8pu8cup",
			accounts: []acc{
				{Address: "cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8ygndsj"},
				{Address: "cosmos1mmlqwyqk7neqegffp99q86eckpm4pjah3ytlpa"},
			},
			want: false,
		}, {
			name:     "empty accounts",
			address:  "cosmos1mmlqwyqk7neqegffp99q86eckpm4pjah3ytlpa",
			accounts: []acc{},
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := ChainGenesis{}
			g.AppState.Auth.Accounts = tt.accounts
			got := g.HasAccount(tt.address)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseGenesis(t *testing.T) {
	tests := []struct {
		name         string
		genesisPath  string
		wantAccounts []acc
		wantErr      bool
	}{
		{
			name:         "parse genesis file 1",
			genesisPath:  "testdata/genesis1.json",
			wantAccounts: []acc{{Address: "cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8ygndsj"}},
		}, {
			name:         "parse genesis file 2",
			genesisPath:  "testdata/genesis2.json",
			wantAccounts: []acc{{Address: "cosmos1mmlqwyqk7neqegffp99q86eckpm4pjah3ytlpa"}},
		}, {
			name:        "parse not found file",
			genesisPath: "testdata/genesis_not_found.json",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGenesis, err := ParseGenesis(tt.genesisPath)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.ElementsMatch(t, tt.wantAccounts, gotGenesis.AppState.Auth.Accounts)
		})
	}
}

func TestParseGentx(t *testing.T) {
	tests := []struct {
		name      string
		gentxPath string
		wantInfo  GentxInfo
		wantErr   bool
	}{
		{
			name:      "parse gentx file 1",
			gentxPath: "testdata/gentx1.json",
			wantInfo: GentxInfo{
				DelegatorAddress: "cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8ygndsj",
				PubKey:           []byte("aeQLCJOjXUyB7evOodI4mbrshIt3vhHGlycJDbUkaMs="),
				SelfDelegation: sdk.Coin{
					Denom:  "stake",
					Amount: sdk.NewInt(95000000),
				},
			},
		}, {
			name:      "parse gentx file 2",
			gentxPath: "testdata/gentx2.json",
			wantInfo: GentxInfo{
				DelegatorAddress: "cosmos1mmlqwyqk7neqegffp99q86eckpm4pjah3ytlpa",
				PubKey:           []byte("OL+EIoo7DwyaBFDbPbgAhwS5rvgIqoUa0x8qWqzfQVQ="),
				SelfDelegation: sdk.Coin{
					Denom:  "stake",
					Amount: sdk.NewInt(95000000),
				},
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
			gotInfo, _, err := ParseGentx(tt.gentxPath)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantInfo, gotInfo)
		})
	}
}
