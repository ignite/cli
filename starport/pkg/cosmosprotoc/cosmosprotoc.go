package cosmosprotoc

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/gomodule"
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
	"github.com/tendermint/starport/starport/pkg/protopath"
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

var (
	protocOuts = []string{
		"--gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:.",
		"--grpc-gateway_out=logtostderr=true:.",
	}
)

// Generate generates code from proto app's proto files.
// make sure that all paths are absolute.
func Generate(
	ctx context.Context,
	projectPath,
	gomodPath,
	protoPath string,
	protoThirdPartyPaths []string,
) error {
	// Cosmos SDK hosts proto files of own x/ modules and some third party ones needed by itself and
	// blockchain apps. Generate should be aware of these and make them available to the blockchain
	// app that wants to generate code for its own proto.
	//
	// blockchain apps may use different versions of the SDK. following code first makes sure that
	// app's dependencies are download by 'go mod' and cached under the local filesystem.
	// and then, it determines which version of the SDK is used by the app and what is the absolute path
	// of its source code.
	if err := cmdrunner.
		New(cmdrunner.DefaultWorkdir(projectPath)).
		Run(ctx, step.New(step.Exec("go", "mod", "download"))); err != nil {
		return err
	}

	modfile, err := gomodule.ParseAt(projectPath)
	if err != nil {
		return err
	}

	// add Google's and SDK's proto paths to third parties list.
	resolved, err := protopath.ResolveDependencyPaths(modfile.Require,
		protopath.NewModule("github.com/cosmos/cosmos-sdk", "proto", "third_party/proto"),
	)
	if err != nil {
		return err
	}

	protoThirdPartyPaths = append(protoThirdPartyPaths, resolved...)

	// created a temporary dir to locate generated code under which later only some of them will be moved to the
	// app's source code. this also prevents having leftover files in the app's source code or its parent dir -when
	// command executed directly there- in case of an interrupt.
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmp)

	// start preparing the protoc command for execution.
	command := []string{
		"protoc",
	}

	// append third party proto locations to the command.
	for _, importPath := range append([]string{protoPath}, protoThirdPartyPaths...) {
		// skip if a third party proto source actually doesn't exist on the filesystem.
		if _, err := os.Stat(importPath); os.IsNotExist(err) {
			continue
		}
		command = append(command, "-I", importPath)
	}

	// find out the list of proto files under the app and generate code for them.
	files, err := protoanalysis.SearchProto(protoPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		// check if the file belongs to a third party proto. if so, skip it since it should
		// only be included via `-I`.
		var includesThirdParty bool
		for _, protoThirdPartyPath := range protoThirdPartyPaths {
			if strings.HasPrefix(file, protoThirdPartyPath) {
				includesThirdParty = true
				break
			}
		}
		if includesThirdParty {
			continue
		}

		// run command for each protocOuts.
		for _, out := range protocOuts {
			command := append(command, out)
			command = append(command, file)

			errb := &bytes.Buffer{}

			err := cmdrunner.
				New(
					cmdrunner.DefaultStderr(errb),
					cmdrunner.DefaultWorkdir(tmp)).
				Run(ctx,
					step.New(step.Exec(command[0], command[1:]...)))

			if err != nil {
				return errors.Wrap(err, errb.String())
			}
		}
	}

	// move generated code for the app under the relative locations in its source code.
	generatedPath := filepath.Join(tmp, gomodPath)
	if err := copy.Copy(generatedPath, projectPath); err != nil {
		return errors.Wrap(err, "cannot copy path")
	}

	return nil
}
