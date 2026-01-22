package xast

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

func TestAppendImports(t *testing.T) {
	existingContent := `package main

import (
	"fmt"
)

func main() {
	// print hello world
	fmt.Println("Hello, world!")
}`

	type args struct {
		fileContent string
		imports     []ImportOptions
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "add single import statement",
			args: args{
				fileContent: existingContent,
				imports: []ImportOptions{
					WithImport("strings"),
				},
			},
			want: `package main

import (
	"fmt"
	"strings"
)

func main() {
	// print hello world
	fmt.Println("Hello, world!")
}
`,
		},
		{
			name: "add multiple import statements",
			args: args{
				fileContent: existingContent,
				imports: []ImportOptions{
					WithNamedImport("st", "strings"),
					WithImport("strconv"),
					WithImport("os"),
				},
			},
			want: `package main

import (
	"fmt"
	"os"
	"strconv"
	st "strings"
)

func main() {
	// print hello world
	fmt.Println("Hello, world!")
}
`,
		},
		{
			name: "add multiple import statements with an existing one",
			args: args{
				fileContent: existingContent,
				imports: []ImportOptions{
					WithNamedImport("st", "strings"),
					WithImport("strconv"),
					WithImport("os"),
				},
			},
			want: `package main

import (
	"fmt"
	"os"
	"strconv"
	st "strings"
)

func main() {
	// print hello world
	fmt.Println("Hello, world!")
}
`,
		},
		{
			name: "add import to specific index",
			args: args{
				fileContent: `package main

import (
	"fmt"
	"os"
	st "strings"
)`,
				imports: []ImportOptions{
					WithImport("strconv"),
				},
			},
			want: `package main

import (
	"fmt"
	"os"
	"strconv"
	st "strings"
)
`,
		},
		{
			name: "add multiple imports to specific index",
			args: args{
				fileContent: `package main

import (
	"fmt"
	"os"
	st "strings"
)`,
				imports: []ImportOptions{
					WithImport("strconv"),
					WithNamedImport("", "testing"),
					WithImport("bytes"),
				},
			},
			want: `package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	st "strings"
	"testing"
)
`,
		},
		{
			name: "add duplicate import statement",
			args: args{
				fileContent: existingContent,
				imports: []ImportOptions{
					WithImport("fmt"),
				},
			},
			want: `package main

import "fmt"

func main() {
	// print hello world
	fmt.Println("Hello, world!")
}
`,
		},
		{
			name: "no import statement",
			args: args{
				fileContent: `package main

func main() {
	// print hello world
	fmt.Println("Hello, world!")
}`,
				imports: []ImportOptions{
					WithImport("fmt"),
				},
			},
			want: `package main

import "fmt"

func main() {
	// print hello world
	fmt.Println("Hello, world!")
}
`,
		},
		{
			name: "no import statement and add two imports",
			args: args{
				fileContent: `package main

func main() {
	// print hello world
	fmt.Println("Hello, world!")
}`,
				imports: []ImportOptions{
					WithImport("fmt"),
					WithImport("os"),
				},
			},
			want: `package main

import (
	"fmt"
	"os"
)

func main() {
	// print hello world
	fmt.Println("Hello, world!")
}
`,
		},
		{
			name: "add invalid import name",
			args: args{
				fileContent: existingContent,
				imports: []ImportOptions{
					WithNamedImport("fmt\"", "fmt\""),
				},
			},
			err: errors.New("format.Node internal error (5:8: expected ';', found fmt (and 2 more errors))"),
		},
		{
			name: "add empty file content",
			args: args{
				fileContent: "",
				imports: []ImportOptions{
					WithImport("fmt"),
				},
			},
			err: errors.New("1:1: expected 'package', found 'EOF'"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AppendImports(tt.args.fileContent, tt.args.imports...)
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

func TestRemoveImports(t *testing.T) {
	type args struct {
		fileContent string
		imports     []ImportOptions
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "remove single import statement",
			args: args{
				fileContent: `package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, world!")
}`,
				imports: []ImportOptions{
					WithImport("strings"),
				},
			},
			want: `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
}
`,
		},
		{
			name: "remove multiple import statements",
			args: args{
				fileContent: `package main

import (
	"fmt"
	"os"
	"strconv"
	st "strings"
)

func main() {
	fmt.Println("Hello, world!")
}`,
				imports: []ImportOptions{
					WithNamedImport("st", "strings"),
					WithImport("strconv"),
					WithImport("os"),
				},
			},
			want: `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
}
`,
		},
		{
			name: "remove all imports",
			args: args{
				fileContent: `package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, world!")
}`,
				imports: []ImportOptions{
					WithImport("fmt"),
					WithImport("strings"),
				},
			},
			want: `package main

func main() {
	fmt.Println("Hello, world!")
}
`,
		},
		{
			name: "remove non-existent import",
			args: args{
				fileContent: `package main

import "fmt"

func main() {
	fmt.Println("Hello, world!")
}`,
				imports: []ImportOptions{
					WithImport("strings"),
				},
			},
			want: `package main

import "fmt"

func main() {
	fmt.Println("Hello, world!")
}
`,
		},
		{
			name: "remove named import",
			args: args{
				fileContent: `package main

import (
	"fmt"
	st "strings"
)

func main() {
	fmt.Println("Hello, world!")
}`,
				imports: []ImportOptions{
					WithNamedImport("st", "strings"),
				},
			},
			want: `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
}
`,
		},
		{
			name: "remove import from file with no imports",
			args: args{
				fileContent: `package main

func main() {
	fmt.Println("Hello, world!")
}`,
				imports: []ImportOptions{
					WithImport("fmt"),
				},
			},
			want: `package main

func main() {
	fmt.Println("Hello, world!")
}
`,
		},
		{
			name: "remove empty file content",
			args: args{
				fileContent: "",
				imports: []ImportOptions{
					WithImport("fmt"),
				},
			},
			err: errors.New("1:1: expected 'package', found 'EOF'"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RemoveImports(tt.args.fileContent, tt.args.imports...)
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
