---
title: In-depth tutorial
---

# In-depth "Hello, World!" tutorial

In this tutorial you will implement "Hello, World!" functionality from
scratch. The functionality of the application you will be building will be
identical to what the one you created in the "Express tutorial" section, but
here you will be doing it manually in order to gain a deeper understanding of
the process.

To begin, let's start with a fresh `hello` blockchain. You can either roll back
the changes you made in the previous section or create a new blockchain using
Ignite. Either way, you will have a blank blockchain that is ready for you to
work with.

```
ignite scaffold chain hello
```

## `SayHello` RPC

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

## `SayHello` keeper method

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

## `CmdSayHello` command

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