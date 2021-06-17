package starportcmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/validation"
	"github.com/tendermint/starport/starport/services/scaffolder"
	"github.com/tendermint/starport/starport/templates/module"
)

const (
	flagIBC                 = "ibc"
	flagIBCOrdering         = "ordering"
	flagRequireRegistration = "require-registration"
)

var ibcRouterPlaceholderInstruction = fmt.Sprintf(`
💬 To enable scaffolding of IBC modules, remove these lines from app/app.go:

%s

💬 Then, find the following line:

%s

💬 Finally, add this block of code below:
%s

`,
	infoColor(`ibcRouter := porttypes.NewRouter()
ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferModule)
app.IBCKeeper.SetRouter(ibcRouter)`),
	infoColor(module.PlaceholderSgAppKeeperDefinition),
	infoColor(`ibcRouter := porttypes.NewRouter()
ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferModule)
`+module.PlaceholderIBCAppRouter+`
app.IBCKeeper.SetRouter(ibcRouter)`),
)

// NewModuleCreate creates a new module create command to scaffold an
// sdk module.
func NewModuleCreate() *cobra.Command {
	c := &cobra.Command{
		Use:   "create [name]",
		Short: "Scaffold a Cosmos SDK module",
		Long:  "Scaffold a new Cosmos SDK module in the `x` directory",
		Args:  cobra.MinimumNArgs(1),
		RunE:  createModuleHandler,
	}
	c.Flags().Bool(flagIBC, false, "scaffold an IBC module")
	c.Flags().String(flagIBCOrdering, "none", "channel ordering of the IBC module [none|ordered|unordered]")
	c.Flags().Bool(flagRequireRegistration, false, "if true command will fail if module can't be registered")
	return c
}

func createModuleHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	var options []scaffolder.ModuleCreationOption

	name := args[0]

	ibcModule, err := cmd.Flags().GetBool(flagIBC)
	if err != nil {
		return err
	}

	ibcOrdering, err := cmd.Flags().GetString(flagIBCOrdering)
	if err != nil {
		return err
	}
	requireRegistration, err := cmd.Flags().GetBool(flagRequireRegistration)
	if err != nil {
		return err
	}

	// Check if the module must be an IBC module
	if ibcModule {
		options = append(options, scaffolder.WithIBCChannelOrdering(ibcOrdering), scaffolder.WithIBC())
	}

	sc, err := scaffolder.New(appPath)
	if err != nil {
		return err
	}
	var msg bytes.Buffer
	fmt.Fprintf(&msg, "\n🎉 Module created %s.\n\n", name)
	sm, err := sc.CreateModule(placeholder.New(), name, options...)
	s.Stop()
	if err != nil {
		// If this is an old scaffolded application that doesn't contain the necessary placeholder
		// We give instruction to the user to modify the application
		if err == scaffolder.ErrNoIBCRouterPlaceholder {
			fmt.Print(ibcRouterPlaceholderInstruction)
		}
		var validationErr validation.Error
		if !requireRegistration && errors.As(err, &validationErr) {
			fmt.Fprintf(&msg, "Can't register module '%s'.\n", name)
			fmt.Fprintln(&msg, validationErr.ValidationInfo())
		} else {
			return err
		}
	} else {
		fmt.Println(sourceModificationToString(sm))
	}

	io.Copy(cmd.OutOrStdout(), &msg)
	return nil
}
