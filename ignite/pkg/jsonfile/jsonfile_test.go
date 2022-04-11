package jsonfile

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestJSONFile_Field(t *testing.T) {
	type (
		invalidStruct struct {
			name   string
			number int
		}
		evidence struct {
			MaxAgeDuration  string `json:"max_age_duration"`
			MaxAgeNumBlocks string `json:"max_age_num_blocks"`
			MaxBytes        int64  `json:"max_bytes"`
		}
	)

	tests := []struct {
		name     string
		filepath string
		key      string
		want     interface{}
		err      error
	}{
		{
			name:     "get string parameter",
			filepath: "testdata/genesis.json",
			key:      "consensus_params.block.max_bytes",
			want:     "22020096",
		},
		{
			name:     "get number parameter",
			filepath: "testdata/genesis.json",
			key:      "consensus_params.block.time_iota_ms",
			want:     1000,
		},
		{
			name:     "get coins parameter",
			filepath: "testdata/genesis.json",
			key:      "app_state.bank.balances.coins",
			want:     sdk.Coins{sdk.NewCoin("stake", sdk.NewInt(95000000))},
		},
		{
			name:     "get custom parameter",
			filepath: "testdata/genesis.json",
			key:      "consensus_params.evidence",
			want: evidence{
				MaxAgeDuration:  "172800000000000",
				MaxAgeNumBlocks: "100000",
				MaxBytes:        1048576,
			},
		},
		{
			name:     "invalid coins parameter",
			filepath: "testdata/genesis.json",
			key:      "app_state.bank.balances.coins",
			want:     invalidStruct{name: "invalid", number: 110},
			err:      ErrInvalidValueType,
		},
		{
			name:     "invalid path",
			filepath: "testdata/genesis.json",
			key:      "invalid.field.path",
			want:     invalidStruct{name: "invalid", number: 110},
			err:      ErrFieldNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := FromPath(tt.filepath)
			require.NoError(t, err)
			out := reflect.New(reflect.TypeOf(tt.want))
			err = f.Field(tt.key, out.Interface())
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, out.Elem().Interface())
		})
	}
}

func TestJSONFile_Update(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		opts     []UpdateFileOption
		err      error
	}{
		{
			name:     "update string field",
			filepath: "testdata/genesis.json",
			opts: []UpdateFileOption{
				WithKeyValue(
					"consensus_params.block.max_bytes",
					"22020096",
				),
			},
		},
		{
			name:     "update string field to number",
			filepath: "testdata/genesis.json",
			opts: []UpdateFileOption{
				WithKeyIntValue(
					"consensus_params.block.max_bytes",
					22020096,
				),
			},
		},
		{
			name:     "update number field",
			filepath: "testdata/genesis.json",
			opts: []UpdateFileOption{
				WithKeyIntValue(
					"consensus_params.block.time_iota_ms",
					1000,
				),
			},
		},
		{
			name:     "update coin field",
			filepath: "testdata/genesis.json",
			opts: []UpdateFileOption{
				WithTime(
					"genesis_time",
					10000000,
				),
			},
		},
		{
			name:     "update all values type",
			filepath: "testdata/genesis.json",
			opts: []UpdateFileOption{
				WithKeyValue(
					"consensus_params.block.max_bytes",
					"3000000",
				),
				WithKeyIntValue(
					"consensus_params.block.time_iota_ms",
					1000,
				),
				WithTime(
					"genesis_time",
					999999999,
				),
			},
		},
		{
			name:     "add non-existing field",
			filepath: "testdata/genesis.json",
			opts: []UpdateFileOption{
				WithKeyValue(
					"app_state.auth.params.sig_verify_cost_ed25519",
					"111",
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := FromPath(tt.filepath)
			require.NoError(t, err)
			err = f.Update(tt.opts...)
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)

			updates := make(map[string][]byte)
			for _, opt := range tt.opts {
				opt(updates)
			}
			for key, value := range updates {
				var out interface{}
				err = f.Field(key, &out)
				switch out := out.(type) {
				case string:
					require.Equal(t, strings.ReplaceAll(string(value), "\"", ""), out)
				case float64:
					v, err := strconv.Atoi(string(value))
					require.NoError(t, err)
					require.Equal(t, v, int(out))
				default:
					require.Equal(t, value, out)
				}
			}
		})
	}
}

func TestJSONFile_Hash(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		want     string
		err      error
	}{
		{
			name:     "file hash",
			filepath: "testdata/genesis.json",
			want:     "f6f6913a1efacc78f61a63041f94672bb903c969ed33c29020f95997f46903dd",
		},
		{
			name:     "not found file",
			filepath: "testdata/genesis_not_found.json",
			err: errors.New(
				"cannot open the file: open testdata/genesis_not_found.json: no such file or directory",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := FromPath(tt.filepath)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, err.Error(), tt.err.Error())
				return
			}
			require.NoError(t, err)
			got, err := f.Hash()
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
