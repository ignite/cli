package placeholder

import (
	"strings"

	"github.com/tendermint/starport/starport/pkg/validation"
)

type iterableStringSet map[string]struct{}

func (set iterableStringSet) Iterate(f func(i int, element string) bool) {
	i := 0
	for key := range set {
		if !f(i, key) {
			return
		}
		i++
	}
}

func (set iterableStringSet) Add(item string) {
	set[item] = struct{}{}
}

var _ validation.Error = (*MissingPlaceholdersError)(nil)

type MissingPlaceholdersError struct {
	missing        iterableStringSet
	additionalInfo string
}

// Is true if both errors have the same list of missing placeholders.
func (e *MissingPlaceholdersError) Is(err error) bool {
	other, ok := err.(*MissingPlaceholdersError)
	if !ok {
		return false
	}
	if len(other.missing) != len(e.missing) {
		return false
	}
	for i := range e.missing {
		if e.missing[i] != other.missing[i] {
			return false
		}
	}
	return true
}

func (e *MissingPlaceholdersError) Error() string {
	var b strings.Builder
	b.WriteString("missing placeholders: ")
	e.missing.Iterate(func(i int, element string) bool {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(element)
		return true
	})
	return b.String()
}

func (e *MissingPlaceholdersError) ValidationInfo() string {
	var b strings.Builder
	b.WriteString("Missing placeholders:\n\n")
	e.missing.Iterate(func(i int, element string) bool {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(element)
		return true
	})
	if e.additionalInfo != "" {
		b.WriteString("\n\n")
		b.WriteString(e.additionalInfo)
	}
	return b.String()
}

// Option for configuring session.
type Option func(*Tracer)

// WithAdditionalInfo will append info to the validation error.
func WithAdditionalInfo(info string) Option {
	return func(s *Tracer) {
		s.additionalInfo = info
	}
}

// New instantiates Session with provided options.
func New(opts ...Option) *Tracer {
	s := &Tracer{missing: iterableStringSet{}}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type Replacer interface {
	Replace(content, placeholder, replacement string) string
}

// Tracer keeps track of missing placeholders.
type Tracer struct {
	missing        iterableStringSet
	additionalInfo string
}

// Replace placeholder in content with replacement string once.
func (t *Tracer) Replace(content, placeholder, replacement string) string {
	// NOTE(dshulyak) we will count twice. once here and second time in strings.Replace
	// if it turns out to be an issue, copy the code from strings.Replace.
	if strings.Count(content, placeholder) == 0 {
		t.missing.Add(placeholder)
		return content
	}
	return strings.Replace(content, placeholder, replacement, 1)
}

// Err if any of the placeholders were missing during execution.
func (t *Tracer) Err() error {
	if len(t.missing) > 0 {
		missing := iterableStringSet{}
		for key := range t.missing {
			missing.Add(key)
		}
		return &MissingPlaceholdersError{missing: missing, additionalInfo: t.additionalInfo}
	}
	return nil
}
