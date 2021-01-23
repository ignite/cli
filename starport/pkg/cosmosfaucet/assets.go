// +build dev

package cosmosfaucet

import "net/http"

// Assets to serve openapi web page.
var Assets http.FileSystem = http.Dir("assets")
