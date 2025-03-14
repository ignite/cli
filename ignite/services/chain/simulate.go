package chain

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/simulation"
)

type simappOptions struct {
	simulationTestName string
	enabled            bool
	config             simulation.Config
	genesisTime        int64
}

func newSimappOptions() simappOptions {
	return simappOptions{
		config: simulation.Config{
			Commit: true,
		},
		enabled:     true,
		genesisTime: 0,
	}
}

// SimappOption provides options for the simapp command.
type SimappOption func(*simappOptions)

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

// SimappWithSimulationTestName allows to set the simulation test name.
func SimappWithSimulationTestName(name string) SimappOption {
	return func(c *simappOptions) {
		c.simulationTestName = name
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
		simappOptions.simulationTestName,
		simappOptions.enabled,
		simappOptions.config,
		simappOptions.genesisTime,
	)
}
