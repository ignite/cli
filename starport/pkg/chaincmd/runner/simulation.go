package chaincmdrunner

import (
	"context"
	"os"
)

// Simulation run the chain simulation.
func (r Runner) Simulation(ctx context.Context, home string) error {
	return r.run(ctx, runOptions{stdout: os.Stdout}, r.chainCmd.SimulationCommand(home))
}
