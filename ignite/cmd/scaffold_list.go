package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/services/scaffolder"
)

// NewScaffoldList returns a new command to scaffold a list.
func NewScaffoldList() *cobra.Command {
	c := &cobra.Command{
		Use: "list NAME [field]...",
		Short: `The "list" scaffolding command is used to generate files that implement the
logic for storing and interacting with data stored as a list in the blockchain
state.

The command accepts a NAME argument that will be used as the name of a new type
of data. It also accepts a list of FIELDs that describe the type.

The interaction with the data follows the create, read, updated, and delete
(CRUD) pattern. For each type three Cosmos SDK messages are defined for writing
data to the blockchain: MsgCreate{Name}, MsgUpdate{Name}, MsgDelete{Name}. For
reading data two queries are defined: {Name} and {Name}All. The type, messages,
and queries are defined in the "proto/" directory as protocol buffer messages.
Messages and queries are mounted in the "Msg" and "Query" services respectively.

When messages are handled, the appropriate keeper methods are called. By
convention, the methods are defined in
"x/{moduleName}/keeper/msg_server_{name}.go". Helpful methods for getting,
setting, removing, and appending are defined in the same "keeper" package in
"{name}.go".

The "list" command essentially allows you to define a new type of data and
provides the logic to create, read, update, and delete instances of the type.
For example, let's review a command that generates the code to handle a list of
posts and each post has "title" and "body" fields:

  ignite scaffold list post title body

This provides you with a "Post" type, MsgCreatePost, MsgUpdatePost,
MsgDeletePost and two queries: Post and PostAll. The compiled CLI, let's say the
binary is "blogd" and the module is "blog", has commands to query the chain (see
"blogd q blog") and broadcast transactions with the messages above (see "blogd
tx blog").

The code generated with the list command is meant to be edited and tailored to
your application needs. Consider the code to be a "skeleton" for the actual
business logic you will implement next.

By default, all fields are assumed to be strings. If you want a field of a
different type, you can specify it after a colon ":". The following types are
supported: string, bool, int, uint, coin, array.string, array.int, array.uint,
array.coin.

  ignite scaffold pool amount:coin tags:array.string height:int
  
Ignite also supports custom types:
  
  ignite scaffold list product-details name description
  
  ignite scaffold list product price:coin details:ProductDetails

In the example above the "ProductDetails" type was defined first, and then used
as a custom type for the "details" field. Ignite doesn't support arrays of
custom types yet.

By default the code will be scaffolded in the module that matches your project's
name. If you have several modules in your project, you might want to specify a
different module:

  ignite scaffold list post title body --module blog

By default, each message comes with a "creator" field that represents the
address of the transaction signer. You can customize the name of this field with
a flag:

  ignite scaffold list post title body --signer author

It's possible to scaffold just the getter/setter logic without the CRUD
messages. This is useful when you want the methods to handle a type, but would
like to scaffold messages manually. Use a flag to skip message scaffolding:

  ignite scaffold list post title body --no-message

The "creator" field is not generated if a list is scaffolded with the
"--no-message" flag.
`,
		Args: cobra.MinimumNArgs(1),
		RunE: scaffoldListHandler,
	}

	flagSetPath(c)
	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagSetScaffoldType())

	return c
}

func scaffoldListHandler(cmd *cobra.Command, args []string) error {
	return scaffoldType(cmd, args, scaffolder.ListType())
}
