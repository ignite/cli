package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the semantic version of Starport.
	Version = ""

	// Date is the build date of Starport.
	Date = ""

	// Head is the HEAD of the current branch.
	Head = ""
)

// Long generates a detailed version info.
func Long() string {
	return fmt.Sprintf("starport version %s %s/%s -build date: %s\ngit object hash: %s",
		Version,
		runtime.GOOS,
		runtime.GOARCH,
		Date,
		Head)
}
