package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/services/scaffolder"
)

const (
	FlagIndexes = "index"
)

// NewScaffoldMap returns a new command to scaffold a map.
func NewScaffoldMap() *cobra.Command {
	c := &cobra.Command{
		Use:   "map NAME [field]...",
		Short: "CRUD for data stored as key-value pairs",
		Long: `The "map" scaffolding command is used to generate files that implement the logic
for storing and interacting with data stored as key-value pairs (or a
dictionary) in the blockchain state.

The "map" command is very similar to "ignite scaffold list" with the main
difference in how values are indexed. With "list" values are indexed by an
incrementing integer, whereas "map" values are indexed by a user-provided value
(or multiple values).

Let's use the same blog post example:

	ignite scaffold map post title body

This command scaffolds a "Post" type and CRUD functionality to create, read,
updated, and delete posts. However, when creating a new post with your chain's
binary (or by submitting a transaction through the chain's API) you will be
required to provide an "index":

	blogd tx blog create-post [index] [title] [body]
	blogd tx blog create-post hello "My first post" "This is the body"

This command will create a post and store it in the blockchain's state under the
"hello" index. You will be able to fetch back the value of the post by querying
for the "hello" key.

	blogd q blog show-post hello

To customize the index, use the "--index" flag. Multiple indices can be
provided, which simplifies querying values. For example:

	ignite scaffold map product price desc --index category,guid

With this command, you would get a "Product" value indexed by both a category
and a GUID (globally unique ID). This will let you programmatically fetch
product values that have the same category but are using different GUIDs.

Since the behavior of "list" and "map" scaffolding is very similar, you can use
the "--no-message", "--module", "--signer" flags as well as the colon syntax for
custom types.
`,
		Args:    cobra.MinimumNArgs(1),
		PreRunE: migrationPreRunHandler,
		RunE:    scaffoldMapHandler,
	}

	flagSetPath(c)
	flagSetClearCache(c)

	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().AddFlagSet(flagSetScaffoldType())
	c.Flags().StringSlice(FlagIndexes, []string{"index"}, "fields that index the value")

	return c
}

func scaffoldMapHandler(cmd *cobra.Command, args []string) error {
	indexes, err := cmd.Flags().GetStringSlice(FlagIndexes)
	if err != nil {
		return err
	}

	return scaffoldType(cmd, args, scaffolder.MapType(indexes...))
}
