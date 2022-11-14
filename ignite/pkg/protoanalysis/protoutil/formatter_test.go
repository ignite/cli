package protoutil

import (
	"fmt"
	"strings"
	"testing"

	"github.com/emicklei/proto"
	"github.com/stretchr/testify/require"
)

// Note that this basically does what parser.go:Printer does
// albeit using the alternative formatter.
func testPrinter(pf *proto.Proto) string {
	output := new(strings.Builder)
	NewFormatter(output, "  ").Format(pf) // 2 spaces

	return output.String()
}

// text in, text out.
type testFormatting struct {
	name     string
	in       string
	out      string
	printFmt bool
}

var formattingTests = []testFormatting{
	{"empty", "", "", false},
	{
		name: "Single syntax stmt",
		in:   "syntax = \"proto3\";",
		out:  "syntax = \"proto3\";\n",
	},
	{
		name: "Syntax with comment.",
		in:   "// a syntax comment\nsyntax = \"proto3\";",
		out:  "// a syntax comment\nsyntax = \"proto3\";\n",
	},
	{
		name: "Syntax with inline comment.",
		in:   `syntax = "proto3"; // a syntax comment`,
		out:  "syntax = \"proto3\"; // a syntax comment\n",
	},
	{
		name: "Syntax with multiple comments",
		in:   "//a syntax comment\n// another syntax comment\nsyntax = \"proto3\";",
		out:  "// a syntax comment\n// another syntax comment\nsyntax = \"proto3\";\n",
	},
	// Package.
	{
		name: "empty",
		in:   "",
		out:  "",
	},
	{
		name: "Single package stmt",
		in:   "package foo;",
		out:  "package foo;\n",
	},
	{
		name: "Package with comment.",
		in:   "// a package comment\npackage foo;",
		out:  "// a package comment\npackage foo;\n",
	},
	{
		name: "Package with multiple comments",
		in:   "//a syntax comment\n// another syntax comment\npackage foo;",
		out:  "// a syntax comment\n// another syntax comment\npackage foo;\n",
	},
	{
		name: "Package with inline comment.",
		in:   `package foo; // a package comment`,
		out:  "package foo; // a package comment\n",
	},
	{
		name: "Package w/o syntax preceeding.",
		in:   "\n\n\n\n// a package comment\npackage foo;",
		out:  "// a package comment\npackage foo;\n",
	},
	{
		// Should ignore all empty lines iff preceeded by syntax
		name: "Package preceeded by syntax.",
		in:   "syntax = \"proto3\";\n\n\npackage foo;",
		out:  "syntax = \"proto3\";\npackage foo;\n",
	},
	{
		// Should leave a new line if preceeded by import/option etc
		name: "Package preceeded by import.",
		in:   "\n\n\nimport \"foo\";\n\n\n\npackage foo;",
		out:  "import \"foo\";\n\npackage foo;\n",
	},
	{
		// Yes, place a new line if a comment is there.
		name: "Package preceeded by syntax && embedded comment.",
		in:   "syntax = \"proto3\";\n\n\n// a package comment\npackage foo;",
		out:  "syntax = \"proto3\";\n\n// a package comment\npackage foo;\n",
	},
	// Imports.
	{
		name: "Single import stmt",
		in:   "import \"foo.proto\";",
		out:  "import \"foo.proto\";\n",
	},
	{
		name: "Import a proto as weak",
		in:   "import weak \"foo.proto\";",
		out:  "import weak \"foo.proto\";\n",
	},
	{
		name: "Import a proto as public",
		in:   "import public \"foo.proto\";",
		out:  "import public \"foo.proto\";\n",
	},
	{
		name: "Import with a longer path.",
		in:   "import \"cosmos/base/query/v1beta1/pagination.proto\";",
		out:  "import \"cosmos/base/query/v1beta1/pagination.proto\";\n",
	},
	{
		name: "Import with a comment.",
		in:   "// a import comment\nimport \"foo.proto\";",
		out:  "// a import comment\nimport \"foo.proto\";\n",
	},
	{
		name: "Import with multiple comments",
		in:   "//a syntax comment\n// another syntax comment\nimport \"foo.proto\";",
		out:  "// a syntax comment\n// another syntax comment\nimport \"foo.proto\";\n",
	},
	{
		name: "Import with inline comment.",
		in:   `import "foo.proto"; // a import comment`,
		out:  "import \"foo.proto\"; // a import comment\n",
	},
	{
		name: "Import preceeded by another import.",
		in:   "import \"foo.proto\";\n\n\n\nimport \"bar.proto\";",
		out:  "import \"foo.proto\";\nimport \"bar.proto\";\n",
	},
	{
		name: "Import preceeded by a package.",
		in:   "package foo;\n\n\n\nimport \"bar.proto\";",
		out:  "package foo;\n\nimport \"bar.proto\";\n",
	},
	{
		name: "Import preceeded by syntax",
		in:   "syntax = \"proto3\";\n\n\n\nimport \"bar.proto\";",
		out:  "syntax = \"proto3\";\n\nimport \"bar.proto\";\n",
	},
	{
		name: "Import with comment preceeded by syntax.",
		in:   "syntax = \"proto3\";\n\n\n\n// a import comment\nimport \"foo.proto\";",
		out:  "syntax = \"proto3\";\n\n// a import comment\nimport \"foo.proto\";\n",
	},
	{
		name: "Import with comment preceeded by import.",
		in:   "import \"foo.proto\";\n\n\n\n// a import comment\nimport \"bar.proto\";",
		out:  "import \"foo.proto\";\n// a import comment\nimport \"bar.proto\";\n",
	},
	// Options.
	{
		name: "Single option stmt",
		in:   "option foo = 1;",
		out:  "option foo = 1;\n",
	},
	{
		name: "Option with a comment",
		in:   "// a option comment\noption foo = 1;",
		out:  "// a option comment\noption foo = 1;\n",
	},
	{
		name: "Option followed by option",
		in:   "option foo = 1;\n\n\noption bar = 2;",
		out:  "option foo = 1;\noption bar = 2;\n",
	},
	{
		name: "Package followed by option",
		in:   "package foo;\n\n\noption bar = 2;",
		out:  "package foo;\n\noption bar = 2;\n",
	},
	{
		name: "Option followed by option with a comment",
		in:   "option foo = 1;\n\n\n// a option comment\noption bar = 2;",
		out:  "option foo = 1;\n// a option comment\noption bar = 2;\n",
	},
	{
		name: "Option followed by option with an inline comment",
		in:   "option foo = 1;\n\n\n\noption bar = 2; // inline",
		out:  "option foo = 1;\noption bar = 2; // inline\n",
	},
	// Enums.
	{
		name: "Single enum stmt",
		in:   "enum Foo {}",
		out:  "enum Foo {}\n",
	},
	{
		// top level elements have extra nl.
		name: "Enum followed by enum",
		in:   "enum Foo {}\n\n\nenum Bar {}",
		out:  "enum Foo {}\n\nenum Bar {}\n",
	},
	{
		// Ok, one newline since we have the comment there.
		name: "Enum followed by an enum with a comment.",
		in:   "enum Foo {}\n\n\n// a enum comment\nenum Bar {}",
		out:  "enum Foo {}\n// a enum comment\nenum Bar {}\n",
	},
	{
		name: "Enum without fields but mucho space.",
		in:   "enum Foo{\n\n\n\n\n\t\t\t\t\n}",
		out:  "enum Foo {}\n",
	},
	{
		name: "Enum with comments inside.",
		in:   "enum Foo{\n\n\n\n\n\t\t\t\t\n// a comment\n\n\n\n\n}",
		out:  "enum Foo {\n  // a comment\n}\n",
	},
	{
		name: "Enum with options",
		in:   "enum Foo {\n\n\n\n\n\t\t\t\t\noption (foo) = 1;\n\n\noption this=\"that\";\n\n}",
		out:  "enum Foo {\n  option (foo) = 1;\n  option this = \"that\";\n}\n",
	},
	{
		name: "Enum with options and comments",
		in:   "enum Foo {\n\n\n\n\n\t\t\t\t\n// a comment\noption (foo) = 1;\n\n\noption this=\"that\";\n\n}",
		out:  "enum Foo {\n  // a comment\n  option (foo) = 1;\n  option this = \"that\";\n}\n",
	},
	// Enum Fields.
	{
		name: "Single enum field",
		in:   "enum Foo { BAR = 1; }",
		out:  "enum Foo {\n  BAR = 1;\n}\n",
	},
	{
		name: "Enum field with inline comment",
		in:   "enum Foo { BAR = 1; // a comment\n}",
		out:  "enum Foo {\n  BAR = 1; // a comment\n}\n",
	},
	{
		name: "Enum field with attached comment",
		in:   "enum Foo {\n  // comment about bar\n  BAR = 1; }",
		out:  "enum Foo {\n  // comment about bar\n  BAR = 1;\n}\n",
	},
	{
		name: "Enum field with a single option",
		in:   "enum Foo { BAR = 1 [foo = 1]; }",
		out:  "enum Foo {\n  BAR = 1 [foo = 1];\n}\n",
	},
	{
		name: "Enum field with a single custom option",
		in:   "enum Foo { BAR = 1 [(fool) = 2]; }",
		out:  "enum Foo {\n  BAR = 1 [(fool) = 2];\n}\n",
	},
	{
		name: "Enum field with a single custom option 2",
		in:   "enum Foo { BAR = 1 [(fool).bar = 2]; }",
		out:  "enum Foo {\n  BAR = 1 [(fool).bar = 2];\n}\n",
	},
	{
		name: "Enum field with multiple options",
		in:   "enum Foo { BAR = 1 [foo = 1, bar = 2]; }",
		out:  "enum Foo {\n  BAR = 1 [\n    foo = 1,\n    bar = 2\n  ];\n}\n",
	},
	{
		name: "Enum fields with comments and single options",
		in:   "enum Foo { BAR = 1 [foo = 1]; // a comment\n  BAZ = 2 [bar = 3];\n}",
		out:  "enum Foo {\n  BAR = 1 [foo = 1]; // a comment\n  BAZ = 2 [bar = 3];\n}\n",
	},
	// Services.
	{
		name: "Single service",
		in:   "service Foo {}",
		out:  "service Foo {}\n",
	},
	{
		name: "Service with a comment",
		in:   "// a service comment\nservice Foo {}",
		out:  "// a service comment\nservice Foo {}\n",
	},
	{
		name: "Service followed by a Service",
		in:   "service Foo {}\n\n\n\nservice Bar {}",
		out:  "service Foo {}\n\nservice Bar {}\n",
	},
	{
		name: "Service with options",
		in:   "service Foo {\n\n\n\noption (foo) = 1;\n\n\noption this=\"that\";\n\n}",
		out:  "service Foo {\n  option (foo) = 1;\n  option this = \"that\";\n}\n",
	},
	// Rpcs
	{
		name: "Single rpc",
		in:   "service Foo {\nrpc Foo(Bar) returns(Baz) {}\n}",
		out:  "service Foo {\n  rpc Foo(Bar) returns (Baz);\n}\n",
	},
	{
		name: "Rpc with streaming input type",
		in:   "service Foo {\nrpc Foo(stream Bar) returns(Baz) {}\n}",
		out:  "service Foo {\n  rpc Foo(stream Bar) returns (Baz);\n}\n",
	},
	{
		name: "Rpc with streaming output type",
		in:   "service Foo {\nrpc Foo(Bar) returns(stream Baz) {}\n}",
		out:  "service Foo {\n  rpc Foo(Bar) returns (stream Baz);\n}\n",
	},
	{
		name: "Rpc with streaming input and output type",
		in:   "service Foo {\nrpc Foo(stream Bar) returns(stream Baz) {}\n}",
		out:  "service Foo {\n  rpc Foo(stream Bar) returns (stream Baz);\n}\n",
	},
	{
		name: "Rpc with a comment",
		in:   "service Foo {\n// a rpc comment\nrpc Foo(Bar) returns(Baz) {}\n}",
		out:  "service Foo {\n  // a rpc comment\n  rpc Foo(Bar) returns (Baz);\n}\n",
	},
	{
		name: "Rpc with a comment and options",
		in:   "service Foo {\n// a rpc comment\nrpc Foo(Bar) returns(Baz) { option (foo) = 1; }}",
		out:  "service Foo {\n  // a rpc comment\n  rpc Foo(Bar) returns (Baz) {\n    option (foo) = 1;\n  }\n}\n",
	},
	{
		name: "A list of options followed by rpcs",
		in:   "service Foo {\noption (foo) = 1;\noption (bar) = 2;\nrpc Foo(Bar) returns(Baz) {}\n rpc Bar(Baz) returns(Bar) {}\n}",
		out:  "service Foo {\n  option (foo) = 1;\n  option (bar) = 2;\n\n  rpc Foo(Bar) returns (Baz);\n  rpc Bar(Baz) returns (Bar);\n}\n",
	},
	// Messages.
	{
		name: "Single message",
		in:   "message Foo {}",
		out:  "message Foo {}\n",
	},
	{
		name: "Message with a comment",
		in:   "// a message comment\nmessage Foo {}",
		out:  "// a message comment\nmessage Foo {}\n",
	},
	{
		name: "Message with a comment and options",
		in:   "// a message comment\nmessage Foo { option (foo) = 1; }",
		out:  "// a message comment\nmessage Foo {\n  option (foo) = 1;\n}\n",
	},
	{
		name: "Extend message.",
		in:   "extend Foo { int32 bar = 1; }",
		out:  "extend Foo {\n  int32 bar = 1;\n}\n",
	},
	{
		name: "Message within a message",
		in:   "message Foo { message Bar { string baz = 1; } }",
		out:  "message Foo {\n  message Bar {\n    string baz = 1;\n  }\n}\n",
	},
	{
		name: "Message within a message with comments",
		in:   "message Foo { // a message comment\n  message Bar { // a nested message comment\n    string baz = 1; // a field comment\n  }\n}",
		out:  "message Foo {\n  // a message comment\n  message Bar {\n    // a nested message comment\n    string baz = 1; // a field comment\n  }\n}\n",
	},
	// Message fields.
	{
		name: "Single message field",
		in:   "message Foo { int32 bar = 1; }",
		out:  "message Foo {\n  int32 bar = 1;\n}\n",
	},
	{
		name: "Simple repeated field",
		in:   "message Foo { repeated int32 bar = 1; }",
		out:  "message Foo {\n  repeated int32 bar = 1;\n}\n",
	},
	{
		name: "Simple optional field",
		in:   "message Foo { optional int32 bar = 1; }",
		out:  "message Foo {\n  optional int32 bar = 1;\n}\n",
	},
	{
		name: "Simple required field",
		in:   "message Foo { required int32 bar = 1; }",
		out:  "message Foo {\n  required int32 bar = 1;\n}\n",
	},
	{
		name: "Message field with a comment",
		in:   "message Foo { // a field comment\n  int32 bar = 1; }",
		out:  "message Foo {\n  // a field comment\n  int32 bar = 1;\n}\n",
	},
	{
		name: "Message field with a comment and options",
		in:   "message Foo { // a field comment\n  int32 bar = 1 [foo = 1]; }",
		out:  "message Foo {\n  // a field comment\n  int32 bar = 1 [foo = 1];\n}\n",
	},
	{
		name: "Messsage with multiple fields",
		in:   "message Foo { int32 bar = 1; int32 baz = 2; }",
		out:  "message Foo {\n  int32 bar = 1;\n  int32 baz = 2;\n}\n",
	},
	{
		name: "Message field with multiple options",
		in:   "message Foo { int32 bar = 1 [foo = 1, bar = 2]; }",
		out:  "message Foo {\n  int32 bar = 1 [\n    foo = 1,\n    bar = 2\n  ];\n}\n",
	},
	{
		name: "Message field with map constant",
		in:   "message Foo { int32 bar = 1 [baz = {1: 2}]; }",
		out:  "message Foo {\n  int32 bar = 1 [baz = {\n    1: 2\n  }];\n}\n",
	},
	// Oneof statement
	{
		name: "Simple empty oneof",
		in:   "message Foo { oneof bar {} }",
		out:  "message Foo {\n  oneof bar {}\n}\n",
	},
	{
		name: "Oneof statement",
		in:   "message Foo { oneof bar { int32 a = 1; string b = 2; } }",
		out:  "message Foo {\n  oneof bar {\n    int32 a = 1;\n    string b = 2;\n  }\n}\n",
	},
	{
		name: "Oneof statement with a comment",
		in:   "message Foo { // a oneof comment\noneof bar { int32 a = 1; string b = 2; } }",
		out:  "message Foo {\n  // a oneof comment\n  oneof bar {\n    int32 a = 1;\n    string b = 2;\n  }\n}\n",
	},
	{
		name: "Oneof statement with a comment and options",
		in:   "message Foo { // a oneof comment\noneof bar { option this = 293.292; } }",
		out:  "message Foo {\n  // a oneof comment\n  oneof bar {\n    option this = 293.292;\n  }\n}\n",
	},
	// Oneof field
	{
		name: "Oneof field",
		in:   "message Foo { oneof bar { int32 a = 1; } }",
		out:  "message Foo {\n  oneof bar {\n    int32 a = 1;\n  }\n}\n",
	},
	{
		name: "Oneof field with a comment",
		in:   "message Foo { oneof bar { // a field comment\n  int32 a = 1; } }",
		out:  "message Foo {\n  oneof bar {\n    // a field comment\n    int32 a = 1;\n  }\n}\n",
	},
	{
		name: "Oneof field with a comment and options",
		in:   "message Foo { oneof bar { // a field comment\n  int32 a = 1 [foo = 1]; } }",
		out:  "message Foo {\n  oneof bar {\n    // a field comment\n    int32 a = 1 [foo = 1];\n  }\n}\n",
	},
	{
		name: "Oneof field with multiple options",
		in:   "message Foo { oneof bar { int32 a = 1 [foo = 1, bar = 2]; } }",
		out:  "message Foo {\n  oneof bar {\n    int32 a = 1 [\n      foo = 1,\n      bar = 2\n    ];\n  }\n}\n",
	},
	// Reserved.
	{
		name: "Simple reserved field",
		in:   "message Foo {\n  reserved 1, 2, 3;\n}\n",
		out:  "message Foo {\n  reserved 1, 2, 3;\n}\n",
	},
	{
		name: "Simple reserved with field names.",
		in:   "message Foo {\n  reserved \"foo\", \"bar\", \"baz\";\n}\n",
		out:  "message Foo {\n  reserved \"foo\", \"bar\", \"baz\";\n}\n",
	},
	{
		// Unfortunately, not Documented, so no visitation of comments.
		name: "Reserved with a preceeding comment",
		in:   "message Foo {\n  // a comment\n  reserved 1, 2, 3;\n}\n",
		out:  "message Foo {\n  reserved 1, 2, 3;\n}\n",
	},
	// Map field
	{
		name: "Simple map field",
		in:   "message Foo {\n  map<string, int32> bar = 1;\n}\n",
		out:  "message Foo {\n  map<string, int32> bar = 1;\n}\n",
	},
	{
		name: "Map field with single option",
		in:   "message Foo {\n  map<string, int32> bar = 1 [default = 2];\n}\n",
		out:  "message Foo {\n  map<string, int32> bar = 1 [default = 2];\n}\n",
	},
	{
		name: "Map field with options",
		in:   "message Foo {\n  map<string, int32> bar = 1 [foo = 1, bar = 2];\n}\n",
		out:  "message Foo {\n  map<string, int32> bar = 1 [\n    foo = 1,\n    bar = 2\n  ];\n}\n",
	},
}

