package module

import (
	"strings"

	"github.com/ignite-hq/cli/ignite/pkg/xstrings"
)

// ProtoPackageName creates a protocol buffer package name for an app module.
func ProtoPackageName(appModulePath, moduleName string) string {
	path := strings.Split(appModulePath, "/")
	path = append(path, moduleName)

	// Make sure that the first path element can be used as proto package name.
	// This is required for app module names like "github.com/username/repo" where
	// "username" might be not be compatible with proto buffer package names.
	path[0] = xstrings.FormatUsername(path[0])

	return strings.Join(path, ".")
}
