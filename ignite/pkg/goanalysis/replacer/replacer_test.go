package replacer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

func TestAppendFunction(t *testing.T) {
	type args struct {
		fileContent string
		function    string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "Append a function after the package declaration",
			args: args{
				fileContent: `package main`,
				function: `func add(a, b int) int {
	return a + b
}`,
			},
			want: `package main

func add(a, b int) int {
	return a + b
}
`,
		},
		{
			name: "Append a function after a var",
			args: args{
				fileContent: `package main

import (
	"fmt"
)

var myIntVar int = 42
`,
				function: `func add(a, b int) int {
	return a + b
}`,
			},
			want: `package main

import (
	"fmt"
)

var myIntVar int = 42

func add(a, b int) int {
	return a + b
}
`,
		},
		{
			name: "Append a function after the import",
			args: args{
				fileContent: `package main

import (
	"fmt"
)
`,
				function: `func add(a, b int) int {
	return a + b
}`,
			},
			want: `package main

import (
	"fmt"
)

func add(a, b int) int {
	return a + b
}
`,
		},
		{
			name: "Append a function after another function",
			args: args{
				fileContent: `package main

import (
	"fmt"
)

var myIntVar int = 42

func myFunction() int {
    return 42
}
`,
				function: `func add(a, b int) int {
	return a + b
}`,
			},
			want: `package main

import (
	"fmt"
)

var myIntVar int = 42

func myFunction() int {
	return 42
}
func add(a, b int) int {
	return a + b
}
`,
		},
		{
			name: "Append a function in an empty file",
			args: args{
				fileContent: ``,
				function: `func add(a, b int) int {
	return a + b
}`,
			},
			err: errors.New("1:1: expected 'package', found 'EOF'"),
		},
		{
			name: "Append a empty function",
			args: args{
				fileContent: `package main`,
				function:    ``,
			},
			err: errors.New("no function declaration found in the provided function body"),
		},
		{
			name: "Append an invalid function",
			args: args{
				fileContent: `package main`,
				function:    `@,.l.e,`,
			},
			err: errors.New("2:1: illegal character U+0040 '@'"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AppendFunction(tt.args.fileContent, tt.args.function)
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

func TestAppendCode(t *testing.T) {
	existingContent := `package main

import (
    "fmt"
)

func main() {
    fmt.Println("Hello, world!")
}

func anotherFunction() bool {
    // Some code here
    fmt.Println("Another function")
    return true
}`

	type args struct {
		fileContent  string
		functionName string
		codeToInsert string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "append code to the end of the function",
			args: args{
				fileContent:  existingContent,
				functionName: "main",
				codeToInsert: "fmt.Println(\"Inserted code here\")",
			},
			want: `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
	fmt.Println("Inserted code here")

}

func anotherFunction() bool {
	// Some code here
	fmt.Println("Another function")
	return true
}
`,
		},
		{
			name: "append code with return statement",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				codeToInsert: "fmt.Println(\"Inserted code here\")",
			},
			want: `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
}

func anotherFunction() bool {
	// Some code here
	fmt.Println("Another function")
	fmt.Println("Inserted code here")

	return true
}
`,
		},
		{
			name: "function not found",
			args: args{
				fileContent:  existingContent,
				functionName: "nonexistentFunction",
				codeToInsert: "fmt.Println(\"Inserted code here\")",
			},
			err: errors.New("function nonexistentFunction not found"),
		},
		{
			name: "invalid code",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				codeToInsert: "%#)(u309f/..\"",
			},
			err: errors.New("1:1: expected operand, found '%' (and 2 more errors)"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AppendCodeToFunction(tt.args.fileContent, tt.args.functionName, tt.args.codeToInsert)
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

func TestReplaceReturn(t *testing.T) {
	existingContent := `package main

import (
    "fmt"
)

func main() {
    x := calculate()
    fmt.Println("Result:", x)
}

func calculate() int {
    return 42
}`

	type args struct {
		fileContent  string
		functionName string
		returnVars   []string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "replace return statement with a single variable",
			args: args{
				fileContent:  existingContent,
				functionName: "calculate",
				returnVars:   []string{"result"},
			},
			want: `package main

import (
	"fmt"
)

func main() {
	x := calculate()
	fmt.Println("Result:", x)
}

func calculate() int {
	return result

}
`,
		},
		{
			name: "replace return statement with multiple variables",
			args: args{
				fileContent:  existingContent,
				functionName: "calculate",
				returnVars:   []string{"result", "err"},
			},
			want: `package main

import (
	"fmt"
)

func main() {
	x := calculate()
	fmt.Println("Result:", x)
}

func calculate() int {
	return result, err

}
`,
		},
		{
			name: "function not found",
			args: args{
				fileContent:  existingContent,
				functionName: "nonexistentFunction",
				returnVars:   []string{"result"},
			},
			err: errors.New("function nonexistentFunction not found"),
		},
		{
			name: "invalid result",
			args: args{
				fileContent:  existingContent,
				functionName: "nonexistentFunction",
				returnVars:   []string{"ae@@of..!\""},
			},
			err: errors.New("1:3: illegal character U+0040 '@' (and 1 more errors)"),
		},
		{
			name: "reserved word",
			args: args{
				fileContent:  existingContent,
				functionName: "nonexistentFunction",
				returnVars:   []string{"range"},
			},
			err: errors.New("1:1: expected operand, found 'range'"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReplaceReturnVars(tt.args.fileContent, tt.args.functionName, tt.args.returnVars...)
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

func TestReplaceCode(t *testing.T) {
	var (
		newFunction     = `fmt.Println("This is the new function.")`
		existingContent = `package main

import (
    "fmt"
)

func main() {
    fmt.Println("Hello, world!")
}

func oldFunction() {
    fmt.Println("This is the old function.")
}`
	)

	type args struct {
		fileContent     string
		oldFunctionName string
		newFunction     string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "replace function implementation",
			args: args{
				fileContent:     existingContent,
				oldFunctionName: "oldFunction",
				newFunction:     newFunction,
			},
			want: `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
}

func oldFunction() { fmt.Println("This is the new function.") }
`,
		},
		{
			name: "replace main function implementation",
			args: args{
				fileContent:     existingContent,
				oldFunctionName: "main",
				newFunction:     newFunction,
			},
			want: `package main

import (
	"fmt"
)

func main() { fmt.Println("This is the new function.") }

func oldFunction() {
	fmt.Println("This is the old function.")
}
`,
		},
		{
			name: "function not found",
			args: args{
				fileContent:     existingContent,
				oldFunctionName: "nonexistentFunction",
				newFunction:     newFunction,
			},
			err: errors.New("function nonexistentFunction not found in file content"),
		},
		{
			name: "invalid new function",
			args: args{
				fileContent:     existingContent,
				oldFunctionName: "nonexistentFunction",
				newFunction:     "ae@@of..!\"",
			},
			err: errors.New("1:25: illegal character U+0040 '@' (and 2 more errors)"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReplaceFunctionContent(tt.args.fileContent, tt.args.oldFunctionName, tt.args.newFunction)
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

func TestAppendParamToFunctionCall(t *testing.T) {
	type args struct {
		fileContent      string
		functionName     string
		functionCallName string
		index            int
		paramToAdd       string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "add new parameter to index 1 in a function extension",
			args: args{
				fileContent: `package main

func myFunction() {
	p := NewParam()
	p.New("param1", "param2")
}`,
				functionName:     "myFunction",
				functionCallName: "New",
				paramToAdd:       "param3",
				index:            1,
			},
			want: `package main

func myFunction() {
	p := NewParam()
	p.New("param1", "param3", "param2")
}
`,
		},
		{
			name: "add new parameter to index 1",
			args: args{
				fileContent: `package main

func myFunction() {
	New("param1", "param2")
}`,
				functionName:     "myFunction",
				functionCallName: "New",
				paramToAdd:       "param3",
				index:            1,
			},
			want: `package main

func myFunction() {
	New("param1", "param3", "param2")
}
`,
		},
		{
			name: "add a new parameter for two functions",
			args: args{
				fileContent: `package main

func myFunction() {
	New("param1", "param2")
	New("param1", "param2", "param4")
}`,
				functionName:     "myFunction",
				functionCallName: "New",
				index:            -1,
				paramToAdd:       "param3",
			},
			want: `package main

func myFunction() {
	New("param1", "param2", "param3")
	New("param1", "param2", "param4", "param3")
}
`,
		},
		{
			name: "function not found",
			args: args{
				fileContent: `package main

func anotherFunction() {
	New("param1", "param2")
}`,
				functionName:     "myFunction",
				functionCallName: "New",
				paramToAdd:       "param3",
				index:            1,
			},
			err: errors.Errorf("function myFunction not found or no calls to New inside the function"),
		},
		{
			name: "function call not found",
			args: args{
				fileContent: `package main

func myFunction() {
	AnotherFunction("param1", "param2")
}`,
				functionName:     "myFunction",
				functionCallName: "New",
				paramToAdd:       "param3",
				index:            1,
			},
			want: `package main

func myFunction() {
	AnotherFunction("param1", "param2")
}`,
			err: errors.Errorf("function myFunction not found or no calls to New inside the function"),
		},
		{
			name: "index out of range",
			args: args{
				fileContent: `package main

func myFunction() {
	New("param1", "param2")
}`,
				functionName:     "myFunction",
				functionCallName: "New",
				paramToAdd:       "param3",
				index:            3,
			},
			err: errors.Errorf("index out of range"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AppendParamToFunctionCall(
				tt.args.fileContent,
				tt.args.functionName,
				tt.args.functionCallName,
				tt.args.paramToAdd,
				tt.args.index,
			)
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
