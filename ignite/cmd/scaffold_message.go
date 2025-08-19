package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/services/scaffolder"
	"github.com/ignite/cli/v29/ignite/templates/field"
)

const flagSigner = "signer"

// NewScaffoldMessage returns the command to scaffold messages.
func NewScaffoldMessage() *cobra.Command {
	c := &cobra.Command{
		Use:   "message [name] [field1:type1] [field2:type2] ...",
		Short: "Message to perform state transition on the blockchain",
		Long: `Message scaffolding is useful for quickly adding functionality to your
blockchain to handle specific Cosmos SDK messages.

Messages are objects whose end goal is to trigger state transitions on the
blockchain. A message is a container for fields of data that affect how the
blockchain's state will change. You can think of messages as "actions" that a
user can perform.

For example, the bank module has a "Send" message for token transfers between
accounts. The send message has three fields: from address (sender), to address
(recipient), and a token amount. When this message is successfully processed,
the token amount will be deducted from the sender's account and added to the
recipient's account.

Ignite's message scaffolding lets you create new types of messages and add them
to your chain. For example:

	ignite scaffold message add-pool amount:coins denom active:bool --module dex

The command above will create a new message MsgAddPool with three fields: amount
(in tokens), denom (a string), and active (a boolean). The message will be added
to the "dex" module.

For detailed type information use ignite scaffold type --help

By default, the message is defined as a proto message in the
"proto/{app}/{module}/tx.proto" and registered in the "Msg" service. A CLI command to
create and broadcast a transaction with MsgAddPool is created in the module's
"cli" package. Additionally, Ignite scaffolds a message constructor and the code
to satisfy the sdk.Msg interface and register the message in the module.

Most importantly in the "keeper" package Ignite scaffolds an "AddPool" function.
Inside this function, you can implement message handling logic.

When successfully processed a message can return data. Use the —response flag to
specify response fields and their types. For example

	ignite scaffold message create-post title body --response id:int,title

The command above will scaffold MsgCreatePost which returns both an ID (an
integer) and a title (a string).

Message scaffolding follows the rules as "ignite scaffold list/map/single" and
supports fields with standard and custom types. See "ignite scaffold list —help"
for details.
`,
		Args:    cobra.MinimumNArgs(1),
		PreRunE: migrationPreRunHandler,
		RunE:    messageHandler,
	}

	flagSetPath(c)
	flagSetClearCache(c)

	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().String(flagModule, "", "module to add the message into. Default: app's main module")
	c.Flags().StringSliceP(flagResponse, "r", []string{}, "response fields")
	c.Flags().Bool(flagNoSimulation, false, "disable CRUD simulation scaffolding")
	c.Flags().StringP(flagDescription, "d", "", "description of the command")
	c.Flags().String(flagSigner, "", "label for the message signer (default: creator)")

	return c
}

func messageHandler(cmd *cobra.Command, args []string) error {
	var (
		module, _         = cmd.Flags().GetString(flagModule)
		resFields, _      = cmd.Flags().GetStringSlice(flagResponse)
		desc, _           = cmd.Flags().GetString(flagDescription)
		signer            = flagGetSigner(cmd)
		appPath           = flagGetPath(cmd)
		withoutSimulation = flagGetNoSimulation(cmd)
	)

	session := cliui.New(
		cliui.StartSpinnerWithText(statusScaffolding),
		cliui.WithoutUserInteraction(getYes(cmd)),
	)
	defer session.End()

	hasMultipleCoinSlice, err := field.MultipleCoins(resFields)
	if err != nil {
		return err
	}
	if hasMultipleCoinSlice {
		session.PauseSpinner()
		_ = session.Print(colors.Info(multipleCoinDisclaimer))
		session.StartSpinner(statusScaffolding)
	}

	cfg, _, err := getChainConfig(cmd)
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	var options []scaffolder.MessageOption

	// Get description
	if desc != "" {
		options = append(options, scaffolder.WithDescription(desc))
	}

	// Get signer
	if signer != "" {
		options = append(options, scaffolder.WithSigner(signer))
	}

	// Skip scaffold simulation
	if withoutSimulation {
		options = append(options, scaffolder.WithoutSimulation())
	}

	sc, err := scaffolder.New(cmd.Context(), appPath, cfg.Build.Proto.Path)
	if err != nil {
		return err
	}

	err = sc.AddMessage(cmd.Context(), module, args[0], args[1:], resFields, options...)
	if err != nil {
		return err
	}

	sm, err := sc.ApplyModifications(xgenny.ApplyPreRun(scaffolder.AskOverwriteFiles(session)))
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
	session.Printf("\n🎉 Created a message `%[1]v`.\n\n", args[0])

	return nil
}
