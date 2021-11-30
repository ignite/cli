package cosmosutil_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/cosmosutil"
)

func TestChainGenesis_HasAccount(t *testing.T) {
	type account struct {
		Address string `json:"address"`
	}
	tests := []struct {
		name     string
		accounts []string
		address  string
		want     bool
	}{
		{
			name:    "found account",
			address: "cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8ygndsj",
			accounts: []string{
				"cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8ygndsj",
				"cosmos1mmlqwyqk7neqegffp99q86eckpm4pjah3ytlpa",
			},
			want: true,
		}, {
			name:    "not found account",
			address: "cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8pu8cup",
			accounts: []string{
				"cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8ygndsj",
				"cosmos1mmlqwyqk7neqegffp99q86eckpm4pjah3ytlpa",
			},
			want: false,
		}, {
			name:     "empty accounts",
			address:  "cosmos1mmlqwyqk7neqegffp99q86eckpm4pjah3ytlpa",
			accounts: []string{},
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := cosmosutil.ChainGenesis{}
			for _, acc := range tt.accounts {
				g.AppState.Auth.Accounts = append(g.AppState.Auth.Accounts, account{Address: acc})
			}
			got := g.HasAccount(tt.address)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseGenesis(t *testing.T) {
	tests := []struct {
		name         string
		genesisPath  string
		wantAccounts []string
		wantErr      bool
	}{
		{
			name:         "parse genesis file 1",
			genesisPath:  "testdata/genesis1.json",
			wantAccounts: []string{"cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8ygndsj"},
		}, {
			name:         "parse genesis file 2",
			genesisPath:  "testdata/genesis2.json",
			wantAccounts: []string{"cosmos1mmlqwyqk7neqegffp99q86eckpm4pjah3ytlpa"},
		}, {
			name:        "parse not found file",
			genesisPath: "testdata/genesis_not_found.json",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGenesis, err := cosmosutil.ParseGenesis(tt.genesisPath)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			gotAddrs := make([]string, 0)
			for _, acc := range gotGenesis.AppState.Auth.Accounts {
				gotAddrs = append(gotAddrs, acc.Address)
			}
			require.ElementsMatch(t, tt.wantAccounts, gotAddrs)
		})
	}
}

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
			gotInfo, _, err := cosmosutil.ParseGentx(tt.gentxPath)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantInfo, gotInfo)
		})
	}
}
