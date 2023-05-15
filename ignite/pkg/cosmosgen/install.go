package cosmosgen

import (
	"context"
	"errors"

	"github.com/ignite/cli/ignite/pkg/gocmd"
)

// UnusedTools deprecated and not necessary tools.
func UnusedTools() []string {
	return []string{
		// regen protoc plugin
		"github.com/regen-network/cosmos-proto/protoc-gen-gocosmos",

		// old ignite repo.
		"github.com/ignite-hq/cli/ignite/pkg/cmdrunner",
		"github.com/ignite-hq/cli/ignite/pkg/cmdrunner/step",
	}
}

// DepTools necessary tools to build and run the chain.
func DepTools() []string {
	return []string{
		// the gocosmos plugin.
		"github.com/cosmos/gogoproto/protoc-gen-gocosmos",

		// Go code generation plugin.
		"github.com/golang/protobuf/protoc-gen-go",

		// grpc-gateway plugins.
		"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway",
		"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger",
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
		return errors.New("unable to install dependency tools, try to run `ignite doctor` and try again")
	}
	return err
}
