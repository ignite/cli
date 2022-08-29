---
sidebar_position: 5
description: Reference list of supported types. 
---

# Ignite CLI Supported Types

Types with CRUD operations are scaffolded with the `ignite scaffold` command. 

## Built-in types

| Type         | Alias    | Index | Code Type   | Description                     |
| ------------ | -------- | ----- | ----------- | ------------------------------- |
| string       | -        | yes   | string      | Text type                       |
| array.string | strings  | no    | []string    | List of text type               |
| bool         | -        | yes   | bool        | Boolean type                    |
| int          | -        | yes   | int32       | Integer type                    |
| array.int    | ints     | no    | []int32     | List of integers types          |
| uint         | -        | yes   | uint64      | Unsigned integer type           |
| array.uint   | uints    | no    | []uint64    | List of unsigned integers types |
| coin         | -        | no    | sdk.Coin    | Cosmos SDK coin type            |
| array.coin   | coins    | no    | sdk.Coins   | List of Cosmos SDK coin types   |

Some types cannot be used an index, like the map and list indexes and module params.

## Custom types

You can create custom types and then use the custom type later.

For example, you can create a `list` type called `user` and then use the `user` type in a subsequent `ignite scaffold` command.

Here's an example of how to scaffold a new `CoordinatorDescription` type that is reusable in the future:

```bash
ignite scaffold list coordinator-description description:string --no-message
```

Now you can scaffold a message using the `CoordinatorDescription` type:

```bash
ignite scaffold message add-coordinator address:string description:CoordinatorDescription
```

Run the chain and then send the message using the CLI.

To pass the custom type in JSON format:

```bash
ignite chain serve
marsd tx mars add-coordinator cosmos1t4jkut0yfnsmqle9vxk3adfwwm9vj9gsj98vqf '{"description":"coordinator description"}' true --from alice --chain-id mars
```

If you try to use a type that is not created yet, the follow error occurs:

```bash
ignite scaffold message validator validator:ValidatorDescription address:string
-> the field type ValidatorDescription doesn't exist
```
