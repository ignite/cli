---
order: 13
---

# Codec File

To [register your types with Amino](https://github.com/tendermint/go-amino#registering-types) so that they can be encoded/decoded, there is a bit of code that needs to be placed in `./x/nameservice/types/codec.go`. Any interface you create and any struct that implements an interface needs to be declared in the `RegisterCodec` function. In this module the three `Msg` implementations (`SetName`, `BuyName` and `DeleteName`) have been registered, but your `Whois` query return type needs to be registered.

<<< @/nameservice/nameservice/x/nameservice/types/codec.go

### Next you need to define CLI interactions with your module.
