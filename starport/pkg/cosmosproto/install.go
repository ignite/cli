package cosmosproto

import (
	"bytes"
	"context"

	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/xexec"
)

// ErrProtocNotInstalled is returned when protoc isn't installed on the system.
var ErrProtocNotInstalled = errors.New("protoc is not installed")

// InstallDependencies installs protoc dependencies needed by Cosmos ecosystem.
func InstallDependencies(ctx context.Context, appPath string) error {
	if !xexec.IsCommandAvailable("protoc") {
		return ErrProtocNotInstalled
	}

	errb := &bytes.Buffer{}
	err := cmdrunner.
		New(
			cmdrunner.DefaultStderr(errb),
			cmdrunner.DefaultWorkdir(appPath),
		).
		Run(ctx,
			// installs the gocosmos plugin with the version specified under the
			// go.mod of the app.
			step.New(
				step.Exec(
					"go",
					"get",
					"github.com/regen-network/cosmos-proto/protoc-gen-gocosmos@v0.3.1",
				),
			),
			// install grpc-gateway.
			step.New(
				step.Exec(
					"go",
					"get",
					"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0",
					"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v1.16.0",
					"github.com/golang/protobuf/protoc-gen-go@v1.4.3",
				),
			),
		)
	return errors.Wrap(err, errb.String())
}
