package ignitecmd

import (
	"github.com/spf13/cobra"
)

// NewScaffoldReact scaffolds a React app for a chain.
func NewScaffoldReact() *cobra.Command {
	c := &cobra.Command{
		Use:        "react",
		Deprecated: "the React scaffolding feature is removed from Ignite CLI.\nPlease use the Ignite CCA app to create a React app.\nFor more information, visit: https://ignite.com/marketplace/CCA",
	}

	return c
}
