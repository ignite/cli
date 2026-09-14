package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCosmos returns the deprecated `ignite cosmos` command tree holding the
// legacy Cosmos SDK tooling. gno.land is now the primary target of Ignite.
func NewCosmos() *cobra.Command {
	c := &cobra.Command{
		Use:   "cosmos [command]",
		Short: "DEPRECATED: Cosmos SDK blockchain tooling",
		Long: `DEPRECATED: the Cosmos SDK tooling has moved under this command.

Ignite is now the developer experience for gno.land smart contracts. The
legacy Cosmos SDK commands remain available here for existing projects but
receive no further feature work. See the gno.land commands at the root of
` + "`ignite --help`" + `.`,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cmd.Printf("⚠️  The Cosmos SDK tooling is deprecated. Migrate your project to gno.land (https://gno.land).\n\n")
			return nil
		},
		Args: cobra.ExactArgs(0),
	}

	c.AddCommand(
		NewScaffold(),
		NewChain(),
		NewGenerate(),
		NewAccount(),
		NewTestnet(),
	)

	return c
}

// cosmosDeprecatedStubs returns hidden deprecated stub commands for the
// previously top-level cosmos commands, pointing users at `ignite cosmos`.
func cosmosDeprecatedStubs() []*cobra.Command {
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
