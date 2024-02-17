package replacer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

func TestInsertGlobal(t *testing.T) {
	type args struct {
		fileContent string
		globalType  GlobalType
		globals     []GlobalOptions
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "Insert global int var",
			args: args{
				fileContent: `package main

import (
	"fmt"
)

`,
				globalType: GlobalTypeVar,
				globals: []GlobalOptions{
					WithGlobal("myIntVar", "int", "42"),
				},
			},
			want: `package main

import (
	"fmt"
)

var myIntVar int = 42
`,
		},
		{
			name: "Insert global int const",
			args: args{
				fileContent: `package main

import (
	"fmt"
)

`,
				globalType: GlobalTypeConst,
				globals: []GlobalOptions{
					WithGlobal("myIntConst", "int", "42"),
				},
			},
			want: `package main

import (
	"fmt"
)

const myIntConst int = 42
`,
		},
		{
			name: "Insert string const",
			args: args{
				fileContent: `package main

import (
    "fmt"
)

`,
				globalType: GlobalTypeConst,
				globals: []GlobalOptions{
					WithGlobal("myStringConst", "string", `"hello"`),
				},
			},
			want: `package main

import (
	"fmt"
)

const myStringConst string = "hello"
`,
		},
		{
			name: "Insert multiples consts",
			args: args{
				fileContent: `package main

import (
	"fmt"
)

`,
				globalType: GlobalTypeConst,
				globals: []GlobalOptions{
					WithGlobal("myStringConst", "string", `"hello"`),
					WithGlobal("myBoolConst", "bool", "true"),
					WithGlobal("myUintConst", "uint64", "40"),
				},
			},
			want: `package main

import (
	"fmt"
)

const myStringConst string = "hello"
const myBoolConst bool = true
const myUintConst uint64 = 40
`,
		},
		{
			name: "Insert global int var with not imports",
			args: args{
				fileContent: `package main
`,
				globalType: GlobalTypeVar,
				globals: []GlobalOptions{
					WithGlobal("myIntVar", "int", "42"),
				},
			},
			want: `package main

var myIntVar int = 42
`,
		},
		{
			name: "Insert global int var int an empty file",
			args: args{
				fileContent: ``,
				globalType:  GlobalTypeVar,
				globals: []GlobalOptions{
					WithGlobal("myIntVar", "int", "42"),
				},
			},
			err: errors.New("1:1: expected 'package', found 'EOF'"),
		},
		{
			name: "Insert a custom var",
			args: args{
				fileContent: `package main`,
				globalType:  GlobalTypeVar,
				globals: []GlobalOptions{
					WithGlobal("fooVar", "foo", "42"),
				},
			},
			want: `package main

var fooVar foo = 42
`,
		},
		{
			name: "Insert an invalid var",
			args: args{
				fileContent: `package main`,
				globalType:  GlobalTypeVar,
				globals: []GlobalOptions{
					WithGlobal("myInvalidVar", "invalid", "AEF#3fa."),
				},
			},
			err: errors.New("1:4: illegal character U+0023 '#'"),
		},
		{
			name: "Insert an invalid type",
			args: args{
				fileContent: `package main`,
				globalType:  102,
				globals: []GlobalOptions{
					WithGlobal("myInvalidVar", "invalid", "AEF#3fa."),
				},
			},
			err: errors.New("unsupported global type: 102"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InsertGlobal(tt.args.fileContent, tt.args.globalType, tt.args.globals...)
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
