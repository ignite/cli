package events

import "fmt"

// Content represents objects suitable for being events content
type Content interface {
	fmt.Stringer
}

// stringContent Content wrapper for string primitive
type stringContent struct {
	str string
}

// StringContent creates Content from a string
func StringContent(content string) Content {
	return stringContent{str: content}
}

// String returns underlying string of stringContent
func (c stringContent) String() string {
	return c.str
}
