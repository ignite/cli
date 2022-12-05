---
sidebar_position: 2
description: Step-by-step guidance to build your first blockchain and your first Cosmos SDK module.
---

# Hello, World!

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

To add a query, run the following command:

```
ignite scaffold query say-hello name --response name
```

The `ignite scaffold query` command is a tool used to quickly create new
queries. When you run this command, it makes changes to your source code to add
the new query and make it available in your API. This command accepts a query
name (`"say-hello"`) and a list of request fields (in our case only `name`). The
optional `--reponse` flag specifies the return values of the query.

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

## Create a query manually

To begin, let's start with a fresh `hello` blockchain. You can either roll back
the changes you made in the previous section or create a new blockchain using
Ignite. Either way, you will have a blank blockchain that is ready for you to
work with.

```
ignite scaffold chain hello
```

### `SayHello` RPC

In Cosmos SDK blockchains, queries are defined as remote procedure calls (RPCs)
in a `Query` service in protocol buffer files. To add a new query, you can add
the following code to the `query.proto` file of your module:

```protobuf title="proto/hello/hello/query.proto"
service Query {
	// highlight-start
	rpc SayHello(QuerySayHelloRequest) returns (QuerySayHelloResponse) {
		option (google.api.http).get = "/hello/hello/say_hello/{name}";
	}
	// highlight-end
}
```

The RPC accepts a request argument of type `QuerySayHelloRequest` and returns a
value of type `QuerySayHelloResponse`. To define these types, you can add the
following code to the `query.proto` file:

```protobuf title="proto/hello/hello/query.proto"
message QuerySayHelloRequest {
  string name = 1;
}

message QuerySayHelloResponse {
  string name = 1;
}
```

To use the types defined in `query.proto`, you must transpile the protocol
buffer files into Go source code. This can be done by running `ignite chain
serve`, which will build and initialize the blockchain and automatically
generate the Go source code from the protocol buffer files. Alternatively, you
can run `ignite generate proto-go` to only generate the Go source code from the
protocol buffer files, without building and initializing the blockchain.

### `SayHello` keeper method

After defining the query, request, and response types in the `query.proto` file,
you will need to implement the logic for the query in your code. This typically
involves writing a function that processes the request and returns the
appropriate response. Create a new file `query_say_hello.go` with the following
contents:

```go title="x/hello/keeper/query_say_hello.go"
package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	
	"hello/x/hello/types"
)

func (k Keeper) SayHello(goCtx context.Context, req *types.QuerySayHelloRequest) (*types.QuerySayHelloResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	// TODO: Process the query
	_ = ctx
	return &types.QuerySayHelloResponse{Name: fmt.Sprintf("hello %s", req.Name)}, nil
}
```

This code defines a `SayHello` function that accepts a request of type
`QuerySayHelloRequest` and returns a value of type `QuerySayHelloResponse`. The
function first checks if the request is valid, and then processes the query by
returning the response message with the provided name as the value for the `%s`
placeholder. You can add additional logic to the function as needed, such as
retrieving data from the blockchain or performing complex operations, to handle
the query and return the appropriate response.

### `CmdSayHello` command

After implementing the query logic, you will need to make the query available to
clients so that they can call it and receive the response. This typically
involves adding the query to the blockchain's application programming interface
(API) and providing a command-line interface (CLI) command that allows users to
easily submit the query and receive the response.

To provide a CLI command for the query, you can create the `query_say_hello.go`
file and implement a `CmdSayHello` command that calls the `SayHello` function
and prints the response to the console.

```go title="x/hello/client/cli/query_say_hello.go"
package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	
	"hello/x/hello/types"
)

var _ = strconv.Itoa(0)

func CmdSayHello() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "say-hello [name]",
		Short: "Query say-hello",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqName := args[0]
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QuerySayHelloRequest{
				Name: reqName,
			}
			res, err := queryClient.SayHello(cmd.Context(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
```

The code defines a `CmdSayHello` command. The command is defined using the
`cobra` library, which is a popular framework for building command-line
applications in Go. The command accepts a `name` as an argument and uses it to
create a `QuerySayHelloRequest` struct that is passed to the `SayHello` function
from the `types.QueryClient`. The `SayHello` function is used to send the
`say-hello` query to the blockchain, and the response is stored in the `res`
variable.

The `QuerySayHelloRequest` struct is defined in the `query.proto` file, which is
a Protocol Buffer file that defines the request and response types for the
query. The `QuerySayHelloRequest` struct includes a `Name` field of type
`string`, which is used to provide the name to be included in the response
message.

After the query has been sent and the response has been received, the code uses
the `clientCtx.PrintProto` function to print the response to the console. The
`clientCtx` variable is obtained using the `client.GetClientQueryContext`
function, which provides access to the client context, including the client's
configuration and connection information. The `PrintProto` function is used to
print the response using the Protocol Buffer format, which allows for efficient
serialization and deserialization of the data.

The `flags.AddQueryFlagsToCmd` function is used to add query-related flags to
the command. This allows users to specify additional options when calling the
command, such as the node URL and other query parameters. These flags are used
to configure the query and provide the necessary information to the `SayHello`
function, allowing it to connect to the blockchain and send the query.

To make the `CmdSayHello` command available to users, you will need to add it to
the chain's binary. This is typically done by modifying the
`x/hello/client/cli/query.go` file and adding the
`cmd.AddCommand(CmdSayHello())` statement. This adds the `CmdSayHello` command
to the list of available commands, allowing users to call it from the
command-line interface (CLI).

```go title="x/hello/client/cli/query.go"
func GetQueryCmd(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(CmdQueryParams())
	// highlight-next-line
	cmd.AddCommand(CmdSayHello())
	return cmd
}
```

Once you have provided a CLI command, users will be able to call the `say-hello`
query and receive the appropriate response.

Save all the changes you made to the source code of your project and run the
following command to start a blockchain node:

```
ignite chain serve
```

Use the following command to submit the query and receive the response:

```
hellod q hello say-hello bob
```

This command will send a "say-hello" query to the blockchain with the name "bob"
and print the response of "Hello, bob!" to the console. You can modify the query
and response as needed to suit your specific requirements and provide the
desired functionality.

Congratulations on completing the "Hello, World!" tutorial! In this tutorial,
you learned how to define a new query in a protocol buffer file, implement the
logic for the query in your code, and make the query available to clients
through the blockchain's API and CLI. By following the steps outlined in the
tutorial, you were able to create a functional query that can be used to
retrieve data from your blockchain or perform other operations as needed.

Now that you have completed the tutorial, you can continue to build on your
knowledge of the Cosmos SDK and explore the many features and capabilities it
offers. You may want to try implementing more complex queries or experiment with
other features of the SDK to see what you can create.