package templates

import "unicode"

// noNumberPrefix adds a underscore at the beginning of the string if it stars with a number
// this is used for package of proto files template because the package name can't start with a string
func NoNumberPrefix(s string) string {
	// Check if it starts with a digit
	if unicode.IsDigit(rune(s[0])) {
		return "_" + s
	}
	return s
}