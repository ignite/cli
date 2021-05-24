package cosmosgen

import (
	"bytes"
	"context"

	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

// InstallDependencies installs protoc dependencies needed by Cosmos ecosystem.
func InstallDependencies(ctx context.Context, appPath string) error {
	errb := &bytes.Buffer{}
	err := cmdrunner.
		New(
			cmdrunner.DefaultStderr(errb),
			cmdrunner.DefaultWorkdir(appPath),
		).
		Run(ctx,
			step.New(
				step.Exec(
					"go",
					"get",
					// installs the gocosmos plugin.
					"github.com/regen-network/cosmos-proto/protoc-gen-gocosmos",

					// install Go code generation plugin.
					"github.com/golang/protobuf/protoc-gen-go",

					// install grpc-gateway plugins.
					"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway",
					"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger",
					"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2",
				),
			),
		)
	return errors.Wrap(err, errb.String())
}
