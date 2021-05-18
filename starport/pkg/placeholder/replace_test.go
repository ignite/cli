package placeholder

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func newErrMissingPlaceholderFromSlice(missing []string) *ErrMissingPlaceholders {
	err := &ErrMissingPlaceholders{missing: iterableStringSet{}}
	for _, placeholder := range missing {
		err.missing[placeholder] = struct{}{}
	}
	return err
}

func TestReplace(t *testing.T) {
	for _, tc := range []struct {
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
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			ctx := EnableTracing(context.Background())
			content := tc.content
			for _, placeholder := range tc.replace {
				content = Replace(ctx, content, placeholder, "")
			}
			err := Validate(ctx)
			if err != nil {
				require.ErrorIs(t, err, newErrMissingPlaceholderFromSlice(tc.missing))
			} else {
				require.Empty(t, tc.missing)
			}
		})
	}
}

func TestReplaceDisabled(t *testing.T) {
	ctx := context.Background()
	_ = Replace(ctx, "", "#one", "")
	require.NoError(t, Validate(ctx))
}
