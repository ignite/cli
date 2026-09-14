# Embedded gno stdlibs

This directory contains a verbatim copy of the Gno standard library sources
(`gnovm/stdlibs`) taken from the gno module version pinned in the root
`go.mod` (`github.com/gnolang/gno`). Ignite embeds and extracts them at
runtime so `ignite chain serve` can boot a dev chain without a local gno
checkout.

The copy is synced with:

```sh
make sync-gno-stdlibs
```

which downloads the pinned gno version from the module proxy and refreshes
the tree. A CI check (`.github/workflows/gno-stdlibs-check.yml`, backed by
`TestStdlibsInSync`) fails when this copy drifts from the pinned version —
e.g. after bumping gno in `go.mod`.

The directory is `_`-prefixed so the go tool ignores the copied `.go` files:
they are gno data, not ignite code. Do not edit by hand; always re-sync.
