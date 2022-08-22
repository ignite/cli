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
		Long: `Scaffolding is a quick way to generate code for major pieces of your
application.

For details on each scaffolding target (chain, module, message, etc.) run the
corresponding command with a "--help" flag, for example, "ignite scaffold chain
--help".

The Ignite team strongly recommends committing the code to a version control
system before running scaffolding commands. This will make it easier to see the
changes to the source code as well as undo the command if you've decided to roll
back the changes.

This blockchain you create with the chain scaffolding command uses the modular
Cosmos SDK framework and imports many standard modules for functionality like
proof of stake, token transfer, inter-blockchain connectivity, governance, and
more. Custom functionality is implemented in modules located by convention in
the "x/" directory. By default, your blockchain comes with an empty custom
module. Use the module scaffolding command to create an additional module.

An empty custom module doesn't do much, it's basically a container for logic
that is responsible for processing transactions and changing the application
state. Cosmos SDK blockchains work by processing user-submitted signed
transactions, which contain one or more messages. A message contains data that
describes a state transition. A module can be responsible for handling any
number of messages.

A message scaffolding command will generate the code for handling a new type of
Cosmos SDK message. Message fields describe the state transition that the
message is intended to produce if processed without errors.

Scaffolding messages is useful to create individual "actions" that your module
can perform. Sometimes, however, you want your blockchain to have the
functionality to create, read, update and delete (CRUD) instances of a
particular type. Depending on how you want to store the data there are three
commands that scaffold CRUD functionality for a type: list, map, and single.
These commands create four messages (one for each CRUD action), and the logic to
add, delete, and fetch the data from the store. If you want to scaffold only the
logic, for example, you've decided to scaffold messages separately, you can do
that as well with the "--no-message" flag.

Reading data from a blockchain happens with a help of queries. Similar to how
you can scaffold messages to write data, you can scaffold queries to read the
data back from your blockchain application.

You can also scaffold a type, which just produces a new protocol buffer file
with a proto message description. Note that proto messages produce (and
correspond with) Go types whereas Cosmos SDK messages correspond to proto "rpc"
in the "Msg" service.

If you're building an application with custom IBC logic, you might need to
scaffold IBC packets. An IBC packet represents the data sent from one blockchain
to another. You can only scaffold IBC packets in IBC-enabled modules scaffolded
with an "--ibc" flag. Note that the default module is not IBC-enabled.
`,
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
