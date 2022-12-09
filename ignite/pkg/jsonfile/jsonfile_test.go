package jsonfile

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/tarball"
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
			name:     "get boolean parameter",
			filepath: "testdata/jsonfile.json",
			key:      "launched",
			want:     true,
		},
		{
			name:     "get array parameter",
			filepath: "testdata/jsonfile.json",
			key:      "consensus_params.block.best_blocks",
			want:     []int{100, 20, 11, 4, 2},
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
			key:      "app_state.bank.balances.[0].coins",
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
			key:      "app_state.bank.balances.[0].coins",
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
			t.Cleanup(func() {
				err = f.Close()
				require.NoError(t, err)
			})
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
	jsonCoins, err := json.Marshal(sdk.NewCoin("bar", sdk.NewInt(500)))
	require.NoError(t, err)

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
				WithKeyValueInt(
					"consensus_params.block.max_bytes",
					22020096,
				),
			},
		},
		{
			name:     "update number field",
			filepath: "testdata/jsonfile.json",
			opts: []UpdateFileOption{
				WithKeyValueInt(
					"consensus_params.block.time_iota_ms",
					1000,
				),
			},
		},
		{
			name:     "update coin field",
			filepath: "testdata/jsonfile.json",
			opts: []UpdateFileOption{
				WithKeyValueTimestamp(
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
				WithKeyValueInt(
					"consensus_params.block.time_iota_ms",
					1000,
				),
				WithKeyValueTimestamp(
					"genesis_time",
					999999999,
				),
			},
		},
		{
			name:     "update bytes",
			filepath: "testdata/jsonfile.json",
			opts: []UpdateFileOption{
				WithKeyValueByte(
					"app_state.crisis.params.constant_fee",
					jsonCoins,
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

			// Rollback files after change
			b, err := f.Bytes()
			require.NoError(t, err)
			t.Cleanup(func() {
				var prettyJSON bytes.Buffer
				err := json.Indent(&prettyJSON, b, "", "  ")
				require.NoError(t, err)

				err = truncate(f.file, 0)
				require.NoError(t, err)
				err = f.Reset()
				require.NoError(t, err)
				_, err = f.file.Write(prettyJSON.Bytes())
				require.NoError(t, err)
				err = f.Close()
				require.NoError(t, err)
			})

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
				newValue := value
				err = f.Field(key, &newValue)
				require.Equal(t, value, newValue)
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
			want:     "036dbc0020f4ab5604f46a8e5a05c368e4cba41f48fcac2864641902c1dfcad5",
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
			t.Cleanup(func() {
				err = f.Close()
				require.NoError(t, err)
			})
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
		filepath        string
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
				filepath: "testdata/jsonfile.json",
			},
			verifyField: "chain_id",
			wantField:   "earth-1",
		},
		{
			name: "tarball URL",
			args: args{
				filepath:        "testdata/example.tar.gz",
				tarballFileName: "example.json",
			},
			verifyField: "chain_id",
			wantField:   "gaia-1",
		},
		{
			name: "invalid tarball file name",
			args: args{
				filepath:        "testdata/example.tar.gz",
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
			url := tt.args.url
			if url == "" {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					file, err := os.ReadFile(tt.args.filepath)
					require.NoError(t, err)
					_, err = w.Write(file)
					require.NoError(t, err)
				}))
				url = ts.URL
			}

			filepath := fmt.Sprintf("%s/jsonfile.json", t.TempDir())
			got, err := FromURL(context.TODO(), url, filepath, tt.args.tarballFileName)
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			t.Cleanup(func() {
				err = got.Close()
				require.NoError(t, err)
			})
			var verificationField string
			err = got.Field(tt.verifyField, &verificationField)
			require.NoError(t, err)
			require.Equal(t, tt.wantField, verificationField)
		})
	}
}
