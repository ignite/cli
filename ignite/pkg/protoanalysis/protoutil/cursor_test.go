package protoutil

import (
	"testing"

	"github.com/emicklei/proto"
	"github.com/stretchr/testify/require"
)

// Make a simple replacement of package -> import.
func TestSimpleReplacement(t *testing.T) {
	f, err := parseStringProto(`package "package"`)
	require.NoError(t, err)
	Apply(f, nil, func(c *Cursor) bool {
		n := c.Node()
		if _, ok := n.(*proto.Package); ok {
			imp := NewImport("that")
			c.Replace(imp)
		}

		return true
	})
	require.True(t, containsElement(f, NewImport("that")))
	require.False(t, containsElement(f, NewPackage("package")))
}

func TestSimpleInsertAfter(t *testing.T) {
	f, err := parseStringProto(`syntax = "proto3"
	
	message Hello {
		message World {}
	}
	`)
	require.NoError(t, err)

	// keep ref for checking containment.
	var msg *proto.Message
	Apply(f, nil, func(c *Cursor) bool {
		n := c.Node()
		if n, ok := n.(*proto.Message); ok {
			if n.Name == "World" {
				msg = NewMessage("WeComeInPeace")
				c.InsertAfter(msg)
			}
		}
		return true
	})
	require.True(t, containsElement(f, msg))
	// check that it is inserted after "World"
	Apply(f, nil, func(c *Cursor) bool {
		n := c.Node()
		if n, ok := n.(*proto.Message); ok {
			if n.Name == "World" {
				next, ok := c.Next()
				require.True(t, ok)
				require.True(t, next.(*proto.Message).Name == "WeComeInPeace")
			}
		}
		return true
	})
}

// Can really only panic with comments since
// other elements in nodes aren't Visitees
func TestInsertAfterPanic(t *testing.T) {
	f, err := parseStringProto(`syntax = "proto3"
	
	// my import
	import "this";
	`)
	require.NoError(t, err)

	// Try calling insertAfter when c is a Comment
	require.Panics(t, func() {
		Apply(f, nil, func(c *Cursor) bool {
			n := c.Node()
			if _, ok := n.(*proto.Comment); ok {
				c.InsertAfter(NewImport("that"))
			}
			return true
		})
	})
}

func TestSimpleInsertBefore(t *testing.T) {
	f, err := parseStringProto(`syntax = "proto3"
	
	message Say {}
	message World {}
	`)
	require.NoError(t, err)

	// keep ref for checking containment.
	var msg *proto.Message
	Apply(f, nil, func(c *Cursor) bool {
		n := c.Node()
		if n, ok := n.(*proto.Message); ok {
			if n.Name == "World" {
				// add hello between say and world
				msg = NewMessage("Hello")
				c.InsertBefore(msg)
			}
		}
		return true
	})
	require.True(t, containsElement(f, msg))

	// check that it is inserted after "Say"
	Apply(f, nil, func(c *Cursor) bool {
		n := c.Node()
		if n, ok := n.(*proto.Message); ok {
			if n.Name == "Say" {
				next, ok := c.Next()
				require.True(t, ok)
				require.True(t, next.(*proto.Message).Name == "Hello")
			}
		}
		return true
	})
}

// Can really only panic with comments since
// other elements in nodes aren't Visitees
func TestInsertBeforePanic(t *testing.T) {
	f, err := parseStringProto(`syntax = "proto3"
	
	// my import
	import "this";
	`)
	require.NoError(t, err)

	// Try calling insertAfter when c is a Comment
	require.Panics(t, func() {
		Apply(f, nil, func(c *Cursor) bool {
			n := c.Node()
			if _, ok := n.(*proto.Comment); ok {
				c.InsertBefore(NewImport("that"))
			}
			return true
		})
	})
}

