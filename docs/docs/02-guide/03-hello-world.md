---
description: Build your first blockchain and your first Cosmos SDK query.
title: Hello World
---

# "Hello world!" Blockchain Tutorial with Ignite CLI

**Introduction**

In this tutorial, you'll build a simple blockchain using Ignite CLI that responds to a custom query with "Hello %s!", where "%s" is a name passed in the query. 
This will enhance your understanding of creating custom queries in a Cosmos SDK blockchain.

## Setup and Scaffold

1. **Create a New Blockchain:**

```bash
ignite scaffold chain hello
```

2. **Navigate to the Blockchain Directory:**

```bash
cd hello
```

## Adding a Custom Query

- **Scaffold the Query:**

```bash
ignite scaffold query say-hello name --response name
```

This command generates code for a new query, `say-hello`, which accepts a name, an input, and returns it in the response.

- **Understanding the Scaffolded Code:**

	- `proto/hello/hello/query.proto`: Defines the request and response structure.
	- `x/hello/client/cli/query_say_hello.go`: Contains the CLI commands for the query.
	- `x/hello/keeper/query_say_hello.go`: Houses the logic for the query response.


## Customizing the Query Response

In the Cosmos SDK, queries are requests for information from the blockchain, used to access data like the ledger's current state or transaction details. While the SDK offers several built-in query methods, developers can also craft custom queries for specific data retrieval or complex operations. 

- **Modify `query_say_hello.go`:**
  
Update the `SayHello` function in `x/hello/keeper/query_say_hello.go` to return a personalized greeting query.

```go title="x/hello/keeper/query_say_hello.go"
package keeper

import (
	"context"
	"fmt"

	"hello/x/hello/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) SayHello(ctx context.Context, req *types.QuerySayHelloRequest) (*types.QuerySayHelloResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Validation and Context unwrapping
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	_ = sdkCtx
	// Custom Response
	return &types.QuerySayHelloResponse{Name: fmt.Sprintf("Hello %s!", req.Name)}, nil
}
```

## Running the Blockchain

1. **Start the Blockchain:**

```bash
ignite chain serve
```

2. **Test the Query:**
   
Use the command-line interface to submit a query.

```
hellod q hello say-hello world
```

Expect a response: `Hello world!`

## Conclusion

Congratulations! 🎉 You've successfully created a blockchain module with a custom query using Ignite CLI. Through this tutorial, you've learned how to scaffold a chain, add a custom query, and modify the logic for personalized responses. This experience illustrates the power of Ignite CLI in streamlining blockchain development and the importance of understanding the underlying code for customization.