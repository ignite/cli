package xgenny

import (
	"context"
	"errors"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/logger"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/validation"
	"os"
)

var _ validation.Error = (*dryRunError)(nil)

type dryRunError struct {
	error
}

// ValidationInfo returns validation info
func (d *dryRunError) ValidationInfo() string {
	return d.Error()
}

// DryRunner is a genny DryRunner with a logger
func DryRunner(ctx context.Context) *genny.Runner {
	runner := genny.DryRunner(context.Background())
	runner.Logger = logger.New(genny.DefaultLogLvl)
	return runner
}

// RunWithValidation checks the generators with a dry run and then execute the wet runner to the generators
func RunWithValidation(
	tracer *placeholder.Tracer,
	gens ...*genny.Generator,
) (sm SourceModification, err error) {
	// run executes the provided runner with the provided generator
	run := func(runner *genny.Runner, gen *genny.Generator) error {
		runner.With(gen)
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
			_, err := os.Stat(file.Name())
			if os.IsNotExist(err) {
				// if the file doesn't exist in the source, it means it has been created by the runner
				sm.AppendCreatedFiles(file.Name())
			} else if err != nil {
				return sm, err
			} else {
				// the file has been modified by the runner
				sm.AppendModifiedFiles(file.Name())
			}
		}

		// execute the modification with a wet runner
		if err := run(genny.WetRunner(context.Background()), gen); err != nil {
			return sm, err
		}
	}
	return sm, nil
}
