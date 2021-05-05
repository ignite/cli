package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

const (
	flagRequest    string = "request"
)

// NewType command creates a new type command to scaffold queries
func NewQuery() *cobra.Command {
	c := &cobra.Command{
		Use:   "query [name] [response_field1] [response_field2] ...",
		Short: "Scaffold a Cosmos SDK query",
		Args:  cobra.MinimumNArgs(1),
		RunE:  queryHandler,
	}
	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	c.Flags().String(flagModule, "", "Module to add the query into. Default: app's main module")
	c.Flags().StringSliceP(flagRequest, "r", []string{}, "Request fields")
	c.Flags().StringP(flagDescription, "d", "", "Description of the command")

	return c
}

func queryHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	// Get the module to add the type into
	module, err := cmd.Flags().GetString(flagModule)
	if err != nil {
		return err
	}

	// Get request fields
	reqFields, err := cmd.Flags().GetStringSlice(flagRequest)
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
		desc = fmt.Sprintf("Query %s", args[0])
	}

	sc, err := scaffolder.New(appPath)
	if err != nil {
		return err
	}
	if err := sc.AddQuery(module, args[0], desc, args[1:], reqFields); err != nil {
		return err
	}

	s.Stop()

	fmt.Printf("\n🎉 Created a query `%[1]v`.\n\n", args[0])
	return nil
}
