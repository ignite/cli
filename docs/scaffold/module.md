---
order: 4
---

# Module Scaffold

Modules are building blocks of Cosmos SDK blockchains. They encapsulate logic and allow sharing functionality between projects. Learn more about [building modules](https://github.com/cosmos/cosmos-sdk/tree/master/docs/building-modules).

Starport supports scaffolding SDK modules.

```
starport module create [name] [flags]
```

`name`

  The name of a new module. This name must be unique within a project.

The following files and directories are created and modified by scaffolding:

* `proto/`: a directory that contains proto files for query and message services.
* `x`: common logic for a module.
* `app/app.go`: imports and initializes your module. 

To scaffold an IBC-enabled module use `--ibc` flag. <!-- Learn more about Starport features related to IBC. -->
