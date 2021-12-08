package chain

import (
	"context"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types/simulation"
)

type simappOptions struct {
	config      simulation.Config
	enabled     bool
	verbose     bool
	period      uint
	genesisTime int64
}

func newSimappOptions() simappOptions {
	return simappOptions{
		config:      simapp.NewConfigFromFlags(),
		enabled:     true,
		verbose:     false,
		period:      0,
		genesisTime: 0,
	}
}

// SimappOption provides options for the simapp command
type SimappOption func(*simappOptions)

// WithVerbose enable the verbose mode
func WithVerbose(verbose bool) SimappOption {
	return func(c *simappOptions) {
		c.verbose = verbose
	}
}

// WithPeriod allows running slow invariants only once every period assertions
func WithPeriod(period uint) SimappOption {
	return func(c *simappOptions) {
		c.period = period
	}
}

// WithGenesisTime allows overriding genesis UNIX time instead of using a random UNIX time
func WithGenesisTime(genesisTime int64) SimappOption {
	return func(c *simappOptions) {
		c.genesisTime = genesisTime
	}
}

// WithConfig allows to add a simulation config
func WithConfig(config simulation.Config) SimappOption {
	return func(c *simappOptions) {
		c.config = config
	}
}

func (c *Chain) Simulate(ctx context.Context, options ...SimappOption) error {
	simappOptions := newSimappOptions()

	// apply the options
	for _, apply := range options {
		apply(&simappOptions)
	}

	commands, err := c.Commands(ctx)
	if err != nil {
		return err
	}
	return commands.Simulation(ctx, c.app.Path)
}
