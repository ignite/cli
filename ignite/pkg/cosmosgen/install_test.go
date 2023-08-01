package cosmosgen_test

import (
	"go/ast"
	"go/token"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cosmosgen"
)

func TestMissingTools(t *testing.T) {
	var (
		tools        = cosmosgen.DepTools()
		someTools    = tools[:2]
		missingTools = tools[2:]
	)
	tests := []struct {
		name    string
		astFile *ast.File
		want    []string
	}{
		{
			name:    "no missing tools",
			astFile: createASTFileWithImports(tools...),
			want:    nil,
		},
		{
			name:    "some missing tools",
			astFile: createASTFileWithImports(someTools...),
			want:    missingTools,
		},
		{
			name:    "all tools missing",
			astFile: createASTFileWithImports(),
			want:    tools,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cosmosgen.MissingTools(tt.astFile)
			require.EqualValues(t, tt.want, got)
		})
	}
}

func TestUnusedTools(t *testing.T) {
	tests := []struct {
		name    string
		astFile *ast.File
		want    []string
	}{
		{
			name: "all unused tools",
			astFile: createASTFileWithImports(
				"fmt",
				"github.com/regen-network/cosmos-proto/protoc-gen-gocosmos",
				"github.com/ignite-hq/cli/ignite/pkg/cmdrunner",
				"github.com/ignite-hq/cli/ignite/pkg/cmdrunner/step",
			),
			want: []string{
				"github.com/regen-network/cosmos-proto/protoc-gen-gocosmos",
				"github.com/ignite-hq/cli/ignite/pkg/cmdrunner",
				"github.com/ignite-hq/cli/ignite/pkg/cmdrunner/step",
			},
		},
		{
			name: "some unused tools",
			astFile: createASTFileWithImports(
				"fmt",
				"github.com/ignite-hq/cli/ignite/pkg/cmdrunner",
			),
			want: []string{"github.com/ignite-hq/cli/ignite/pkg/cmdrunner"},
		},
		{
			name:    "no tools unused",
			astFile: createASTFileWithImports("fmt"),
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cosmosgen.UnusedTools(tt.astFile)
			require.EqualValues(t, tt.want, got)
		})
	}
}

// createASTFileWithImports helper function to create an AST file with given imports.
func createASTFileWithImports(imports ...string) *ast.File {
	f := &ast.File{Imports: make([]*ast.ImportSpec, len(imports))}
	for i, imp := range imports {
		f.Imports[i] = &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote(imp),
			},
		}
	}
	return f
}
