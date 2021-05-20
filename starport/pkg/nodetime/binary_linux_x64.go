// +build linux amd64

package nodetime

import _ "embed" // embed is required for binary embedding.

//go:embed nodetime-linux-x64.tar.gz
var binaryCompressed []byte
