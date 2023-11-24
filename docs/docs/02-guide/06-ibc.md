---
sidebar_position: 7
description: Explore the essentials of blockchain interoperability with our tutorial on creating and transmitting packets between networks using the Inter-Blockchain Communication protocol.
title: "Mastering IBC: Inter-Blockchain Connectivity Essentials"
---

# Mastering IBC: Inter-Blockchain Connectivity Essentials

The Inter-Blockchain Communication protocol (IBC) is an important part of the Cosmos SDK ecosystem. The Hello World tutorial is a time-honored tradition in computer programming. This tutorial builds an understanding of how to create and send packets across blockchain. This foundational knowledge helps you navigate between blockchains with the Cosmos SDK.

**You will learn how to**

- Master IBC for packet creation and inter-blockchain transmission.
- Navigate and link blockchains with Cosmos SDK and Ignite App Hermes Relayer.
- Create and manage basic blog posts, and transfer them to another blockchain.

## What is IBC?

IBC stands for [Inter-Blockchain Communication protocol]((https://ibc.cosmos.network/main/ibc/overview.html)), a core component for enabling blockchains to communicate with each other. This protocol manages transport and authentication, ensuring reliable and ordered data exchange between heterogeneous blockchains.

## Tutorial Overview

- **Build and Connect Blockchains:** Learn to set up two interconnected blockchains and understand the mechanics of sending and acknowledging packets via IBC.
- **Modules and Lifecycle of IBC Packets:** Explore the modules, packet lifecycle, and essential elements like IBC packets, relayer, and more.
- **Practical Application:** Create a simple blog module capable of posting "Hello World" messages, demonstrating the practical use of IBC in sending data across blockchains.

## Steps to Implement IBC

1. **Scaffold the Blockchain:**

Use Ignite CLI to create the blockchain app named `planet`

```bash
ignite scaffold chain planet --no-module && cd planet
```

2. **Build the Blog Module:**

```bash
ignite scaffold module blog --ibc
```

A new directory with the code for an IBC module is created in `planet/x/blog`.
Modules scaffolded with the `--ibc` flag include all the necessary logic for IBC to work.

3. **Generate CRUD Actions:**

- **Create Blog Posts:**

  ```bash
  ignite scaffold list post title content creator --no-message --module blog
  ```

- **Process Acknowledgments Posts:**

  ```bash
  ignite scaffold list sentPost postID title chain creator --no-message --module blog
  ```

- **Manage Post Timeouts:**

  ```bash
  ignite scaffold list timedoutPost title chain creator --no-message --module blog
  ```

Learn more about the [`ignite scaffold list` command](https://docs.ignite.com/nightly/references/cli#ignite-scaffold-list).

### Craft IBC Packets

1. **Scaffold a packet named `ibcPost`:**

The `title` and `content` are stored on the target chain.
The `postID` is acknowledged on the sending chain.

```bash
ignite scaffold packet ibcPost title content --ack postID --module blog
```


2. **Add creator to the blog post packet:**

Start with the proto file that defines the structure of the IBC packet.

```protobuf title="proto/planet/blog/packet.proto"
message IbcPostPacketData {
  string title = 1;
  string content = 2;
  // highlight-next-line
  string creator = 3;
}
```

To make sure the receiving chain has content on the creator of a blog post, add the `msg.Creator` value to the IBC `packet`.

```go title="x/blog/keeper/msg_server_ibc_post.go"
package keeper

import (
	// ...
	"planet/x/blog/types"
)

func (k msgServer) SendIbcPost(goCtx context.Context, msg *types.MsgSendIbcPost) (*types.MsgSendIbcPostResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: logic before transmitting the packet

	// Construct the packet
	var packet types.IbcPostPacketData

	packet.Title = msg.Title
	packet.Content = msg.Content
	// highlight-next-line
	packet.Creator = msg.Creator

	// Transmit the packet
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

3. **Receive the Post:**

In the `x/blog/keeper/ibc_post.go` file, make sure to import `"strconv"` below
`"errors"`:

```go title="x/blog/keeper/ibc_post.go"
import (
	//...

	"strconv"

// ...
)
```

Then modify the `OnRecvIbcPostPacket` keeper function with the following code:

```go title="x/blog/keeper/ibc_post.go"
func (k Keeper) OnRecvIbcPostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcPostPacketData) (packetAck types.IbcPostPacketAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	id := k.AppendPost(
		ctx,
		types.Post{
			Creator: packet.SourcePort + "-" + packet.SourceChannel + "-" + data.Creator,
			Title:   data.Title,
			Content: data.Content,
		},
	)

	packetAck.PostId = strconv.FormatUint(id, 10)

	return packetAck, nil
}
```

4. **Receive the Post Acknowledgement:**

On the sending blockchain, store a `sentPost` so you know that the post has been
received on the target chain.

```go title="x/blog/keeper/ibc_post.go"
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

		k.AppendSentPost(
			ctx,
			types.SentPost{
				Creator: data.Creator,
				PostId:  packetAck.PostId,
				Title:   data.Title,
				Chain:   packet.DestinationPort + "-" + packet.DestinationChannel,
			},
		)

		return nil
	default:
		return errors.New("the counter-party module does not implement the correct acknowledgment format")
	}
}
```

5. **Store Information About the Timed-Out Packet:**

Store posts that have not been received by target chains in `timedoutPost`
posts. This logic follows the same format as `sentPost`.

```go title="x/blog/keeper/ibc_post.go"
func (k Keeper) OnTimeoutIbcPostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcPostPacketData) error {
	k.AppendTimedoutPost(
		ctx,
		types.TimedoutPost{
			Creator: data.Creator,
			Title:   data.Title,
			Chain:   packet.DestinationPort + "-" + packet.DestinationChannel,
		},
	)

	return nil
}

```

This last step completes the basic `blog` module setup. The blockchain is now
ready!

## Relayer Configuration

Start two blockchain networks on the same machine. Both blockchains use the same source code. Each blockchain has a unique chain ID.

One blockchain is named `earth` and the other blockchain is named `mars`.

Create and setup the `earth.yml` and `mars.yml` files:

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

Start the `earth` blockchain:

```bash
ignite chain serve -c earth.yml
```

Start the `mars` blockchain:

```bash
ignite chain serve -c mars.yml
```

- **Install the Hermes Relayer App:**

```bash
ignite app add -g ($GOPATH)/src/github.com/ignite/apps/hermes
```

- **Configure the Relayer:**

```bash
ignite relayer hermes configure "earth" "http://localhost:26657" "http://localhost:9090" "mars" "http://localhost:26659" "http://localhost:9092"
```

### Interact and Test

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
planetd q blog list-timedout-post
```

Results:

```yaml
TimedoutPost:
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

By completing this tutorial, you've gained valuable insights into using IBC with Ignite CLI. You've learned how to set up interconnected blockchains, manage data transmission, and understand the practical applications of IBC in the Cosmos ecosystem.

Congratulations on embarking on this journey into the world of blockchain interoperability with Ignite CLI and Cosmos SDK! 🎉