package cosmosproto

import (
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
	"github.com/tendermint/starport/starport/pkg/nodetime/protobufjs"
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
	"github.com/tendermint/starport/starport/pkg/protoc"
	"github.com/tendermint/starport/starport/pkg/protopath"
	"golang.org/x/mod/modfile"
)

var (
	protocOuts = []string{
		"--gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:.",
		"--grpc-gateway_out=logtostderr=true:.",
	}

	sdkImport          = "github.com/cosmos/cosmos-sdk"
	sdkProto           = "proto"
	sdkProtoThirdParty = "third_party/proto"

	fileTypes = "types"
)

type generateOptions struct {
	gomodPath string
	jsOut     func(protoanalysis.Package, string) string
}

// Target adds a new code generation target to Generate.
type Target func(*generateOptions)

// WithJSGeneration adds JS code generation.
func WithJSGeneration(out func(pkg protoanalysis.Package, moduleName string) (path string)) Target {
	return func(o *generateOptions) {
		o.jsOut = out
	}
}

// WithGoGeneration adds Go code generation.
func WithGoGeneration(gomodPath string) Target {
	return func(o *generateOptions) {
		o.gomodPath = gomodPath
	}
}

// generator generates code for sdk and sdk apps.
type generator struct {
	ctx          context.Context
	projectPath  string
	protoPath    string
	includePaths []string
	o            *generateOptions
	modfile      *modfile.File
}

// Generate generates code from proto app's proto files.
// make sure that all paths are absolute.
func Generate(
	ctx context.Context,
	projectPath,
	protoPath string,
	includePaths []string,
	target Target,
	otherTargets ...Target,
) error {
	g := &generator{
		ctx:          ctx,
		projectPath:  projectPath,
		protoPath:    protoPath,
		includePaths: includePaths,
		o:            &generateOptions{},
	}

	for _, target := range append(otherTargets, target) {
		target(g.o)
	}

	if err := g.setup(); err != nil {
		return err
	}

	if g.o.jsOut != nil {
		if err := g.generateJS(); err != nil {
			return err
		}
	}

	if g.o.gomodPath != "" {
		if err := g.generateGo(); err != nil {
			return err
		}
	}

	return nil
}

func (g *generator) setup() (err error) {
	// Cosmos SDK hosts proto files of own x/ modules and some third party ones needed by itself and
	// blockchain apps. Generate should be aware of these and make them available to the blockchain
	// app that wants to generate code for its own proto.
	//
	// blockchain apps may use different versions of the SDK. following code first makes sure that
	// app's dependencies are download by 'go mod' and cached under the local filesystem.
	// and then, it determines which version of the SDK is used by the app and what is the absolute path
	// of its source code.
	if err := cmdrunner.
		New(cmdrunner.DefaultWorkdir(g.projectPath)).
		Run(g.ctx, step.New(step.Exec("go", "mod", "download"))); err != nil {
		return err
	}

	// parse the go.mod of the app.
	g.modfile, err = gomodule.ParseAt(g.projectPath)

	return
}

func (g *generator) generateGo() error {
	includePaths, err := g.resolveInclude(protopath.NewModule(sdkImport, sdkProto, sdkProtoThirdParty))
	if err != nil {
		return err
	}

	// created a temporary dir to locate generated code under which later only some of them will be moved to the
	// app's source code. this also prevents having leftover files in the app's source code or its parent dir -when
	// command executed directly there- in case of an interrupt.
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	if err := protoc.Generate(g.ctx, tmp, g.protoPath, append(g.includePaths, includePaths...), protocOuts); err != nil {
		return err
	}

	// move generated code for the app under the relative locations in its source code.
	generatedPath := filepath.Join(tmp, g.o.gomodPath)
	err = copy.Copy(generatedPath, g.projectPath)
	return errors.Wrap(err, "cannot copy path")
}

func (g *generator) generateJS() error {
	includePaths, err := g.resolveInclude(protopath.NewModule(sdkImport, sdkProto))
	if err != nil {
		return err
	}

	pkgs, err := protoanalysis.DiscoverPackages(g.protoPath)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		msp := strings.Split(pkg.Name, ".")
		moduleName := msp[len(msp)-1]

		if err := protobufjs.Generate(
			g.ctx,
			g.o.jsOut(pkg, moduleName),
			fileTypes,
			pkg.Path,
			append(g.includePaths, includePaths...),
		); err != nil {
			return err
		}
	}

	return nil
}

func (g *generator) resolveInclude(modules ...protopath.Module) (paths []string, err error) {
	return protopath.ResolveDependencyPaths(g.modfile.Require, modules...)
}
