package env

import "os"

const (
	debug = "IGN_DEBUG"
)

func DebugEnabled() bool {
	return os.Getenv(debug) == "1"
}
