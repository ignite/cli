# Module Scaffold

Modules are building blocks of Cosmos SDK blockchains. They encapsulate logic and allow sharing functionality between projects. Learn more about [building modules](https://github.com/cosmos/cosmos-sdk/tree/master/docs/building-modules).

Starport supports scaffolding SDK modules.

```
starport module create [name] [flags]
```

- **name**

     The name of a new module. This name must be unique within a project.

Files and directories created and modified by scaffolding:

- **proto** 

     A directory that contains proto files for query and message services.

- **x**

    Common logic for a module.
    
- **app/app.go** 

    Imports and initializes your module. 

To scaffold an IBC-enabled module use `--ibc` flag. <!-- Learn more about Starport features related to IBC. -->
