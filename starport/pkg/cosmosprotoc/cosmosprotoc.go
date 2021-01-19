package cosmosprotoc

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-zglob"
	"github.com/otiai10/copy"
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
// make sure that all paths are absolute.
func Generate(
	ctx context.Context,
	projectPath,
	gomodPath,
	protoPath string,
	protoThirdPartyPaths []string,
) error {
	// define protoc command with proto paths(-I).
	command := []string{
		"protoc",
	}

	for _, importPath := range append([]string{protoPath}, protoThirdPartyPaths...) {
		command = append(command, "-I", importPath)
	}

	pattern := func(path string) string {
		return path + "/**/*.proto"
	}

	// get a list of proto dirs under path and run protoc for each individually to all protocOuts.
	dirs, err := xos.DirList(pattern(protoPath))
	if err != nil {
		return err
	}

	// put generated files under the tmp dir and move from here to the source code.
	// this prevents having leftover generated files in the app's source code or its parent dir in case of an interrupt.
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	for _, dir := range dirs {
		// check if the third party proto files are located under the same dir of
		// app's proto files, if so skip them since they should only be included with `-I`.
		var includesThirdParty bool
		for _, protoThirdPartyPath := range protoThirdPartyPaths {
			if strings.HasPrefix(dir, protoThirdPartyPath) {
				includesThirdParty = true
			}
		}
		if includesThirdParty {
			continue
		}

		// find out the list of proto files under the dir and code generate for them.
		files, err := zglob.Glob(pattern(dir))
		if err != nil {
			return err
		}

		for _, out := range protocOuts {
			command := append(command, out)
			command = append(command, files...)

			errb := &bytes.Buffer{}

			err := cmdrunner.
				New(cmdrunner.DefaultStderr(errb),
					cmdrunner.DefaultWorkdir(tmp)).
				Run(ctx, step.New(
					step.Exec(command[0], command[1:]...)))

			if err != nil {
				return errors.Wrap(err, errb.String())
			}
		}

	}

	// move generated files to the proper places.
	generatedPath := filepath.Join(tmp, gomodPath)
	if err := copy.Copy(generatedPath, projectPath); err != nil {
		return errors.Wrap(err, "cannot copy path")
	}

	return nil
}