// Build a skeleton of a file by continuous appends on the file.
func TestAppendFile(t *testing.T) {
	f, err := parseStringProto(`syntax = "proto3"`)
	require.NoError(t, err)

	i := NewImport("importpath")
	Append(f, i)
	require.True(t, containsElement(f, i))

	p := NewPackage("package")
	Append(f, p)
	require.True(t, containsElement(f, p))

	o := NewOption("this", "that")
	Append(f, o)
	require.True(t, containsElement(f, o))

	oneof_f := NewOneofField("this", "string", 2)
	// Can directly append an option if required:
	opt := NewOption("this", "that")
	Append(oneof_f, opt)
	require.True(t, containsElement(oneof_f, opt))

	oneof := NewOneof("myoneof")
	Append(oneof, oneof_f)
	require.True(t, containsElement(oneof, oneof_f))

	normalfield := NewField("that", "string", 3)

	m := NewMessage("Hello")
	Append(m, oneof)
	require.True(t, containsElement(m, oneof))
	Append(m, normalfield)
	require.True(t, containsElement(m, normalfield))

	Append(f, m)
	require.True(t, containsElement(f, m))

	// Append an empty service
	s := NewService("Hey")
	Append(f, s)
	require.True(t, containsElement(f, s))

	// An empty enum
	e := NewEnum("Hey")
	// Add an enum field to it:
	e_f := NewEnumField("HEY", 1)
	Append(e, e_f)
	require.True(t, containsElement(e, e_f))

	Append(f, e)
	require.True(t, containsElement(f, e))
}

// Append to a node w/o elements panics.
func TestAppendEdges(t *testing.T) {
	f, err := parseStringProto(`syntax = "proto3"`)
	require.NoError(t, err)

	// Can't append to a Syntax node, panic.
	require.Panics(t, func() {
		Apply(f, nil, func(c *Cursor) bool {
			n := c.Node()
			if n, ok := n.(*proto.Syntax); ok {
				Append(n, NewImport("that"))
			}
			return true
		})
	})

	// Empty append does nothing.
	elems := len(f.Elements)
	Append(f)
	require.True(t, len(f.Elements) == elems)

	// Appending a non-option to NormalField/OneOfField panics.
	require.Panics(t, func() {
		f := NewField("that", "string", 3)
		Append(f, NewImport("that"))
	})
}

func TestCursorOps(t *testing.T) {
	f, err := parseStringProto(`syntax = "proto3"
	
	message Hello {}
	message World {
		message Hey {}
		enum E {}
	}
	`)
	require.NoError(t, err)

	Apply(f, nil, func(c *Cursor) bool {
		n := c.Node()
		if n, ok := n.(*proto.Message); ok {
			if n.Name == "Hello" {
				require.False(t, c.IsLast())
				n, ok := c.Next()
				require.True(t, ok)
				require.NotNil(t, n)

				parent, ok := c.Parent().(*proto.Proto)
				require.True(t, ok)
				require.True(t, parent.Filename == "")
				// currently useless.
				require.True(t, c.Name() == "Elements")
			}
			if n.Name == "World" {
				require.True(t, c.IsLast())
				n, ok := c.Next()
				require.False(t, ok)
				require.Nil(t, n)

				parent, ok := c.Parent().(*proto.Proto)
				require.True(t, ok)
				require.True(t, parent.Filename == "")
				// currently useless.
				require.True(t, c.Name() == "Elements")
			}

			if n.Name == "Hey" {
				require.False(t, c.IsLast())
				n, ok := c.Next()
				require.True(t, ok)
				require.NotNil(t, n)

				// parent is the message
				parent, ok := c.Parent().(*proto.Message)
				require.True(t, ok)
				require.True(t, parent.Name == "World")
				// currently useless.
				require.True(t, c.Name() == "Elements")
			}
		}

		if _, ok := n.(*proto.Enum); ok {
			require.True(t, c.IsLast())
			n, ok := c.Next()
			require.False(t, ok)
			require.Nil(t, n)

			// parent is the message
			parent, ok := c.Parent().(*proto.Message)
			require.True(t, ok)
			require.True(t, parent.Name == "World")
			// currently useless.
			require.True(t, c.Name() == "Elements")
		}

		// Don't make sense for elements not contained in a slice (currently
		// proto.Proto or comments)
		if _, ok := n.(*proto.Proto); ok {
			require.Panics(t, func() { c.IsLast() })
			require.Panics(t, func() { c.Next() })
		}
		return true
	})
}

