package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/services/scaffolder"
)

const (
	flagPaginated = "paginated"
)

// NewScaffoldQuery command creates a new type command to scaffold queries.
func NewScaffoldQuery() *cobra.Command {
	c := &cobra.Command{
		Use:   "query [name] [field1:type1] [field2:type2] ...",
		Short: "Query for fetching data from a blockchain",
		Long: `Query for fetching data from a blockchain.
		
For detailed type information use ignite scaffold type --help.`,
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
	session := cliui.New(
		cliui.StartSpinnerWithText(statusScaffolding),
		cliui.WithoutUserInteraction(getYes(cmd)),
	)
	defer session.End()

	cfg, _, err := getChainConfig(cmd)
	if err != nil {
		return err
	}

	// Get the module to add the type into
	module, _ := cmd.Flags().GetString(flagModule)

	// Get request fields
	resFields, _ := cmd.Flags().GetStringSlice(flagResponse)

	// Get description
	desc, _ := cmd.Flags().GetString(flagDescription)
	if desc == "" {
		// Use a default description
		desc = fmt.Sprintf("Query %s", args[0])
	}

	var (
		paginated, _ = cmd.Flags().GetBool(flagPaginated)
		appPath      = flagGetPath(cmd)
	)

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	sc, err := scaffolder.New(cmd.Context(), appPath, cfg.Build.Proto.Path)
	if err != nil {
		return err
	}

	err = sc.AddQuery(cmd.Context(), module, args[0], desc, args[1:], resFields, paginated)
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
	session.Printf("\n🎉 Created a query `%[1]v`.\n\n", args[0])

	return nil
}
