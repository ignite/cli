package xast_test

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/xast"
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
			f: func(n ast.Node) error {
				return errors.New("oups")
			},
			expectedError: "oups",
		},
		{
			name: "stop error",
			f: func(n ast.Node) error {
				calls++
				return xast.ErrStop
			},
			expectedCalls: 1,
		},
		{
			name: "no error",
			f: func(n ast.Node) error {
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
