package goanalysis_test

import (
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/goanalysis"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
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
			if tt.err != nil {
				require.Error(t, err)
				require.True(t, errors.Is(tt.err, err))
				return
			}
			require.NoError(t, err)
			require.Equal(t, 1, len(want))
			require.Equal(t, want[0], actual)
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

func TestFuncVarExists(t *testing.T) {
	tests := []struct {
		name            string
		testFile        string
		goImport        string
		methodSignature string
		want            bool
	}{
		{
			name:            "test a declaration inside a method success",
			testFile:        "testdata/varexist",
			methodSignature: "Background",
			goImport:        "context",
			want:            true,
		},
		{
			name:            "test global declaration success",
			testFile:        "testdata/varexist",
			methodSignature: "Join",
			goImport:        "path/filepath",
			want:            true,
		},
		{
			name:            "test a declaration inside an if and inside a method success",
			testFile:        "testdata/varexist",
			methodSignature: "SplitList",
			goImport:        "path/filepath",
			want:            true,
		},
		{
			name:            "test global variable success assign",
			testFile:        "testdata/varexist",
			methodSignature: "New",
			goImport:        "errors",
			want:            true,
		},
		{
			name:            "test invalid import",
			testFile:        "testdata/varexist",
			methodSignature: "Join",
			goImport:        "errors",
			want:            false,
		},
		{
			name:            "test invalid case sensitive assign",
			testFile:        "testdata/varexist",
			methodSignature: "join",
			goImport:        "context",
			want:            false,
		},
		{
			name:            "test invalid struct assign",
			testFile:        "testdata/varexist",
			methodSignature: "fooStruct",
			goImport:        "context",
			want:            false,
		},
		{
			name:            "test invalid method signature",
			testFile:        "testdata/varexist",
			methodSignature: "fooMethod",
			goImport:        "context",
			want:            false,
		},
		{
			name:            "test not found name",
			testFile:        "testdata/varexist",
			methodSignature: "Invalid",
			goImport:        "context",
			want:            false,
		},
		{
			name:            "test invalid assign with wrong",
			testFile:        "testdata/varexist",
			methodSignature: "invalid.New",
			goImport:        "context",
			want:            false,
		},
		{
			name:            "test invalid assign with wrong",
			testFile:        "testdata/varexist",
			methodSignature: "SplitList",
			goImport:        "path/filepath",
			want:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appPkg, _, err := xast.ParseFile(tt.testFile)
			require.NoError(t, err)

			got := goanalysis.FuncVarExists(appPkg, tt.goImport, tt.methodSignature)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFindBlankImports(t *testing.T) {
	tests := []struct {
		name     string
		testfile string
		want     []string
	}{
		{
			name:     "test a declaration inside a method success",
			testfile: "testdata/varexist",
			want:     []string{"embed", "mvdan.cc/gofumpt"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appPkg, _, err := xast.ParseFile(tt.testfile)
			require.NoError(t, err)

			got := goanalysis.FindBlankImports(appPkg)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFormatImports(t *testing.T) {
	tests := []struct {
		name  string
		input *ast.File
		want  map[string]string
	}{
		{
			name: "Test one import",
			input: &ast.File{
				Imports: []*ast.ImportSpec{
					{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"fmt\"",
						},
					},
				},
			},
			want: map[string]string{
				"fmt": "fmt",
			},
		},
		{
			name: "Test underscore import",
			input: &ast.File{
				Imports: []*ast.ImportSpec{
					{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"net/http\"",
						},
					},
					{
						Name: &ast.Ident{
							Name: "_",
						},
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"github.com/example/pkg\"",
						},
					},
				},
			},
			want: map[string]string{
				"http": "net/http",
				"pkg":  "github.com/example/pkg",
			},
		},
		{
			name: "Test dot import",
			input: &ast.File{
				Imports: []*ast.ImportSpec{
					{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"net/http\"",
						},
					},
					{
						Name: &ast.Ident{
							Name: ".",
						},
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"github.com/example/pkg\"",
						},
					},
					{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"fmt\"",
						},
					},
				},
			},
			want: map[string]string{
				"http": "net/http",
				"pkg":  "github.com/example/pkg",
				"fmt":  "fmt",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, goanalysis.FormatImports(tt.input))
		})
	}
}

func TestReplaceCode(t *testing.T) {
	var (
		newFunction = `package test
func NewMethod1() {
	n := "test new method"
	bla := fmt.Sprintf("test new - %s", n)
	fmt.Println(bla)
}`
		rollback = `package test
func NewMethod1() {
	foo := 100
	bar := fmt.Sprintf("test - %d", foo)
	fmt.Println(bar)
}`
	)

	type args struct {
		path            string
		oldFunctionName string
		newFunction     string
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "function fooTest",
			args: args{
				path:            "testdata",
				oldFunctionName: "fooTest",
				newFunction:     newFunction,
			},
		},
		{
			name: "function BazTest",
			args: args{
				path:            "testdata",
				oldFunctionName: "BazTest",
				newFunction:     newFunction,
			},
		},
		{
			name: "function invalidFunction",
			args: args{
				path:            "testdata",
				oldFunctionName: "invalidFunction",
				newFunction:     newFunction,
			},
		},
		{
			name: "invalid path",
			args: args{
				path:            "invalid_path",
				oldFunctionName: "invalidPath",
				newFunction:     newFunction,
			},
			err: os.ErrNotExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := goanalysis.ReplaceCode(tt.args.path, tt.args.oldFunctionName, tt.args.newFunction)
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, goanalysis.ReplaceCode(tt.args.path, tt.args.oldFunctionName, rollback))
		})
	}
}

func TestHasStructFieldsInPkg(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		structName string
		fields     []string
		err        error
		want       bool
	}{
		{
			name:       "test a value with an empty struct",
			path:       "testdata",
			structName: "emptyStruct",
			fields:     []string{"name"},
			want:       false,
		},
		{
			name:       "test no value with an empty struct",
			path:       "testdata",
			structName: "emptyStruct",
			fields:     []string{""},
			want:       false,
		},
		{
			name:       "test a valid field into single field struct",
			path:       "testdata",
			structName: "fooStruct",
			fields:     []string{"name"},
			want:       true,
		},
		{
			name:       "test a not valid field into single field struct",
			path:       "testdata",
			structName: "fooStruct",
			fields:     []string{"baz"},
			want:       false,
		},
		{
			name:       "test a not valid field into struct",
			path:       "testdata",
			structName: "bazStruct",
			fields:     []string{"baz"},
			want:       false,
		},
		{
			name:       "test a valid field into struct",
			path:       "testdata",
			structName: "bazStruct",
			fields:     []string{"name"},
			want:       true,
		},
		{
			name:       "test two valid fields into struct",
			path:       "testdata",
			structName: "bazStruct",
			fields:     []string{"name", "title"},
			want:       true,
		},
		{
			name:       "test a valid and a not valid fields into struct",
			path:       "testdata",
			structName: "bazStruct",
			fields:     []string{"foo", "title"},
			want:       true,
		},
		{
			name:       "test three not valid fields into struct",
			path:       "testdata",
			structName: "bazStruct",
			fields:     []string{"foo", "baz", "bla"},
			want:       false,
		},
		{
			name: "invalid path",
			path: "invalid_path",
			err:  os.ErrNotExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := goanalysis.HasAnyStructFieldsInPkg(tt.path, tt.structName, tt.fields)
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
