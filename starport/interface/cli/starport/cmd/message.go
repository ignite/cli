package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

const (
	responseFlag string = "res"
)

// NewType command creates a new type command to scaffold types.
func NewMessage() *cobra.Command {
	c := &cobra.Command{
		Use:   "message [name] [description] [field1] [field2] ...",
		Short: "Generates an empty message",
		Args:  cobra.MinimumNArgs(2),
		RunE:  messageHandler,
	}
	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	addSdkVersionFlag(c)

	c.Flags().String(moduleFlag, "", "Module to add the message into. Default: app's main module")
	c.Flags().StringSlice(responseFlag, []string{}, "Response fields")

	return c
}

func messageHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	// Get the module to add the type into
	module, err := cmd.Flags().GetString(moduleFlag)
	if err != nil {
		return err
	}

	// Get response fields
	resFields, err := cmd.Flags().GetStringSlice(responseFlag)
	if err != nil {
		return err
	}

	sc := scaffolder.New(appPath)
	if err := sc.AddMessage(module, args[0], args[1], args[2:], resFields); err != nil {
		return err
	}

	s.Stop()

	fmt.Printf("\n🎉 Created a message `%[1]v`.\n\n", args[0])
	return nil
}
