package v1

import (
	"strings"

	"github.com/spf13/cobra"
)

const igniteBinaryName = "ignite"

// Path returns the absolute command path including the binary name as prefix.
func (c *Command) Path() string {
	return ensureFullCommandPath(c.PlaceCommandUnder)
}

// ToCobraCommand returns a new Cobra command that matches the current command.
func (c *Command) ToCobraCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     c.Use,
		Aliases: c.Aliases,
		Short:   c.Short,
		Long:    c.Long,
		Hidden:  c.Hidden,
	}

	for _, f := range c.Flags {
		if err := f.exportFlags(cmd); err != nil {
			return nil, err
		}
	}

	return cmd, nil
}

// ImportFlags imports flags from a Cobra command.
func (c *ExecutedCommand) ImportFlags(cmd *cobra.Command) {
	c.Flags = extractCobraFlags(cmd)
}

func ensureFullCommandPath(path string) string {
	if !strings.HasPrefix(path, igniteBinaryName) {
		path = igniteBinaryName + " " + path
	}
	return strings.TrimSpace(path)
}
