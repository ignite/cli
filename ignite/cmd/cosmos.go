package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCosmos returns the `ignite cosmos` command tree holding the Cosmos SDK
// tooling. gno.land is the primary target of Ignite.
func NewCosmos() *cobra.Command {
	c := &cobra.Command{
		Use:   "cosmos [command]",
		Short: "Cosmos SDK blockchain tooling",
		Long: `Cosmos SDK blockchain tooling.

Ignite's focus is gno.land smart contracts (see the commands at the root of
` + "`ignite --help`" + `), while the Cosmos SDK commands for sovereign
blockchains live here.`,
		Args: cobra.ExactArgs(0),
	}

	c.AddCommand(
		NewScaffold(),
		NewChain(),
		NewGenerate(),
		NewAccount(),
		NewTestnet(),
		NewDoctor(),
	)

	return c
}

// cosmosLegacyStubs returns hidden stub commands for the previously
// top-level cosmos commands, pointing users at `ignite cosmos`.
func cosmosLegacyStubs() []*cobra.Command {
	stub := func(use, useLine, instead string) *cobra.Command {
		return &cobra.Command{
			Use:        use,
			Hidden:     true,
			Deprecated: fmt.Sprintf("use `ignite cosmos %s` instead.", instead),
			RunE: func(cmd *cobra.Command, args []string) error {
				return fmt.Errorf("use `ignite cosmos %s` instead", instead)
			},
		}
	}

	return []*cobra.Command{
		stub("chain [command]", "chain [flags]", "chain"),
		stub("generate [command]", "generate [flags]", "generate"),
		stub("testnet [command]", "testnet [flags]", "testnet"),
	}
}
