package v1

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
		var fs *pflag.FlagSet
		if f.Persistent {
			fs = cmd.PersistentFlags()
		} else {
			fs = cmd.Flags()
		}

		if err := f.exportToFlagSet(fs); err != nil {
			return nil, err
		}
	}

	return cmd, nil
}

// ImportFlags imports flags from a Cobra command.
func (c *ExecutedCommand) ImportFlags(cmd *cobra.Command) {
	c.Flags = extractCobraFlags(cmd)
}

// NewFlags creates a new flags set initialized with the executed command's flags.
func (c *ExecutedCommand) NewFlags() (*pflag.FlagSet, error) {
	fs := pflag.NewFlagSet(igniteBinaryName, pflag.ContinueOnError)

	for _, f := range c.Flags {
		if f.Persistent {
			continue
		}

		if err := f.exportToFlagSet(fs); err != nil {
			return nil, err
		}
	}

	return fs, nil
}

// NewPersistentFlags creates a new flags set initialized with the executed command's persistent flags.
func (c *ExecutedCommand) NewPersistentFlags() (*pflag.FlagSet, error) {
	fs := pflag.NewFlagSet(igniteBinaryName, pflag.ContinueOnError)

	for _, f := range c.Flags {
		if !f.Persistent {
			continue
		}

		if err := f.exportToFlagSet(fs); err != nil {
			return nil, err
		}
	}

	return fs, nil
}

func ensureFullCommandPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return igniteBinaryName
	}

	if !strings.HasPrefix(path, igniteBinaryName) {
		path = igniteBinaryName + " " + path
	}
	return path
}
