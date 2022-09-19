---
sidebar_position: 2
description: Blockchain client in Go
---

# Create a blockchain client in Go

Learn how to connect your blockchain to an independent application with RPC.

After creating the blog blockchain in this tutorial you will learn how to connect to your blockchain from a separate client.

## Use the blog blockchain

Navigate to a separate directory right next to the `blog` blockchain you built in the [Build a Blog](index.md) tutorial.

## Creating a blockchain client

Create a new directory called `blogclient` on the same level as `blog` directory. As the name suggests, `blogclient` will contain a standalone Go program that acts as a client to your `blog` blockchain.

The command:

```bash
ls
```

Shows just `blog` now. More results are listed when you have more directories here.

Create your `blogclient` directory first, change your current working directory, and initialize the new Go module.

```bash
mkdir blogclient
cd blogclient
go mod init blogclient
touch main.go
```

The `go.mod` file is created inside your `blogclient` directory.

Your blockchain client has only two dependencies: 

- The `blog` blockchain `types` for message types and a query client
- `ignite` for the `cosmosclient` blockchain client

```go-module
module blogclient

go 1.18

require (
	blog v0.0.0-00010101000000-000000000000
	github.com/ignite/cli v0.23.0
)

replace blog => ../blog
replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
```

The `replace` directive uses the package from the local `blog` directory and is specified as a relative path to the `blogclient` directory.

Cosmos SDK uses a custom version of the `protobuf` package, so use the `replace` directive to specify the correct dependency.

The `blogclient` will eventually have only two files: 

- `main.go` for the main logic of the client
- `go.mod` for specifying dependencies.

### Main logic of the client in `main.go`

Add the following code to your `main.go` file to make a connection to your blockchain from a separate app.

```go
package main

import (
	"context"
	"fmt"
	"log"

	// Importing the general purpose Cosmos blockchain client
	"github.com/ignite/cli/ignite/pkg/cosmosclient"

	// Importing the types package of your blog blockchain
	"blog/x/blog/types"
)

func main() {
	// Prefix to use for account addresses.
	// The address prefix was assigned to the blog blockchain
	// using the `--address-prefix` flag during scaffolding.
	addressPrefix := "blog"

	// Create a Cosmos client instance
	cosmos, err := cosmosclient.New(
		context.Background(),
		cosmosclient.WithAddressPrefix(addressPrefix),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Account `alice` was initialized during `ignite chain serve`
	accountName := "alice"

	// Get account from the keyring
	account, err := cosmos.Account(accountName)
	if err != nil {
		log.Fatal(err)
	}

	addr, err := account.Address(addressPrefix)
	if err != nil {
		log.Fatal(err)
	}

	// Define a message to create a post
	msg := &types.MsgCreatePost{
		Creator: addr,
		Title:   "Hello!",
		Body:    "This is the first post",
	}

	// Broadcast a transaction from account `alice` with the message
	// to create a post store response in txResp
	txResp, err := cosmos.BroadcastTx(account, msg)
	if err != nil {
		log.Fatal(err)
	}

	// Print response from broadcasting a transaction
	fmt.Print("MsgCreatePost:\n\n")
	fmt.Println(txResp)

	// Instantiate a query client for your `blog` blockchain
	queryClient := types.NewQueryClient(cosmos.Context())

	// Query the blockchain using the client's `Posts` method
	// to get all posts store all posts in queryResp
	queryResp, err := queryClient.Posts(context.Background(), &types.QueryPostsRequest{})
	if err != nil {
		log.Fatal(err)
	}

	// Print response from querying all the posts
	fmt.Print("\n\nAll posts:\n\n")
	fmt.Println(queryResp)
}
```

Read the comments in the code carefully to learn details about each line of code.

To learn more about the `cosmosclient` package, see the Go 
[cosmosclient](https://pkg.go.dev/github.com/ignite/cli/ignite/pkg/cosmosclient) package documentation. Details are provided to learn how to use the `Client` type with `Options` and `KeyringBackend`.

## Run the blockchain and the client

Make sure your blog blockchain is still running with `ignite chain serve`.

Install dependencies for your `blogclient`:

```bash
go mod tidy
```

Run the blockchain client:

```bash
go run main.go
```

If successful, the results of running the command are printed to the terminal:

```
# github.com/keybase/go-keychain
### Some warnings might be displayed which can be ignored
MsgCreatePost:

Response:
  Height: 3222
  TxHash: AFCA76B0FEE5113382C068967B610180C105FCE045FF8C7943EA45EF4B7A1E69
  Data: 0A280A222F636F736D6F6E6175742E626C6F672E626C6F672E4D7367437265617465506F737412020801
  Raw Log: [{"events":[{"type":"message","attributes":[{"key":"action","value":"CreatePost"}]}]}]
  Logs: [{"events":[{"type":"message","attributes":[{"key":"action","value":"CreatePost"}]}]}]
  GasWanted: 300000
  GasUsed: 45805


All posts:

Post:<creator:"blog1j8d8pyjr5vynjvcq7xgzme0ny6ha30rpakxk3n" title:"foo" body:"bar" > Post:<creator:"blog1j8d8pyjr5vynjvcq7xgzme0ny6ha30rpakxk3n" id:1 title:"Hello!" body:"This is the first post" > pagination:<total:2 > 
```

You can confirm the new post with using the `blogd query blog posts` command that you learned about in the previous chapter.
The result looks similar to:

```yaml
Post:
- body: bar
  creator: blog1j8d8pyjr5vynjvcq7xgzme0ny6ha30rpakxk3n
  id: "0"
  title: foo
- body: This is the first post
  creator: blog1j8d8pyjr5vynjvcq7xgzme0ny6ha30rpakxk3n
  id: "1"
  title: Hello!
pagination:
  next_key: null
  total: "2"
```

Congratulations, you have just created a post using a separate app.
