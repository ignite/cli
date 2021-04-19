---
order: 2
---

# IBC Scaffold

Starport supports IBC-specifc scaffolding.

To create a Cosmos SDK module with IBC logic:

```
starport module create ibcdex --ibc
```

To create a custom packet:

```
starport packet buyOrder amountDenom amount:int priceDenom price:int --ack remainingAmount:int,purchase:int --module ibcdex
```
