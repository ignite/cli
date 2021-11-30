---
order: 15
description: Module Parameters
---

# Module Parameters

Sometimes the developer needs to set some default parameters for a module. The [cosmos-sdk params package](https://docs.cosmos.network/master/modules/params) provides a globally available parameter saved into the store.
The starport can scaffold parameters to be accessible for the module. These parameters have default values but can change along the chain is alive.

you can scaffold a module with params using the `--params` flag:
```shell
starport scaffold module launch --params minLaunch:uint,maxLaunch:int
```

After scaffolding, the developer should change the `x/<module>/types/params.go` file to set the default values and validate the field. The params support all built-in starport types