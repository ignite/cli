package placeholder

import (
	"strings"
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
	ReplaceAll(content, placeholder, replacement string) string
	ReplaceOnce(content, placeholder, replacement string) string
	AppendMiscError(miscError string)
}

// Tracer keeps track of missing placeholders or other issues related to file modification.
type Tracer struct {
	missing        iterableStringSet
	miscErrors     []string
	additionalInfo string
}

// ReplaceAll replace all placeholders in content with replacement string.
func (t *Tracer) ReplaceAll(content, placeholder, replacement string) string {
	if strings.Count(content, placeholder) == 0 {
		t.missing.Add(placeholder)
		return content
	}
	return strings.ReplaceAll(content, placeholder, replacement)
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

// ReplaceOnce will replace placeholder in content only if replacement is not already found in content.
func (t *Tracer) ReplaceOnce(content, placeholder, replacement string) string {
	if !strings.Contains(content, replacement) {
		return t.Replace(content, placeholder, replacement)
	}
	return content
}

// AppendMiscError allows to track errors not related to missing placeholders during file modification.
func (t *Tracer) AppendMiscError(miscError string) {
	t.miscErrors = append(t.miscErrors, miscError)
}

// Err if any of the placeholders were missing during execution.
func (t *Tracer) Err() error {
	// miscellaneous errors represent errors preventing source modification not related to missing placeholder
	var miscErrors error
	if len(t.miscErrors) > 0 {
		miscErrors = &ValidationMiscError{
			errors: t.miscErrors,
		}
	}

	if len(t.missing) > 0 {
		missing := iterableStringSet{}
		for key := range t.missing {
			missing.Add(key)
		}
		return &MissingPlaceholdersError{
			missing:          missing,
			additionalInfo:   t.additionalInfo,
			additionalErrors: miscErrors,
		}
	}

	// if not missing placeholder but still miscellaneous errors, return them
	return miscErrors
}
