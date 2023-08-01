package protoutil

import (
	"errors"
	"fmt"

	"github.com/emicklei/proto"
)

// AddAfterSyntax tries to add the given Visitee after the 'syntax' statement.
// If no syntax statement is found, returns an error.
func AddAfterSyntax(f *proto.Proto, v proto.Visitee) error {
	// return false to immediately stop
	inserted := false
	Apply(f, nil, func(c *Cursor) bool {
		if _, ok := c.Node().(*proto.Syntax); ok {
			c.InsertAfter(v)
			inserted = true
			return false
		}
		// Continue until we insert.
		return true
	})
	if inserted {
		return nil
	}
	return errors.New("could not find syntax statement")
}

// AddAfterPackage tries to add the given Visitee after the 'package' statement.
// If no package statement is found, returns an error.
func AddAfterPackage(f *proto.Proto, v proto.Visitee) error {
	inserted := false
	Apply(f, nil, func(c *Cursor) bool {
		if _, ok := c.Node().(*proto.Package); ok {
			c.InsertAfter(v)
			inserted = true
			return false
		}
		// Continue until we insert.
		return true
	})
	if inserted {
		return nil
	}
	return errors.New("could not find package statement")
}

// Fallback logic, try and use import after a package and if that fails
// attempts to use it after a syntax statement.
// If that fails, returns an error.
func importFallback(f *proto.Proto, imp *proto.Import) error {
	if err := AddAfterPackage(f, imp); err != nil {
		if err = AddAfterSyntax(f, imp); err != nil {
			return err
		}
	}
	return nil
}

// AddImports attempts to add the given import *after* any other imports
// in the file.
//
// If fallback is supplied, attempts to add it after the 'package'
// statement and then the 'syntax' statement are made.
//
// If none of the attempts are successful, returns an error.
func AddImports(f *proto.Proto, fallback bool, imports ...*proto.Import) (err error) {
	// No effect.
	if len(imports) == 0 {
		return nil
	}
	importMap, inserted := make(map[string]*proto.Import), false
	for _, i := range imports {
		importMap[i.Filename] = i
	}

	Apply(f, nil, func(c *Cursor) bool {
		if i, ok := c.Node().(*proto.Import); ok {
			delete(importMap, i.Filename)
			if next, ok := c.Next(); ok {
				if _, ok := next.(*proto.Import); ok {
					return true
				}
				for _, imp := range importMap {
					c.InsertAfter(imp)
				}
				inserted = true
				return false
			}
			// We're at the end (no Next())
			for _, imp := range importMap {
				c.InsertAfter(imp)
			}
			inserted = true
			return false
		}
		return true
	})
	// return if inserted.
	if inserted {
		return nil
	}
	// else fallback if defined.
	if fallback {
		// if the number of imports is > 1, we can try and insert the first after
		// the package/syntax and then recurse into AddImport with the rest (which we'll)
		// know that we can insert after an import since we just added it.
		imports = []*proto.Import{}
		for _, imp := range importMap {
			imports = append(imports, imp)
		}
		if len(imports) == 0 {
			return nil
		}
		if err := importFallback(f, imports[0]); err != nil {
			return err
		}
		// recurse with the rest. (might be empty)
		return AddImports(f, false, imports[1:]...)
	}
	return errors.New("unable to add import, no import statements found")
}

// NextUniqueID goes through the fields of the given Message and returns
// an id > max(fieldIds). It does not try to 'plug the holes' by selecting the
// least available id.
//
//	 // In 'example.proto' file
//	 syntax = "proto3"
//
//		message Hello {
//			string g = 1;
//			string foo = 2;
//			int32 bar = 3;
//			int64 baz = 5;
//		}
//	 f := ParseProtoPath("example.proto")
//	 m := GetMessageByName(f, "Hello")
//	 NextUniqueID(m) // 6
func NextUniqueID(m *proto.Message) int {
	// Best to recurse through elements directly here since
	// messages can embed other messages and the Apply could get
	// hairy.
	// if no elements exist => 1.
	max := 0
	for _, el := range m.Elements {
		if f, ok := el.(*proto.NormalField); ok {
			if f.Sequence > max {
				max = f.Sequence
			}
		}
	}
	return max + 1
}

