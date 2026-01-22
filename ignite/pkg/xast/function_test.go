package xast

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

func TestModifyFunction(t *testing.T) {
	existingContent := `package main

import (
	"fmt"
)

// main function
func main() {
	// print hello world
	fmt.Println("Hello, world!")
	// call new param function
	New(param1, param2)
}

// anotherFunction another function
func anotherFunction() bool {
	// init param
	p := bla.NewParam()
	// start to call something
	p.CallSomething("Another call")
	// return always true
	return true
}

// TestValidate test the validations
func TestValidate(t *testing.T) {
	tests := []struct {
		desc     string
		genState types.GenesisState
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
		},
		{
			desc:     "valid genesis state",
			genState: types.GenesisState{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			require.NoError(t, err)
		})
	}
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
			name: "add a case to switch statement",
			args: args{
				fileContent: `package test

func processPacket(packet interface{}) error {
    switch packet := packet.(type) {
    default:
        return fmt.Errorf("unknown packet type: %T", packet)
    }
}`,
				functionName: "processPacket",
				functions: []FunctionOptions{
					AppendSwitchCase(
						"packet := packet.(type)",
						"*types.FooPacket",
						"return handleFooPacket(packet)",
					),
				},
			},
			want: `package test

func processPacket(packet interface{}) error {
	switch packet := packet.(type) {
	case *types.FooPacket:
		return handleFooPacket(packet)

	default:
		return fmt.Errorf("unknown packet type: %T", packet)
	}
}`,
		},
		{
			name: "add multiple cases to switch statement",
			args: args{
				fileContent: `package test

func handlePacket(data interface{}) error {
    switch v := data.(type) {
    case string:
        return processString(v)
    default:
        return fmt.Errorf("unsupported type: %T", v)
    }
}`,
				functionName: "handlePacket",
				functions: []FunctionOptions{
					AppendSwitchCase(
						"v := data.(type)",
						"int",
						"return processInt(v)",
					),
					AppendSwitchCase(
						"v := data.(type)",
						"bool",
						"return processBool(v)",
					),
				},
			},
			want: `package test

func handlePacket(data interface{}) error {
	switch v := data.(type) {
	case string:
		return processString(v)
	case int:
		return processInt(v)
	case bool:
		return processBool(v)

	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}`,
		},
		{
			name: "add multiple cases to two switch statement",
			args: args{
				fileContent: `package test

func handlePacket(data interface{}) error {
    switch v := data.(type) {
    case string:
        return processString(v)
    default:
        return fmt.Errorf("unsupported type: %T", v)
    }

    switch x {
    case 1:
        return "one"
    default:
        return "unknown"
    }
}`,
				functionName: "handlePacket",
				functions: []FunctionOptions{
					AppendSwitchCase(
						"v := data.(type)",
						"int",
						"return processInt(v)",
					),
					AppendSwitchCase(
						"x",
						"2",
						`return "two"`,
					),
				},
			},
			want: `package test

func handlePacket(data interface{}) error {
	switch v := data.(type) {
	case string:
		return processString(v)
	case int:
		return processInt(v)

	default:
		return fmt.Errorf("unsupported type: %T", v)
	}

	switch x {
	case 1:
		return "one"
	case 2:
		return "two"

	default:
		return "unknown"
	}
}`,
		},
		{
			name: "add case to switch with non-matching condition",
			args: args{
				fileContent: `package test

func process(x int) string {
    switch x {
    case 1:
        return "one"
    default:
        return "unknown"
    }
}`,
				functionName: "process",
				functions: []FunctionOptions{
					AppendSwitchCase(
						"wrongCondition",
						"2",
						`return "two"`,
					),
				},
			},
			err: errors.New("function switch not found: map[wrongCondition:[{wrongCondition 2 return \"two\"}]]"),
		},

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
					AppendFuncAtLine(`if param1 == "" {
						return false
					}`, 2),
					AppendFuncCode(`fmt.Println("Appended code.")`),
					AppendFuncCode(`Param{
						Baz: baz,
						Foo: foo,
					}`),
					NewFuncReturn("1"),
					AppendInsideFuncCall("SimpleCall", "baz", 0),
					AppendInsideFuncCall("SimpleCall", "bla", -1),
					AppendInsideFuncCall("Println", strconv.Quote("test"), -1),
					AppendFuncStruct("Param", "Bar", strconv.Quote("bar")),
					AppendFuncTestCase(`{
								desc:     "valid first genesis state",
								genState: GenesisState{},
					}`),
				},
			},
			want: `package main

import (
	"fmt"
)

// main function
func main() {
	// print hello world
	fmt.Println("Hello, world!")
	// call new param function
	New(param1, param2)
}

// anotherFunction another function
func anotherFunction(param1 string) bool {
	fmt.Println("Appended at line 0.", "test")
	SimpleCall(baz, foo, bar, bla)
	if param1 == "" {
		return false
	}
	fmt.Println("Appended code.", "test")
	Param{
		Baz: baz,
		Foo: foo,
		Bar: "bar",
	}
	return 1
}

// TestValidate test the validations
func TestValidate(t *testing.T) {
	tests := []struct {
		desc     string
		genState types.GenesisState
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
		},
		{
			desc:     "valid genesis state",
			genState: types.GenesisState{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			require.NoError(t, err)
		})
	}
}`,
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

// main function
func main() {
	// print hello world
	fmt.Println("Hello, world!")
	// call new param function
	New(param1, param2)
}

// anotherFunction another function
func anotherFunction() bool { return false }

// TestValidate test the validations
func TestValidate(t *testing.T) {
	tests := []struct {
		desc     string
		genState types.GenesisState
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
		},
		{
			desc:     "valid genesis state",
			genState: types.GenesisState{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			require.NoError(t, err)
		})
	}
}`,
		},
		{
			name: "add a new test case",
			args: args{
				fileContent:  existingContent,
				functionName: "TestValidate",
				functions: []FunctionOptions{
					AppendFuncTestCase(`{
	desc: "valid genesis state",
	genState: GenesisState{},
}`),
				},
			},
			want: `package main

import (
	"fmt"
)

// main function
func main() {
	// print hello world
	fmt.Println("Hello, world!")
	// call new param function
	New(param1, param2)
}

// anotherFunction another function
func anotherFunction() bool {
	// init param
	p := bla.NewParam()
	// start to call something
	p.CallSomething("Another call")
	// return always true
	return true
}

// TestValidate test the validations
func TestValidate(t *testing.T) {
	tests := []struct {
		desc     string
		genState types.GenesisState
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
		},
		{
			desc:     "valid genesis state",
			genState: types.GenesisState{},
		}, {
			desc:     "valid genesis state",
			genState: GenesisState{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			require.NoError(t, err)
		})
	}
}`,
		},
		{
			name: "add two test cases",
			args: args{
				fileContent:  existingContent,
				functionName: "TestValidate",
				functions: []FunctionOptions{
					AppendFuncTestCase(`
{
	desc:     "valid first genesis state",
	genState: GenesisState{},
}`),
					AppendFuncTestCase(`
{
	desc:     "valid second genesis state",
	genState: GenesisState{},
}`),
				},
			},
			want: `package main

import (
	"fmt"
)

// main function
func main() {
	// print hello world
	fmt.Println("Hello, world!")
	// call new param function
	New(param1, param2)
}

// anotherFunction another function
func anotherFunction() bool {
	// init param
	p := bla.NewParam()
	// start to call something
	p.CallSomething("Another call")
	// return always true
	return true
}

// TestValidate test the validations
func TestValidate(t *testing.T) {
	tests := []struct {
		desc     string
		genState types.GenesisState
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
		},
		{
			desc:     "valid genesis state",
			genState: types.GenesisState{},
		}, {
			desc:     "valid first genesis state",
			genState: GenesisState{},
		}, {
			desc:     "valid second genesis state",
			genState: GenesisState{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			require.NoError(t, err)
		})
	}
}`,
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

// main function
func main() {
	// print hello world
	fmt.Println("Hello, world!")
	// call new param function
	New(param1, param2)
}

// anotherFunction another function
func anotherFunction() bool {
	fmt.Println("Appended at line 0.")
	SimpleCall(foo, bar)

	// init param
	p := bla.NewParam()
	// start to call something
	p.CallSomething("Another call")
	fmt.Println("Appended code.")

	// return always true
	return true
}

// TestValidate test the validations
func TestValidate(t *testing.T) {
	tests := []struct {
		desc     string
		genState types.GenesisState
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
		},
		{
			desc:     "valid genesis state",
			genState: types.GenesisState{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			require.NoError(t, err)
		})
	}
}`,
		},
		{
			name: "add all modifications type",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{NewFuncReturn("1")},
			},
			want: strings.ReplaceAll(existingContent, "return true", "return 1\n"),
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

// main function
func main() {
	// print hello world
	fmt.Println("Hello, world!")
	// call new param function
	New(param1, param2)
}

// anotherFunction another function
func anotherFunction() bool {
	// init param
	p := bla.NewParam(baz, bla)
	// start to call something
	p.CallSomething("test2", "Another call", "test1")
	// return always true
	return true
}

// TestValidate test the validations
func TestValidate(t *testing.T) {
	tests := []struct {
		desc     string
		genState types.GenesisState
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
		},
		{
			desc:     "valid genesis state",
			genState: types.GenesisState{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			require.NoError(t, err)
		})
	}
}`,
		},
		{
			name: "add inside call modifications with qualified package name",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions: []FunctionOptions{
					AppendInsideFuncCall("bla.NewParam", "baz", 0),
					AppendInsideFuncCall("bla.NewParam", "bla", -1),
					AppendInsideFuncCall("CallSomething", strconv.Quote("test1"), -1),
				},
			},
			want: `package main

import (
	"fmt"
)

// main function
func main() {
	// print hello world
	fmt.Println("Hello, world!")
	// call new param function
	New(param1, param2)
}

// anotherFunction another function
func anotherFunction() bool {
	// init param
	p := bla.NewParam(baz, bla)
	// start to call something
	p.CallSomething("Another call", "test1")
	// return always true
	return true
}

// TestValidate test the validations
func TestValidate(t *testing.T) {
	tests := []struct {
		desc     string
		genState types.GenesisState
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
		},
		{
			desc:     "valid genesis state",
			genState: types.GenesisState{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			require.NoError(t, err)
		})
	}
}`,
		},
		{
			name: "add inside call modifications with mixed qualified and unqualified names",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions: []FunctionOptions{
					AppendInsideFuncCall("bla.NewParam", "ctx", 0),
					AppendInsideFuncCall("NewParam", "baz", -1),
					AppendInsideFuncCall("p.CallSomething", strconv.Quote("test1"), 0),
					AppendInsideFuncCall("CallSomething", strconv.Quote("test2"), -1),
				},
			},
			want: `package main

import (
	"fmt"
)

// main function
func main() {
	// print hello world
	fmt.Println("Hello, world!")
	// call new param function
	New(param1, param2)
}

// anotherFunction another function
func anotherFunction() bool {
	// init param
	p := bla.NewParam(ctx, baz)
	// start to call something
	p.CallSomething("test1", "Another call", "test2")
	// return always true
	return true
}

// TestValidate test the validations
func TestValidate(t *testing.T) {
	tests := []struct {
		desc     string
		genState types.GenesisState
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
		},
		{
			desc:     "valid genesis state",
			genState: types.GenesisState{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			require.NoError(t, err)
		})
	}
}`,
		},
		{
			name: "add inside struct modifications",
			args: args{
				fileContent: `package main

import (
	"fmt"
)

// anotherFunction another function
func anotherFunction() bool {
	Param{
		Baz: baz,
		Foo: foo,
	}
	Client{baz, foo}
	// return always true
	return true
}

// TestValidate test the validations
func TestValidate(t *testing.T) {
	tests := []struct {
		desc     string
		genState types.GenesisState
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
		},
		{
			desc:     "valid genesis state",
			genState: types.GenesisState{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			require.NoError(t, err)
		})
	}
}`,
				functionName: "anotherFunction",
				functions: []FunctionOptions{
					AppendFuncStruct("Param", "Bar", "bar"),
					AppendFuncStruct("Param", "Bla", "bla"),
					AppendFuncStruct("Client", "", "bar"),
				},
			},
			want: `package main

import (
	"fmt"
)

// anotherFunction another function
func anotherFunction() bool {
	Param{
		Baz: baz,
		Foo: foo,
		Bar: bar,
		Bla: bla,
	}
	Client{baz, foo, bar}
	// return always true
	return true
}

// TestValidate test the validations
func TestValidate(t *testing.T) {
	tests := []struct {
		desc     string
		genState types.GenesisState
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
		},
		{
			desc:     "valid genesis state",
			genState: types.GenesisState{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			require.NoError(t, err)
		})
	}
}`,
		},
		{
			name: "function without test case assertion",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions: []FunctionOptions{
					AppendFuncTestCase(`{
								desc:     "valid second genesis state",
								genState: GenesisState{},
					}`),
				},
			},
			want: existingContent,
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
			err: errors.New("format.Node internal error (16:22: expected ')', found 9 (and 1 more errors))"),
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
			err: errors.New("line number 4 out of range (max 2)"),
		},
		{
			name: "invalid code for append at line",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				functions:    []FunctionOptions{AppendFuncAtLine("9#.(c", 0)},
			},
			err: errors.New("1:24: illegal character U+0023 '#'"),
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
			err: errors.New("format.Node internal error (18:21: illegal character U+0023 '#' (and 4 more errors))"),
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
			want: existingContent,
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

