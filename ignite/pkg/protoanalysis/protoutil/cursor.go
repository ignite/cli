package protoutil

import (
	"fmt"
	"reflect"

	"github.com/emicklei/proto"
)

// Note: The traversing can also be done with proto.Walk but there's some reasons
// why I chose the cursor instead:
//
// 1. We can abort traversing deeper in the tree at any point. (a post
//    ApplyFunc returning false)
// 2. We keep track of the parent to have finer grained control over where we
//    are in the tree.
// 3. We can use pre/post handling.

// Modeled heavily after the Apply/Cursor logic in astutil, using proto.Visitee as
// the common interface, abilities for reflection aren't as rich but can still
// manage in order to get the job done.
// Cursor has been augmented to add a couple more methods that can make
// our life easier.

// An ApplyFunc is invoked by Apply for each Visitee n, even if n is nil,
// before and/or after the node's children, using a Cursor describing
// the current node and providing operations on it.
//
// The return value of ApplyFunc controls the syntax tree traversal.
// See Apply for details.
type ApplyFunc func(*Cursor) bool

// Apply traverses a syntax tree recursively, starting with root,
// and calling pre and post for each node as described below.
// Apply returns the syntax tree, possibly modified.
//
// If pre is not nil, it is called for each node before the node's
// children are traversed (pre-order). If pre returns false, no
// children are traversed, and post is not called for that node.
//
// If post is not nil, and a prior call of pre didn't return false,
// post is called for each node after its children are traversed
// (post-order). If post returns false, traversal is terminated and
// Apply returns immediately.
func Apply(root proto.Visitee, pre, post ApplyFunc) (result proto.Visitee) {
	parent := &struct{ proto.Visitee }{root}
	defer func() {
		if r := recover(); r != nil && r != abort {
			panic(r)
		}
		result = parent.Visitee
	}()
	a := &application{pre: pre, post: post}
	a.apply(parent, "Visitee", nil, root)
	return
}

var abort = new(int) // singleton, to signal termination of Apply

// A Cursor describes a node encountered during Apply.
// Information about the node and its parent is available
// from the Node, Parent, Name, and Index methods.
type Cursor struct {
	parent proto.Visitee // parent (containing a []proto.Visitee slice)
	name   string
	iter   *iterator
	node   proto.Visitee // current node we're applying over
}

type iterator struct {
	index, step int
}

// Index reports the index >= 0 of the current Visitee in the slice of Visitees that
// contains it, or a value < 0 if the current Visitee is not part of a slice.
// The index of the current node changes if InsertBefore is called while
// processing the current node.
func (c *Cursor) Index() int {
	if c.iter != nil {
		return c.iter.index
	}
	return -1
}

// field returns the current node's parent field value.
func (c *Cursor) field() reflect.Value {
	return reflect.Indirect(reflect.ValueOf(c.parent)).FieldByName(c.name)
}

// Node returns the current Node.
func (c *Cursor) Node() proto.Visitee { return c.node }

// Parent returns the parent of the current Node.
func (c *Cursor) Parent() proto.Visitee { return c.parent }

// Name returns the name of the parent Node field that contains the current Node.
// If the parent is a *ast.Package and the current Node is a *ast.File, Name returns
// the filename for the current Node.
func (c *Cursor) Name() string { return c.name }

// IsLast returns if the current node being traversed is the final node in the
// slice of nodes. Can be used to determine if a node is the last one.
func (c *Cursor) IsLast() bool {
	i := c.Index()
	if i < 0 {
		panic("IsLast node not contained in slice")
	}
	v := c.field()
	return i == v.Len()-1
}

// Next returns the next Visitee. Can be used to check the next value
// before deciding to continue.
func (c *Cursor) Next() (proto.Visitee, bool) {
	i := c.Index()
	if i < 0 {
		panic("Next node not contained in slice")
	}
	v := c.field()
	if i == v.Len()-1 {
		return nil, false
	}
	var x proto.Visitee
	if e := v.Index(i + 1); e.IsValid() {
		x = e.Interface().(proto.Visitee)
	}
	return x, true
}

// Replace replaces the current Node with n.
// The replacement node is not walked by Apply.
func (c *Cursor) Replace(n proto.Visitee) {
	v := c.field()
	if i := c.Index(); i >= 0 {
		v = v.Index(i)
	}
	v.Set(reflect.ValueOf(n))
}

// InsertAfter inserts n after the current Node in its containing slice.
// If the current Node is not part of a slice, InsertAfter panics.
// Apply does not walk n.
func (c *Cursor) InsertAfter(n proto.Visitee) {
	i := c.Index()
	if i < 0 {
		panic("InsertAfter node not contained in slice")
	}
	v := c.field()
	v.Set(reflect.Append(v, reflect.Zero(v.Type().Elem())))
	l := v.Len()
	reflect.Copy(v.Slice(i+2, l), v.Slice(i+1, l))
	v.Index(i + 1).Set(reflect.ValueOf(n))
	c.iter.step++
}

// InsertBefore inserts n before the current Node in its containing slice.
// If the current Node is not part of a slice, InsertBefore panics.
// Apply will not walk n.
func (c *Cursor) InsertBefore(n proto.Visitee) {
	i := c.Index()
	if i < 0 {
		panic("InsertBefore node not contained in slice")
	}
	v := c.field()
	v.Set(reflect.Append(v, reflect.Zero(v.Type().Elem())))
	l := v.Len()
	reflect.Copy(v.Slice(i+1, l), v.Slice(i, l))
	v.Index(i).Set(reflect.ValueOf(n))
	c.iter.index++
}

