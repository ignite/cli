package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

const flagSigner = "signer"

// NewScaffoldMessage returns the command to scaffold messages
func NewScaffoldMessage() *cobra.Command {
	c := &cobra.Command{
		Use:   "message [name] [field1] [field2] ...",
		Short: "Message to perform state transition on the blockchain",
		Args:  cobra.MinimumNArgs(1),
		RunE:  messageHandler,
	}
	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	c.Flags().String(flagModule, "", "Module to add the message into. Default: app's main module")
	c.Flags().StringSliceP(flagResponse, "r", []string{}, "Response fields")
	c.Flags().StringP(flagDescription, "d", "", "Description of the command")
	c.Flags().String(flagSigner, "", "Label for the message signer (default: creator)")

	return c
}

func messageHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	var (
		module, _    = cmd.Flags().GetString(flagModule)
		resFields, _ = cmd.Flags().GetStringSlice(flagResponse)
		desc, _      = cmd.Flags().GetString(flagDescription)
		signer, _    = cmd.Flags().GetString(flagSigner)
	)

	var options []scaffolder.MessageOption

	// Get description
	if desc != "" {
		options = append(options, scaffolder.WithDescription(desc))
	}

	// Get signer
	if signer != "" {
		options = append(options, scaffolder.WithSigner(signer))
	}

	sc, err := scaffolder.New(appPath)
	if err != nil {
		return err
	}
	sm, err := sc.AddMessage(placeholder.New(), module, args[0], args[1:], resFields, options...)
	if err != nil {
		return err
	}

	s.Stop()

	fmt.Println(sourceModificationToString(sm))
	fmt.Printf("\n🎉 Created a message `%[1]v`.\n\n", args[0])
	return nil
}
