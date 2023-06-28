---
description: Step-by-step guidance to build your first blockchain and your first Cosmos SDK module.
title: Express tutorial
---

# "Hello, World!" in 5 minutes

In this tutorial, you will create a simple blockchain with a custom query that
responds with `"Hello, %s!"`, where `%s` is a name provided in the query. To do
this, you will use the Ignite CLI to generate most of the code, and then modify
the query to return the desired response. After completing the tutorial, you
will have a better understanding of how to create custom queries in a
blockchain.

First, create a new `hello` blockchain with Ignite CLI:

```
ignite scaffold chain hello
```

Let's add a query to the blockchain we just created.

In the Cosmos SDK, a query is a request for information from the blockchain.
Queries are used to retrieve data from the blockchain, such as the current state
of the ledger or the details of a specific transaction. The Cosmos SDK provides
a number of built-in query methods that can be used to retrieve data from the
blockchain, and developers can also create custom queries to access specific
data or perform complex operations. Queries are processed by the blockchain's
nodes and the results are returned to the querying client.

## Create a query with Ignite

To add a query, run the following command inside the `hello` directory:

```
ignite scaffold query say-hello name --response name
```

The `ignite scaffold query` command is a tool used to quickly create new
queries. When you run this command, it makes changes to your source code to add
the new query and make it available in your API. This command accepts a query
name (`"say-hello"`) and a list of request fields (in our case only `name`). The
optional `--response` flag specifies the return values of the query.

This command made the following changes to the source code.

The `proto/hello/hello/query.proto` file was modified to define the request and
response for a query, as well as to add the `SayHello` query in the `Query`
service. 

The `x/hello/client/cli/query_say_hello.go` file was created and added to the
project. This file contains a CLI command `CmdSayHello` that allows users to
submit a "say hello" query to the blockchain. This command allows users to
interact with the blockchain in a more user-friendly way, allowing them to
easily submit queries and receive responses from the blockchain.

The `x/hello/client/cli/query.go` was modified to add the `CmdSayHello` command
to the CLI of the blockchain.

The `x/hello/keeper/query_say_hello.go` file was created with a keeper method
called `SayHello`. This method is responsible for handling the "say hello"
query, which can be called by a client using the command-line interface (CLI) or
an API. When the "say hello" query is executed, the `SayHello` method is called
to perform the necessary actions and return a response to the client. The
`SayHello` method may retrieve data from the application's database, process the
data, and return a result to the client in a specific format, such as a string
of text or a data structure.

To change the source code so that the query returns the `"Hello, %s!"` string,
modify the return statement in `query_say_hello.go` to return
`fmt.Sprintf("hello %s", req.Name)`.

```go title="x/hello/keeper/query_say_hello.go"
func (k Keeper) SayHello(goCtx context.Context, req *types.QuerySayHelloRequest) (*types.QuerySayHelloResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Process the query
	_ = ctx
	// highlight-next-line
	return &types.QuerySayHelloResponse{Name: fmt.Sprintf("Hello, %s!", req.Name)}, nil
}
```

The function now returns a `QuerySayHelloResponse` struct with the `Name` field
set to the string `"Hello, %s!"` with `req.Name` as the value for the `%s`
placeholder. It also returns a nil error to indicate success.

Now that you have added a query your blockchain and modified it return the value
you want, you can start your blockchain with Ignite:

```
ignite chain serve
```

After starting your blockchain, you can use its command-line interface (CLI) to
interact with it and perform various actions such as querying the blockchain's
state, sending transactions, and more.

You can use the `hellod` binary to run the `say-hello` query:

```
hellod q hello say-hello bob
```

Once you run this command, the `hellod` binary will send a `say-hello` query to
your blockchain with the argument `bob`. The blockchain will process the query
and return the result, which will be printed by the `hellod` binary. In this
case, the expected result is a string containing the message `Hello, bob!`.

```
name: Hello, bob!
```

Congratulations! 🎉 You have successfully created a new Cosmos SDK module called
`hello` with a custom query functionality. This allows users to query the
blockchain and receive a response with a personalized greeting. This tutorial
demonstrated how to use Ignite CLI to create a custom query in a blockchain.

Ignite is an incredibly convenient tool for developers because it automatically
generates much of the code required for a project. This saves developers time
and effort by reducing the amount of code they need to write manually. With
Ignite, developers can quickly and easily set up the basic structure of their
project, allowing them to focus on the more complex and unique aspects of their
work.

However, it is also important for developers to understand how the code
generated by Ignite works under the hood. One way to do this is to implement the
same functionality manually, without using Ignite. For example, in this tutorial
Ignite was used to generate query functionality, now could try implementing the
same functionality manually to see how it works and gain a deeper understanding
of the code.

Implementing the same functionality manually can be time-consuming and
challenging, but it can also be a valuable learning experience. By seeing how
the code works at a low level, developers can gain a better understanding of how
different components of their project fit together and how they can be
customized and optimized.
