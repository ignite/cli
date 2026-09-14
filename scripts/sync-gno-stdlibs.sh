#!/usr/bin/env sh
# Sync the embedded gno stdlibs from the gno module version pinned in go.mod.
# The embedded copy must match the module's gnovm/stdlibs exactly: the dev
# chain loads stdlibs from it, and a CI check fails on drift.
set -eu

MOD=github.com/gnolang/gno

# Resolve the module source dir at the pinned version (downloads it if needed).
GNO_DIR=$(go mod download -json "$MOD" | sed -n 's/^\t"Dir": "\(.*\)",$/\1/p')
[ -n "$GNO_DIR" ] || { echo "error: unable to resolve $MOD source dir" >&2; exit 1; }

SRC="$GNO_DIR/gnovm/stdlibs"
DEST="ignite/services/gno/_stdlibs/gnovm/stdlibs"

[ -d "$SRC" ] || { echo "error: $SRC not found in module" >&2; exit 1; }

VERSION=$(go list -m -f '{{.Version}}' "$MOD")
chmod -R u+w "$DEST" 2>/dev/null || true
rm -rf "$DEST"
mkdir -p "$DEST"
cp -R "$SRC/." "$DEST/"
chmod -R u+w "$DEST"

echo "synced stdlibs from $MOD $VERSION ($GNO_DIR)"
