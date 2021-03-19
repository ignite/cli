# IBC Scaffold

Starport supports IBC-specifc scaffolding.

Creating a Cosmos SDK module with IBC logic:

```
starport module create ibcdex --ibc
```

Creating a custom packet:

```
starport packet buyOrder amountDenom amount:int priceDenom price:int --ack remainingAmount:int,purchase:int --module ibcdex
```