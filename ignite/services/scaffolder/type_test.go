package scaffolder

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/multiformatname"
	"github.com/ignite/cli/ignite/templates/field"
	"github.com/ignite/cli/ignite/templates/field/datatype"
)

func TestParseTypeFields(t *testing.T) {
	const (
		testModuleName = "test"
		testSigner     = "creator"
	)

	tests := []struct {
		name            string
		addKind         AddTypeKind
		addOptions      []AddTypeOption
		expectedOptions addTypeOptions
		shouldError     bool
		expectedFields  field.Fields
	}{
		{
			name:    "list type with fields",
			addKind: ListType(),
			addOptions: []AddTypeOption{
				TypeWithFields("foo", "bar"),
			},
			expectedOptions: addTypeOptions{
				moduleName: testModuleName,
				fields:     []string{"foo", "bar"},
				isList:     true,
				signer:     testSigner,
			},
			shouldError: false,
			expectedFields: field.Fields{
				{
					Name: multiformatname.Name{
						Original:   "foo",
						LowerCamel: "foo",
						UpperCamel: "Foo",
						LowerCase:  "foo",
						UpperCase:  "FOO",
						Kebab:      "foo",
						Snake:      "foo",
					},
					DatatypeName: "string",
					Datatype:     "",
				},
				{
					Name: multiformatname.Name{
						Original:   "bar",
						LowerCamel: "bar",
						UpperCamel: "Bar",
						LowerCase:  "bar",
						UpperCase:  "BAR",
						Kebab:      "bar",
						Snake:      "bar",
					},
					DatatypeName: "string",
					Datatype:     "",
				},
			},
		},
		{
			name:    "singleton type with module",
			addKind: SingletonType(),
			addOptions: []AddTypeOption{
				TypeWithModule("module"),
			},
			expectedOptions: addTypeOptions{
				moduleName:  "module",
				isSingleton: true,
				signer:      testSigner,
			},
			shouldError:    false,
			expectedFields: nil,
		},
		{
			name:    "map type without simulation",
			addKind: MapType("foo", "bar"),
			addOptions: []AddTypeOption{
				TypeWithoutSimulation(),
			},
			expectedOptions: addTypeOptions{
				moduleName:        testModuleName,
				indexes:           []string{"foo", "bar"},
				isMap:             true,
				withoutSimulation: true,
				signer:            testSigner,
			},
			shouldError:    false,
			expectedFields: nil,
		},
		{
			name:    "dry type with signer, without message",
			addKind: DryType(),
			addOptions: []AddTypeOption{
				TypeWithoutMessage(),
				TypeWithSigner("signer"),
				TypeWithFields("FieldFoo"),
			},
			expectedOptions: addTypeOptions{
				moduleName:     testModuleName,
				withoutMessage: true,
				fields:         []string{"FieldFoo"},
				signer:         "signer",
			},
			shouldError: false,
			expectedFields: field.Fields{
				{
					Name: multiformatname.Name{
						Original:   "FieldFoo",
						LowerCamel: "fieldFoo",
						UpperCamel: "FieldFoo",
						LowerCase:  "fieldfoo",
						UpperCase:  "FIELDFOO",
						Kebab:      "field-foo",
						Snake:      "field_foo",
					},
					DatatypeName: "string",
					Datatype:     "",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			o := newAddTypeOptions(testModuleName)
			for _, apply := range append(tc.addOptions, AddTypeOption(tc.addKind)) {
				apply(&o)
			}

			require.Equal(t, tc.expectedOptions, o)
			fields, err := parseTypeFields(o)
			if tc.shouldError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedFields, fields)
		})
	}
}

// indirectly tests checkForbiddenTypeField()
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
		{
			name:  "should pass - typed index",
			index: "blogID:uint",
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

func TestAddType(t *testing.T) {
}
