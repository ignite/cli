package goanalysis_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/goanalysis"
)

var (
	MainFile   = []byte(`package main`)
	ImportFile = []byte(`
package app

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	queryonlymodmodule "github.com/tendermint/testchain/x/queryonlymod"
	queryonlymodmodulekeeper "github.com/tendermint/testchain/x/queryonlymod/keeper"
	queryonlymodmoduletypes "github.com/tendermint/testchain/x/queryonlymod/types"
)
`)
)

func TestDiscoverMain(t *testing.T) {
	tests := []struct {
		name       string
		mainFiles  []string
		expectFind bool
	}{
		{
			name:       "single main",
			mainFiles:  []string{"main.go"},
			expectFind: true,
		},
		{
			name:       "no mains",
			mainFiles:  []string{},
			expectFind: false,
		},
		{
			name:       "single main in sub-folder",
			mainFiles:  []string{"sub/main.go"},
			expectFind: true,
		},
		{
			name:       "single main with different name",
			mainFiles:  []string{"sub/somethingelse.go"},
			expectFind: true,
		},
		{
			name: "multiple mains",
			mainFiles: []string{
				"main.go",
				"sub/main.go",
				"diffSub/alsomain.go",
			},
			expectFind: true,
		},
		{
			name:       "single main with wrong extension",
			mainFiles:  []string{"main.ogg"},
			expectFind: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			want, err := createMainFiles(tmpDir, tt.mainFiles)
			require.NoError(t, err)

			actual, err := goanalysis.DiscoverMain(tmpDir)
			require.NoError(t, err)
			if !tt.expectFind {
				want = []string{}
			}
			require.ElementsMatch(t, actual, want)
		})
	}
}

func TestDiscoverOneMain(t *testing.T) {
	tests := []struct {
		name      string
		mainFiles []string
		err       error
	}{
		{
			name:      "single main",
			mainFiles: []string{"main.go"},
			err:       nil,
		},
		{
			name: "multiple mains",
			mainFiles: []string{
				"main.go",
				"sub/main.go",
			},
			err: goanalysis.ErrMultipleMainPackagesFound,
		},
		{
			name:      "no mains",
			mainFiles: []string{},
			err:       errors.New("main package cannot be found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			want, err := createMainFiles(tmpDir, tt.mainFiles)
			require.NoError(t, err)

			actual, err := goanalysis.DiscoverOneMain(tmpDir)

			require.Equal(t, tt.err, err)

			if tt.err == nil {
				require.Equal(t, 1, len(want))
				require.Equal(t, want[0], actual)
			}
		})
	}
}

func createMainFiles(tmpDir string, mainFiles []string) (pathsWithMain []string, err error) {
	for _, mf := range mainFiles {
		mainFile := filepath.Join(tmpDir, mf)
		dir := filepath.Dir(mainFile)

		if err = os.MkdirAll(dir, 0770); err != nil {
			return nil, err
		}

		if err = os.WriteFile(mainFile, MainFile, 0644); err != nil {
			return nil, err
		}

		pathsWithMain = append(pathsWithMain, dir)
	}

	return pathsWithMain, nil
}

func TestFindImportedPackages(t *testing.T) {
	tmpDir := t.TempDir()

	tmpFile := filepath.Join(tmpDir, "app.go")
	err := os.WriteFile(tmpFile, ImportFile, 0644)
	require.NoError(t, err)

	packages, err := goanalysis.FindImportedPackages(tmpFile)
	require.NoError(t, err)
	require.EqualValues(t, packages, map[string]string{
		"io":                       "io",
		"http":                     "net/http",
		"filepath":                 "path/filepath",
		"baseapp":                  "github.com/cosmos/cosmos-sdk/baseapp",
		"client":                   "github.com/cosmos/cosmos-sdk/client",
		"types":                    "github.com/cosmos/cosmos-sdk/codec/types",
		"api":                      "github.com/cosmos/cosmos-sdk/server/api",
		"config":                   "github.com/cosmos/cosmos-sdk/server/config",
		"servertypes":              "github.com/cosmos/cosmos-sdk/server/types",
		"simapp":                   "github.com/cosmos/cosmos-sdk/simapp",
		"queryonlymodmodule":       "github.com/tendermint/testchain/x/queryonlymod",
		"queryonlymodmodulekeeper": "github.com/tendermint/testchain/x/queryonlymod/keeper",
		"queryonlymodmoduletypes":  "github.com/tendermint/testchain/x/queryonlymod/types",
	})
}
