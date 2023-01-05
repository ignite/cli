package cosmosgen

import (
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
	"github.com/pkg/errors"

	"github.com/ignite/cli/ignite/pkg/protoanalysis"
	"github.com/ignite/cli/ignite/pkg/protoc"
)

var goOuts = []string{
	"--gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:.",
	"--grpc-gateway_out=logtostderr=true:.",
}

func (g *generator) generateGo() error {
	includePaths, err := g.resolveInclude(g.appPath)
	if err != nil {
		return err
	}

	// create a temporary dir to locate generated code under which later only some of them will be moved to the
	// app's source code. this also prevents having leftover files in the app's source code or its parent dir - when
	// command executed directly there - in case of an interrupt.
	tmp, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	// discover proto packages in the app.
	pp := filepath.Join(g.appPath, g.protoDir)
	pkgs, err := protoanalysis.Parse(g.ctx, nil, pp)
	if err != nil {
		return err
	}

	// code generate for each module.
	for _, pkg := range pkgs {
		if err := protoc.Generate(g.ctx, tmp, pkg.Path, includePaths, goOuts); err != nil {
			return err
		}
	}

	// move generated code for the app under the relative locations in its source code.
	generatedPath := filepath.Join(tmp, g.o.gomodPath)

	_, err = os.Stat(generatedPath)
	if err == nil {
		err = copy.Copy(generatedPath, g.appPath)
		if err != nil {
			return errors.Wrap(err, "cannot copy path")
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	return nil
}
