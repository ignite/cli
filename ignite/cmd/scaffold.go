package ignitecmd

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/ignite/cli/ignite/pkg/cliui/clispinner"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/xgit"
	"github.com/ignite/cli/ignite/services/scaffolder"
)

// flags related to component scaffolding
const (
	flagModule       = "module"
	flagNoMessage    = "no-message"
	flagNoSimulation = "no-simulation"
	flagResponse     = "response"
	flagDescription  = "desc"
)

// NewScaffold returns a command that groups scaffolding related sub commands.
func NewScaffold() *cobra.Command {
	c := &cobra.Command{
		Use:   "scaffold [command]",
		Short: "Scaffold a new blockchain, module, message, query, and more",
		Long: `Scaffold commands create and modify the source code files to add functionality.

CRUD stands for "create, read, update, delete".`,
		Aliases: []string{"s"},
		Args:    cobra.ExactArgs(1),
	}

	c.AddCommand(NewScaffoldChain())
	c.AddCommand(addGitChangesVerifier(NewScaffoldModule()))
	c.AddCommand(addGitChangesVerifier(NewScaffoldList()))
	c.AddCommand(addGitChangesVerifier(NewScaffoldMap()))
	c.AddCommand(addGitChangesVerifier(NewScaffoldSingle()))
	c.AddCommand(addGitChangesVerifier(NewScaffoldType()))
	c.AddCommand(addGitChangesVerifier(NewScaffoldMessage()))
	c.AddCommand(addGitChangesVerifier(NewScaffoldQuery()))
	c.AddCommand(addGitChangesVerifier(NewScaffoldPacket()))
	c.AddCommand(addGitChangesVerifier(NewScaffoldBandchain()))
	c.AddCommand(addGitChangesVerifier(NewScaffoldVue()))
	c.AddCommand(addGitChangesVerifier(NewScaffoldFlutter()))
	// c.AddCommand(NewScaffoldWasm())

	return c
}

func scaffoldType(
	cmd *cobra.Command,
	args []string,
	kind scaffolder.AddTypeKind,
) error {
	var (
		typeName          = args[0]
		fields            = args[1:]
		moduleName        = flagGetModule(cmd)
		withoutMessage    = flagGetNoMessage(cmd)
		withoutSimulation = flagGetNoSimulation(cmd)
		signer            = flagGetSigner(cmd)
		appPath           = flagGetPath(cmd)
	)

	var options []scaffolder.AddTypeOption

	if len(fields) > 0 {
		options = append(options, scaffolder.TypeWithFields(fields...))
	}
	if moduleName != "" {
		options = append(options, scaffolder.TypeWithModule(moduleName))
	}
	if withoutMessage {
		options = append(options, scaffolder.TypeWithoutMessage())
	} else {
		if signer != "" {
			options = append(options, scaffolder.TypeWithSigner(signer))
		}
		if withoutSimulation {
			options = append(options, scaffolder.TypeWithoutSimulation())
		}
	}

	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	sc, err := newApp(appPath)
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	sm, err := sc.AddType(cmd.Context(), cacheStorage, typeName, placeholder.New(), kind, options...)
	if err != nil {
		return err
	}

	s.Stop()

	modificationsStr, err := sourceModificationToString(sm)
	if err != nil {
		return err
	}

	fmt.Println(modificationsStr)
	fmt.Printf("\n🎉 %s added. \n\n", typeName)

	return nil
}

func addGitChangesVerifier(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().AddFlagSet(flagSetYes())

	preRunFun := cmd.PreRunE
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if preRunFun != nil {
			if err := preRunFun(cmd, args); err != nil {
				return err
			}
		}

		appPath := flagGetPath(cmd)

		changesCommitted, err := xgit.AreChangesCommitted(appPath)
		if err != nil {
			return err
		}

		if !getYes(cmd) && !changesCommitted {
			var confirmed bool
			prompt := &survey.Confirm{
				Message: "Your saved project changes have not been committed. To enable reverting to your current state, commit your saved changes. Do you want to proceed with scaffolding without committing your saved changes",
			}
			if err := survey.AskOne(prompt, &confirmed); err != nil || !confirmed {
				return errors.New("said no")
			}
		}
		return nil
	}
	return cmd
}

func flagSetScaffoldType() *flag.FlagSet {
	f := flag.NewFlagSet("", flag.ContinueOnError)
	f.String(flagModule, "", "Module to add into. Default is app's main module")
	f.Bool(flagNoMessage, false, "Disable CRUD interaction messages scaffolding")
	f.Bool(flagNoSimulation, false, "Disable CRUD simulation scaffolding")
	f.String(flagSigner, "", "Label for the message signer (default: creator)")
	return f
}

func flagGetModule(cmd *cobra.Command) string {
	module, _ := cmd.Flags().GetString(flagModule)
	return module
}

func flagGetNoSimulation(cmd *cobra.Command) bool {
	noMessage, _ := cmd.Flags().GetBool(flagNoSimulation)
	return noMessage
}

func flagGetNoMessage(cmd *cobra.Command) bool {
	noMessage, _ := cmd.Flags().GetBool(flagNoMessage)
	return noMessage
}

func flagGetSigner(cmd *cobra.Command) string {
	signer, _ := cmd.Flags().GetString(flagSigner)
	return signer
}
