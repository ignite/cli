---
order: 14
description: Starport scaffolding types
---

# Starport scaffolding types

Starport can support different types, like a number, string, bool, etc. This document shows all current types supported into the `starport`.

## Built-in Types

| Type         | Alias    | Code Type   | Description                      |
| ------------ | -------- | ----------- | -------------------------------- |
| string       | -        | string      | Text type                        |
| array.string | strings  | []string    | List of text type                |
| bool         | -        | bool        | Boolean type                     |
| int          | -        | int32       | Integer numbers                  |
| array.int    | ints     | []int32     | List of integer numbers          |
| uint         | -        | uint64      | Unsigned integer numbers         |
| array.uint   | uints    | []uint64    | List of unsigned integer numbers |
| coin         | -        | sdk.Coin    | Cosmos-sdk coin type             |
| array.coin   | coins    | sdk.Coins   | List of Cosmos-sdk coin type     |

## Custom Type Scaffold

Starport gives the to use custom field scaffolded before into the chain. The developer can create a `list` type called `user` and use this type in the next scaffold type.

### Scaffolding

We will scaffold a new `Coordinator` type to be reusable in the future:
```shell
starport scaffold list coordinator address:string id:uint --no-message
```
Now we can scaffold a message using the `Coordinator` type:
```shell
starport scaffold message blocks owner:Coordinator approve:bool
```
To send the message using the CLI, we should pass the custom type as a JSON:
```shell
testd tx test settle-request '{"id":100,"address":"cosmos1t4jkut0yfnsmqle9vxk3adfwwm9vj9gsj98vqf","validatorID":33}' true --from alice --chain-id test
```
If the developer tries to use another type not created yet, the starport fails:
```shell
starport scaffold message settle-request validator:Validator approve:bool
-> the field type Validator doesn't exist
```
