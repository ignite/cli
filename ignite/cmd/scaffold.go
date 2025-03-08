package ignitecmd

import (
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gocmd"
	"github.com/ignite/cli/v29/ignite/pkg/xgit"
	"github.com/ignite/cli/v29/ignite/services/scaffolder"
	"github.com/ignite/cli/v29/ignite/version"
)

// flags related to component scaffolding.
const (
	flagModule       = "module"
	flagNoMessage    = "no-message"
	flagNoSimulation = "no-simulation"
	flagResponse     = "response"
	flagDescription  = "desc"
	flagProtoDir     = "proto-dir"

	msgCommitPrefix = "Your saved project changes have not been committed.\nTo enable reverting to your current state, commit your saved changes."
	msgCommitPrompt = "Do you want to proceed without committing your saved changes"

	statusScaffolding = "Scaffolding..."

	supportFieldTypes = `
Currently supports: 

| Type         | Alias   | Index | Code Type | Description                     |
|--------------|---------|-------|-----------|---------------------------------|
| string       | -       | yes   | string    | Text type                       |
| array.string | strings | no    | []string  | List of text type               |
| bool         | -       | yes   | bool      | Boolean type                    |
| int          | -       | yes   | int64     | Integer type                    |
| array.int    | ints    | no    | []int64   | List of integers types          |
| uint         | -       | yes   | uint64    | Unsigned integer type           |
| array.uint   | uints   | no    | []uint64  | List of unsigned integers types |
| coin         | -       | no    | sdk.Coin  | Cosmos SDK coin type            |
| array.coin   | coins   | no    | sdk.Coins | List of Cosmos SDK coin types   |

Field Usage:
    - fieldName
    - fieldName:fieldType

If no :fieldType, default (string) is used
`
)

// NewScaffold returns a command that groups scaffolding related sub commands.
func NewScaffold() *cobra.Command {
	c := &cobra.Command{
		Use:   "scaffold [command]",
		Short: "Create a new blockchain, module, message, query, and more",
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

	c.AddCommand(
		NewScaffoldChain(),
		NewScaffoldModule(),
		NewScaffoldList(),
		NewScaffoldMap(),
		NewScaffoldSingle(),
		NewScaffoldType(),
		NewScaffoldParams(),
		NewScaffoldConfigs(),
		NewScaffoldMessage(),
		NewScaffoldQuery(),
		NewScaffoldPacket(),
		NewScaffoldVue(),
		NewScaffoldReact(),
		NewScaffoldChainRegistry(),
	)

	return c
}

func migrationPreRunHandler(cmd *cobra.Command, args []string) error {
	if err := gitChangesConfirmPreRunHandler(cmd, args); err != nil {
		return err
	}

	session := cliui.New()
	defer session.End()

	cfg, _, err := getChainConfig(cmd)
	if err != nil {
		return err
	}

	appPath, err := goModulePath(cmd)
	if err != nil {
		return err
	}

	ver, err := cosmosver.Detect(appPath)
	if err != nil {
		return err
	}

	if err := version.AssertSupportedCosmosSDKVersion(ver); err != nil {
		return err
	}

	if err := toolsMigrationPreRunHandler(cmd, session, appPath); err != nil {
		return err
	}

	if err := bufMigrationPreRunHandler(cmd, session, appPath, cfg.Build.Proto.Path); err != nil {
		return err
	}

	// we go mod tidy in case new dependencies were added or removed
	if err := gocmd.ModTidy(cmd.Context(), appPath); err != nil {
		return err
	}

	return nil
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

	session := cliui.New(cliui.StartSpinnerWithText(statusScaffolding))
	defer session.End()

	cfg, _, err := getChainConfig(cmd)
	if err != nil {
		return err
	}

	sc, err := scaffolder.New(cmd.Context(), appPath, cfg.Build.Proto.Path)
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	err = sc.AddType(cmd.Context(), typeName, kind, options...)
	if err != nil {
		return err
	}

	sm, err := sc.ApplyModifications()
	if err != nil {
		return err
	}

	if err := sc.PostScaffold(cmd.Context(), cacheStorage, false); err != nil {
		return err
	}

	modificationsStr, err := sm.String()
	if err != nil {
		return err
	}

	session.Println(modificationsStr)
	session.Printf("\n🎉 %s added. \n\n", typeName)

	return nil
}

func gitChangesConfirmPreRunHandler(cmd *cobra.Command, _ []string) error {
	// Don't confirm when the "--yes" flag is present
	if getYes(cmd) {
		return nil
	}

	appPath := flagGetPath(cmd)
	session := cliui.New()

	defer session.End()

	return confirmWhenUncommittedChanges(session, appPath)
}

func confirmWhenUncommittedChanges(session *cliui.Session, appPath string) error {
	cleanState, err := xgit.AreChangesCommitted(appPath)
	if err != nil {
		return err
	}

	if !cleanState {
		session.Println(msgCommitPrefix)
		if err := session.AskConfirm(msgCommitPrompt); err != nil {
			if errors.Is(err, promptui.ErrAbort) {
				return errors.New("No")
			}

			return err
		}
	}

	return nil
}

func flagSetScaffoldType() *flag.FlagSet {
	f := flag.NewFlagSet("", flag.ContinueOnError)
	f.String(flagModule, "", "specify which module to generate code in")
	f.Bool(flagNoMessage, false, "skip generating message handling logic")
	f.Bool(flagNoSimulation, false, "skip simulation logic")
	f.String(flagSigner, "", "label for the message signer (default: creator)")
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
