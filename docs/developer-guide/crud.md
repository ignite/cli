---
order: 8
description: Implementing CRUD manually
---

# Type Scaffold Reference

The `starport scaffold list` command scaffolds files that implement create, read, update, and delete (CRUD) functionality for a custom new type.

To add a custom type with create, read, update, and delete (CRUD) functionality:

```sh
starport scaffold type post title body
```

The `type` command scaffolds functionality a custom type.

## Command Overview

```sh
starport scaffold list [typeName] [field1] [field2] ... [flags]
```

`typeName` string

The name of a new type. Must be unique within a module.

`field1`, `field2`, and so on

Fields of the type. Define fields with a compact notation colon (`:`) syntax. For example, for an `amount` field that accepts an integer, use: `amount:int`. Supported types: `string`, `bool`, `int`, `uint`. By default, fields are `string`.

A type is scaffolded in a module.

`--module`

The name of the custom module in which a type is scaffolded. By default, a type is scaffolded in a module name that matches the project name.

## Files and Directories

The following files and directories are created and modified by scaffolding:

- `proto`: services for SDK messages and queries, HTTP endpoints
- `x/module_name/keeper`: gRPC message server and query handler
- `x/module_name/types`: message types, keys
- `x/module_name/client/cli`: CRUD actions on the CLI
- `x/module_name/client/rest`: legacy HTTP endpoints
- `vue/src/views`: Vue component, a CRUD form for interacting with the type

## CLI Commands

CLI commands are created for CRUD interactions with the type.

For example, if the binary is named `appd`, the module is `blog`, and the type is `post`, then the following transaction commands become available:

```sh
appd tx blog create-post [title] [content]
appd tx blog delete-post [id]
appd tx blog update-post [id] [title] [content]
```

Commands for querying:

```sh
appd q blog list-post
appd q blog show-post [id]
```

## Create Type Example

This example command creates a `post` type with four fields: `title` and `body` strings, boolean `comments` and integer `count`. This type is created in a module called blog.

```sh
starport scaffold list post title body comments:bool count:int32 --module blog
```