// Also test the utilities here.

func TestAddImports(t *testing.T) {
	f, err := parseStringProto(`syntax = "proto3"`)
	require.NoError(t, err)

	// Add an import
	err = AddImports(f, true, NewImport("this.proto"))
	require.NoError(t, err)
	require.True(t, HasImport(f, "this.proto"))
	// Note: added in reverse order.
	err = AddImports(f, true,
		NewImport("that.proto"),
		NewImport("the.other.proto"),
		NewImport("and.another.proto"),
	)
	require.NoError(t, err)
	require.True(t, HasImport(f, "that.proto"))
	require.True(t, HasImport(f, "the.other.proto"))
	require.True(t, HasImport(f, "and.another.proto"))

	// Empty import is no-op.
	require.NoError(t, AddImports(f, true))
	// Importing on empty file is currently an error.
	require.Error(t, AddImports(
		&proto.Proto{},
		true,
		NewImport("this.proto"),
	))

	// Exercise the recursive case:
	f, err = parseStringProto(`syntax = "proto3"`)
	require.NoError(t, err)
	err = AddImports(f, true,
		NewImport("this.proto"),
		NewImport("that.proto"),
	)
	require.NoError(t, err)
	require.True(t, HasImport(f, "this.proto"))
	require.True(t, HasImport(f, "that.proto"))

	f, err = parseStringProto(`syntax = "proto3";
package cosmonaut.chainname.chainname;

import "gogoproto/gogo.proto";
import "google/api/annotations.proto";
import "cosmos/base/query/v1beta1/pagination.proto";
import "chainname/params.proto";
`)
	require.NoError(t, err)
	err = AddImports(f, true, NewImport("chainname/bleep.proto"))
	require.NoError(t, err)
	// Add dupes:
	err = AddImports(f, true, NewImport("chainname/bleep.proto"))
	require.NoError(t, err)
	err = AddImports(f, true, NewImport("chainname/bleep.proto"))
	require.NoError(t, err)
	err = AddImports(f, true, NewImport("chainname/params.proto"))
	require.NoError(t, err)
	// just checking that is added last.
	// fmt.Print(Printer(f))

	// Check that adding duplicates does nothing.
	f, err = parseStringProto(`syntax = "proto3";
package cosmonaut.chainname.chainname;

import "gogoproto/gogo.proto";
import "google/api/annotations.proto";
import "cosmos/base/query/v1beta1/pagination.proto";
import "chainname/params.proto";
`)
	require.NoError(t, err)
	imports := []*proto.Import{
		NewImport("chainname/params.proto"),
		NewImport("gogoproto/gogo.proto"),
	}
	err = AddImports(f, true, imports...)
	require.NoError(t, err)
	require.Equal(t, len(f.Elements), 6, "The number of elements shouldn't have changed")

	// Pass an empty import list.
	f, err = parseStringProto(`syntax = "proto3";`)
	require.NoError(t, err)
	ret := AddImports(f, true)
	require.Nil(t, ret)

	// No imports, no fallback.
	f, err = parseStringProto(`syntax = "proto3";`)
	require.NoError(t, err)
	err = AddImports(f, false, NewImport("this.proto"))
	require.Error(t, err)
}

func TestHasImport(t *testing.T) {
	f, err := parseStringProto(`syntax = "proto3"

	import "this.proto";
	import "that.proto";
	import "the.other.proto"
	`)
	require.NoError(t, err)
	require.True(t, HasImport(f, "this.proto"))
	require.True(t, HasImport(f, "that.proto"))
	require.True(t, HasImport(f, "the.other.proto"))
	require.False(t, HasImport(f, "this.proto.proto"))
}

