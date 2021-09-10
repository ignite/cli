---
description: Blockchain Client in Go
order: 10
---

# Creating a Blockchain Client in Go

**Learn how to connect your Blockchain to an independent application with RPC**

## Creating a Blockchain
Scaffold a new blockchain using `starport`:
```zsh
starport scaffold chain github.com/cosmonaut/blog
```

Scaffold create, read, update, delete functionality for a type `post` with two fields: `title` and `body`. Use `starport scaffold list` to scaffold code for storing posts in a list-like data structure.
```zsh
starport scaffold list post title body
```

Start a blockchain node in development:
```zsh
starport chain serve
```

## Creating a Blockchain Client
**In a new directory (outside the blockchain directory) create two files - `go.main` and `go.mod`**
```zsh
touch main.go && touch go.mod
```
**`go.main`**

- **Import files**

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cosmonaut/blog/x/blog/types"
	blogtypes "github.com/alijnmerchant21/blog/x/blog/types"
	"github.com/tendermint/starport/starport/pkg/cosmosclient"
)
```

- **Instance of cosmosclient.**

```go
func main() {

	cosmos, err := cosmosclient.New(context.Background())
	if err != nil {
		log.Fatal(err)	
	}
```
- **Using the account name Alice (obtained while scaffolding the blockchain), define the account address**

```go
accountName := "alice"
address, err := cosmos.Address(accountName)
	if err != nil {
		log.Fatal(err)
	}
```

- **`Create` a new post, `broadcast` it on the chain and `display` the results**

```go
// Create a post
msg := &types.MsgCreatePost {
	Creator: address.String(),
	Title: "Hello Ali",
	Body: "This is the first post",
}

// Broadcast the post
txResp, err := cosmos.BroadcastTx(accountName, msg)
if err != nil {
	log.Fatal(err)
}
fmt.Print("MsgCreatePost:\n\n")
fmt.Println(txResp)

// Display the results
queryClient := blogtypes.NewQueryClient(cosmos.Context)
queryResp, err := queryClient.Post(context.Background(), &types.QueryGetPostRequest{})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\n\nUser:\n\n")
	fmt.Println(queryResp)

```
**Tip:** To display the results of all users, replace `POST` with `POSTALL` as shown in this example code:

```go
queryResp, err := queryClient.PostAll(context.Background(), &types.QueryAllPostRequest{})
```

- **Add `go.mod` file**
```go
module github.com/tendermint/starport/local_test/client

go 1.17

require (
	github.com/alijnmerchant21/blog v0.0.0-20210825193134-b6859adfa282
	github.com/tendermint/starport v0.15.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/tendermint/starport => github.com/ilgooz/starport v0.0.500
```
***Tip:** Instead of `github.com/alijnmerchant21/blog v0.0.0-20210825193134-b6859adfa282` using your own `blog` package*

## Step 3
To run this program, you need follow these steps:
- Populate the file `main.go` with the code given above
- Run `go mod tidy`
- Run the main file `go run main.go`

**Output**
```zsh
MsgCreatePost:

Response:
  Height: 3522
  TxHash: E6904A959A4F62B361CAD9C1F8A6976E1B48EB4DA77A1FBB0D41DE921F1E8028
  Data: 0A100A0A437265617465506F737412020801
  Raw Log: [{"events":[{"type":"message","attributes":[{"key":"action","value":"CreatePost"}]}]}]
  Logs: [{"events":[{"type":"message","attributes":[{"key":"action","value":"CreatePost"}]}]}]
  GasWanted: 300000
  GasUsed: 45764


User:

Post:<creator:"cosmos1u27xu76zamjzgus2py87r8kmafea8cp6rvgtv9" 
title:"Hello Ali" body:"This is the first post" >
```

## Step 4
Similaryly, you can `update` and `delete` the POST
```go
// Update a Post
	msgUpd := &types.MsgUpdatePost{
		Creator: address.String(),
		Id:      0,
		Title:   "Hello Ali",
		Body:    "This is a modified post",
	}
```

## Delete a Post

To delete a post:
	msgDel := &types.MsgDeletePost{
		Creator: address.String(),
		Id:      1,
	}
```

## Final Code in main.go 

After completing the steps in this article, your `main.go` file looks like:


```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alijnmerchant21/blog/x/blog/types"
	blogtypes "github.com/alijnmerchant21/blog/x/blog/types"
	"github.com/tendermint/starport/starport/pkg/cosmosclient"
)

func main() {

	cosmos, err := cosmosclient.New(context.Background())
	if err != nil {
		log.Fatal(err)	
	}

	accountName := "alice"
	address, err := cosmos.Address(accountName)
	if err != nil {
		log.Fatal(err)
	}

	// Create a post
	msg := &types.MsgCreatePost {
		Creator: address.String(),
		Title: "Hello Ali",
		Body: "This is the first post",
	}

	// Broadcast the post
	txResp, err := cosmos.BroadcastTx(accountName, msg)
	if err != nil {
	log.Fatal(err)
	}
	fmt.Print("MsgCreatePost:\n\n")
	fmt.Println(txResp)

	// Display the results
	queryClient := blogtypes.NewQueryClient(cosmos.Context)
	queryResp, err := queryClient.Post(context.Background(), &types.QueryGetPostRequest{})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\n\nUser:\n\n")
	fmt.Println(queryResp)

	// Tip: To display result of all users, replace POST with POSTALL
	queryResp, err := queryClient.PostAll(context.Background(), &types.QueryAllPostRequest{})

	// Delete a Post
	msgDel := &types.MsgDeletePost{
		Creator: address.String(),
		Id:      1,
	}

	// Broadcast Delete Post
	txResp, err := cosmos.BroadcastTx(accountName, msgDel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("MsgDeletePost:\n\n")
	fmt.Println(txResp)

	// Display the Results
	queryClient := blogtypes.NewQueryClient(cosmos.Context)
	queryRespAll, err := queryClient.PostAll(context.Background(), &types.QueryAllPostRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\n\nUser All:\n\n")
	fmt.Println(queryRespAll)

	// Update a Post
	msgUpd := &types.MsgUpdatePost{
		Creator: address.String(),
		Id:      0,
		Title:   "Hello Ali",
		Body:    "This is a modified post",
	}

	// Broadcast Update Post
	txResp, err := cosmos.BroadcastTx(accountName, msgUpd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("MsgUpdatePost:\n\n")
	fmt.Println(txResp)

	// Display the Results
	queryClient := blogtypes.NewQueryClient(cosmos.Context)
	queryRespAll, err := queryClient.PostAll(context.Background(), &types.QueryAllPostRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\n\nUser All:\n\n")
	fmt.Println(queryRespAll)
}
```
