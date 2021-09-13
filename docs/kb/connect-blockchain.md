---
description: Blockchain Client in Go
order: 10
---

# Creating a Blockchain Client in Go

Learn how to connect your Blockchain to an independent application with RPC.

## Creating a Blockchain

Scaffold a new blockchain using `starport`:

```
starport scaffold chain github.com/cosmonaut/blog
```

Scaffold create, read, update, delete functionality for a type `post` with two fields: `title` and `body`. Use `starport scaffold list` to scaffold code for storing posts in a list-like data structure.

```
starport scaffold list post title body
```

Start a blockchain node in development:

```
starport chain serve
```

## Creating a Blockchain Client

Create a new directory called `blogclient` on the same level as the `blog` directory. As the name suggests, `blogclient` will contain a standalone Go program that will act as a client to your `blog` blockchain.

`blogclient` will have two files: `main.go` for the main logic of the client and `go.mod` for specifying dependencies.

### Main Logic of the Client in main.go

```go
package main

import (
	"context"
	"fmt"
	"log"

	// importing the types package of your blog blockchain
	"github.com/cosmonaut/blog/x/blog/types"
	// importing the general purpose Cosmos blockchain client
	"github.com/tendermint/starport/starport/pkg/cosmosclient"
)

func main() {
	// create an instance of cosmosclient
	cosmos, err := cosmosclient.New(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// account `alice` was initialized during `starport chain serve`
	accountName := "alice"
	// get account from the keyring by account name and return a bech32 address
	address, err := cosmos.Address(accountName)
	if err != nil {
		log.Fatal(err)
	}
	// define a message to create a post
	msg := &types.MsgCreatePost{
		Creator: address.String(),
		Title:   "Hello!",
		Body:    "This is the first post",
	}
	// broadcast a transaction from account `alice` with the message to create a post
	// store response in txResp
	txResp, err := cosmos.BroadcastTx(accountName, msg)
	if err != nil {
		log.Fatal(err)
	}
	// print response from broadcasting a transaction
	fmt.Print("MsgCreatePost:\n\n")
	fmt.Println(txResp)
	// instantiate a query client for your `blog` blockchain
	queryClient := types.NewQueryClient(cosmos.Context)
	// query the blockchain using the client's `PostAll` method to get all posts
	// store all posts in queryResp
	queryResp, err := queryClient.PostAll(context.Background(), &types.QueryAllPostRequest{})
	if err != nil {
		log.Fatal(err)
	}
	// print response from querying all the posts
	fmt.Print("\n\nAll posts:\n\n")
	fmt.Println(queryResp)
}
```

### Specifying Dependencies in go.mod

Your blockchain client has only two dependencies: the `blog` blockchain (`types` for message types and a query client) and `starport` (for the `cosmosclient` blockchain client).

```go
module blogclient

go 1.17

require (
	github.com/cosmonaut/blog v0.0.0
	github.com/tendermint/starport v0.15.0
)

replace github.com/cosmonaut/blog => ../blog
replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
// DON'T FORGET TO REMOVE BEFORE THE STARPORT RELEASE!
replace github.com/tendermint/starport => github.com/ilgooz/starport v0.0.500
```

Use the `replace` directive to use the package from the local `blog` directory (specified as a relative path). Skip this if you've pushed the source code for your blockchain to a location accessible online (like GitHub).

Cosmos SDK uses a custom version of the `protobuf` package, so use the `replace` directive to specify the correct dependency.

## Running the Blockchain and the Client

Start a blockchain node in development (if you haven't already):

```
cd blog

starport chain serve
```

Install dependencies:

```
cd blogclient

go mod tidy
```

Run the blockchain client:

```
go run main.go
```

If successful, the results of running the command are printed to the terminal:

```
MsgCreatePost:

Response:
  Height: 3522
  TxHash: E6904A959A4F62B361CAD9C1F8A6976E1B48EB4DA77A1FBB0D41DE921F1E8028
  Data: 0A100A0A437265617465506F737412020801
  Raw Log: [{"events":[{"type":"message","attributes":[{"key":"action","value":"CreatePost"}]}]}]
  Logs: [{"events":[{"type":"message","attributes":[{"key":"action","value":"CreatePost"}]}]}]
  GasWanted: 300000
  GasUsed: 45764


All posts:

Post:<title:"Hello!" body:"This is the first post" creator:"cosmos1v8xuglg3xjj50nt2k253nrppzrtvp2yjw40eae" >
```

### Updating and Deleting Posts

To update a post modify the message type to `MsgUpdatePost` and provide the correct `Id`:

```go
	msg := &types.MsgUpdatePost{
		Creator: address.String(),
		Id:      0,
		Title:   "Hello cosmonaut",
		Body:    "You can change the world by building blockchains.",
	}
```

To delete a post use the `MsgDeletePost` message type.

## Delete a Post

```go
	msg := &types.MsgDeletePost{
		Creator: address.String(),
		Id:      1,
	}
```