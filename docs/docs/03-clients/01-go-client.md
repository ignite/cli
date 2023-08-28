---
description: Blockchain client in Go
title: Go client
---

# A client in the Go programming language

In this tutorial, we will show you how to create a standalone Go program that
serves as a client for a blockchain. We will use the Ignite CLI to set up a
standard blockchain. To communicate with the blockchain, we will utilize the
`cosmosclient` package, which provides an easy-to-use interface for interacting
with the blockchain. You will learn how to use the `cosmosclient` package to
send transactions and query the blockchain. By the end of this tutorial, you
will have a good understanding of how to build a client for a blockchain using
Go and the `cosmosclient` package.

## Create a blockchain

To create a blockchain using the Ignite CLI, use the following command:

```
ignite scaffold chain blog
```

This will create a new Cosmos SDK blockchain called "blog".

Once the blockchain has been created, you can generate code for a "blog" model
that will enable you to perform create, read, update, and delete (CRUD)
operations on blog posts. To do this, you can use the following command:

```
cd blog
ignite scaffold list post title body
```

This will generate the necessary code for the "blog" model, including functions
for creating, reading, updating, and deleting blog posts. With this code in
place, you can now use your blockchain to perform CRUD operations on blog posts.
You can use the generated code to create new blog posts, retrieve existing ones,
update their content, and delete them as needed. This will give you a fully
functional Cosmos SDK blockchain with the ability to manage blog posts.

Start your blockchain node with the following command:

```
ignite chain serve
```

## Creating a blockchain client

Create a new directory called `blogclient` on the same level as `blog`
directory. As the name suggests, `blogclient` will contain a standalone Go
program that acts as a client to your `blog` blockchain.

```bash
mkdir blogclient
```

This command will create a new directory called `blogclient` in your current
location. If you type `ls` in your terminal window, you should see both the
`blog` and `blogclient` directories listed.

To initialize a new Go package inside the `blogclient` directory, you can use
the following command:

```
cd blogclient
go mod init blogclient
```

This will create a `go.mod` file in the `blogclient` directory, which contains
information about the package and the Go version being used.

To import dependencies for your package, you can add the following code to the
`go.mod` file:

```text title="blogclient/go.mod"
module blogclient

go 1.20

require (
	blog v0.0.0-00010101000000-000000000000
	github.com/ignite/cli v0.25.2
)

replace blog => ../blog
replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
```

Your package will import two dependencies:

* `blog`, which contains `types` of messages and a query client
* `ignite` for the `cosmosclient` package

The `replace` directive uses the package from the local `blog` directory and is
specified as a relative path to the `blogclient` directory.

Cosmos SDK uses a custom version of the `protobuf` package, so use the `replace`
directive to specify the correct dependency.

Finally, install dependencies for your `blogclient`:

```bash
go mod tidy
```

### Main logic of the client in `main.go`

Create a `main.go` file inside the `blogclient` directory and add the following
code:

```go title="blogclient/main.go"
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
	ctx := context.Background()
	addressPrefix := "cosmos"

	// Create a Cosmos client instance
	client, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix))
	if err != nil {
		log.Fatal(err)
	}

	// Account `alice` was initialized during `ignite chain serve`
	accountName := "alice"

	// Get account from the keyring
	account, err := client.Account(accountName)
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
	txResp, err := client.BroadcastTx(ctx, account, msg)
	if err != nil {
		log.Fatal(err)
	}

	// Print response from broadcasting a transaction
	fmt.Print("MsgCreatePost:\n\n")
	fmt.Println(txResp)

	// Instantiate a query client for your `blog` blockchain
	queryClient := types.NewQueryClient(client.Context())

	// Query the blockchain using the client's `PostAll` method
	// to get all posts store all posts in queryResp
	queryResp, err := queryClient.PostAll(ctx, &types.QueryAllPostRequest{})
	if err != nil {
		log.Fatal(err)
	}

	// Print response from querying all the posts
	fmt.Print("\n\nAll posts:\n\n")
	fmt.Println(queryResp)
}
```

