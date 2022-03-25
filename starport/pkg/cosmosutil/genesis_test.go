package cosmosutil_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/starport/starport/pkg/cosmosutil"
)

const (
	genesisSample = `
{
	"foo": "bar",
	"genesis_time": "foobar"
}
`
	unixTime = 1600000000
	rfcTime  = "2020-09-13T12:26:40Z"
)

func TestSetGenesisTime(t *testing.T) {
	tmp := t.TempDir()
	tmpGenesis := filepath.Join(tmp, "genesis.json")

	// fails with no file
	require.Error(t, cosmosutil.SetGenesisTime(tmpGenesis, 0))

	require.NoError(t, os.WriteFile(tmpGenesis, []byte(genesisSample), 0644))
	require.NoError(t, cosmosutil.SetGenesisTime(tmpGenesis, unixTime))

	// check genesis modified value
	var actual struct {
		Foo         string `json:"foo"`
		GenesisTime string `json:"genesis_time"`
	}
	actualBytes, err := os.ReadFile(tmpGenesis)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(actualBytes, &actual))
	require.Equal(t, "bar", actual.Foo)
	require.Equal(t, rfcTime, actual.GenesisTime)
}

func TestChainGenesis_HasAccount(t *testing.T) {
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
			g := cosmosutil.Genesis{Accounts: tt.accounts}
			got := g.HasAccount(tt.address)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseChainGenesis(t *testing.T) {
	genesis1 := cosmosutil.ChainGenesis{ChainID: "earth-1"}
	genesis1.AppState.Auth.Accounts = []struct {
		Address string `json:"address"`
	}{{Address: "cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8ygndsj"}}
	genesis1.AppState.Staking.Params.BondDenom = "stake"

	genesis2 := cosmosutil.ChainGenesis{ChainID: "earth-1"}
	genesis2.AppState.Auth.Accounts = []struct {
		Address string `json:"address"`
	}{{Address: "cosmos1mmlqwyqk7neqegffp99q86eckpm4pjah3ytlpa"}}
	genesis2.AppState.Staking.Params.BondDenom = "stake"

	tests := []struct {
		name        string
		genesisPath string
		want        cosmosutil.ChainGenesis
		wantErr     bool
	}{
		{
			name:        "parse genesis file 1",
			genesisPath: "testdata/genesis1.json",
			want:        genesis1,
		}, {
			name:        "parse genesis file 2",
			genesisPath: "testdata/genesis2.json",
			want:        genesis2,
		}, {
			name:        "parse not found file",
			genesisPath: "testdata/genesis_invalid.json",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			genesisFile, err := os.ReadFile(tt.genesisPath)
			require.NoError(t, err)

			got, err := cosmosutil.ParseChainGenesis(genesisFile)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.want, got)
		})
	}
}

func TestParseGenesis(t *testing.T) {
	tests := []struct {
		name        string
		genesisPath string
		want        cosmosutil.Genesis
		wantErr     bool
	}{
		{
			name:        "parse genesis file 1",
			genesisPath: "testdata/genesis1.json",
			want: cosmosutil.Genesis{
				Accounts:   []string{"cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8ygndsj"},
				StakeDenom: "stake",
			},
		}, {
			name:        "parse genesis file 2",
			genesisPath: "testdata/genesis2.json",
			want: cosmosutil.Genesis{
				Accounts:   []string{"cosmos1mmlqwyqk7neqegffp99q86eckpm4pjah3ytlpa"},
				StakeDenom: "stake",
			},
		}, {
			name:        "parse not found file",
			genesisPath: "testdata/genesis_invalid.json",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			genesisFile, err := os.ReadFile(tt.genesisPath)
			require.NoError(t, err)

			got, err := cosmosutil.ParseGenesis(genesisFile)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.ElementsMatch(t, tt.want.Accounts, got.Accounts)
			require.Equal(t, tt.want.StakeDenom, got.StakeDenom)
		})
	}
}

func TestParseGenesisFromPath(t *testing.T) {
	tests := []struct {
		name        string
		genesisPath string
		want        cosmosutil.Genesis
		wantErr     bool
	}{
		{
			name:        "parse genesis file 1",
			genesisPath: "testdata/genesis1.json",
			want: cosmosutil.Genesis{
				Accounts:   []string{"cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8ygndsj"},
				StakeDenom: "stake",
			},
		}, {
			name:        "parse genesis file 2",
			genesisPath: "testdata/genesis2.json",
			want: cosmosutil.Genesis{
				Accounts:   []string{"cosmos1mmlqwyqk7neqegffp99q86eckpm4pjah3ytlpa"},
				StakeDenom: "stake",
			},
		}, {
			name:        "parse not found file",
			genesisPath: "testdata/genesis_not_found.json",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cosmosutil.ParseGenesisFromPath(tt.genesisPath)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.ElementsMatch(t, tt.want.Accounts, got.Accounts)
			require.Equal(t, tt.want.StakeDenom, got.StakeDenom)
		})
	}
}
