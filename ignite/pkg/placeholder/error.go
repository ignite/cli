package placeholder

import (
	"fmt"
	"strings"

	"github.com/ignite/cli/ignite/pkg/validation"
)

var _ validation.Error = (*MissingPlaceholdersError)(nil)

// MissingPlaceholdersError is used as an error when a source file is missing placeholder.
type MissingPlaceholdersError struct {
	missing          iterableStringSet
	additionalInfo   string
	additionalErrors error
}

// Is true if both errors have the same list of missing placeholders.
func (e *MissingPlaceholdersError) Is(err error) bool {
	other, ok := err.(*MissingPlaceholdersError) //nolint:errorlint
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

// Error implements error interface.
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

// ValidationInfo implements validation.Error interface.
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
	if e.additionalErrors != nil {
		b.WriteString("\n\nAdditional errors: ")
		b.WriteString(e.additionalErrors.Error())
	}
	return b.String()
}

var _ validation.Error = (*ValidationMiscError)(nil)

// ValidationMiscError is used as a miscellaneous error related to validation.
type ValidationMiscError struct {
	errors []string
}

// Error implements error interface.
func (e *ValidationMiscError) Error() string {
	return fmt.Sprintf("validation errors: %v", e.errors)
}

// ValidationInfo implements validation.Error interface.
func (e *ValidationMiscError) ValidationInfo() string {
	return fmt.Sprintf("Validation errors:\n\n%v", strings.Join(e.errors, "\n"))
}
