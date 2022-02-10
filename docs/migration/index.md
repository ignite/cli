---
order: 1
title: v0.19.2
parent:
  title: Migration
  order: 3
description: For chains that were scaffolded with Starport versions lower than v0.19.2, changes are required to use Starport v0.19.2. 
---

# Upgrading a blockchain to use Starport v0.19.2

Starport v0.19.2 comes with IBC v2.0.2.

With Starport v0.19.2, the contents of the deprecated Starport Modules `tendermint/spm` repo are moved to the official Starport repo which introduces breaking changes.

To migrate your chain that was scaffolded with Starport versions lower than v0.19.2: 

1. IBC upgrade: Use the [IBC migration documents](https://github.com/cosmos/ibc-go/blob/main/docs/migrations/v1-to-v2.md)
   
2. In your chain's `go.mod` file, remove `tendermint/spm` and add the v0.19.2 version of `tendermint/starport`. If your chain uses these packages, change the import paths as shown: 

    - `github.com/tendermint/spm/ibckeeper` moved to `github.com/tendermint/starport/starport/pkg/cosmosibckeeper`
    - `github.com/tendermint/spm/cosmoscmd` moved to `github.com/tendermint/starport/starport/pkg/cosmoscmd` 
    - `github.com/tendermint/spm/openapiconsole` moved to `github.com/tendermint/starport/starport/pkg/openapiconsole`
    - `github.com/tendermint/spm/testutil/sample` moved to `github.com/tendermint/starport/starport/pkg/cosmostestutil/sample`