The code above creates a standalone Go program that acts as a client to the
`blog` blockchain. It begins by importing the required packages, including the
general purpose Cosmos blockchain client and the `types` package of the `blog`
blockchain.

In the `main` function, the code creates a Cosmos client instance and sets the
address prefix to "cosmos". It then retrieves an account named `"alice"` from
the keyring and gets the address of the account using the address prefix.

Next, the code defines a message to create a blog post with the title "Hello!"
and body "This is the first post". It then broadcasts a transaction from the
account "alice" with the message to create the post, and stores the response in
the variable `txResp`.

The code then instantiates a query client for the blog blockchain and uses it to
query the blockchain to retrieve all the posts. It stores the response in the
variable `queryResp` and prints it to the console.

Finally, the code prints the response from broadcasting the transaction to the
console. This allows the user to see the results of creating and querying a blog
post on the `blog` blockchain using the client.

To find out more about the `cosmosclient` package, you can refer to the Go
package documentation for
[`cosmosclient`](https://pkg.go.dev/github.com/ignite/cli/ignite/pkg/cosmosclient).
This documentation provides information on how to use the `Client` type with
`Options` and `KeyringBackend`.

## Run the blockchain and the client

Make sure your blog blockchain is still running with `ignite chain serve`.

Run the blockchain client:

```bash
go run main.go
```

If the command is successful, the results of running the command will be printed
to the terminal. The output may include some warnings, which can be ignored.

```yml
MsgCreatePost:

code: 0
codespace: ""
data: 12220A202F626C6F672E626C6F672E4D7367437265617465506F7374526573706F6E7365
events:
- attributes:
  - index: true
    key: ZmVl
    value: null
  - index: true
    key: ZmVlX3BheWVy
    value: Y29zbW9zMWR6ZW13NzZ3enQ3cDBnajd3MzQyN2E0eHg3MjRkejAzd3hnOGhk
  type: tx
- attributes:
  - index: true
    key: YWNjX3NlcQ==
    value: Y29zbW9zMWR6ZW13NzZ3enQ3cDBnajd3MzQyN2E0eHg3MjRkejAzd3hnOGhkLzE=
  type: tx
- attributes:
  - index: true
    key: c2lnbmF0dXJl
    value: UWZncUJCUFQvaWxWVzJwNUJNTngzcDlvRzVpSXp0elhXdE9yMHcwVE00OEtlSkRqR0FEdU9VNjJiY1ZRNVkxTHdEbXNuYUlsTmc3VE9uMnJ2ZWRHSlE9PQ==
  type: tx
- attributes:
  - index: true
    key: YWN0aW9u
    value: L2Jsb2cuYmxvZy5Nc2dDcmVhdGVQb3N0
  type: message
gas_used: "52085"
gas_wanted: "300000"
height: "20"
info: ""
logs:
- events:
  - attributes:
    - key: action
      value: /blog.blog.MsgCreatePost
    type: message
  log: ""
  msg_index: 0
raw_log: '[{"msg_index":0,"events":[{"type":"message","attributes":[{"key":"action","value":"/blog.blog.MsgCreatePost"}]}]}]'
timestamp: ""
tx: null
txhash: 4F53B75C18254F96EF159821DDD665E965DBB576A5AC2B94CE863EB62E33156A

All posts:

Post:<title:"Hello!" body:"This is the first post" creator:"cosmos1dzemw76wzt7p0gj7w3427a4xx724dz03wxg8hd" > pagination:<total:1 >
```

As you can see the client has successfully broadcasted a transaction and queried
the chain for blog posts.

Please note, that some values in the output on your terminal (like transaction
hash and block height) might be different from the output above.

You can confirm the new post with using the `blogd q blog list-post` command:

```yaml
Post:
- body: This is the first post
  creator: cosmos1dzemw76wzt7p0gj7w3427a4xx724dz03wxg8hd
  id: "0"
  title: Hello!
pagination:
  next_key: null
  total: "0"
```

Great job! You have successfully completed the process of creating a Go client
for your Cosmos SDK blockchain, submitting a transaction, and querying the
chain.