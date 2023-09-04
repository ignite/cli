package cosmosgen

import (
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
	"github.com/pkg/errors"
)

func (g *generator) gogoTemplate() string {
	return filepath.Join(g.appPath, g.protoDir, "buf.gen.gogo.yaml")
}

func (g *generator) pulsarTemplate() string {
	return filepath.Join(g.appPath, g.protoDir, "buf.gen.pulsar.yaml")
}

func (g *generator) generateGo() error {
	// create a temporary dir to locate generated code under which later only some of them will be moved to the
	// app's source code. this also prevents having leftover files in the app's source code or its parent dir - when
	// command executed directly there - in case of an interrupt.
	tmp, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	protoPath := filepath.Join(g.appPath, g.protoDir)

	// code generate for each module.
	if err := g.buf.Generate(
		g.ctx,
		protoPath,
		tmp,
		g.gogoTemplate(),
		"module.proto",
	); err != nil {
		return err
	}

	// move generated code for the app under the relative locations in its source code.
	generatedPath := filepath.Join(tmp, g.gomodPath)

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

func (g *generator) generatePulsar() error {
	// create a temporary dir to locate generated code under which later only some of them will be moved to the
	// app's source code. this also prevents having leftover files in the app's source code or its parent dir - when
	// command executed directly there - in case of an interrupt.
	tmp, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	protoPath := filepath.Join(g.appPath, g.protoDir)

	// code generate for each module.
	if err := g.buf.Generate(
		g.ctx,
		protoPath,
		tmp,
		g.pulsarTemplate(),
	); err != nil {
		return err
	}

	// move generated code for the app under the relative locations in its source code.
	_, err = os.Stat(tmp)
	if err == nil {
		err = copy.Copy(tmp, g.appPath)
		if err != nil {
			return errors.Wrap(err, "cannot copy path")
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	return nil
}
