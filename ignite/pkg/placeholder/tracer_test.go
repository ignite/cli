package placeholder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func newErrMissingPlaceholder(missing []string) *MissingPlaceholdersError {
	err := &MissingPlaceholdersError{missing: iterableStringSet{}}
	for _, placeholder := range missing {
		err.missing[placeholder] = struct{}{}
	}
	return err
}

func TestReplace(t *testing.T) {
	tests := []struct {
		desc    string
		content string
		replace []string
		missing []string
	}{
		{
			desc:    "FoundAll",
			content: "#one #two",
			replace: []string{"#one", "#two"},
		},
		{
			desc:    "MissingAll",
			content: "",
			replace: []string{"#one", "#two"},
			missing: []string{"#one", "#two"},
		},
		{
			desc:    "MissingOne",
			content: "#two",
			replace: []string{"#one", "#two"},
			missing: []string{"#one"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			tr := New()
			content := tc.content
			for _, placeholder := range tc.replace {
				content = tr.Replace(content, placeholder, "")
			}
			err := tr.Err()
			if err != nil {
				require.ErrorIs(t, err, newErrMissingPlaceholder(tc.missing))
			} else {
				require.Empty(t, tc.missing)
			}
		})
	}
}

func TestReplaceAll(t *testing.T) {
	tests := []struct {
		desc    string
		content string
		replace []string
		missing []string
	}{
		{
			desc:    "FoundAll",
			content: "#one #one #two",
			replace: []string{"#one", "#two"},
		},
		{
			desc:    "MissingAll",
			content: "",
			replace: []string{"#one", "#two"},
			missing: []string{"#one", "#two"},
		},
		{
			desc:    "MissingOne",
			content: "#two #two",
			replace: []string{"#one", "#two"},
			missing: []string{"#one"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			tr := New()
			content := tc.content
			for _, placeholder := range tc.replace {
				content = tr.ReplaceAll(content, placeholder, "")
			}
			err := tr.Err()
			if err != nil {
				require.ErrorIs(t, err, newErrMissingPlaceholder(tc.missing))
			} else {
				require.Empty(t, tc.missing)
			}
		})
	}
}
