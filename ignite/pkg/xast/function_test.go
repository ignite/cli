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
					AppendFuncParams("param1", "string", 0),
					ReplaceFuncBody(`return false`),
					AppendFuncAtLine(`fmt.Println("Appended at line 0.")`, 0),
					AppendFuncAtLine(`SimpleCall(foo, bar)`, 1),
					AppendFuncCode(`fmt.Println("Appended code.")`),
					AppendFuncCode(`Param{Baz: baz, Foo: foo}`),
					NewFuncReturn("1"),
					AppendInsideFuncCall("SimpleCall", "baz", 0),
					AppendInsideFuncCall("SimpleCall", "bla", -1),
					AppendInsideFuncCall("Println", strconv.Quote("test"), -1),
					AppendInsideFuncStruct("Param", "Bar", strconv.Quote("bar"), -1),
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
	Param{Baz: baz, Foo: foo, Bar: "bar"}
	return 1
}
`,
		},
		{
			name: "add the replace body",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{ReplaceFuncBody(`return false`)},
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
					AppendFuncAtLine(`fmt.Println("Appended at line 0.")`, 0),
					AppendFuncAtLine(`SimpleCall(foo, bar)`, 1),
					AppendFuncCode(`fmt.Println("Appended code.")`),
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
				functions:    []FunctionOptions{NewFuncReturn("1")},
			},
			want: strings.ReplaceAll(existingContent, "return true", "return 1\n") + "\n",
		},
		{
			name: "add inside call modifications",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions: []FunctionOptions{
					AppendInsideFuncCall("NewParam", "baz", 0),
					AppendInsideFuncCall("NewParam", "bla", -1),
					AppendInsideFuncCall("CallSomething", strconv.Quote("test1"), -1),
					AppendInsideFuncCall("CallSomething", strconv.Quote("test2"), 0),
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
			name: "add inside struct modifications",
			args: args{
				fileContent: `package main

import (
	"fmt"
)

func anotherFunction() bool {
	Param{Baz: baz, Foo: foo}
	Client{baz, foo}
	return true
}`,
				functionName: "anotherFunction",
				functions: []FunctionOptions{
					AppendInsideFuncStruct("Param", "Bar", "bar", -1),
					AppendInsideFuncStruct("Param", "Bla", "bla", 1),
					AppendInsideFuncStruct("Client", "", "bar", 0),
				},
			},
			want: `package main

import (
	"fmt"
)

func anotherFunction() bool {
	Param{Baz: baz, Bla: bla, Foo: foo, Bar: bar}
	Client{bar, baz, foo}
	return true
}
`,
		},
		{
			name: "params out of range",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{AppendFuncParams("param1", "string", 1)},
			},
			err: errors.New("params index 1 out of range"),
		},
		{
			name: "invalid params",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{AppendFuncParams("9#.(c", "string", 0)},
			},
			err: errors.New("format.Node internal error (12:22: expected ')', found 9 (and 1 more errors))"),
		},
		{
			name: "invalid content for replace body",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{ReplaceFuncBody("9#.(c")},
			},
			err: errors.New("1:24: illegal character U+0023 '#'"),
		},
		{
			name: "line number out of range",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{AppendFuncAtLine(`fmt.Println("")`, 4)},
			},
			err: errors.New("line number 4 out of range"),
		},
		{
			name: "invalid code for append at line",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{AppendFuncAtLine("9#.(c", 0)},
			},
			err: errors.New("1:2: illegal character U+0023 '#'"),
		},
		{
			name: "invalid code append",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{AppendFuncCode("9#.(c")},
			},
			err: errors.New("1:24: illegal character U+0023 '#'"),
		},
		{
			name: "invalid new return",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{NewFuncReturn("9#.(c")},
			},
			err: errors.New("1:2: illegal character U+0023 '#'"),
		},
		{
			name: "call name not found",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{AppendInsideFuncCall("FooFunction", "baz", 0)},
			},
			err: errors.New("function calls not found: map[FooFunction:[{FooFunction baz 0}]]"),
		},
		{
			name: "invalid call param",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{AppendInsideFuncCall("NewParam", "9#.(c", 0)},
			},
			err: errors.New("format.Node internal error (13:21: illegal character U+0023 '#' (and 2 more errors))"),
		},
		{
			name: "call params out of range",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{AppendInsideFuncCall("NewParam", "baz", 1)},
			},
			err: errors.New("function call index 1 out of range"),
		},
		{
			name: "empty modifications",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{},
			},
			want: existingContent + "\n",
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

func TestModifyCaller(t *testing.T) {
	existingContent := `package main

