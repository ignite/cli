package xgenny

import (
	"context"
	"errors"
	"os"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/logger"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/validation"
)

var _ validation.Error = (*dryRunError)(nil)

type dryRunError struct {
	error
}

func (d *dryRunError) ValidationInfo() string {
	return d.Error()
}

func DryRunner(ctx context.Context) *genny.Runner {
	runner := genny.DryRunner(context.Background())
	runner.Logger = logger.New(genny.DefaultLogLvl)
	return runner
}

func RunWithValidation(tracer *placeholder.Tracer, gens ...*genny.Generator) error {
	run := func(runner *genny.Runner, gen *genny.Generator) error {
		runner.With(gen)
		return runner.Run()
	}
	for _, gen := range gens {
		if err := run(DryRunner(context.Background()), gen); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return &dryRunError{err}
			}
			return err
		}
		if err := tracer.Err(); err != nil {
			return err
		}
		if err := run(genny.WetRunner(context.Background()), gen); err != nil {
			return err
		}
	}
	return nil
}
