package params

import (
	"testing"
)

// XXX move elsewhere in the gno.land/ dir.

func TestValidate(t *testing.T) {
	t.Skip("this test isn't testing useful things. we need to move it to the good location and make sure to use an appropriate validate() function.")
	tests := []struct {
		module    string
		submodule string
		name      string
		type_     string
		wantErr   bool
	}{
		// Valid cases
		{"module", "p", "valid_key", "string", false},
		{"module", "p", "valid_key.string", "string", false}, // is a string
		{"module1", "p", "validKey", "int64", false},
		{"module", "p", "1invalid", "string", false}, // Starts with a number (see IsASCII)
		{"p_", "p", "_valid123", "bool", false},
		{"module", "_", "valid_key", "string", false},                         // Underscore as submodule
		{"module", "gno.land/r/myuser/myrealm", "valid_key", "string", false}, // Realm path as submodule
		{"1module", "p", "valid_key.string", "string", false},                 // Module starts with a number

		// Invalid key cases
		{"module", "p", "", "string", true},              // Empty key
		{"module", "p", "-invalid", "string", true},      // Starts with an invalid character
		{"module", "p", "invalid-123", "string", true},   // Contains invalid character (-)
		{"module", "p", "valid/path.key", "bool", true},  // Contains invalid character (/)
		{"module", "p", "invalid.string", "int64", true}, // Not a string
		{"module", "p", "valid:key", "string", true},     // ":" in name
		{"module", "p", "valid:", "string", true},        // ":" in name

		// Invalid submodule cases
		{"module", "", "valid_key", "string", true},    // Empty submodule
		{"module", "p:q", "valid_key", "string", true}, // ":" in submodule
		{"module", "p:", "valid_key", "string", true},  // ":" in submodule

		// Invalid module cases
		{"module!", "p", "valid_key", "string", true},          // Module contains invalid character (!)
		{"-prefix", "p", "valid_key", "string", true},          // Module starts with an invalid character
		{"module/submodule", "p", "valid_key", "string", true}, // Module contains invalid character (/)
		{"module:submodule", "p", "valid_key", "string", true}, // Module contains invalid character (/)
	}

	for _, tt := range tests {
		validate := func(module, submodule, name, type_ string) error {
			// FIXME: use a real module key validator.
			return nil
		}
		err := validate(tt.module, tt.submodule, tt.name, tt.type_)
		if (err != nil) != tt.wantErr {
			t.Errorf("validate(%q, %q, %q, %q) = %v, wantErr %v", tt.module, tt.submodule, tt.name, tt.type_, err, tt.wantErr)
		}
	}
}
