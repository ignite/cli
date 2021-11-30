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
| coin         | -        | sdk.Coin    | Cosmos-SDK coin type             |
| array.coin   | coins    | sdk.Coins   | List of Cosmos-SDK coin type     |

## Custom Type Scaffold

Starport allows to use previously scaffolded fields. The developer can create a `list` type called `user` and use this type in the next scaffold type.

### Scaffolding

We scaffold a new `CoordinatorDescription` type to be reusable in the future:
```shell
starport scaffold list coordinator-description description:string --no-message
```
Now we can scaffold a message using the `CoordinatorDescription` type:
```shell
starport scaffold map coordinator description:CoordinatorDescription address:string --no-message
```
To send the message using the CLI, we should pass the custom type as a JSON:
```shell
testd tx test settle-request '{"description":"coordinator description"}' true --from alice --chain-id mars
```
If the developer tries to use another type not created yet, the starport fails:
```shell
starport scaffold message validator validator:ValidatorDescription address:string
-> the field type ValidatorDescription doesn't exist
```
