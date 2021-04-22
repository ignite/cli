---
order: 6
description: IBC packet data structure and packet semantic reference.
---

# Packet Scaffold

Packets are sent to other blockchains using Inter-Blockchain Communication [(IBC) channels](https://docs.cosmos.network/master/ibc/overview.html). An IBC packet is a data structure with sequence-related metadata and an opaque value field referred to as the packet data. The packet semantics are defined by the application layer, for example, token amount and denomination.

## IBC Module Packet Scaffold

Packets can be be scaffolded only in IBC modules.

To scaffold a packet:

```
starport packet [packetName] [field1] [field2] --module [module_name] [flags]
```

### Custom acknowledgement types

Define acknowledgement fields with the `--ack` flag and a comma-separated list (no spaces) of fields that describe the acknowledgement fields.

## Files and Directories

When you scaffold a packet, the following files and directories are created and modified:

- `proto`: packet data and acknowledgement type and message type
- `x/module_name/keeper`: IBC hooks, gRPC message server
- `x/module_name/types`: message types, IBC events
- `x/module_name/client/cli`: CLI command to broadcast a transaction containing a message with a packet

## Packet Scaffold Example

The following command scaffolds the IBC-enabled `buyOrder` packet for the `amountDenom` and `remainingAmount` fields with custom acknowledgements for the `remainingAmount:int` and `purchase:int` fields.

```
starport packet buyOrder amountDenom amount:int priceDenom price:int --ack remainingAmount:int,purchase:int --module ibcdex
```
