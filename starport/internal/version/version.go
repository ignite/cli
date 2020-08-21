package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the semantic version of Starport.
	Version = "dev"

	// Date is the build date of Starport.
	Date = ""
)

// Long generates a detailed version info.
func Long() string {
	return fmt.Sprintf("starport version %s %s/%s -build date: %s",
		Version,
		runtime.GOOS,
		runtime.GOARCH,
		Date)
}
