// Package prefixgen is a prefix generation helper for log messages
// and any other kind.
package prefixgen

import (
	"fmt"
	"strings"

	"github.com/ignite/cli/ignite/pkg/cliui/colors"
)

// Prefixer generates prefixes.
type Prefixer struct {
	format           string
	color            string
	left, right      string
	convertUppercase bool
}

// Option configures Prefixer.
type Option func(p *Prefixer)

// Color sets color to the prefix.
func Color(color string) Option {
	return func(p *Prefixer) {
		p.color = color
	}
}

// SquareBrackets adds square brackets to the prefix.
func SquareBrackets() Option {
	return func(p *Prefixer) {
		p.left = "["
		p.right = "]"
	}
}

// SpaceRight adds rights space to the prefix.
func SpaceRight() Option {
	return func(p *Prefixer) {
		p.right += " "
	}
}

// Uppercase formats the prefix to uppercase.
func Uppercase() Option {
	return func(p *Prefixer) {
		p.convertUppercase = true
	}
}

// Common holds some common prefix options and extends those
// options by given options.
func Common(options ...Option) []Option {
	return append([]Option{
		SquareBrackets(),
		SpaceRight(),
		Uppercase(),
	}, options...)
}

// New creates a new Prefixer with format and options.
// Format is an fmt.Sprintf() like format to dynamically create prefix texts
// as needed.
func New(format string, options ...Option) *Prefixer {
	p := &Prefixer{
		format: format,
	}
	for _, o := range options {
		o(p)
	}
	return p
}

// Gen generates a new prefix by applying s to format given during New().
func (p *Prefixer) Gen(s ...interface{}) string {
	format := p.format
	format = p.left + format
	format += p.right
	prefix := fmt.Sprintf(format, s...)
	if p.convertUppercase {
		prefix = strings.ToUpper(prefix)
	}
	if p.color != "" {
		return colors.SprintFunc(p.color)(prefix)
	}
	return prefix
}