import (
	"context"
	"fmt"
)

func main() {
	// Simple function call
	fmt.Println("Hello, world!")
	
	// Call with multiple arguments
	server.Foo(param1, param2, 42)
	
	// Call with no arguments
	EmptyFunc()
	
	// Call with complex arguments
	ComplexFunc([]string{"a", "b"}, map[string]int{"a": 1})
	
	// Multiple calls to the same function
	fmt.Println("First call")
	fmt.Println("Second call")
}
`

	tests := []struct {
		name          string
		content       string
		callerExpr    string
		modifierFunc  func([]string) ([]string, error)
		expected      string
		expectedError string
	}{
		{
			name:       "replace arguments in fmt.Println",
			content:    existingContent,
			callerExpr: "fmt.Println",
			modifierFunc: func(args []string) ([]string, error) {
				return []string{`"Modified output"`}, nil
			},
			expected: `package main

import (
	"context"
	"fmt"
)

func main() {
	// Simple function call
	fmt.Println("Modified output")
	
	// Call with multiple arguments
	server.Foo(param1, param2, 42)
	
	// Call with no arguments
	EmptyFunc()
	
	// Call with complex arguments
	ComplexFunc([]string{"a", "b"}, map[string]int{"a": 1})
	
	// Multiple calls to the same function
	fmt.Println("Modified output")
	fmt.Println("Modified output")
}
`,
		},
		{
			name:       "replace server.Foo arguments",
			content:    existingContent,
			callerExpr: "server.Foo",
			modifierFunc: func(args []string) ([]string, error) {
				return []string{"context.Background()", "newParam", "123"}, nil
			},
			expected: `package main

import (
	"context"
	"fmt"
)

func main() {
	// Simple function call
	fmt.Println("Hello, world!")
	
	// Call with multiple arguments
	server.Foo(context.Background(), newParam, 123)
	
	// Call with no arguments
	EmptyFunc()
	
	// Call with complex arguments
	ComplexFunc([]string{"a", "b"}, map[string]int{"a": 1})
	
	// Multiple calls to the same function
	fmt.Println("First call")
	fmt.Println("Second call")
}
`,
		},
		{
			name:       "add argument to EmptyFunc",
			content:    existingContent,
			callerExpr: "EmptyFunc",
			modifierFunc: func(args []string) ([]string, error) {
				return []string{`"new argument"`}, nil
			},
			expected: `package main

import (
	"context"
	"fmt"
)

func main() {
	// Simple function call
	fmt.Println("Hello, world!")
	
	// Call with multiple arguments
	server.Foo(param1, param2, 42)
	
	// Call with no arguments
	EmptyFunc("new argument")
	
	// Call with complex arguments
	ComplexFunc([]string{"a", "b"}, map[string]int{"a": 1})
	
	// Multiple calls to the same function
	fmt.Println("First call")
	fmt.Println("Second call")
}
`,
		},
		{
			name:       "modify complex arguments",
			content:    existingContent,
			callerExpr: "ComplexFunc",
			modifierFunc: func(args []string) ([]string, error) {
				return []string{`[]string{"x", "y", "z"}`, `map[string]int{"x": 10}`}, nil
			},
			expected: `package main

import (
	"context"
	"fmt"
)

func main() {
	// Simple function call
	fmt.Println("Hello, world!")
	
	// Call with multiple arguments
	server.Foo(param1, param2, 42)
	
	// Call with no arguments
	EmptyFunc()
	
	// Call with complex arguments
	ComplexFunc([]string{"x", "y", "z"}, map[string]int{"x": 10})
	
	// Multiple calls to the same function
	fmt.Println("First call")
	fmt.Println("Second call")
}
`,
		},
		{
			name:       "function not found",
			content:    existingContent,
			callerExpr: "NonExistentFunc",
			modifierFunc: func(args []string) ([]string, error) {
				return []string{`"test"`}, nil
			},
			expectedError: "function call NonExistentFunc not found in file content",
		},
		{
			name:       "error in modifier function",
			content:    existingContent,
			callerExpr: "fmt.Println",
			modifierFunc: func(args []string) ([]string, error) {
				return nil, errors.New("custom error in modifier")
			},
			expectedError: "custom error in modifier",
		},
		{
			name:       "invalid caller expression",
			content:    existingContent,
			callerExpr: "pkg.sub.Function",
			modifierFunc: func(args []string) ([]string, error) {
				return []string{`"test"`}, nil
			},
			expectedError: "invalid caller expression format, use 'pkgname.FuncName' or 'FuncName'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ModifyCaller(tt.content, tt.callerExpr, tt.modifierFunc)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}
