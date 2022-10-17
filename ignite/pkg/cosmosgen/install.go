package cosmosgen

import (
	"context"

	"github.com/ignite/cli/ignite/pkg/gocmd"
)

func DepTools() []string {
	return []string{
		// the gocosmos plugin.
		"github.com/regen-network/cosmos-proto/protoc-gen-gocosmos",

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
	// NOTE(tb): go get is not required thanks to tools.go, but for retrocompat
	// with older scaffolded chains that doesn't have this file, we need to run
	// it anyway.
	if err := gocmd.Get(ctx, appPath, DepTools()); err != nil {
		return err
	}
	return gocmd.Install(ctx, appPath, DepTools())
}
