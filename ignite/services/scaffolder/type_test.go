package scaffolder

import (
	"github.com/ignite/cli/ignite/templates/field/datatype"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddTypeOptions(t *testing.T) {
	const (
		testModuleName = "test"
		testSigner     = "creator"
	)

	tests := []struct {
		name       string
		addKind    AddTypeKind
		addOptions []AddTypeOption
		expected   addTypeOptions
	}{
		{
			name:    "list type with fields",
			addKind: ListType(),
			addOptions: []AddTypeOption{
				TypeWithFields("foo", "bar"),
			},
			expected: addTypeOptions{
				moduleName: testModuleName,
				fields:     []string{"foo", "bar"},
				isList:     true,
				signer:     testSigner,
			},
		},
		{
			name:    "singleton type with module",
			addKind: SingletonType(),
			addOptions: []AddTypeOption{
				TypeWithModule("module"),
			},
			expected: addTypeOptions{
				moduleName:  "module",
				isSingleton: true,
				signer:      testSigner,
			},
		},
		{
			name:    "map type without simulation",
			addKind: MapType("foo", "bar"),
			addOptions: []AddTypeOption{
				TypeWithoutSimulation(),
			},
			expected: addTypeOptions{
				moduleName:        testModuleName,
				indexes:           []string{"foo", "bar"},
				isMap:             true,
				withoutSimulation: true,
				signer:            testSigner,
			},
		},
		{
			name:    "dry type with signer, without message",
			addKind: DryType(),
			addOptions: []AddTypeOption{
				TypeWithoutMessage(),
				TypeWithSigner("signer"),
			},
			expected: addTypeOptions{
				moduleName:     testModuleName,
				withoutMessage: true,
				signer:         "signer",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			o := newAddTypeOptions(testModuleName)
			for _, apply := range append(tc.addOptions, AddTypeOption(tc.addKind)) {
				apply(&o)
			}

			require.Equal(t, tc.expected, o)
		})
	}
}

func TestCheckForbiddenTypeIndexField(t *testing.T) {
	tests := []struct {
		name        string
		index       string
		shouldError bool
	}{
		{
			name:        "should fail with empty index",
			index:       "",
			shouldError: true,
		},
		{
			name:        "should fail with reserved Go keyword",
			index:       "uint",
			shouldError: true,
		},
		{
			name:        "should fail with forbidden ignite keyword - id",
			index:       "id",
			shouldError: true,
		},
		{
			name:        "should fail with forbidden ignite keyword - ID",
			index:       "id",
			shouldError: true,
		},
		{
			name:        "should fail with forbidden ignite keyword - params",
			index:       "params",
			shouldError: true,
		},
		{
			name:        "should fail with forbidden ignite keyword - appendedvalue",
			index:       "appendedvalue",
			shouldError: true,
		},
		{
			name:        "should fail with forbidden ignite keyword - customtype keyword",
			index:       datatype.TypeCustom,
			shouldError: true,
		},
		{
			name:  "should pass - blog",
			index: "blog",
		},
		{
			name:  "should pass - post",
			index: "post",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := checkForbiddenTypeIndex(tc.index)
			if tc.shouldError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
