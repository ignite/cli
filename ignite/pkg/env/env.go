package env

import "os"

const (
	debug = "IGNT_DEBUG"
)

func DebugEnabled() bool {
	return os.Getenv(debug) == "1"
}
