package xast

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
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
			name: "add all modifications type",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions: []FunctionOptions{
					AppendParams("param1", "string", 0),
					ReplaceBody(`return false`),
					AppendAtLine(`fmt.Println("Appended at line 0.")`, 0),
					AppendAtLine(`SimpleCall(foo, bar)`, 1),
					AppendCode(`fmt.Println("Appended code.")`),
					NewReturn("1"),
					InsideCall("SimpleCall", "baz", 0),
					InsideCall("SimpleCall", "bla", -1),
					InsideCall("Println", strconv.Quote("test"), -1),
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
		{
			name: "add the replace body",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{ReplaceBody(`return false`)},
			},
			want: `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
	New(param1, param2)
}

func anotherFunction() bool { return false }
`,
		},
		{
			name: "add append line and code modification",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions: []FunctionOptions{
					AppendAtLine(`fmt.Println("Appended at line 0.")`, 0),
					AppendAtLine(`SimpleCall(foo, bar)`, 1),
					AppendCode(`fmt.Println("Appended code.")`),
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

func anotherFunction() bool {
	fmt.Println("Appended at line 0.")
	SimpleCall(foo, bar)

	p := bla.NewParam()
	p.CallSomething("Another call")
	fmt.Println("Appended code.")

	return true
}
`,
		},
		{
			name: "add all modifications type",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{NewReturn("1")},
			},
			want: strings.ReplaceAll(existingContent, "return true", "return 1\n") + "\n",
		},
		{
			name: "add inside call modifications",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions: []FunctionOptions{
					InsideCall("NewParam", "baz", 0),
					InsideCall("NewParam", "bla", -1),
					InsideCall("CallSomething", strconv.Quote("test1"), -1),
					InsideCall("CallSomething", strconv.Quote("test2"), 0),
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

func anotherFunction() bool {
	p := bla.NewParam(baz, bla)
	p.CallSomething("test2", "Another call", "test1")
	return true
}
`,
		},
		{
			name: "params out of range",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{AppendParams("param1", "string", 1)},
			},
			err: errors.New("params index 1 out of range"),
		},
		{
			name: "invalid params",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{AppendParams("9#.(c", "string", 0)},
			},
			err: errors.New("format.Node internal error (12:22: expected ')', found 9 (and 1 more errors))"),
		},
		{
			name: "invalid content for replace body",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{ReplaceBody("9#.(c")},
			},
			err: errors.New("1:24: illegal character U+0023 '#'"),
		},
		{
			name: "line number out of range",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{AppendAtLine(`fmt.Println("")`, 4)},
			},
			err: errors.New("line number 4 out of range"),
		},
		{
			name: "invalid code for append at line",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{AppendAtLine("9#.(c", 0)},
			},
			err: errors.New("1:2: illegal character U+0023 '#'"),
		},
		{
			name: "invalid code append",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{AppendCode("9#.(c")},
			},
			err: errors.New("1:2: illegal character U+0023 '#'"),
		},
		{
			name: "invalid new return",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{NewReturn("9#.(c")},
			},
			err: errors.New("1:2: illegal character U+0023 '#'"),
		},
		{
			name: "call name not found",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{InsideCall("FooFunction", "baz", 0)},
			},
			err: errors.New("function calls not found: map[FooFunction:[{FooFunction baz 0}]]"),
		},
		{
			name: "invalid call param",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{InsideCall("NewParam", "9#.(c", 0)},
			},
			err: errors.New("format.Node internal error (13:21: illegal character U+0023 '#' (and 2 more errors))"),
		},
		{
			name: "call params out of range",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{InsideCall("NewParam", "baz", 1)},
			},
			err: errors.New("function call index 1 out of range"),
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
