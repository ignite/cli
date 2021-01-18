package cosmosprotoc

import (
	"bytes"
	"context"

	"github.com/mattn/go-zglob"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/xexec"
	"github.com/tendermint/starport/starport/pkg/xos"
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
					"github.com/regen-network/cosmos-proto/protoc-gen-gocosmos",
				),
			),
			// install grpc-gateway.
			step.New(
				step.Exec(
					"go",
					"install",
					"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway",
					"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger",
					"github.com/golang/protobuf/protoc-gen-go",
				),
			),
		)
	return errors.Wrap(err, errb.String())
}

var (
	protocOuts = []string{
		"--gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:.",
		"--grpc-gateway_out=logtostderr=true:.",
	}
)

// Generate generates source code from proto files residing under dirs.
func Generate(ctx context.Context, path string, importPaths ...string) error {
	// define protoc command with proto paths(-I).
	command := []string{
		"protoc",
	}

	for _, importPath := range append([]string{path}, importPaths...) {
		command = append(command, "-I", importPath)
	}

	pattern := func(path string) string {
		return path + "/**/*.proto"
	}

	// get a list of proto dirs under path and run protoc for each individually to all protocOuts.
	dirs, err := xos.DirList(pattern(path))
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		files, err := zglob.Glob(pattern(dir))
		if err != nil {
			return err
		}

		for _, out := range protocOuts {
			command := append(command, out)
			command = append(command, files...)

			errb := &bytes.Buffer{}

			err := cmdrunner.
				New(
					cmdrunner.DefaultStderr(errb),
				).
				Run(ctx,
					step.New(
						step.Exec(command[0], command[1:]...),
					),
				)

			if err != nil {
				return errors.Wrap(err, errb.String())
			}
		}

	}

	return nil
}
