package cosmosgen

import (
	"context"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

func (g *generator) gogoTemplate() string {
	return filepath.Join(g.appPath, g.protoDir, "buf.gen.gogo.yaml")
}

func (g *generator) protoPath() string {
	return filepath.Join(g.appPath, g.protoDir)
}

func (g *generator) generateGoGo(ctx context.Context) error {
	// create a temporary dir to locate generated code under which later only some of them will be moved to the
	// app's source code. this also prevents having leftover files in the app's source code or its parent dir - when
	// command executed directly there - in case of an interrupt.
	tmp, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	// code generate for each module.
	if err := g.buf.Generate(
		ctx,
		g.protoPath(),
		tmp,
		g.gogoTemplate(),
	); err != nil {
		return err
	}

	// move generated code for the app under the relative locations in its source code.
	path := filepath.Join(tmp, g.goModPath)
	if _, err := os.Stat(path); err == nil {
		err = copy.Copy(path, g.appPath)
		if err != nil {
			return errors.Wrap(err, "cannot copy path")
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	return nil
}
