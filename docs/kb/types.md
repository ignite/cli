---
order: 6
description: Reference list of supported types. 
---

# Starport Supported Types

Types with CRUD operations are scaffolded with the `starport scaffold` command. 

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
| coin         | -        | sdk.Coin    | Cosmos SDK coin type             |
| array.coin   | coins    | sdk.Coins   | List of Cosmos SDK coin types     |

## Custom Types

You can create custom types and then use the custom type later. 

For example, you can create a `list` type called `user` and then use the `user` type in a subsequent `starport scaffold` command.

Here's an example of how to scaffold a new `CoordinatorDescription` type that is reusable in the future:

```shell
starport scaffold list coordinator-description description:string --no-message
```

Now you can scaffold a message using the `CoordinatorDescription` type:

```shell
starport scaffold message add-coordinator address:string description:CoordinatorDescription
```

Run the chain and then send the message using the CLI. 

To pass the custom type in JSON format:

```shell
starport chain serve
marsd tx mars add-coordinator cosmos1t4jkut0yfnsmqle9vxk3adfwwm9vj9gsj98vqf '{"description":"coordinator description"}' true --from alice --chain-id mars
```

If you try to use a type that is not created yet, the follow error occurs:

```shell
starport scaffold message validator validator:ValidatorDescription address:string
-> the field type ValidatorDescription doesn't exist
```
