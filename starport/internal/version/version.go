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
	output := fmt.Sprintf("starport version %s %s/%s -build date: %s",
		Version,
		runtime.GOOS,
		runtime.GOARCH,
		Date)

	if Head != "" {
		output = fmt.Sprintf("%s\ngit object hash: %s", output, Head)
	}
	return output
}
