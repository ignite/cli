package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/services/gno"
)

// NewGnoScaffold returns the `ignite scaffold` command for gno.land
// realms and packages.
func NewGnoScaffold() *cobra.Command {
	c := &cobra.Command{
		Use:     "scaffold [command]",
		Aliases: []string{"s"},
		Short:   "Scaffold gno.land realms and packages",
		Long:    `Scaffold new gno.land smart contracts: stateful realms (r/) and stateless packages (p/).`,
		Args:    cobra.ExactArgs(0),
	}

	c.AddCommand(
		newGnoScaffoldRealm(),
		newGnoScaffoldPackage(),
	)

	return c
}

func newGnoScaffoldRealm() *cobra.Command {
	return &cobra.Command{
		Use:   "realm <name>",
		Short: "Scaffold a new gno.land realm (stateful smart contract)",
		Long: `Scaffold a new gno.land realm.

A realm is a stateful smart contract: package-level variables are persisted
on-chain. <name> is either a bare name ("counter", deployed as
gno.land/r/counter) or a full path ("gno.land/r/demo/counter").`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			session := cliui.New(cliui.StartSpinnerWithText(statusCreating))
			defer session.End()

			dir, modulePath, err := gno.Scaffold(gno.KindRealm, args[0], gno.ScaffoldOptions{})
			if err != nil {
				return err
			}
			return session.Printf("%s Realm scaffolded.\n\n  Directory: %s\n  Module:   %s\n", icons.OK, dir, modulePath)
		},
	}
}

func newGnoScaffoldPackage() *cobra.Command {
	return &cobra.Command{
		Use:   "package <name>",
		Short: "Scaffold a new gno.land package (stateless library)",
		Long: `Scaffold a new gno.land package.

A package is a stateless library of pure functions (gno.land/p/...).`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			session := cliui.New(cliui.StartSpinnerWithText(statusCreating))
			defer session.End()

			dir, modulePath, err := gno.Scaffold(gno.KindPackage, args[0], gno.ScaffoldOptions{})
			if err != nil {
				return err
			}
			return session.Printf("%s Package scaffolded.\n\n  Directory: %s\n  Module:   %s\n", icons.OK, dir, modulePath)
		},
	}
}
