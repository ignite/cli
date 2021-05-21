package xgenny

import (
	"context"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/logger"
	"github.com/tendermint/starport/starport/pkg/placeholder"
)

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
