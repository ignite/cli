package jsonfile

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/ignite-hq/cli/ignite/pkg/tarball"
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
			filepath: "testdata/jsonfile.json",
			key:      "consensus_params.block.max_bytes",
			want:     "22020096",
		},
		{
			name:     "get number parameter",
			filepath: "testdata/jsonfile.json",
			key:      "consensus_params.block.time_iota_ms",
			want:     1000,
		},
		{
			name:     "get coins parameter",
			filepath: "testdata/jsonfile.json",
			key:      "app_state.bank.balances.coins",
			want:     sdk.Coins{sdk.NewCoin("stake", sdk.NewInt(95000000))},
		},
		{
			name:     "get custom parameter",
			filepath: "testdata/jsonfile.json",
			key:      "consensus_params.evidence",
			want: evidence{
				MaxAgeDuration:  "172800000000000",
				MaxAgeNumBlocks: "100000",
				MaxBytes:        1048576,
			},
		},
		{
			name:     "invalid coins parameter",
			filepath: "testdata/jsonfile.json",
			key:      "app_state.bank.balances.coins",
			want:     invalidStruct{name: "invalid", number: 110},
			err:      ErrInvalidValueType,
		},
		{
			name:     "invalid path",
			filepath: "testdata/jsonfile.json",
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
			filepath: "testdata/jsonfile.json",
			opts: []UpdateFileOption{
				WithKeyValue(
					"consensus_params.block.max_bytes",
					"22020096",
				),
			},
		},
		{
			name:     "update string field to number",
			filepath: "testdata/jsonfile.json",
			opts: []UpdateFileOption{
				WithKeyIntValue(
					"consensus_params.block.max_bytes",
					22020096,
				),
			},
		},
		{
			name:     "update number field",
			filepath: "testdata/jsonfile.json",
			opts: []UpdateFileOption{
				WithKeyIntValue(
					"consensus_params.block.time_iota_ms",
					1000,
				),
			},
		},
		{
			name:     "update coin field",
			filepath: "testdata/jsonfile.json",
			opts: []UpdateFileOption{
				WithTime(
					"genesis_time",
					10000000,
				),
			},
		},
		{
			name:     "update all values type",
			filepath: "testdata/jsonfile.json",
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
			filepath: "testdata/jsonfile.json",
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
			filepath: "testdata/jsonfile.json",
			want:     "4d685d9cb6f9fb9815a33f10a75cd9970f162bbc6ebc8c5c0e3fd166d1b3ee93",
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

func TestFromURL(t *testing.T) {
	type args struct {
		url             string
		tarballFileName string
	}
	tests := []struct {
		name        string
		args        args
		verifyField string
		wantField   string
		err         error
	}{
		{
			name: "JSON URL",
			args: args{
				url: "https://raw.githubusercontent.com/ignite-hq/cli/4d4f6b436d15aa3fb8aeca4bf91b6a557f897f9b/ignite/pkg/tarball/testdata/example.json",
			},
			verifyField: "chain_id",
			wantField:   "gaia-1",
		},
		{
			name: "tarball URL",
			args: args{
				url:             "https://github.com/ignite-hq/cli/raw/4d4f6b436d15aa3fb8aeca4bf91b6a557f897f9b/ignite/pkg/tarball/testdata/example-subfolder.tar.gz",
				tarballFileName: "example.json",
			},
			verifyField: "chain_id",
			wantField:   "gaia-1",
		},
		{
			name: "invalid tarball file name",
			args: args{
				url:             "https://github.com/ignite-hq/cli/raw/4d4f6b436d15aa3fb8aeca4bf91b6a557f897f9b/ignite/pkg/tarball/testdata/example-subfolder.tar.gz",
				tarballFileName: "invalid.json",
			},
			err: tarball.ErrGzipFileNotFound,
		},
		{
			name: "invalid link",
			args: args{
				url: "https://github.com/invalid_example.json",
			},
			err: ErrInvalidURL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filepath := fmt.Sprintf("%s/jsonfile.json", os.TempDir())
			got, err := FromURL(context.TODO(), tt.args.url, filepath, tt.args.tarballFileName)
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			var verificationField string
			err = got.Field(tt.verifyField, &verificationField)
			require.NoError(t, err)
			require.Equal(t, tt.wantField, verificationField)
		})
	}
}
