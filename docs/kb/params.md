---
order: 15
description: Module Parameters
---

# Module Parameters

Sometimes you need to set default parameters for a module. The Cosmos SDK [params package](https://docs.cosmos.network/master/modules/params) provides a globally available parameter saved into the store. Params are managed and centralized by the Cosmos SDK `params` module and are updated with a governance proposal.
The starport can scaffold parameters to be accessible for the module. These parameters have default values but can change along the chain is alive. Since they are managed and centralized by the params SDK module, they can be easily updated through a Governance proposal.

To scaffold a module with params using the `--params` flag:
```shell
starport scaffold module launch --params minLaunch:uint,maxLaunch:int
```

After the parameters are scaffolded, change the `x/<module>/types/params.go` file to set the default values and validate the field. The params support all built-in Starport types.
