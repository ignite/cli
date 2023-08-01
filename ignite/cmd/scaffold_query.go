package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/services/scaffolder"
)

const (
	flagPaginated = "paginated"
)

// NewScaffoldQuery command creates a new type command to scaffold queries.
func NewScaffoldQuery() *cobra.Command {
	c := &cobra.Command{
		Use:     "query [name] [request_field1] [request_field2] ...",
		Short:   "Query for fetching data from a blockchain",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: migrationPreRunHandler,
		RunE:    queryHandler,
	}

	flagSetPath(c)
	flagSetClearCache(c)

	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().String(flagModule, "", "module to add the query into. Default: app's main module")
	c.Flags().StringSliceP(flagResponse, "r", []string{}, "response fields")
	c.Flags().StringP(flagDescription, "d", "", "description of the CLI to broadcast a tx with the message")
	c.Flags().Bool(flagPaginated, false, "define if the request can be paginated")

	return c
}

func queryHandler(cmd *cobra.Command, args []string) error {
	appPath := flagGetPath(cmd)

	session := cliui.New(cliui.StartSpinnerWithText(statusScaffolding))
	defer session.End()

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

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	sc, err := scaffolder.New(appPath)
	if err != nil {
		return err
	}

	sm, err := sc.AddQuery(cmd.Context(), cacheStorage, placeholder.New(), module, args[0], desc, args[1:], resFields, paginated)
	if err != nil {
		return err
	}

	modificationsStr, err := sourceModificationToString(sm)
	if err != nil {
		return err
	}

	session.Println(modificationsStr)
	session.Printf("\n🎉 Created a query `%[1]v`.\n\n", args[0])

	return nil
}
