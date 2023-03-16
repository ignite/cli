package module

import (
	"regexp"
	"strings"

	"github.com/ignite/cli/ignite/pkg/xstrings"
)

// ProtoPackageName creates a protocol buffer package name for an app module.
func ProtoPackageName(appModulePath, moduleName string) string {
	pathArray := strings.Split(appModulePath, "/")
	path := []string{pathArray[len(pathArray)-1], moduleName}

	// Make sure that the first path element can be used as proto package name.
	// This is required for app module names like "github.com/username/repo" where
	// "username" might be not be compatible with proto buffer package names.
	path[0] = xstrings.NoNumberPrefix(path[0])

	return cleanProtoPackageName(strings.Join(path, "."))
}

func cleanProtoPackageName(name string) string {
	r := regexp.MustCompile("[^a-zA-Z0-9_.]+")
	return strings.ToLower(r.ReplaceAllString(name, ""))
}
