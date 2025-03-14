package xast

import (
	"fmt"
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

func main() {
	fmt.Println("Hello, world!")
	New(param1, param2)
}

func anotherFunction() bool {
	p := bla.NewParam()
	p.CallSomething("Another call")
	return true
}

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
					AppendFuncStruct("Param", "Bar", strconv.Quote("bar"), -1),
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
}
`,
		},
		{
			name: "add a new test case",
			args: args{
				fileContent:  existingContent,
				functionName: "TestValidate",
				functions: []FunctionOptions{
					AppendFuncTestCase(`{
								desc:     "valid genesis state",
								genState: GenesisState{},
					}`),
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
	p := bla.NewParam()
	p.CallSomething("Another call")
	return true
}

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

			desc: "valid genesis state",

			genState: GenesisState{}},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			require.NoError(t, err)
		})
	}
}
`,
		},
		{
			name: "add two test cases",
			args: args{
				fileContent:  existingContent,
				functionName: "TestValidate",
				functions: []FunctionOptions{
					AppendFuncTestCase(`{
								desc:     "valid first genesis state",
								genState: GenesisState{},
					}`),
					AppendFuncTestCase(`{
								desc:     "valid second genesis state",
								genState: GenesisState{},
					}`),
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
	p := bla.NewParam()
	p.CallSomething("Another call")
	return true
}

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

			desc: "valid first genesis state",

			genState: GenesisState{},
		}, {

			desc: "valid second genesis state",

			genState: GenesisState{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			require.NoError(t, err)
		})
	}
}
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
}

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
}
`,
				functionName: "anotherFunction",
				functions: []FunctionOptions{
					AppendFuncStruct("Param", "Bar", "bar", -1),
					AppendFuncStruct("Param", "Bla", "bla", 1),
					AppendFuncStruct("Client", "", "bar", 0),
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
}
`,
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
			want: fmt.Sprintln(existingContent),
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
			err: errors.New("format.Node internal error (13:21: illegal character U+0023 '#' (and 4 more errors))"),
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
