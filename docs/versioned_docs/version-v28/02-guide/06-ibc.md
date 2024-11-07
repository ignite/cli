---
sidebar_position: 7
description: Build an understanding of how to create and send packets across blockchains and navigate between blockchains.
title: "Inter-Blockchain Communication: Basics"
---

# Inter-Blockchain Communication: Basics

The Inter-Blockchain Communication protocol (IBC) is an important part of the
Cosmos SDK ecosystem. The Hello World tutorial is a time-honored tradition in
computer programming. This tutorial builds an understanding of how to create and
send packets across blockchain. This foundational knowledge helps you navigate
between blockchains with the Cosmos SDK.

**You will learn how to**

- Use IBC to create and send packets between blockchains.
- Navigate between blockchains using the Cosmos SDK and the Ignite CLI Relayer.
- Create a basic blog post and save the post on another blockchain.

## What is IBC?

The Inter-Blockchain Communication protocol (IBC) allows blockchains to talk to
each other. IBC handles transport across different sovereign blockchains. This
end-to-end, connection-oriented, stateful protocol provides reliable, ordered,
and authenticated communication between heterogeneous blockchains.

The [IBC protocol in the Cosmos
SDK](https://ibc.cosmos.network/main/ibc/overview) is the standard for the
interaction between two blockchains. The IBCmodule interface defines how packets
and messages are constructed to be interpreted by the sending and the receiving
blockchain.

The IBC relayer lets you connect between sets of IBC-enabled chains. This
tutorial teaches you how to create two blockchains and then start and use the
relayer with Ignite CLI to connect two blockchains.

This tutorial covers essentials like modules, IBC packets, relayer, and the
lifecycle of packets routed through IBC.

## Create a blockchain

Create a blockchain app with a blog module to write posts on other blockchains
that contain the Hello World message. For this tutorial, you can write posts for
the Cosmos SDK universe that contain Hello Mars, Hello Cosmos, and Hello Earth
messages.

For this simple example, create an app that contains a blog module that has a
post transaction with title and text.

After you define the logic, run two blockchains that have this module installed.

- The chains can send posts between each other using IBC.

- On the sending chain, save the `acknowledged` and `timed out` posts.

After the transaction is acknowledged by the receiving chain, you know that the
post is saved on both blockchains.

- The sending chain has the additional data `postID`.

- Sent posts that are acknowledged and timed out contain the title and the
  target chain of the post. These identifiers
- are visible on the parameter `chain`. The following chart shows the lifecycle
  of a packet that travels through IBC.

![The Lifecycle of an IBC packet](./images/packet_sendpost.png)

## Build your blockchain app

Use Ignite CLI to scaffold the blockchain app and the blog module.

### Build a new blockchain

To scaffold a new blockchain named `planet`:

```bash
ignite scaffold chain planet --no-module
cd planet
```

A new directory named `planet` is created in your home directory. The `planet`
directory contains a working blockchain app.

### Scaffold the blog module inside your blockchain

Next, use Ignite CLI to scaffold a blog module with IBC capabilities. The blog
module contains the logic for creating blog posts and routing them through IBC
to the second blockchain.

To scaffold a module named `blog`:

```bash
ignite scaffold module blog --ibc
```

A new directory with the code for an IBC module is created in `planet/x/blog`.
Modules scaffolded with the `--ibc` flag include all the logic for the
scaffolded IBC module.

### Generate CRUD actions for types

Next, create the CRUD actions for the blog module types.

Use the `ignite scaffold list` command to scaffold the boilerplate code for the
create, read, update, and delete (CRUD) actions.

These `ignite scaffold list` commands create CRUD code for the following
transactions:

- Creating blog posts

```bash
ignite scaffold list post title content creator --no-message --module blog
```

- Processing acknowledgments for sent posts

```bash
ignite scaffold list sentPost postID:uint title chain creator --no-message --module blog
```

- Managing post timeouts

```bash
ignite scaffold list timeoutPost title chain creator --no-message --module blog
```

The scaffolded code includes proto files for defining data structures, messages,
messages handlers, keepers for modifying the state, and CLI commands.

### Ignite CLI Scaffold List Command Overview

```
ignite scaffold list [typeName] [field1] [field2] ... [flags]
```

The first argument of the `ignite scaffold list [typeName]` command specifies
the name of the type being created. For the blog app, you created `post`,
`sentPost`, and `timeoutPost` types.

The next arguments define the fields that are associated with the type. For the
blog app, you created `title`, `content`, `postID`, and `chain` fields.

The `--module` flag defines which module the new transaction type is added to.
This optional flag lets you manage multiple modules within your Ignite CLI app.
When the flag is not present, the type is scaffolded in the module that matches
the name of the repo.

When a new type is scaffolded, the default behavior is to scaffold messages that
can be sent by users for CRUD operations. The `--no-message` flag disables this
feature. Disable the messages option for the app since you want the posts to be
created upon reception of IBC packets and not directly created from a user's
messages.

### Scaffold a sendable and interpretable IBC packet

You must generate code for a packet that contains the title and the content of
the blog post.

The `ignite packet` command creates the logic for an IBC packet that can be sent
to another blockchain.

- The `title` and `content` are stored on the target chain.

- The `postID` is acknowledged on the sending chain.

To scaffold a sendable and interpretable IBC packet:

```bash
ignite scaffold packet ibcPost title content --ack postID:uint --module blog
```

Notice the fields in the `ibcPost` packet match the fields in the `post` type
that you created earlier.

- The `--ack` flag defines which identifier is returned to the sending
  blockchain.

- The `--module` flag specifies to create the packet in a particular IBC module.

The `ignite packet` command also scaffolds the CLI command that is capable of
sending an IBC packet:

```bash
planetd tx blog send-ibcPost [portID] [channelID] [title] [content]
```

## Modify the source code

After you create the types and transactions, you must manually insert the logic
to manage updates in the database. Modify the source code to save the data as
specified earlier in this tutorial.

### Add creator to the blog post packet

Start with the proto file that defines the structure of the IBC packet.

To identify the creator of the post in the receiving blockchain, add the
`creator` field inside the packet. This field was not specified directly in the
command because it would automatically become a parameter in the `SendIbcPost`
CLI command.

```protobuf title="proto/planet/blog/packet.proto"
message IbcPostPacketData {
  string title = 1;
  string content = 2;
  // highlight-next-line
  string creator = 3;
}
```

To make sure the receiving chain has content on the creator of a blog post, add
the `msg.Creator` value to the IBC `packet`.

- The content of the `sender` of the message is automatically included in
  `SendIbcPost` message.
- The sender is verified as the signer of the message, so you can add the
  `msg.Sender` as the creator to the new packet
- before it is sent over IBC.

```go title="x/blog/keeper/msg_server_ibc_post.go"
package keeper

func (k msgServer) SendIbcPost(goCtx context.Context, msg *types.MsgSendIbcPost) (*types.MsgSendIbcPostResponse, error) {
	// validate incoming message
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	if msg.Port == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid packet port")
	}

	if msg.ChannelID == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid packet channel")
	}

	if msg.TimeoutTimestamp == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid packet timeout")
	}

	// TODO: logic before transmitting the packet

	// Construct the packet
	var packet types.IbcPostPacketData

	packet.Title = msg.Title
	packet.Content = msg.Content
	// highlight-next-line
	packet.Creator = msg.Creator

	// Transmit the packet
	ctx := sdk.UnwrapSDKContext(goCtx)
	_, err := k.TransmitIbcPostPacket(
		ctx,
		packet,
		msg.Port,
		msg.ChannelID,
		clienttypes.ZeroHeight(),
		msg.TimeoutTimestamp,
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendIbcPostResponse{}, nil
}
```

### Receive the post

The methods for primary transaction logic are in the `x/blog/keeper/ibc_post.go`
file. Use these methods to manage IBC packets:

- `TransmitIbcPostPacket` is called manually to send the packet over IBC. This
  method also defines the logic before the packet is sent over IBC to another
  blockchain app.
- `OnRecvIbcPostPacket` hook is automatically called when a packet is received
  on the chain. This method defines the packet reception logic.
- `OnAcknowledgementIbcPostPacket` hook is called when a sent packet is
  acknowledged on the source chain. This method defines the logic when the
  packet has been received.
- `OnTimeoutIbcPostPacket` hook is called when a sent packet times out. This
  method defines the logic when the packet is not received on the target chain

You must modify the source code to add the logic inside those functions so that
the data tables are modified accordingly.

On reception of the post message, create a new post with the title and the
content on the receiving chain.

To identify the blockchain app that a message is originating from and who
created the message, use an identifier in the following format:

`<portID>-<channelID>-<creatorAddress>`

Finally, the Ignite CLI-generated AppendPost function returns the ID of the new
appended post. You can return this value to the source chain through
acknowledgment.

Append the type instance as `PostId` on receiving the packet:

- The context `ctx` is an [immutable data
  structure](https://docs.cosmos.network/main/core/context#go-context-package)
  that has header data from the transaction. See [how the context is
  initiated](https://github.com/cosmos/cosmos-sdk/blob/main/types/context.go#L71)
- The identifier format that you defined earlier
- The `title` is the Title of the blog post
- The `content` is the Content of the blog post

Then modify the `OnRecvIbcPostPacket` keeper function with the following code:

```go title="x/blog/keeper/ibc_post.go"
package keeper

func (k Keeper) OnRecvIbcPostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcPostPacketData) (packetAck types.IbcPostPacketAck, err error) {
	packetAck.PostId, err = k.PostSeq.Next(ctx)
	if err != nil {
		return packetAck, err
	}
	return packetAck, k.Post.Set(ctx, packetAck.PostId, types.Post{Title: data.Title, Content: data.Content})
}
```

### Receive the post acknowledgement

On the sending blockchain, store a `sentPost` so you know that the post has been
received on the target chain.

Store the title and the target to identify the post.

When a packet is scaffolded, the default type for the received acknowledgment
data is a type that identifies if the packet treatment has failed. The
`Acknowledgement_Error` type is set if `OnRecvIbcPostPacket` returns an error
from the packet.

```go title="x/blog/keeper/ibc_post.go"
package keeper

func (k Keeper) OnAcknowledgementIbcPostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcPostPacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		// We will not treat acknowledgment error in this tutorial
		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.IbcPostPacketAck
		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		seq, err := k.SentPostSeq.Next(ctx)
		if err != nil {
			return err
		}

		return k.SentPost.Set(ctx, seq,
			types.SentPost{
				PostId: packetAck.PostId,
				Title:  data.Title,
				Chain:  packet.DestinationPort + "-" + packet.DestinationChannel,
			},
		)
	default:
		return errors.New("the counter-party module does not implement the correct acknowledgment format")
	}
}
```

### Store information about the timed-out packet

Store posts that have not been received by target chains in `timeoutPost`
posts. This logic follows the same format as `sentPost`.

```go title="x/blog/keeper/ibc_post.go"
func (k Keeper) OnTimeoutIbcPostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcPostPacketData) error {
	seq, err := k.TimeoutPostSeq.Next(ctx)
	if err != nil {
		return err
	}

	return k.TimeoutPost.Set(ctx, seq,
		types.TimeoutPost{
			Title: data.Title,
			Chain: packet.DestinationPort + "-" + packet.DestinationChannel,
		},
	)
}
```

This last step completes the basic `blog` module setup. The blockchain is now
ready!

## Use the IBC modules

You can now spin up the blockchain and send a blog post from one blockchain app
to the other. Multiple terminal windows are required to complete these next
steps.

### Test the IBC modules

To test the IBC module, start two blockchain networks on the same machine. Both
blockchains use the same source code. Each blockchain has a unique chain ID.

One blockchain is named `earth` and the other blockchain is named `mars`.

The `earth.yml` and `mars.yml` files are required in the project directory:

```yaml title="earth.yml"
version: 1
build:
  proto:
    path: proto
    third_party_paths:
    - third_party/proto
    - proto_vendor
accounts:
- name: alice
  coins:
  - 1000token
  - 100000000stake
- name: bob
  coins:
  - 500token
  - 100000000stake
faucet:
  name: bob
  coins:
  - 5token
  - 100000stake
  host: 0.0.0.0:4500
genesis:
  chain_id: earth
validators:
- name: alice
  bonded: 100000000stake
  home: $HOME/.earth
```

```yaml title="mars.yml"
version: 1
build:
  proto:
    path: proto
    third_party_paths:
    - third_party/proto
    - proto_vendor
accounts:
- name: alice
  coins:
  - 1000token
  - 1000000000stake
- name: bob
  coins:
  - 500token
  - 100000000stake
faucet:
  name: bob
  coins:
  - 5token
  - 100000stake
  host: :4501
genesis:
  chain_id: mars
validators:
- name: alice
  bonded: 100000000stake
  app:
    api:
      address: :1318
    grpc:
      address: :9092
    grpc-web:
      address: :9093
  config:
    p2p:
      laddr: :26658
    rpc:
      laddr: :26659
      pprof_laddr: :6061
  home: $HOME/.mars
```

Open a terminal window and run the following command to start the `earth`
blockchain:

```bash
ignite chain serve -c earth.yml
```

Open a different terminal window and run the following command to start the
`mars` blockchain:

```bash
ignite chain serve -c mars.yml
```

### Remove Existing Relayer and Ignite CLI Configurations

If you previously used the relayer, follow these steps to remove exiting relayer
and Ignite CLI configurations:

- Stop your blockchains and delete previous configuration files:

  ```bash
  rm -rf ~/.ignite/relayer
  ```

If existing relayer configurations do not exist, the command returns `no matches
found` and no action is taken.

### Configure and start the relayer

First, add the Hermes relayer app.

```bash
ignite app install -g github.com/ignite/apps/hermes
```

and after configure the relayer.

```bash
ignite relayer hermes configure \                        
"earth" "http://localhost:26657" "http://localhost:9090" \
"mars" "http://localhost:26659" "http://localhost:9092" \
--chain-a-faucet "http://0.0.0.0:4500" \
--chain-b-faucet "http://0.0.0.0:4501" \
--chain-a-port-id "blog" \
--chain-b-port-id "blog" \
--channel-version "blog-1"
```

When prompted, press Enter to accept the default values for `Chain A Account` and
`Chain B Account`.

The output looks like:

```
Hermes config created at /Users/danilopantani/.ignite/relayer/hermes/earth_mars
? Chain earth doesn't have a default Hermes key. Type your mnemonic to continue or type enter to generate a new one: (optional) 
New mnemonic generated: danger plate flavor twist chimney myself sketch assist copy expand core tattoo ignore ensure quote mean forum carbon enroll gadget immense grab early maze
Chain earth key created
Chain earth relayer wallet: cosmos1jk6wmyl880j6t9vw6umy9v8ex0yhrfwgx0vv2d
New balance from faucet: 100000stake,5token
? Chain mars doesn't have a default Hermes key. Type your mnemonic to continue or type enter to generate a new one: (optional) 
New mnemonic generated: invest box icon session lens demise purse link boss dwarf give minimum jazz eye vocal seven sunset coach express want ask version anger ranch
Chain mars key created
Chain mars relayer wallet: cosmos1x9kt37c0sutanaqwy9gxpvq5990yt0qnpqntmp
New balance from faucet: 100000stake,5token
Client '07-tendermint-0' created (earth -> mars)
Client 07-tendermint-0' created (mars -> earth)
Connection 'earth (connection-0) <-> mars (connection-0)' created
Channel 'earth (channel-0) <-> mars (channel-0)' created
```

Now start the relayer:

```bash
ignite relayer hermes start "earth" "mars"
```

### Send packets

You can now send packets and verify the received posts:

```bash
planetd tx blog send-ibc-post blog channel-0 "Hello" "Hello Mars, I'm Alice from Earth" --from alice --chain-id earth --home ~/.earth
```

To verify that the post has been received on Mars:

```bash
planetd q blog list-post --node tcp://localhost:26659
```

The packet has been received:

```yaml
Post:
  - content: Hello Mars, I'm Alice from Earth
    creator: blog-channel-0-cosmos1aew8dk9cs3uzzgeldatgzvm5ca2k4m98xhy20x
    id: "0"
    title: Hello
pagination:
  next_key: null
  total: "1"
```

To check if the packet has been acknowledged on Earth:

```bash
planetd q blog list-sent-post
```

Output:

```yaml
SentPost:
  - chain: blog-channel-0
    creator: cosmos1aew8dk9cs3uzzgeldatgzvm5ca2k4m98xhy20x
    id: "0"
    postID: "0"
    title: Hello
pagination:
  next_key: null
  total: "1"
```

To test timeout, set the timeout time of a packet to 1 nanosecond, verify that
the packet is timed out, and check the timed-out posts:

```bash
planetd tx blog send-ibc-post blog channel-0 "Sorry" "Sorry Mars, you will never see this post" --from alice --chain-id earth --home ~/.earth --packet-timeout-timestamp 1
```

Check the timed-out posts:

```bash
planetd q blog list-timeout-post
```

Results:

```yaml
TimeoutPost:
  - chain: blog-channel-0
    creator: cosmos1fhpcsxn0g8uask73xpcgwxlfxtuunn3ey5ptjv
    id: "0"
    title: Sorry
pagination:
  next_key: null
  total: "2"
```

You can also send a post from Mars:

```bash
planetd tx blog send-ibc-post blog channel-0 "Hello" "Hello Earth, I'm Alice from Mars" --from alice --chain-id mars --home ~/.mars --node tcp://localhost:26659
```

List post on Earth:

```bash
planetd q blog list-post
```

Results:

```yaml
Post:
  - content: Hello Earth, I'm Alice from Mars
    creator: blog-channel-0-cosmos1xtpx43l826348s59au24p22pxg6q248638q2tf
    id: "0"
    title: Hello
pagination:
  next_key: null
  total: "1"
```

## Congratulations 🎉

By completing this tutorial, you've learned to use the Inter-Blockchain
Communication protocol (IBC).

Here's what you accomplished in this tutorial:

- Built two Hello blockchain apps as IBC modules
- Modified the generated code to add CRUD action logic
- Configured and used the Ignite CLI relayer to connect two blockchains with
  each other
- Transferred IBC packets from one blockchain to another
