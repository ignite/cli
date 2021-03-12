# Type Scaffold Reference

You can scaffold types within Starport by running a command:
```
starport type [type-name] [field1:type1] [field2:type2] ...
```
<<!-- why do we scaffold types? what is the module string? is this the Cosmos SDK module? and this is how we add modules to a blockchain? I am sure this is explained somewhere, where can I learn? -->


## Stargate

A project that is scaffolded with the Stargate version of Cosmos SDK creates the following files: <!-- how does this tie in to types? -->

```
.
├── proto
│   └── module
│       └── v1beta
│           └── {{typeName}}.proto
└── x
    └── module
        ├── client
        │   ├── cli
        │   │   ├── query{{TypeName}}.go
        │   │   └── tx{{TypeName}}.go
        │   └── rest
        │       ├── query{{TypeName}}.go
        │       └── tx{{TypeName}}.go
        ├── keeper
        │   ├── grpc_query_{{typeName}}.go
        │   ├── querier_{{typeName}}.go
        │   └── {{typeName}}.go
        ├── types
        │   ├── MsgCreate{{TypeName}}.go
        │   └── {{typeName}}.pb.go
        └── handlerMsgCreate{{TypeName}}.go
```

The following existing files are updated:

```
.
├── proto
│   └── module
│       └── v1beta
│           └── querier.proto
└── x
    └── module
        ├── client
        │   ├── cli
        │   │   ├── query.go
        │   │   └── tx.go
        │   └── rest
        │       └── rest.go
        ├── handler.go
        ├── keeper
        │   └── querier.go
        └── types
            ├── codec.go
            ├── keys.go
            ├── querier.go
            └── querier.pb.go
```

# Launchpad

A project that is scaffolded with the Launchpad version of Cosmos SDK creates the following files:

```
.
└── x
    └── module
        ├── client
        │   ├── cli
        │   │   ├── query{{TypeName}}.go
        │   │   └── tx{{TypeName}}.go
        │   └── rest
        │       ├── query{{TypeName}}.go
        │       └── tx{{TypeName}}.go
        ├── handlerMsgCreate{{TypeName}}.go
        ├── handlerMsgDelete{{TypeName}}.go
        ├── handlerMsgSet{{TypeName}}.go
        ├── keeper
        │   └── {{typeName}}.go
        └── types
            ├── MsgCreate{{TypeName}}.go
            ├── MsgDelete{{TypeName}}.go
            ├── MsgSet{{TypeName}}.go
            └── Type{{TypeName}}.go
```

The following existing files are updated:

```
.
├── vue
│   └── src
│       └── views
│           └── Index.vue
└── x
    └── module
        ├── client
        │   ├── cli
        │   │   ├── query.go
        │   │   └── tx.go
        │   └── rest
        │       └── rest.go
        ├── handler.go
        ├── keeper
        │   └── querier.go
        └── types
            ├── codec.go
            ├── key.go
            └── querier.go
```
