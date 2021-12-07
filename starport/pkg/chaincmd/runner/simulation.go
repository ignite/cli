package chaincmdrunner

import (
	"context"
)

// Simulation run the chain simulation.
func (r Runner) Simulation(ctx context.Context, home string) error {
	return r.run(ctx, runOptions{}, r.chainCmd.SimulationCommand(home))
}
