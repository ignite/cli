package version

import (
	"fmt"
	"runtime"
	"time"
)

var (
	// Version is the semantic version of Starport.
	Version = "0.12.0-develop"

	// Date is the build date of Starport.
	Date = time.Now()
)

// Long generates a detailed version info.
func Long() string {
	return fmt.Sprintf("starport version %s %s/%s -build date: %s",
		Version,
		runtime.GOOS,
		runtime.GOARCH,
		Date.Format(time.RFC3339))
}
