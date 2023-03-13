package module

import (
	"regexp"
	"strings"
)

// ProtoPackageName creates a protocol buffer package name for an app module.
func ProtoPackageName(appModulePath, moduleName string) string {
	pathArray := strings.Split(appModulePath, "/")
	path := []string{pathArray[len(pathArray)-1], moduleName}
	return cleanProtoPackageName(strings.Join(path, "."))
}

func cleanProtoPackageName(name string) string {
	r := regexp.MustCompile("[^a-zA-Z0-9_.]+")
	return strings.ToLower(r.ReplaceAllString(name, ""))
}
