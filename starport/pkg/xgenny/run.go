package xgenny

import (
	"context"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/placeholder"
)

// Run provide generator with validation.
func Run(ctx context.Context, gen *genny.Generator) error {
	run := func(runner *genny.Runner, gen *genny.Generator) error {
		if err := runner.With(gen); err != nil {
			return err
		}
		return runner.Run()
	}
	if err := run(genny.DryRunner(ctx), gen); err != nil {
		return err
	}
	if err := placeholder.Validate(ctx); err != nil {
		return err
	}
	return run(genny.WetRunner(ctx), gen)
}
