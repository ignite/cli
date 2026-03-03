---
sidebar_position: 3
title: Chain Registry Types (chainregistry)
slug: /packages/chainregistry
---

# Chain Registry Types (chainregistry)

The `chainregistry` package provides Go structs for Cosmos chain registry files such as
`chain.json` and `assetlist.json`, plus helpers to save them as JSON.

For full API details, see the
[`chainregistry` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/chainregistry).

## Example: Generate chain.json and assetlist.json

```go
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ignite/cli/v29/ignite/pkg/chainregistry"
)

func main() {
	outDir := "./chain-registry"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatal(err)
	}

	chain := chainregistry.Chain{
		ChainName:    "ignite",
		PrettyName:   "Ignite",
		ChainID:      "ignite-1",
		Bech32Prefix: "cosmos",
		DaemonName:   "ignited",
		NodeHome:     ".ignite",
		Status:       chainregistry.ChainStatusActive,
		NetworkType:  chainregistry.NetworkTypeMainnet,
		ChainType:    chainregistry.ChainTypeCosmos,
	}

	if err := chain.SaveJSON(filepath.Join(outDir, "chain.json")); err != nil {
		log.Fatal(err)
	}

	assetList := chainregistry.AssetList{
		ChainName: "ignite",
		Assets: []chainregistry.Asset{
			{
				Name:      "Ignite Token",
				Symbol:    "IGNT",
				Base:      "uignite",
				Display:   "ignite",
				TypeAsset: "sdk.coin",
			},
		},
	}

	if err := assetList.SaveJSON(filepath.Join(outDir, "assetlist.json")); err != nil {
		log.Fatal(err)
	}
}
```