// application carries all the shared data, so we can pass it around cheaply.
type application struct {
	pre, post ApplyFunc
	cursor    Cursor
	iter      iterator
}

func (a *application) apply(parent proto.Visitee, name string, iter *iterator, n proto.Visitee) {
	// don't walk into nil's
	if v := reflect.ValueOf(n); v.Kind() == reflect.Ptr && v.IsNil() {
		return
	}

	// avoid heap-allocating a new cursor for each apply call; reuse a.cursor instead
	saved := a.cursor
	a.cursor.parent = parent
	a.cursor.name = name
	a.cursor.iter = iter
	a.cursor.node = n

	if a.pre != nil && !a.pre(&a.cursor) {
		a.cursor = saved
		return
	}

	// Walk the children.
	// This is the issue with proto. Structure isn't really here in order to be able to
	// visit every component using a distinct interface. They are all Visitee's.
	// Ideally, we could wrap the proto nodes into interfaces that enforce a structure,
	// i.e. Nodes, MessageNodes, ServiceNodes, ProtoNodes, etc.
	// this way, inserting into a slice would be guarded (by error-ing) by the type of the slice.
	//
	// An alternative would be to reflect in the insertion methods for cursor and only allow
	// specific elements per type.
	switch n := n.(type) {
	case *proto.Proto:
		a.applyList(n, "Elements")
	case *proto.Service:
		a.apply(n, "Comment", nil, n.Comment)
		a.applyList(n, "Elements")
	case *proto.RPC:
		a.apply(n, "Comment", nil, n.Comment)
		a.apply(n, "Inline Comment", nil, n.InlineComment)
		a.applyList(n, "Elements")
	case *proto.Message:
		a.apply(n, "Comment", nil, n.Comment)
		a.applyList(n, "Elements")
	case *proto.NormalField:
		a.apply(n, "Comment", nil, n.Comment)
		a.apply(n, "Inline Comment", nil, n.InlineComment)
		a.applyList(n, "Options")
	case *proto.Oneof:
		a.apply(n, "Comment", nil, n.Comment)
		a.applyList(n, "Elements")
	case *proto.OneOfField:
		a.apply(n, "Comment", nil, n.Comment)
		a.apply(n, "Inline Comment", nil, n.InlineComment)
		a.applyList(n, "Options")
	case *proto.Enum:
		a.apply(n, "Comment", nil, n.Comment)
		a.applyList(n, "Elements")
	case *proto.EnumField:
		a.apply(n, "Comment", nil, n.Comment)
		a.apply(n, "Inline Comment", nil, n.InlineComment)
		a.applyList(n, "Elements")
	case *proto.Import:
		a.apply(n, "Comment", nil, n.Comment)
		a.apply(n, "Inline Comment", nil, n.InlineComment)
	case *proto.Option:
		a.apply(n, "Comment", nil, n.Comment)
		a.apply(n, "Inline Comment", nil, n.InlineComment)
	case *proto.Package:
		a.apply(n, "Comment", nil, n.Comment)
		a.apply(n, "Inline Comment", nil, n.InlineComment)
	default:
		// Probably a comment, ignore it.
	}

	if a.post != nil && !a.post(&a.cursor) {
		panic(abort)
	}
	a.cursor = saved
}

// applyList calls apply on each of the elements of the Visitee.
func (a *application) applyList(parent proto.Visitee, name string) {
	// avoid heap-allocating a new iterator for each applyList call; reuse a.iter instead
	saved := a.iter
	a.iter.index = 0
	for {
		// must reload parent.name each time, since cursor modifications might change it
		v := reflect.Indirect(reflect.ValueOf(parent)).FieldByName(name)
		if a.iter.index >= v.Len() {
			break
		}

		// element x may be nil in a bad AST - be cautious
		var x proto.Visitee
		if e := v.Index(a.iter.index); e.IsValid() {
			x = e.Interface().(proto.Visitee)
		}

		// reset step on each iteration.
		a.iter.step = 1
		a.apply(parent, name, &a.iter, x)
		a.iter.index += a.iter.step
	}
	a.iter = saved
}

// Append appends the elements provided to the node `n`. `n` must be
// a node that can accept elements, such as a proto.File or a proto.Message.
//
// Append panics if `n` is not a node that can accept elements or if the type
// of the elements provided is not compatible with the type of the elements
// contained by `n`. (Basically, this applies to NormalFields and OneOfFields,).
func Append(n proto.Visitee, elems ...proto.Visitee) {
	// return early if the slice is empty.
	if len(elems) == 0 {
		return
	}
	switch n.(type) {
	case *proto.Proto, *proto.Message, *proto.Enum, *proto.Oneof,
		*proto.Service, *proto.EnumField, *proto.RPC:
		// Can just append directly.
		v := reflect.Indirect(reflect.ValueOf(n)).FieldByName("Elements")
		v.Set(reflect.AppendSlice(v, reflect.ValueOf(elems)))
	case *proto.NormalField, *proto.OneOfField:
		// Make into options, panic on failure of one of the objects to do
		// so.
		var elements []*proto.Option
		for _, e := range elems {
			o, ok := e.(*proto.Option)
			if !ok {
				panic(fmt.Sprintf("Tried to append %T to a slice of Options", e))
			}
			elements = append(elements, o)
		}
		// append
		v := reflect.Indirect(reflect.ValueOf(n)).FieldByName("Options")
		v.Set(reflect.AppendSlice(v, reflect.ValueOf(elements)))
		return
	default:
		panic("Append: node not a slice")
	}
}
