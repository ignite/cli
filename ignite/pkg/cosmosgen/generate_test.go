package cosmosgen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindInnerProtoFolder(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "proto-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// create dummy files
	create := func(path string) {
		dir := filepath.Dir(path)
		err := os.MkdirAll(dir, 0o755)
		require.NoError(t, err)
		_, err = os.Create(path)
		require.NoError(t, err)
	}

	tests := []struct {
		name         string
		setup        func(root string)
		expectedPath string
		expectError  bool
	}{
		{
			name: "no proto files",
			setup: func(root string) {
				// No files created
			},
			expectError: true,
		},
		{
			name: "single proto file in root",
			setup: func(root string) {
				create(filepath.Join(root, "a.proto"))
			},
			expectedPath: ".",
		},
		{
			name: "single proto file in proto dir",
			setup: func(root string) {
				create(filepath.Join(root, "proto", "a.proto"))
			},
			expectedPath: "proto",
		},
		{
			name: "multiple proto files in same proto dir",
			setup: func(root string) {
				create(filepath.Join(root, "proto", "a.proto"))
				create(filepath.Join(root, "proto", "b.proto"))
			},
			expectedPath: "proto",
		},
		{
			name: "nested proto directories",
			setup: func(root string) {
				create(filepath.Join(root, "proto", "a.proto"))
				create(filepath.Join(root, "proto", "api", "v1", "b.proto"))
			},
			expectedPath: "proto",
		},
		{
			name: "highest level proto directory",
			setup: func(root string) {
				create(filepath.Join(root, "proto", "a.proto"))
				create(filepath.Join(root, "foo", "proto", "b.proto"))
			},
			expectedPath: "proto",
		},
		{
			name: "no dir named proto",
			setup: func(root string) {
				create(filepath.Join(root, "api", "a.proto"))
			},
			expectedPath: "api",
		},
		{
			name: "deeply nested with no proto dir name",
			setup: func(root string) {
				create(filepath.Join(root, "foo", "bar", "a.proto"))
			},
			expectedPath: "foo/bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			caseRoot := filepath.Join(tmpDir, tt.name)
			err := os.MkdirAll(caseRoot, 0o755)
			require.NoError(t, err)

			tt.setup(caseRoot)

			result, err := findInnerProtoFolder(caseRoot)

			if tt.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			expected := filepath.Join(caseRoot, tt.expectedPath)
			require.Equal(t, expected, result)
		})
	}
}