func TestFmts(t *testing.T) {
	for _, tt := range formattingTests {
		n, err := parseStringProto(tt.in)
		require.NoError(t, err, tt.name)
		out := testPrinter(n)
		if tt.printFmt {
			t.Log(tt.name, "\n"+out)
		}
		require.Equal(t, tt.out, out, "%s expected %q, got %q", tt.name, tt.out, out)

		// Check that whatever we got still parses:
		_, err = parseStringProto(out)
		require.NoError(t, err, tt.name, "\n"+out)

	}
}

// ============ Misc Old =======================

func TestCommentedFmt(t *testing.T) {
	t.Skip("Run manually by commenting this out.")
	pf := &proto.Proto{
		Elements: []proto.Visitee{},
		Filename: "foo.proto",
	}

	imp := &proto.Message{
		Name: "Bar",
		Comment: &proto.Comment{
			Lines: strings.Split(" this is a comment", "\n"),
		},
	}

	pf.Elements = append(pf.Elements, imp)
	output := new(strings.Builder)
	f := Formatter{w: output, indentSeparator: "  "}
	f.Format(pf)

	fmt.Println(output.String())
}

func Test_dump_out(t *testing.T) {
	t.Skip("Run manually by commenting this out.")
	pf, err := parseStringProto(`//syntax comment
	syntax = "proto3";
	
	// this comment
	import "this.proto";
	// that comment
	import "that.proto";
	// package comment
	package foo;
	`)
	if err != nil {
		fmt.Println(err)
		return
	}

	output := new(strings.Builder)
	f := Formatter{w: output, indentSeparator: "  "}
	f.Format(pf)

	fmt.Println(output.String())
}
