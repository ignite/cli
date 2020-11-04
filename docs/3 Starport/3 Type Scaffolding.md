# Type Scaffolding

You can scaffold types within Starport by running a command:
```
starport type [type-name] [field1:type1] [field2:type2] ...
```



## Stargate

This will create the following files:

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

As well as update the following files:

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

Using `starport type` on a Launchpad application will create the following files:
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

As well as update the following files:

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
