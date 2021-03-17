// +build linux arm64

package nodetime

import _ "embed" // embed is required for binary embedding.

//go:embed nodetime-linux-arm64.tar.gz
var binaryCompressed []byte
