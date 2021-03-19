# Packet Scaffold

An IBC packet is a data structure with sequence-related metadata and an opaque value field referred to as the packet data (with semantics defined by the application layer, e.g. token amount and denomination). Packets are sent through IBC channels.

Starport knows how to scaffold IBC packets.

```
starport packet [packetName] [field1] [field2]
```

Optional `--ack`: list of comma-separated (no spaces) fields that describe acknowledgement fields.

Files and directories created and modified by scaffolding:

* `proto`: packet data and acknowledgement type and message type
* `x/module_name/keeper`: IBC hooks, gRPC message server
* `x/module_name/types`: message types, IBC events
* `x/module_name/client/cli`: CLI command to broadcast a transaction containing a message with a packet

Packets can only be scaffolded in IBC modules.

## Example

```
starport packet buyOrder amountDenom amount:int priceDenom price:int --ack remainingAmount:int,purchase:int --module ibcdex
```