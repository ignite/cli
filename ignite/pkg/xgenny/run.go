package xgenny

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/packd"

	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/validation"
)

var _ validation.Error = (*dryRunError)(nil)

type dryRunError struct {
	error
}

// ValidationInfo returns validation info.
func (d *dryRunError) ValidationInfo() string {
	return d.Error()
}

// DryRunner is a genny DryRunner with a logger.
func DryRunner(ctx context.Context) *genny.Runner {
	runner := genny.DryRunner(ctx)
	runner.Logger = logger.New(genny.DefaultLogLvl)
	return runner
}

// RunWithValidation checks the generators with a dry run and then execute the wet runner to the generators.
func RunWithValidation(
	tracer *placeholder.Tracer,
	gens ...*genny.Generator,
) (sm SourceModification, err error) {
	// run executes the provided runner with the provided generator
	run := func(runner *genny.Runner, gen *genny.Generator) error {
		err := runner.With(gen)
		if err != nil {
			return err
		}
		return runner.Run()
	}
	for _, gen := range gens {
		// check with a dry runner the generators
		dryRunner := DryRunner(context.Background())
		if err := run(dryRunner, gen); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return sm, &dryRunError{err}
			}
			return sm, err
		}
		if err := tracer.Err(); err != nil {
			return sm, err
		}

		// fetch the source modification
		sm = NewSourceModification()
		for _, file := range dryRunner.Results().Files {
			fileName := file.Name()
			_, err := os.Stat(fileName)

			//nolint:gocritic
			if os.IsNotExist(err) {
				// if the file doesn't exist in the source, it means it has been created by the runner
				sm.AppendCreatedFiles(fileName)
			} else if err != nil {
				return sm, err
			} else {
				// the file has been modified by the runner
				sm.AppendModifiedFiles(fileName)
			}
		}

		// execute the modification with a wet runner
		if err := run(genny.WetRunner(context.Background()), gen); err != nil {
			return sm, err
		}
	}
	return sm, nil
}

// Box will mount each file in the Box and wrap it, already existing files are ignored.
func Box(g *genny.Generator, box packd.Walker) error {
	return box.Walk(func(path string, bf packd.File) error {
		f := genny.NewFile(path, bf)
		f, err := g.Transform(f)
		if err != nil {
			return err
		}
		filePath := strings.TrimSuffix(f.Name(), ".plush")
		_, err = os.Stat(filePath)
		if os.IsNotExist(err) {
			// path doesn't exist. move on.
			g.File(f)
			return nil
		}
		return err
	})
}
