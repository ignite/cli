package v1

import "github.com/spf13/cobra"

// CommandPath returns the absolute command path including the binary name as prefix.
func (h *Hook) CommandPath() string {
	return ensureFullCommandPath(h.PlaceHookOn)
}

// ImportFlags imports flags from a Cobra command.
func (h *Hook) ImportFlags(cmd *cobra.Command) {
	h.Flags = extractCobraFlags(cmd)
}
