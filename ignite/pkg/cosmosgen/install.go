package cosmosgen

import (
	"context"
	"errors"
	"go/ast"

	"github.com/ignite/cli/ignite/pkg/goanalysis"
	"github.com/ignite/cli/ignite/pkg/gocmd"
)

// DepTools necessary tools to build and run the chain.
func DepTools() []string {
	return []string{
		// buf build code generation.
		"github.com/bufbuild/buf/cmd/buf",
		"github.com/cosmos/gogoproto/protoc-gen-gocosmos",

		// Go code generation plugin.
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc",
		"google.golang.org/protobuf/cmd/protoc-gen-go",
		"github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar",

		// grpc-gateway plugins.
		"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway",
		"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2",
	}
}

// InstallDepTools installs protoc dependencies needed by Cosmos ecosystem.
func InstallDepTools(ctx context.Context, appPath string) error {
	if err := gocmd.ModTidy(ctx, appPath); err != nil {
		return err
	}
	err := gocmd.Install(ctx, appPath, DepTools())
	if gocmd.IsInstallError(err) {
		return errors.New("unable to install dependency tools, run `ignite doctor` and try again")
	}
	return err
}

// MissingTools find missing tools import indo a *ast.File.
func MissingTools(f *ast.File) (missingTools []string) {
	imports := make(map[string]string)
	for name, imp := range goanalysis.FormatImports(f) {
		imports[imp] = name
	}

	for _, tool := range DepTools() {
		if _, ok := imports[tool]; !ok {
			missingTools = append(missingTools, tool)
		}
	}
	return
}

// UnusedTools find unused tools import indo a *ast.File.
func UnusedTools(f *ast.File) (unusedTools []string) {
	unused := []string{
		// regen protoc plugin
		"github.com/regen-network/cosmos-proto/protoc-gen-gocosmos",

		// old ignite repo.
		"github.com/ignite-hq/cli/ignite/pkg/cmdrunner",
		"github.com/ignite-hq/cli/ignite/pkg/cmdrunner/step",
	}

	imports := make(map[string]string)
	for name, imp := range goanalysis.FormatImports(f) {
		imports[imp] = name
	}

	for _, tool := range unused {
		if _, ok := imports[tool]; ok {
			unusedTools = append(unusedTools, tool)
		}
	}
	return
}
