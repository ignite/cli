package scaffolder

import (
	"context"
	"os"
	"path/filepath"
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.contains, containsCustomTypes(tc.fields))
		})
	}
}

func TestCheckTypeProtoCreated(t *testing.T) {
	t.Run("should fail when proto type already exists", func(t *testing.T) {
		tmp := t.TempDir()
		protoFile := filepath.Join(tmp, "proto", "blog", "blog", "v1", "post.proto")
		require.NoError(t, os.MkdirAll(filepath.Dir(protoFile), 0o755))

		content := `syntax = "proto3";
package blog.blog.v1;

message Post {}
`
		require.NoError(t, os.WriteFile(protoFile, []byte(content), 0o644))

		name, err := multiformatname.NewName("post")
		require.NoError(t, err)

		err = checkTypeProtoCreated(context.Background(), tmp, "blog", "proto", "blog", name)
		require.EqualError(t, err, "component type with name post is already created (type Post exists)")
	})

	t.Run("should pass when proto type does not exist", func(t *testing.T) {
		tmp := t.TempDir()
		protoFile := filepath.Join(tmp, "proto", "blog", "blog", "v1", "comment.proto")
		require.NoError(t, os.MkdirAll(filepath.Dir(protoFile), 0o755))

		content := `syntax = "proto3";
package blog.blog.v1;

message Comment {}
`
		require.NoError(t, os.WriteFile(protoFile, []byte(content), 0o644))

		name, err := multiformatname.NewName("post")
		require.NoError(t, err)

		require.NoError(t, checkTypeProtoCreated(context.Background(), tmp, "blog", "proto", "blog", name))
	})
}
