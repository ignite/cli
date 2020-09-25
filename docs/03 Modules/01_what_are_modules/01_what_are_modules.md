# What are modules

In the Cosmos SDK modules are the basis for the logic of your blockchain. Each module serves specific purposes and functions. The Cosmos SDK offers a variety of native modules to make a blockchain work. These modules handle authentification for users, token transfers, governmental functions, staking of tokens, supply of tokens and many more.

If you want to change the default functionality of a module or just change certain hardcoded parameter, you can fork a module and change it, therefore owning your own logic for your blockchain. While forking and editing a module should be done carefully, this approach marks the Cosmos SDK as especially powerful, as you can experiment with different parameters as the standard implementation suggests.

Modules do not need to be created by a specific company or invididual. They can be created by anyone and offered for general use to the public. Although there do exist standards that projects look into before integrating a module to their blockchain. It is recommended that a module has understandable specification, handles one thing good and is well tested - optimally battle-tested on a live blockchain.
When growing more complex, sometimes it makes more sense to have two modules instead of one module trying to "solve-it-all", this consideration can make it more attractive for other projects to use a module on their blockchain project.

## Default modules

When creating a new project with starport, the following modules are activated by default

- auth
- bank
- staking
- params
- supply

In the next Chapter (Basic modules) we will look at each of those modules in more detail.

## Using modules

With starport you can add a module with the command `starport add modulename`. When adding a module manually to a blockchain application, it requires to edit the `app/app.go` and the `myappcli/main.go` with the according entries. Starport manages the code edits and additions for you conviniently.

## Summary

- Importing modules in a Cosmos SDK built blockchain exposes new functionality for the blockchain.
- Any combination of modules is allowed.
- The modules define what can be done on the blockchain.
- Modules are editable, but the success of your blockchain will be dependend on choosing the correct modules for your blockchain, for functionality and security sake.
- `starport add modulename` lets you import modules into your blockchain application.

[◀️ Previous - Configuration](../../02%20using_starport/03_configuration/03_configuration.md) | [▶️ Next - Basic Modules](../../03%20modules/02_basic_modules/02_basic_modules.md)  