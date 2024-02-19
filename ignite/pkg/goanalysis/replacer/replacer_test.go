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

func TestAppendImports(t *testing.T) {
	existingContent := `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
}`

	type args struct {
		fileContent      string
		importStatements []string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "Add single import statement",
			args: args{
				fileContent:      existingContent,
				importStatements: []string{"strings"},
			},
			want: `package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, world!")
}
`,
			err: nil,
		},
		{
			name: "Add multiple import statements",
			args: args{
				fileContent:      existingContent,
				importStatements: []string{"st strings", "strconv", "os"},
			},
			want: `package main

import (
	"fmt"
	"os"
	"strconv"
	st "strings"
)

func main() {
	fmt.Println("Hello, world!")
}
`,
			err: nil,
		},
		{
			name: "Add multiple import statements with an existing one",
			args: args{
				fileContent:      existingContent,
				importStatements: []string{"st strings", "strconv", "os", "fmt"},
			},
			want: `package main

import (
	"fmt"
	"os"
	"strconv"
	st "strings"
)

func main() {
	fmt.Println("Hello, world!")
}
`,
			err: nil,
		},
		{
			name: "Add duplicate import statement",
			args: args{
				fileContent:      existingContent,
				importStatements: []string{"fmt"},
			},
			want: existingContent + "\n",
			err:  nil,
		},
		{
			name: "No import statement",
			args: args{
				fileContent: `package main

func main() {
	fmt.Println("Hello, world!")
}`,
				importStatements: []string{"fmt"},
			},
			want: `package main

import  "fmt"

func main() {
	fmt.Println("Hello, world!")
}
`,
			err: nil,
		},
		{
			name: "No import statement and add two imports",
			args: args{
				fileContent: `package main

func main() {
	fmt.Println("Hello, world!")
}`,
				importStatements: []string{"fmt", "os"},
			},
			want: `package main

import (
	 "fmt"
	 "os"
)

func main() {
	fmt.Println("Hello, world!")
}
`,
			err: nil,
		},
		{
			name: "Add invalid import statement",
			args: args{
				fileContent:      existingContent,
				importStatements: []string{"fmt\""},
			},
			err: errors.New("format.Node internal error (5:8: string literal not terminated (and 1 more errors))"),
		},
		{
			name: "Add empty file content",
			args: args{
				fileContent:      "",
				importStatements: []string{"fmt"},
			},
			err: errors.New("1:1: expected 'package', found 'EOF'"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AppendImports(tt.args.fileContent, tt.args.importStatements...)
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
			name: "Append code to the end of the function",
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
			name: "Append code with return statement",
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
			name: "Function not found",
			args: args{
				fileContent:  existingContent,
				functionName: "nonexistentFunction",
				codeToInsert: "fmt.Println(\"Inserted code here\")",
			},
			err: errors.New("function nonexistentFunction not found"),
		},
		{
			name: "Invalid code",
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
			name: "Replace return statement with a single variable",
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
			name: "Replace return statement with multiple variables",
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
			name: "Function not found",
			args: args{
				fileContent:  existingContent,
				functionName: "nonexistentFunction",
				returnVars:   []string{"result"},
			},
			err: errors.New("function nonexistentFunction not found"),
		},
		{
			name: "Invalid result",
			args: args{
				fileContent:  existingContent,
				functionName: "nonexistentFunction",
				returnVars:   []string{"ae@@of..!\""},
			},
			err: errors.New("1:3: illegal character U+0040 '@' (and 1 more errors)"),
		},
		{
			name: "Reserved word",
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
			name: "Replace function implementation",
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
			name: "Replace main function implementation",
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
			name: "Function not found",
			args: args{
				fileContent:     existingContent,
				oldFunctionName: "nonexistentFunction",
				newFunction:     newFunction,
			},
			err: errors.New("function nonexistentFunction not found in file content"),
		},
		{
			name: "Invalid new function",
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
	tests := []struct {
		name             string
		fileContent      string
		functionName     string
		functionCallName string
		index            int
		paramToAdd       string
		want             string
		err              error
	}{
		{
			name: "add new parameter to index 1 in a function extension",
			fileContent: `package main

func myFunction() {
	p := NewParam()
	p.New("param1", "param2")
}`,
			functionName:     "myFunction",
			functionCallName: "New",
			paramToAdd:       "param3",
			index:            1,
			want: `package main

func myFunction() {
	p := NewParam()
	p.New("param1", "param3", "param2")
}
`,
		},
		{
			name: "add new parameter to index 1",
			fileContent: `package main

func myFunction() {
	New("param1", "param2")
}`,
			functionName:     "myFunction",
			functionCallName: "New",
			paramToAdd:       "param3",
			index:            1,
			want: `package main

func myFunction() {
	New("param1", "param3", "param2")
}
`,
		},
		{
			name: "add a new parameter for two functions",
			fileContent: `package main

func myFunction() {
	New("param1", "param2")
	New("param1", "param2", "param4")
}`,
			functionName:     "myFunction",
			functionCallName: "New",
			index:            -1,
			paramToAdd:       "param3",
			want: `package main

func myFunction() {
	New("param1", "param2", "param3")
	New("param1", "param2", "param4", "param3")
}
`,
		},
		{
			name: "FunctionNotFound",
			fileContent: `package main

func anotherFunction() {
	New("param1", "param2")
}`,
			functionName:     "myFunction",
			functionCallName: "New",
			paramToAdd:       "param3",
			index:            1,
			err:              errors.Errorf("function myFunction not found or no calls to New inside the function"),
		},
		{
			name: "FunctionCallNotFound",
			fileContent: `package main

func myFunction() {
	AnotherFunction("param1", "param2")
}`,
			functionName:     "myFunction",
			functionCallName: "New",
			paramToAdd:       "param3",
			index:            1,
			want: `package main

func myFunction() {
	AnotherFunction("param1", "param2")
}`,
			err: errors.Errorf("function myFunction not found or no calls to New inside the function"),
		},
		{
			name: "IndexOutOfRange",
			fileContent: `package main

func myFunction() {
	New("param1", "param2")
}`,
			functionName:     "myFunction",
			functionCallName: "New",
			paramToAdd:       "param3",
			index:            3,
			err:              errors.Errorf("index out of range"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AppendParamToFunctionCall(tt.fileContent, tt.functionName, tt.functionCallName, tt.paramToAdd, tt.index)
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
