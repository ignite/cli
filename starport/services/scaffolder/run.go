package scaffolder

import (
	"context"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/placeholder"
)

func runWithValidation(tracer *placeholder.Tracer, gens ...*genny.Generator) error {
	run := func(runner *genny.Runner, gen *genny.Generator) error {
		runner.With(gen)
		return runner.Run()
	}
	for _, gen := range gens {
		if err := run(genny.DryRunner(context.Background()), gen); err != nil {
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