// GetMessageByName returns the message with the given name or nil if not found.
// Only traverses in proto.Proto and proto.Message since they are the only nodes
// that contain messages:
//
//	f, _ := ParseProtoPath("foo.proto")
//	m := GetMessageByName(f, "Foo")
//	m.Name // "Foo"
func GetMessageByName(f *proto.Proto, name string) (node *proto.Message, err error) {
	node, err = nil, nil
	found := false
	Apply(f,
		func(c *Cursor) bool {
			if m, ok := c.Node().(*proto.Message); ok {
				if m.Name == name {
					found = true
					node = m
					return false
				}
				// keep looking if we're in a Message
				return true
			}
			// keep looking while we're in a proto.Proto.
			_, ok := c.Node().(*proto.Proto)
			return ok
		},
		// return immediately iff found.
		func(c *Cursor) bool { return !found })
	if found {
		return
	}
	return nil, fmt.Errorf("message %s not found", name)
}

// GetServiceByName returns the service with the given name or nil if not found.
// Only traverses in proto.Proto since it is the only node that contain services:
//
//	f, _ := ParseProtoPath("foo.proto")
//	s := GetServiceByName(f, "FooSrv")
//	s.Name // "FooSrv"
func GetServiceByName(f *proto.Proto, name string) (node *proto.Service, err error) {
	node, err = nil, nil
	found := false
	Apply(f,
		func(c *Cursor) bool {
			if s, ok := c.Node().(*proto.Service); ok {
				if s.Name == name {
					found = true
					node = s
				}
				// No nested services
				return false
			}
			// keep looking while we're in a proto.Proto.
			_, ok := c.Node().(*proto.Proto)
			return ok
		},
		// return immediately iff found.
		func(c *Cursor) bool { return !found })
	if found {
		return
	}
	return nil, fmt.Errorf("service %s not found", name)
}

// GetImportByPath returns the import with the given path or nil if not found.
// Only traverses in proto.Proto since it is the only node that contain imports:
//
//	f, _ := ParseProtoPath("foo.proto")
//	s := GetImportByPath(f, "other.proto")
//	s.FileName // "other.proto"
func GetImportByPath(f *proto.Proto, path string) (node *proto.Import, err error) {
	found := false
	node, err = nil, nil
	Apply(f,
		func(c *Cursor) bool {
			if i, ok := c.Node().(*proto.Import); ok {
				if i.Filename == path {
					found = true
					node = i
				}
				// No nested imports
				return false
			}
			// keep looking while we're in a proto.Proto.
			_, ok := c.Node().(*proto.Proto)
			return ok
		},
		// return immediately iff found.
		func(c *Cursor) bool { return !found })
	if found {
		return
	}
	return nil, fmt.Errorf("import %s not found", path)
}

// HasMessage returns true if the given message is found in the given file.
//
//	f, _ := ParseProtoPath("foo.proto")
//	// true if 'foo.proto' contains message Foo { ... }
//	r := HasMessage(f, "Foo")
func HasMessage(f *proto.Proto, name string) bool {
	_, err := GetMessageByName(f, name)
	return err == nil
}

// HasService returns true if the given service is found in the given file.
//
//	f, _ := ParseProtoPath("foo.proto")
//	// true if 'foo.proto' contains service FooSrv { ... }
//	r := HasService(f, "FooSrv")
func HasService(f *proto.Proto, name string) bool {
	_, err := GetServiceByName(f, name)
	return err == nil
}

// HasImport returns true if the given import (by path) is found in the given file.
//
//	f, _ := ParseProtoPath("foo.proto")
//	// true if 'foo.proto' contains import "path.to.other.proto"
//	r := HasImport(f, "path.to.other.proto")
func HasImport(f *proto.Proto, path string) bool {
	_, err := GetImportByPath(f, path)
	return err == nil
}
