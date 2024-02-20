package replacer

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModifyFunction(t *testing.T) {
	existingContent := `package main

import (
    "fmt"
)

func main() {
    fmt.Println("Hello, world!")
	New(param1, param2)
}

func anotherFunction() bool {
	p := bla.NewParam()
    p.CallSomething("Another call")
    return true
}`

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
		{
			name: "",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions: []FunctionOptions{
					AppendParams("param1", "string", 0),
					ReplaceBody(`return false`),
					AppendAtLine(`fmt.Println("Appended at line 0.")`, 0),
					AppendCode(`fmt.Println("Appended code.")`),
					AppendAtLine(`SimpleCall(foo, bar)`, 1),
					InsideCall("SimpleCall", "baz", 0),
					InsideCall("SimpleCall", "bla", -1),
					InsideCall("Println", strconv.Quote("test"), -1),
					NewReturn("1"),
				},
			},
			want: `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
	New(param1, param2)
}

func anotherFunction(param1 string) bool {
	fmt.Println("Appended at line 0.", "test")
	SimpleCall(baz, foo, bar, bla)
	fmt.Println("Appended code.", "test")
	return 1
}
`,
		},
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
