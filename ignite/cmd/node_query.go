package ignitecmd

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const (
	flagPage       = "page"
	flagLimit      = "limit"
	flagPageKey    = "page-key"
	flagOffset     = "offset"
	flagCountTotal = "count-total"
	flagReverse    = "reverse"
)

func NewNodeQuery() *cobra.Command {
	c := &cobra.Command{
		Use:     "query",
		Short:   "Querying subcommands",
		Aliases: []string{"q"},
	}

	c.AddCommand(
		NewNodeQueryBank(),
		NewNodeQueryTx(),
	)

	return c
}

func flagSetPagination(query string) *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.Uint64(flagPage, 1, fmt.Sprintf("pagination page of %s to query for. This sets offset to a multiple of limit", query))
	fs.String(flagPageKey, "", fmt.Sprintf("pagination page-key of %s to query for", query))
	fs.Uint64(flagOffset, 0, fmt.Sprintf("pagination offset of %s to query for", query))
	fs.Uint64(flagLimit, 100, fmt.Sprintf("pagination limit of %s to query for", query))
	fs.Bool(flagCountTotal, false, fmt.Sprintf("count total number of records in %s to query for", query))
	fs.Bool(flagReverse, false, "results are sorted in descending order")

	return fs
}

func getPagination(cmd *cobra.Command) (*query.PageRequest, error) {
	var (
		pageKey, _    = cmd.Flags().GetString(flagPageKey)
		offset, _     = cmd.Flags().GetUint64(flagOffset)
		limit, _      = cmd.Flags().GetUint64(flagLimit)
		countTotal, _ = cmd.Flags().GetBool(flagCountTotal)
		page, _       = cmd.Flags().GetUint64(flagPage)
		reverse, _    = cmd.Flags().GetBool(flagReverse)
	)

	if page > 1 && offset > 0 {
		return nil, errors.New("page and offset cannot be used together")
	}

	if page > 1 {
		offset = (page - 1) * limit
	}

	return &query.PageRequest{
		Key:        []byte(pageKey),
		Offset:     offset,
		Limit:      limit,
		CountTotal: countTotal,
		Reverse:    reverse,
	}, nil
}
