package cosmosgen

import (
	"golang.org/x/mod/modfile"
)

// DepTools necessary tools to build and run the chain.
func DepTools() []string {
	return []string{
		// buf build code generation.
		"github.com/bufbuild/buf/cmd/buf",
		"github.com/cosmos/gogoproto/protoc-gen-gocosmos",
		"github.com/cosmos/gogoproto/protoc-gen-gogo",
		"github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar",

		// Go code generation plugin.
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc",
		"google.golang.org/protobuf/cmd/protoc-gen-go",

		// grpc-gateway plugins.
		"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway",
		"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2",

		// code style
		"golang.org/x/tools/cmd/goimports",
		"github.com/golangci/golangci-lint/cmd/golangci-lint",
	}
}

// MissingTools find missing tools imports from a given go.mod.
func MissingTools(f *modfile.File) (missingTools []string) {
	imports := make(map[string]struct{})
	for _, imp := range f.Tool {
		imports[imp.Path] = struct{}{}
	}

	for _, tool := range DepTools() {
		if _, ok := imports[tool]; !ok {
			missingTools = append(missingTools, tool)
		}
	}

	return missingTools
}

// UnusedTools find unused tools imports from a given go.mod.
func UnusedTools(f *modfile.File) (unusedTools []string) {
	unused := []string{
		// regen protoc plugin
		"github.com/regen-network/cosmos-proto/protoc-gen-gocosmos",

		// old ignite repo.
		"github.com/ignite-hq/cli/ignite/pkg/cmdrunner",
		"github.com/ignite-hq/cli/ignite/pkg/cmdrunner/step",
	}

	imports := make(map[string]struct{})
	for _, imp := range f.Tool {
		imports[imp.Path] = struct{}{}
	}

	for _, tool := range unused {
		if _, ok := imports[tool]; ok {
			unusedTools = append(unusedTools, tool)
		}
	}
	return
}
