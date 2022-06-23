package cosmosutil_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cosmosutil"
)

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

func TestUpdateGenesis(t *testing.T) {
	genesisSample := `
{
  "number": 33,
  "foo": "bar",
  "genesis_time": "foobar",
  "app_state": {
    "monitoring": {
      "chain-id": "ignite-1"
    },
    "foobar": "baz",
	"debug": "false",
	"height": "100",
	"time": "100",
    "staling": {
      "params": []
    }
  }
}
`
	type args struct {
		genesis string
		options []cosmosutil.GenesisField
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "with key and value",
			args: args{
				genesis: genesisSample,
				options: []cosmosutil.GenesisField{
					cosmosutil.WithKeyValue("foo", "foobar"),
				},
			},
		},
		{
			name: "with key and bool value",
			args: args{
				genesis: genesisSample,
				options: []cosmosutil.GenesisField{
					cosmosutil.WithKeyValueBoolean("app_state.debug", true),
				},
			},
		},
		{
			name: "with key and int value",
			args: args{
				genesis: genesisSample,
				options: []cosmosutil.GenesisField{
					cosmosutil.WithKeyValueInt("app_state.height", 199),
				},
			},
		},
		{
			name: "with key and uint value",
			args: args{
				genesis: genesisSample,
				options: []cosmosutil.GenesisField{
					cosmosutil.WithKeyValueUint("app_state.height", 438),
				},
			},
		},
		{
			name: "with key and timestamp value",
			args: args{
				genesis: genesisSample,
				options: []cosmosutil.GenesisField{
					cosmosutil.WithKeyValueTimestamp("app_state.time", 3000),
				},
			},
		},
		{
			name: "with all key and value types",
			args: args{
				genesis: genesisSample,
				options: []cosmosutil.GenesisField{
					cosmosutil.WithKeyValue("foo", "baz"),
					cosmosutil.WithKeyValue("app_state.monitoring.chain-id", "spn-1"),
					cosmosutil.WithKeyValueBoolean("app_state.debug", false),
					cosmosutil.WithKeyValueInt("app_state.height", 123),
					cosmosutil.WithKeyValueUint("app_state.height", 343),
					cosmosutil.WithKeyValueTimestamp("app_state.time", 999999),
				},
			},
		},
		{
			name: "casting key value",
			args: args{
				genesis: genesisSample,
				options: []cosmosutil.GenesisField{
					cosmosutil.WithKeyValue("number", "number_value"),
				},
			},
		},
		{
			name: "with wrong key path",
			args: args{
				genesis: genesisSample,
				options: []cosmosutil.GenesisField{
					cosmosutil.WithKeyValue("wrong.path", "foobar"),
					cosmosutil.WithKeyValue("app_state.monitoring.wrong", "baz"),
				},
			},
		},
		{
			name: "with file path",
			args: args{
				genesis: "",
				options: []cosmosutil.GenesisField{
					cosmosutil.WithKeyValue("foo", "foobar"),
				},
			},
			err: errors.New("Key path not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpGenesis := filepath.Join(t.TempDir(), "genesis.json")
			if tt.args.genesis != "" {
				require.Error(t,
					cosmosutil.UpdateGenesis(
						tmpGenesis,
						cosmosutil.WithKeyValue("test", "test"),
					),
				)
			}
			require.NoError(t, os.WriteFile(tmpGenesis, []byte(tt.args.genesis), 0644))
			err := cosmosutil.UpdateGenesis(tmpGenesis, tt.args.options...)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), tt.err.Error())
				return
			}
			require.NoError(t, err)

			genesisBytes, err := os.ReadFile(tmpGenesis)

			f := map[string]string{}
			for _, applyField := range tt.args.options {
				applyField(f)
			}

			for key, value := range f {
				val, err := jsonparser.GetString(genesisBytes, strings.Split(key, ".")...)
				require.NoError(t, err)
				require.Equal(t, val, value)
			}
		})
	}
}
