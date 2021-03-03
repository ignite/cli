package gocmd

import "os"

// Name returns the name of Go binary to use.
func Name() string {
	custom := os.Getenv("GONAME")
	if custom != "" {
		return custom
	}
	return "go"
}
