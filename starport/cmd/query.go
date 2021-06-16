package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

const (
	flagPaginated = "paginated"
)

// NewType command creates a new type command to scaffold queries
func NewQuery() *cobra.Command {
	c := &cobra.Command{
		Use:   "query [name] [request_field1] [request_field2] ...",
		Short: "Scaffold a Cosmos SDK query",
		Args:  cobra.MinimumNArgs(1),
		RunE:  queryHandler,
	}
	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	c.Flags().String(flagModule, "", "Module to add the query into. Default: app's main module")
	c.Flags().StringSliceP(flagResponse, "r", []string{}, "Response fields")
	c.Flags().StringP(flagDescription, "d", "", "Description of the command")
	c.Flags().Bool(flagPaginated, false, "Define if the request can be paginated")

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
		desc = fmt.Sprintf("Query %s", args[0])
	}

	paginated, err := cmd.Flags().GetBool(flagPaginated)
	if err != nil {
		return err
	}

	sc, err := scaffolder.New(appPath)
	if err != nil {
		return err
	}
	sm, err := sc.AddQuery(placeholder.New(), module, args[0], desc, args[1:], resFields, paginated)
	if err != nil {
		return err
	}

	s.Stop()

	fmt.Println(sourceModificationToString(sm))
	fmt.Printf("\n🎉 Created a query `%[1]v`.\n\n", args[0])
	return nil
}
