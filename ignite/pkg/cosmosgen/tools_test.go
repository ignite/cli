package cosmosgen

import (
	"go/ast"
	"go/token"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMissingTools(t *testing.T) {
	tests := []struct {
		name    string
		astFile *ast.File
		want    []string
	}{
		{
			name:    "no missing tools",
			astFile: createASTFileWithImports(DepTools()...),
			want:    nil,
		},
		{
			name:    "some missing tools",
			astFile: createASTFileWithImports("github.com/golang/protobuf/protoc-gen-go", "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway", "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"),
			want:    []string{"github.com/cosmos/gogoproto/protoc-gen-gocosmos", "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"},
		},
		{
			name:    "all tools missing",
			astFile: createASTFileWithImports(),
			want:    DepTools(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MissingTools(tt.astFile)
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
