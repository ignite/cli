package chain

import "context"

func (c *Chain) Simulate(ctx context.Context) error {
	commands, err := c.Commands(ctx)
	if err != nil {
		return err
	}
	return commands.Simulation(ctx, c.app.Path)
}
