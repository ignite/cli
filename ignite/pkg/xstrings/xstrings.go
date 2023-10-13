package xstrings

import (
	"strings"
	"unicode"

	"golang.org/x/exp/slices" // TODO: replace with slices.Contains when it will be available in stdlib (1.21)
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// AllOrSomeFilter filters elems out from the list as they  present in filterList and
// returns the remaining ones.
// if filterList is empty, all elems from list returned.
func AllOrSomeFilter(list, filterList []string) []string {
	if len(filterList) == 0 {
		return list
	}

	var elems []string

	for _, elem := range list {
		if !slices.Contains(filterList, elem) {
			elems = append(elems, elem)
		}
	}

	return elems
}

// List returns a slice of strings captured after the value returned by do which is
// called n times.
func List(n int, do func(i int) string) []string {
	var list []string

	for i := 0; i < n; i++ {
		list = append(list, do(i))
	}

	return list
}

// FormatUsername formats a username to make it usable as a variable.
func FormatUsername(s string) string {
	return NoDash(NoNumberPrefix(s))
}

// NoDash removes dash from the string.
func NoDash(s string) string {
	return strings.ReplaceAll(s, "-", "")
}

// NoNumberPrefix adds an underscore at the beginning of the string if it stars with a number
// this is used for package of proto files template because the package name can't start with a number.
func NoNumberPrefix(s string) string {
	// Check if it starts with a digit
	if unicode.IsDigit(rune(s[0])) {
		return "_" + s
	}
	return s
}

// Title returns a copy of the string s with all Unicode letters that begin words
// mapped to their Unicode title case.
func Title(s string) string {
	return cases.Title(language.English).String(s)
}

// ToUpperFirst returns a copy of the string with the first unicode letter in upper case.
func ToUpperFirst(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}
