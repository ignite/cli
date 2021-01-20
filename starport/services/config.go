package services

import "os"

var (
	StarportConfDir = os.ExpandEnv("$HOME/.starport")
)
