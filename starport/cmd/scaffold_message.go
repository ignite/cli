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
	c.Flags().String(flagSigner, "", "Name of the message signer")

	return c
}

func messageHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	// Get the module to add the type into
	module, err := cmd.Flags().GetString(flagModule)
	if err != nil {
		return err
	}

	// Get response fields
	resFields, err := cmd.Flags().GetStringSlice(flagResponse)
	if err != nil {
		return err
	}

	var options []scaffolder.MessageOption

	// Get description
	desc, err := cmd.Flags().GetString(flagDescription)
	if err != nil {
		return err
	}
	if desc == "" {
		options = append(options, scaffolder.WithDescription(desc))
	}

	// Get signer
	signer, err := cmd.Flags().GetString(flagSigner)
	if err != nil {
		return err
	}
	if signer == "" {
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
