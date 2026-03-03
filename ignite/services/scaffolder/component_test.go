package scaffolder

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
)

func TestCheckForbiddenComponentName(t *testing.T) {
	tests := []struct {
		name        string
		compName    string
		shouldError bool
	}{
		{
			name:        "should allow valid case",
			compName:    "valid",
			shouldError: false,
		},
		{
			name:        "should prevent forbidden name",
			compName:    "genesis",
			shouldError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mfName, err := multiformatname.NewName(tc.compName)
			require.NoError(t, err)

			err = checkForbiddenComponentName(mfName)
			if tc.shouldError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestCheckGoReservedWord(t *testing.T) {
	tests := []struct {
		name        string
		word        string
		shouldError bool
	}{
		{
			name:        "should allow valid case",
			word:        "valid",
			shouldError: false,
		},
		{
			name:        "should prevent forbidden go identifier",
			word:        "panic",
			shouldError: true,
		},
		{
			name:        "should prevent forbidden go keyword",
			word:        "for",
			shouldError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := checkGoReservedWord(tc.word)
			if tc.shouldError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestContainsCustomTypes(t *testing.T) {
	tests := []struct {
		name     string
		fields   []string
		contains bool
	}{
		{
			name:     "contains no custom types",
			fields:   []string{"foo", "bar"},
			contains: false,
		},
		{
			name:     "contains one non-custom type",
			fields:   []string{"foo", "bar:coin"},
			contains: false,
		},
		{
			name:     "contains one custom type",
			fields:   []string{"foo", "bar:CustomType"},
			contains: true,
		},
		{
			name:     "contains one custom array type",
			fields:   []string{"foo", "bar:array.CustomType"},
			contains: true,
		},
		{
			name:     "contains one built-in array type",
			fields:   []string{"foo", "bar:array.string"},
			contains: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.contains, containsCustomTypes(tc.fields))
		})
	}
}

func TestCustomFieldType(t *testing.T) {
	tests := []struct {
		name      string
		fieldType string
		wantType  string
		isCustom  bool
	}{
		{
			name:      "built-in scalar type",
			fieldType: "string",
			isCustom:  false,
		},
		{
			name:      "built-in array type",
			fieldType: "array.string",
			isCustom:  false,
		},
		{
			name:      "custom scalar type",
			fieldType: "ProductDetails",
			wantType:  "ProductDetails",
			isCustom:  true,
		},
		{
			name:      "custom array type",
			fieldType: "array.ProductDetails",
			wantType:  "ProductDetails",
			isCustom:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotType, isCustom := customFieldType(tc.fieldType)
			require.Equal(t, tc.isCustom, isCustom)
			require.Equal(t, tc.wantType, gotType)
		})
	}
}
