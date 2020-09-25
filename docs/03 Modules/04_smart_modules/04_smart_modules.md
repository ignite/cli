# Smart Modules

Modules in the Cosmos SDK expand the blockchain application by a set of features. These features should be well defined and must be deterministic in outcome. Just like smart contracts, smart modules enhances a blockchain application with defined functions that the users can refer to and use.
The difference is that the Virtual Machine for smart contracts allows every user to upload a contract with their defined behavior to the blockchain. These contracts can eventually be bugged (The famous DAO hack), malicious functions ([as for exmaple described by consensys](https://consensys.github.io/smart-contract-best-practices/known_attacks/)) or not tested enough (like ERC20, which got a predecessor ERC777). Smart modules allow a certain functionality to the Cosmos SDK blockchain that can be well tested or upgraded and allows for all users to use a certain set of functions on your blockchain application.

With using Starport, you are already using a few smart modules like `IBC`, `slashing` or `upgrade`. This can be enhanced by installing other modules such as before described the Ethereum Virtual Machine `evm`, CosmWasm `wasm` or other modules.

## Summary

- Smart modules can add a set of functions to your blockchain.

[◀️ Previous - Advanced Modules](../../03%20modules/03_advanced_modules/03_advanced_modules.md) | [▶️ Next - Your own Modules](../../03%20modules/05_your_own_module/05_your_own_module.md)  