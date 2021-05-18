package placeholder

import (
	"context"
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

var _ validation.Error = (*ErrMissingPlaceholders)(nil)

type ErrMissingPlaceholders struct {
	missing iterableStringSet
}

// Is true if both errors have the same list of missing placeholders.
func (e *ErrMissingPlaceholders) Is(err error) bool {
	other, ok := err.(*ErrMissingPlaceholders)
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

func (e *ErrMissingPlaceholders) Error() string {
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

func (e *ErrMissingPlaceholders) ValidationInfo() string {
	var b strings.Builder
	b.WriteString("Missing placeholders:\n\n")
	e.missing.Iterate(func(i int, element string) bool {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(element)
		return true
	})
	b.WriteString("\n\n")
	b.WriteString("Visit https://docs.starport.network/troubleshooting/placeholders.html.")
	return b.String()
}

type contextKey struct {
	string
}

func (c contextKey) String() string {
	return c.string
}

type tracer struct {
	placeholders iterableStringSet
}

func (t *tracer) add(placeholder string) {
	t.placeholders[placeholder] = struct{}{}
}

var missingPlaceholdersKey = contextKey{"missing placeholders"}

// EnableTracing will enable placeholder tracing when Replace is used.
func EnableTracing(ctx context.Context) context.Context {
	return context.WithValue(ctx, missingPlaceholdersKey, &tracer{placeholders: iterableStringSet{}})
}

func traceIfEnabled(ctx context.Context, placeholder string) {
	tc, ok := ctx.Value(missingPlaceholdersKey).(*tracer)
	if !ok {
		return
	}
	tc.add(placeholder)
}

// Replace placeholder in content with replacement string once.
func Replace(ctx context.Context, content, placeholder, replacement string) string {
	// NOTE(dshulyak) we will count twice. once here and second time in strings.Replace
	// if it turns out to be an issue, copy the code from strings.Replace.
	if strings.Count(content, placeholder) == 0 {
		traceIfEnabled(ctx, placeholder)
		return content
	}
	return strings.Replace(content, placeholder, replacement, 1)
}

// Validate if any of the placeholders were missing during execution.
func Validate(ctx context.Context) error {
	tc, ok := ctx.Value(missingPlaceholdersKey).(*tracer)
	if !ok {
		return nil
	}
	if len(tc.placeholders) > 0 {
		return &ErrMissingPlaceholders{missing: tc.placeholders}
	}
	return nil
}