// main function
func main() {
	// Simple function call
	// print hello world
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

// main function
func main() {
	// Simple function call
	// print hello world
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

// main function
func main() {
	// Simple function call
	// print hello world
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

// main function
func main() {
	// Simple function call
	// print hello world
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

// main function
func main() {
	// Simple function call
	// print hello world
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

func TestRemoveFunction(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		funcName    string
		expected    string
		expectError bool
	}{
		{
			name: "remove a simple function",
			content: `package main

func main() {
	println("hello")
}

func anotherFunction() {
	println("another")
}

func thirdFunction() {
	println("third")
}
`,
			funcName: "anotherFunction",
			expected: `package main

func main() {
	println("hello")
}

func thirdFunction() {
	println("third")
}`,
		},
		{
			name: "remove first function",
			content: `package main

func first() {
	println("first")
}

func second() {
	println("second")
}
`,
			funcName: "first",
			expected: `package main

func second() {
	println("second")
}`,
		},
		{
			name: "remove last function",
			content: `package main

func first() {
	println("first")
}

func second() {
	println("second")
}
`,
			funcName: "second",
			expected: `package main

func first() {
	println("first")
}`,
		},
		{
			name: "remove function with comments",
			content: `package main

// main is the entry point
func main() {
	println("main")
}

// helperFunc does something
func helperFunc() {
	println("helper")
}
`,
			funcName: "helperFunc",
			expected: `package main

// main is the entry point
func main() {
	println("main")
}`,
		},
		{
			name: "function not found",
			content: `package main

func main() {
	println("hello")
}
`,
			funcName:    "notFound",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RemoveFunction(tt.content, tt.funcName)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestRemoveFuncCall(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		funcName string
		callName string
		expected string
	}{
		{
			name: "remove a function call",
			content: `package main

func main() {
	fmt.Println("before")
	doSomething()
	fmt.Println("after")
}
`,
			funcName: "main",
			callName: "doSomething",
			expected: `package main

func main() {
	fmt.Println("before")

	fmt.Println("after")
}`,
		},
		{
			name: "remove qualified function call",
			content: `package main

func main() {
	fmt.Println("hello")
	pkg.DoSomething()
	fmt.Println("world")
}
`,
			funcName: "main",
			callName: "pkg.DoSomething",
			expected: `package main

func main() {
	fmt.Println("hello")

	fmt.Println("world")
}`,
		},
		{
			name: "remove multiple calls to same function",
			content: `package main

func main() {
	doSomething()
	fmt.Println("middle")
	doSomething()
}
`,
			funcName: "main",
			callName: "doSomething",
			expected: `package main

func main() {

	fmt.Println("middle")

}`,
		},
		{
			name: "remove call with arguments",
			content: `package main

func process() {
	validate(arg1, arg2)
	execute()
}
`,
			funcName: "process",
			callName: "validate",
			expected: `package main

func process() {

	execute()
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ModifyFunction(tt.content, tt.funcName, RemoveFuncCall(tt.callName))
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}
