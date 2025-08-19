package xast

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
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

// This is a comment
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

// This is a comment
`,
		},
		{
			name: "Insert global int var without type",
			args: args{
				fileContent: `package main

import (
	"fmt"
)

`,
				globalType: GlobalTypeVar,
				globals: []GlobalOptions{
					WithGlobal("myIntVar", "", "42"),
				},
			},
			want: `package main

import (
	"fmt"
)

var myIntVar = 42
`,
		},
		{
			name: "Insert global int const",
			args: args{
				fileContent: `package main

import (
	"fmt"
)

// This is a comment
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

// This is a comment
`,
		},
		{
			name: "Insert string const",
			args: args{
				fileContent: `package main

import (
    "fmt"
)

// This is a comment
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

// This is a comment
`,
		},
		{
			name: "Insert string const when already exist one",
			args: args{
				fileContent: `package main

import (
    "fmt"
)

// myIntConst is my const int
const myIntConst int = 42

// This is a comment
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

// myIntConst is my const int
const myIntConst int = 42

// This is a comment
`,
		},
		{
			name: "Insert multiples consts",
			args: args{
				fileContent: `package main

import (
	"fmt"
)

// This is a comment
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

// This is a comment
`,
		},
		{
			name: "Insert global int var with not imports",
			args: args{
				fileContent: `package main

// This is a comment
`,
				globalType: GlobalTypeVar,
				globals: []GlobalOptions{
					WithGlobal("myIntVar", "int", "42"),
				},
			},
			want: `package main

var myIntVar int = 42

// This is a comment
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
				globalType:  "invalid",
				globals: []GlobalOptions{
					WithGlobal("fooVar", "foo", "42"),
				},
			},
			err: errors.New("unsupported global type: invalid"),
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

func TestModifyStruct(t *testing.T) {
	type args struct {
		fileContent string
		structName  string
		options     []StructOpts
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "Add field to existing struct",
			args: args{
				fileContent: `package main

type MyStruct struct {
	ExistingField int
}
`,
				structName: "MyStruct",
				options:    []StructOpts{AppendStructValue("NewField", "string")},
			},
			want: `package main

type MyStruct struct {
	ExistingField int
	NewField      string
}
`,
		},
		{
			name: "Add field to empty struct",
			args: args{
				fileContent: `package main

type EmptyStruct struct {
}
`,
				structName: "EmptyStruct",
				options:    []StructOpts{AppendStructValue("NewField", "string")},
			},
			want: `package main

type EmptyStruct struct {
	NewField string
}
`,
		},
		{
			name: "Struct not found",
			args: args{
				fileContent: `package main

type AnotherStruct struct {
	ExistingField int
}
`,
				structName: "NonExistentStruct",
				options:    []StructOpts{AppendStructValue("NewField", "string")},
			},
			err: errors.New(`struct "NonExistentStruct" not found in file content`),
		},
		{
			name: "Invalid Go code",
			args: args{
				fileContent: `package main

type MyStruct`,
				structName: "MyStruct",
				options:    []StructOpts{AppendStructValue("NewField", "string")},
			},
			err: errors.New("3:14: expected type, found newline"),
		},
		{
			name: "Add field after multiple existing fields",
			args: args{
				fileContent: `package main

type MyStruct struct {
	Field1 int
	Field2 string
}
`,
				structName: "MyStruct",
				options:    []StructOpts{AppendStructValue("Field3", "bool")},
			},
			want: `package main

type MyStruct struct {
	Field1 int
	Field2 string
	Field3 bool
}
`,
		},
		{
			name: "Empty file input",
			args: args{
				fileContent: ``,
				structName:  "MyStruct",
				options:     []StructOpts{AppendStructValue("NewField", "string")},
			},
			err: errors.New("1:1: expected 'package', found 'EOF'"),
		},
		{
			name: "Add field with pointer type",
			args: args{
				fileContent: `package main

type MyStruct struct {
	ExistingField int
}
`,
				structName: "MyStruct",
				options:    []StructOpts{AppendStructValue("PointerField", "*int")},
			},
			want: `package main

type MyStruct struct {
	ExistingField int
	PointerField  *int
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ModifyStruct(tt.args.fileContent, tt.args.structName, tt.args.options...)
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

func TestModifyGlobalArrayVar(t *testing.T) {
	type args struct {
		fileContent string
		globalName  string
		options     []GlobalArrayOpts
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "Add field to custom variable array",
			args: args{
				fileContent: `package app
var (
	moduleAccPerms = []*authmodulev1.ModuleAccountPermission{
		{Account: nft.ModuleName},
		{Account: ibctransfertypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}},
	}
)
`,
				globalName: "moduleAccPerms",
				options:    []GlobalArrayOpts{AppendGlobalArrayValue("{Account: icatypes.ModuleName}")},
			},
			want: `package app

var (
	moduleAccPerms = []*authmodulev1.ModuleAccountPermission{
		{Account: nft.ModuleName},
		{Account: ibctransfertypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}},
		{Account: icatypes.ModuleName},
	}
)
`,
		},
		{
			name: "Add field to string variable array",
			args: args{
				fileContent: `package app

var (
	blockAccAddrs = []string{
		authtypes.FeeCollectorName,
		distrtypes.ModuleName,
		minttypes.ModuleName,
		stakingtypes.BondedPoolName,
		stakingtypes.NotBondedPoolName,
	}
)
`,
				globalName: "blockAccAddrs",
				options:    []GlobalArrayOpts{AppendGlobalArrayValue("nft.ModuleName")},
			},
			want: `package app

var (
	blockAccAddrs = []string{
		authtypes.FeeCollectorName,
		distrtypes.ModuleName,
		minttypes.ModuleName,
		stakingtypes.BondedPoolName,
		stakingtypes.NotBondedPoolName,
		nft.ModuleName,
	}
)
`,
		},
		{
			name: "name not found",
			args: args{
				fileContent: `package app

var (
	blockAccAddrs = []string{
		authtypes.FeeCollectorName,
		distrtypes.ModuleName,
		minttypes.ModuleName,
		stakingtypes.BondedPoolName,
		stakingtypes.NotBondedPoolName,
	}
)
`,
				globalName: "notFound",
				options:    []GlobalArrayOpts{AppendGlobalArrayValue("nft.ModuleName")},
			},
			err: errors.New("global array \"notFound\" not found in file content"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ModifyGlobalArrayVar(tt.args.fileContent, tt.args.globalName, tt.args.options...)
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
