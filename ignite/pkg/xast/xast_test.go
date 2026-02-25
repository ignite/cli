package xast_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
)

func TestInspect(t *testing.T) {
	fset := token.NewFileSet()
	n, err := parser.ParseFile(fset, "testdata/inspect/test.go", nil, 0)
	require.NoError(t, err)
	var calls int
	tests := []struct {
		name          string
		f             func(n ast.Node) error
		expectedError string
		expectedCalls int
	}{
		{
			name: "random error",
			f: func(ast.Node) error {
				return errors.New("oups")
			},
			expectedError: "oups",
		},
		{
			name: "stop error",
			f: func(ast.Node) error {
				calls++
				return xast.ErrStop
			},
			expectedCalls: 1,
		},
		{
			name: "no error",
			f: func(ast.Node) error {
				calls++
				return nil
			},
			expectedCalls: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls = 0
			err = xast.Inspect(n, tt.f)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectedCalls, calls)
		})
	}
}

func TestParseDir(t *testing.T) {
	pkg, fileSet, err := xast.ParseDir("testdata/parseDir")

	require.NoError(t, err)
	require.NotNil(t, fileSet)
	require.Equal(t, "file", pkg.Name)
}

func TestParseFile(t *testing.T) {
	dir := t.TempDir()

	t.Run("parse valid file", func(t *testing.T) {
		filePath := filepath.Join(dir, "valid.go")
		content := `package sample

func hello() {}
`
		require.NoError(t, os.WriteFile(filePath, []byte(content), 0o644))

		file, fileSet, err := xast.ParseFile(filePath)
		require.NoError(t, err)
		require.NotNil(t, fileSet)
		require.Equal(t, "sample", file.Name.Name)
		require.Len(t, file.Decls, 1)
	})

	t.Run("invalid file path", func(t *testing.T) {
		_, _, err := xast.ParseFile(filepath.Join(dir, "missing.go"))
		require.Error(t, err)
	})

	t.Run("invalid go file", func(t *testing.T) {
		filePath := filepath.Join(dir, "invalid.go")
		require.NoError(t, os.WriteFile(filePath, []byte("package sample\nfunc"), 0o644))

		_, _, err := xast.ParseFile(filePath)
		require.Error(t, err)
	})
}
