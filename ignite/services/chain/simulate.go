package chain

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/simulation"
)

type simappOptions struct {
	enabled     bool
	config      simulation.Config
	period      uint
	genesisTime int64
}

func newSimappOptions() simappOptions {
	return simappOptions{
		config: simulation.Config{
			Commit: true,
		},
		enabled:     true,
		period:      0,
		genesisTime: 0,
	}
}

// SimappOption provides options for the simapp command.
type SimappOption func(*simappOptions)

// SimappWithPeriod allows running slow invariants only once every period assertions.
func SimappWithPeriod(period uint) SimappOption {
	return func(c *simappOptions) {
		c.period = period
	}
}

// SimappWithGenesisTime allows overriding genesis UNIX time instead of using a random UNIX time.
func SimappWithGenesisTime(genesisTime int64) SimappOption {
	return func(c *simappOptions) {
		c.genesisTime = genesisTime
	}
}

// SimappWithConfig allows to add a simulation config.
func SimappWithConfig(config simulation.Config) SimappOption {
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
	return commands.Simulation(ctx,
		c.app.Path,
		simappOptions.enabled,
		simappOptions.config,
		simappOptions.period,
		simappOptions.genesisTime,
	)
}
