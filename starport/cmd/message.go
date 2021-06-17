package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

const (
	flagResponse    string = "response"
	flagDescription string = "desc"
)

// NewType command creates a new type command to scaffold messages
func NewMessage() *cobra.Command {
	c := &cobra.Command{
		Use:   "message [name] [field1] [field2] ...",
		Short: "Scaffold a Cosmos SDK message",
		Args:  cobra.MinimumNArgs(1),
		RunE:  messageHandler,
	}
	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	c.Flags().String(flagModule, "", "Module to add the message into. Default: app's main module")
	c.Flags().StringSliceP(flagResponse, "r", []string{}, "Response fields")
	c.Flags().StringP(flagDescription, "d", "", "Description of the command")

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

	// Get description
	desc, err := cmd.Flags().GetString(flagDescription)
	if err != nil {
		return err
	}
	if desc == "" {
		// Use a default description
		desc = fmt.Sprintf("Broadcast message %s", args[0])
	}

	sc, err := scaffolder.New(appPath)
	if err != nil {
		return err
	}
	sm, err := sc.AddMessage(placeholder.New(), module, args[0], desc, args[1:], resFields)
	if err != nil {
		return err
	}

	s.Stop()

	fmt.Println(sourceModificationToString(sm))
	fmt.Printf("\n🎉 Created a message `%[1]v`.\n\n", args[0])
	return nil
}
