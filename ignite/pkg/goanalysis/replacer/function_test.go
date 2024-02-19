package replacer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModifyFunction(t *testing.T) {
	type args struct {
		fileContent  string
		functionName string
		functions    []FunctionOptions
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ModifyFunction(tt.args.fileContent, tt.args.functionName, tt.args.functions...)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
