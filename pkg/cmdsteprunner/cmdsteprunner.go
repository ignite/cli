package cmdsteprunner

import (
	"context"
)

type Runner struct {
	o runnerOptions
}

type RunnerOption func(*runnerOptions)

type runnerOptions struct {
}

func New(options ...RunnerOption) *Runner {
	var opts runnerOptions
	for _, o := range options {
		o(&opts)
	}
	r := &Runner{
		o: opts,
	}
	return r
}

func (r *Runner) Run(ctx context.Context, steps ...Step) error {
	return nil
}