func TestGetMessage(t *testing.T) {
	f, err := parseStringProto(`syntax = "proto3"
	
	message Hello {
		message World {
			message WeComeInPeace {
				message TheAnswerToLifeTheUniverseAndEverything {
					message IsActuallyFortyTwo {}
				}
			}
		}
	}
	`)
	require.NoError(t, err)
	m, err := GetMessageByName(f, "Hello")
	require.NoError(t, err)
	require.Equal(t, "Hello", m.Name)

	m, err = GetMessageByName(f, "World")
	require.NoError(t, err)
	require.Equal(t, "World", m.Name)

	m, err = GetMessageByName(f, "WeComeInPeace")
	require.NoError(t, err)
	require.Equal(t, "WeComeInPeace", m.Name)

	m, err = GetMessageByName(f, "TheAnswerToLifeTheUniverseAndEverything")
	require.NoError(t, err)
	require.Equal(t, "TheAnswerToLifeTheUniverseAndEverything", m.Name)

	m, err = GetMessageByName(f, "IsActuallyFortyTwo")
	require.NoError(t, err)
	require.Equal(t, "IsActuallyFortyTwo", m.Name)

	_, err = GetMessageByName(f, "DoesNotExist")
	require.Error(t, err)
}

func TestHasMessage(t *testing.T) {
	f, err := parseStringProto(`syntax = "proto3"
	
	message Hello {
		message World {
			message WeComeInPeace {
				message TheAnswerToLifeTheUniverseAndEverything {
					message IsActuallyFortyTwo {}
				}
			}
		}
	}
	`)
	require.NoError(t, err)
	require.True(t, HasMessage(f, "Hello"))
	require.True(t, HasMessage(f, "World"))
	require.True(t, HasMessage(f, "WeComeInPeace"))
	require.True(t, HasMessage(f, "TheAnswerToLifeTheUniverseAndEverything"))
	require.True(t, HasMessage(f, "IsActuallyFortyTwo"))
	require.False(t, HasMessage(f, "DoesNotExist"))
	require.False(t, HasMessage(f, "Hello.World"))
}

func TestGetService(t *testing.T) {
	f, err := parseStringProto(`syntax = "proto3"
	
	service Msg {
	}
	service AnotherMsg {}
	service YetAnotherMsg {
		rpc Foo(Bar) returns (Bar) {}
	}
	`)
	require.NoError(t, err)
	s, err := GetServiceByName(f, "Msg")
	require.NoError(t, err)
	require.Equal(t, "Msg", s.Name)

	s, err = GetServiceByName(f, "AnotherMsg")
	require.NoError(t, err)
	require.Equal(t, "AnotherMsg", s.Name)

	s, err = GetServiceByName(f, "YetAnotherMsg")
	require.NoError(t, err)
	require.Equal(t, "YetAnotherMsg", s.Name)

	_, err = GetServiceByName(f, "DoesNotExist")
	require.Error(t, err)
}

func TestHasService(t *testing.T) {
	f, err := parseStringProto(`syntax = "proto3"

	service Msg {}
	service AnotherMsg {}
	service YetAnotherMsg {}
	`)
	require.NoError(t, err)
	require.True(t, HasService(f, "Msg"))
	require.True(t, HasService(f, "AnotherMsg"))
	require.True(t, HasService(f, "YetAnotherMsg"))
	require.False(t, HasService(f, "DoesNotExist"))
}

func TestGetNextId(t *testing.T) {
	f, err := parseStringProto(`syntax = "proto3"

	message Hello {
		string g = 1;
		message World {
			message WeComeInPeace {
				message TheAnswerToLifeTheUniverseAndEverything {
					message IsActuallyFortyTwo {
						string foo = 1;
						int32 bar = 2;
						int64 baz = 3;
					}
				}
			}
		}
	}
	`)
	require.NoError(t, err)

	m, err := GetMessageByName(f, "IsActuallyFortyTwo")
	require.NoError(t, err)
	require.Equal(t, 4, NextUniqueID(m))

	m, err = GetMessageByName(f, "Hello")
	require.NoError(t, err)
	require.Equal(t, 2, NextUniqueID(m))

	f, err = parseStringProto(`syntax = "proto3"

	message Hello {
		string g = 1;
		string foo = 2;
		int32 bar = 3;
		int64 baz = 5;
	}`)
	require.NoError(t, err)
	m, err = GetMessageByName(f, "Hello")
	require.NoError(t, err)
	require.Equal(t, 6, NextUniqueID(m))
}
