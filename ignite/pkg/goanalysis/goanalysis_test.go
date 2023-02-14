package goanalysis_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/goanalysis"
	"github.com/ignite/cli/ignite/pkg/xast"
)

var MainFile = []byte(`package main`)

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

		if err = os.MkdirAll(dir, 0o770); err != nil {
			return nil, err
		}

		if err = os.WriteFile(mainFile, MainFile, 0o644); err != nil {
			return nil, err
		}

		pathsWithMain = append(pathsWithMain, dir)
	}

	return pathsWithMain, nil
}

func TestGenVarExists(t *testing.T) {
	testFile := `
		package goanalysis
		
		import (
			"context"
			"errors"
			"path/filepath"
		)
		
		const (
			fooConst = "foo"
		)
		
		type (
			fooStruct struct {
				name string
			}
		)
		
		var (
			fooVar       = filepath.Join("test", "join")
			contextVar   = context.Background()
			fooStructVar = fooStruct{}
		)
		
		var (
			errorFooVar  = errors.New("error foo")
			bazStructVar = fooStruct{}
			errorBarVar  = errors.New("error bar")
		)
		
		func fooMethod(foo string) error {
			return nil
		}
`
	filename := filepath.Join(t.TempDir(), "var.go")
	require.NoError(t, os.WriteFile(filename, []byte(testFile), 0o644))

	tests := []struct {
		name            string
		goImport        string
		methodSignature string
		want            bool
	}{
		{
			name:            "test success assign",
			methodSignature: "Background",
			goImport:        "context",
			want:            true,
		},
		{
			name:            "test success assign",
			methodSignature: "Join",
			goImport:        "path/filepath",
			want:            true,
		},
		{
			name:            "test success assign",
			methodSignature: "New",
			goImport:        "errors",
			want:            true,
		},
		{
			name:            "test invalid import",
			methodSignature: "Join",
			goImport:        "errors",
			want:            false,
		},
		{
			name:            "test invalid case sensitive assign",
			methodSignature: "join",
			goImport:        "context",
			want:            false,
		},
		{
			name:            "test invalid struct assign",
			methodSignature: "fooStruct",
			goImport:        "context",
			want:            false,
		},
		{
			name:            "test invalid method signature",
			methodSignature: "fooMethod",
			goImport:        "context",
			want:            false,
		},
		{
			name:            "test not found name",
			methodSignature: "Invalid",
			goImport:        "context",
			want:            false,
		},
		{
			name:            "test invalid assign with wrong",
			methodSignature: "invalid.New",
			goImport:        "context",
			want:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appPkg, _, err := xast.ParseFile(filename)
			require.NoError(t, err)

			got := goanalysis.GenVarExists(appPkg, tt.goImport, tt.methodSignature)
			require.Equal(t, tt.want, got)
		})
	}
}
