# Migration from Cosmos SDK 0.39 to 0.40

The Cosmos SDK has major changes coming with version 0.40. To make the process as smooth as possible in a decentralized environment, the version 0.39 - also known as launchpad - is backward compatible with the previous versions and created as a LTS (Long Term Support) release that projects might upgrade to.

The version 0.40 - also called the Stargate - changes important specification like decoding and encoding of transactions and adds among others the Inter-Blockchain Communication (IBC).

## Major code differences

There are a variety of things that have changed from 0.39 to 0.40. 
On version 0.39 developer and maintainer of a blockchain are used to have the gaiad and gaiacli tools available to interact with the blockchain. On version 0.40 these have been merged together into the gaiad tool that now supports all commands.

// Version 0.4x still not publically released. Currently supported 0.39, these are preparations
// TODO: list changes in app.go, genesis, ...
...