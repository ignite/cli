# Vendored gno stdlibs

This directory contains a verbatim copy of the Gno standard library sources
(`gnovm/stdlibs`) from https://github.com/gnolang/gno (BSD-3-Clause, by the
Gnoland maintainers). Ignite embeds and extracts them at runtime so
`ignite chain serve` can boot a dev chain without a local gno checkout.

The directory is `_`-prefixed so the go tool ignores the copied `.go` files:
they are gno data, not ignite code. Do not edit; re-vendor from the gno repo.
