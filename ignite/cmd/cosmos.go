package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCosmos returns the `ignite cosmos` command tree holding the legacy
// Cosmos SDK tooling, now in maintenance mode. gno.land is the primary
// target of Ignite.
func NewCosmos() *cobra.Command {
	c := &cobra.Command{
		Use:   "cosmos [command]",
		Short: "Cosmos SDK blockchain tooling (maintenance mode)",
		Long: `The Cosmos SDK tooling lives under this command and is in maintenance mode.

Ignite is now the tool for building gno.land smart contracts. The Cosmos SDK
commands remain fully functional here for existing projects: they receive
bug fixes and compatibility updates, but no new features. See the gno.land
commands at the root of ` + "`ignite --help`" + `.`,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cmd.Printf("ℹ️  The Cosmos SDK tooling is in maintenance mode. For new projects, consider gno.land (https://gno.land).\n\n")
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

// cosmosMaintenanceStubs returns hidden stub commands for the previously
// top-level cosmos commands, pointing users at `ignite cosmos`.
func cosmosMaintenanceStubs() []*cobra.Command {
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
