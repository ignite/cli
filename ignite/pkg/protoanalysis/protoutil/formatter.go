package protoutil

// An alternative formatter that doesn't try to be as smart as the one
// currently used.

// NOTES:
//  - Doesn't handle C-style comments (prints as regular) or inline ones.
//  - Doesn't handle proto2 syntax.
//  - Can make it detect duplicate imports/double syntax/package and fields
//	  with weird sequencing.

import (
	"fmt"
	"io"
	"strings"

	"github.com/emicklei/proto"
)

const (
	srv     = "service"
	enum    = "enum"
	message = "message"
	oneof   = "oneof"
)

// Formatter visits a Proto and writes formatted source.
type Formatter struct {
	w                  io.Writer
	indentSeparator    string
	lastStmt, currStmt string
	indentLevel        int
}

// NewFormatter returns a new Formatter. Only the indentation separator is configurable.
func NewFormatter(writer io.Writer, indentSeparator string) *Formatter {
	return &Formatter{w: writer, indentSeparator: indentSeparator}
}

// Format visits all proto elements and writes formatted source.
func (f *Formatter) Format(p *proto.Proto) {
	f.lastStmt = "proto"
	for _, each := range p.Elements {
		// Visit all comments as you go.
		f.VisitCommentable(each)
		each.Accept(f)
	}
}

// VisitCommentable visits the comment preceding each element and returns true if
// one was visited. This is used to automatically print comments if they exist without
// needing to explicitly visit them in each Visit* function.
func (f *Formatter) VisitCommentable(c proto.Visitee) {
	// Note that Comment nodes are not Documented.
	if com, ok := c.(proto.Documented); ok {
		comment := com.Doc()
		if comment != nil {
			f.currStmt = currentNode(c)
			comment.Accept(f)
			f.lastStmt = "embedded-comment"
		}
	}
}

// VisitComment formats a Comment and writes a newline.
func (f *Formatter) VisitComment(c *proto.Comment) {
	// Note: don't need to guard against nil, we only visit if we
	// see it in the tree or if VisitCommentable is called.
	if f.lastStmt != f.currStmt {
		// containers + comment.
		f.allowedLastStatements(
			"proto", "comment", "enum", "message", "service", "rpc",
			"oneof",
		)
	}
	last := len(c.Lines) - 1
	for i, each := range c.Lines {
		if each != "" && each[0] == ' ' {
			f.printWithIndent(fmt.Sprintf("//%s", each))
		} else {
			f.printWithIndent(fmt.Sprintf("// %s", each))
		}
		if i != last {
			f.newline()
		}
	}
	f.newline()
}

func (f *Formatter) formatLiteral(l *proto.Literal) {
	if len(l.OrderedMap) == 0 && len(l.Array) == 0 {
		fmt.Fprintf(f.w, "%s", l.SourceRepresentation())
		return
	}
	fmt.Fprintf(f.w, "{\n")
	for _, other := range l.OrderedMap {
		f.indent()
		// yes, print with ident here.
		f.printWithIndent(other.Name)
		if other.PrintsColon {
			fmt.Fprintf(f.w, ": ")
		}
		f.formatLiteral(other.Literal)
		f.dedent()
		f.newline()
	}
	// and here.
	f.printWithIndent("}")
}

// VisitSyntax formats a Syntax.
func (f *Formatter) VisitSyntax(s *proto.Syntax) {
	// Syntax can't be indented, we don't care about indent in here.
	f.allowedLastStatements("proto", "comment")

	fmt.Fprintf(f.w, "syntax = %q", s.Value)
	f.endWithComment(s.InlineComment)

	f.lastStmt = "syntax"
	f.newline()
}

// VisitPackage formats a Package.
func (f *Formatter) VisitPackage(p *proto.Package) {
	// Packages can't be indented, we don't care about indent in here.
	f.allowedLastStatements("syntax", "proto")

	fmt.Fprintf(f.w, "package %s", p.Name)
	f.endWithComment(p.InlineComment)
	f.lastStmt = "package"
	f.newline()
}

// VisitImport formats a Import.
func (f *Formatter) VisitImport(i *proto.Import) {
	// Imports can't be indented, we don't care about indent in here.
	f.allowedLastStatements("import", "proto")

	kind := i.Kind
	if kind != "" {
		kind += " "
	}
	fmt.Fprintf(f.w, "import %s%q", kind, i.Filename)
	f.endWithComment(i.InlineComment)
	f.lastStmt = "import"
	f.newline()
}

// VisitOption formats a Option.
func (f *Formatter) VisitOption(o *proto.Option) {
	// try and keep options together.
	f.allowedLastStatements(
		"option", "normalfield", "proto",
		"rpc", "enum", "message", "service",
		"oneof",
	)

	f.printWithIndent(fmt.Sprintf("option %s = ", o.Name))
	f.formatLiteral(&o.Constant)
	f.endWithComment(o.InlineComment)
	f.lastStmt = "option"
	f.newline()
}

// VisitEnum formats a Enum.
func (f *Formatter) VisitEnum(e *proto.Enum) {
	// Always place a newline between top-level elements. (except comments!)
	f.allowedLastStatements("proto")

	// Technically, we're in it.
	f.lastStmt = enum
	f.printWithIndent(fmt.Sprintf("%s %s {", enum, e.Name))
	if len(e.Elements) > 0 {
		f.currStmt = enum
		f.newline()
		for _, each := range e.Elements {
			f.indent()
			// Don't forget to visit comment.
			f.VisitCommentable(each)
			each.Accept(f)
			f.dedent()
		}
	}
	io.WriteString(f.w, "}")
	f.lastStmt = enum
	f.newline()
}

// VisitEnumField formats a EnumField.
func (f *Formatter) VisitEnumField(ef *proto.EnumField) {
	f.allowedLastStatements("enumfield", "enum")
	f.printWithIndent(fmt.Sprintf("%s = %d", ef.Name, ef.Integer))

	numElements := len(ef.Elements)
	if numElements == 1 {
		if e, ok := ef.Elements[0].(*proto.Option); ok {
			io.WriteString(f.w, " [")
			f.VisitInlineOption(e)
			io.WriteString(f.w, "]")
		} else {
			ef.Elements[0].Accept(f) // Shouldn't occur! (Non option)
		}
	} else if numElements > 1 {
		// Format as an array.
		io.WriteString(f.w, " [")
		f.indent()
		for idx, each := range ef.Elements {
			io.WriteString(f.w, "\n")
			if e, ok := each.(*proto.Option); ok {
				f.printWithIndent("") // move f.w to right spot.
				f.VisitInlineOption(e)
				if idx < numElements-1 {
					io.WriteString(f.w, ",")
				}
			} else {
				each.Accept(f) // Shouldn't occur! (Non option)
			}
		}
		io.WriteString(f.w, "\n")
		f.dedent()
		f.printWithIndent("]")
	}
	f.endWithComment(ef.InlineComment)
	f.lastStmt = "enumfield"
	f.newline()
}

// VisitService formats a Service.
func (f *Formatter) VisitService(s *proto.Service) {
	// Services are never indented.
	f.allowedLastStatements("proto")

	// Technically, we're in it.
	f.lastStmt = srv
	fmt.Fprintf(f.w, "%s %s {", srv, s.Name)
	if len(s.Elements) > 0 {
		f.currStmt = srv
		f.newline()
		for _, each := range s.Elements {
			f.indent()
			// Don't forget to visit comment.
			f.VisitCommentable(each)
			each.Accept(f)
			f.dedent()
		}
	}
	io.WriteString(f.w, "}")
	f.lastStmt = srv
	f.newline()
}

// VisitRPC formats a RPC.
func (f *Formatter) VisitRPC(r *proto.RPC) {
	f.allowedLastStatements("rpc", "service")
	streamsReq, streamsResp := "", ""
	if r.StreamsRequest {
		streamsReq = "stream "
	}
	if r.StreamsReturns {
		streamsResp = "stream "
	}
	header := fmt.Sprintf(
		"rpc %s(%s%s) returns (%s%s)",
		r.Name, streamsReq, r.RequestType, streamsResp, r.ReturnsType)
	f.printWithIndent(header)
	if len(r.Elements) > 0 {
		fmt.Fprintf(f.w, " {")
		f.newline()
		for _, each := range r.Elements {
			f.indent()
			each.Accept(f)
			f.dedent()
		}
	}
	if len(r.Elements) > 0 {
		f.printWithIndent("}")
	} else {
		f.endWithComment(r.InlineComment)
	}
	f.lastStmt = "rpc"
	f.newline()
}

// VisitMessage formats a Message.
func (f *Formatter) VisitMessage(m *proto.Message) {
	// Always place a newline between top-level elements.
	// unless we're currently in a message.
	if f.currStmt != message {
		f.allowedLastStatements("proto")
	}
	prefix := message
	if m.IsExtend {
		prefix = "extend"
	}
	f.printWithIndent(fmt.Sprintf("%s %s {", prefix, m.Name))
	f.lastStmt = message
	if len(m.Elements) > 0 {
		f.newline()
		f.currStmt = message
		for _, each := range m.Elements {
			f.indent()
			// Don't forget to visit comment.
			f.VisitCommentable(each)
			each.Accept(f)
			f.dedent()
		}
	}
	f.printWithIndent("}")
	// reset it.
	f.lastStmt = message
	f.newline()
}

// VisitNormalField formats a NormalField.
func (f *Formatter) VisitNormalField(nf *proto.NormalField) {
	f.allowedLastStatements("normalfield", "message")
	// don't support optional, required which are proto2.
	prefix := ""
	switch {
	case nf.Repeated:
		prefix = "repeated "
	case nf.Optional:
		prefix = "optional "
	case nf.Required:
		prefix = "required "
	}
	f.printWithIndent(fmt.Sprintf("%s%s %s = %d", prefix, nf.Type, nf.Name, nf.Sequence))
	f.VisitInlineOptions(nf.Options)

	f.endWithComment(nf.InlineComment)
	f.lastStmt = "normalfield"
	f.newline()
}

// VisitOneof formats a Oneof.
func (f *Formatter) VisitOneof(o *proto.Oneof) {
	f.allowedLastStatements("message")

	f.printWithIndent(fmt.Sprintf("%s %s {", oneof, o.Name))
	f.lastStmt = oneof
	if len(o.Elements) > 0 {
		f.newline()
		f.currStmt = oneof
		for _, each := range o.Elements {
			f.indent()
			// Don't forget to visit comment.
			f.VisitCommentable(each)
			each.Accept(f)
			f.dedent()
		}
		f.printWithIndent("}")
	} else {
		io.WriteString(f.w, "}")
	}
	f.newline()
	f.lastStmt = oneof
}

// VisitOneofField formats a OneofField.
func (f *Formatter) VisitOneofField(o *proto.OneOfField) {
	f.allowedLastStatements("oneoffield", "oneof")

	f.printWithIndent(fmt.Sprintf("%s %s = %d", o.Type, o.Name, o.Sequence))
	f.VisitInlineOptions(o.Options)

	f.endWithComment(o.InlineComment)
	f.lastStmt = "oneoffield"
	f.newline()
}

// VisitReserved formats a Reserved.
func (f *Formatter) VisitReserved(r *proto.Reserved) {
	f.printWithIndent("reserved ")
	if len(r.Ranges) > 0 {
		for i, each := range r.Ranges {
			if i > 0 {
				io.WriteString(f.w, ", ")
			}
			fmt.Fprintf(f.w, "%s", each.SourceRepresentation())
		}
	} else {
		for i, each := range r.FieldNames {
			if i > 0 {
				io.WriteString(f.w, ", ")
			}
			fmt.Fprintf(f.w, "%q", each)
		}
	}
	f.endWithComment(r.InlineComment)
	f.lastStmt = "reserved"
	f.newline()
}

// VisitMapField formats a MapField.
func (f *Formatter) VisitMapField(m *proto.MapField) {
	f.allowedLastStatements("message")

	header := fmt.Sprintf("map<%s, %s> %s = %d", m.KeyType, m.Type, m.Name, m.Sequence)
	f.printWithIndent(header)
	f.VisitInlineOptions(m.Options)

	f.endWithComment(m.InlineComment)
	f.lastStmt = "mapfield"
	f.newline()
}

// Stub, not implemented for proto3.
func (f *Formatter) VisitExtensions(_ *proto.Extensions) {}

// Stub, not implemented for proto3.
func (f *Formatter) VisitGroup(_ *proto.Group) {}

// Small helpers for formatting.

// printWithIndent writes indentation based on the current indent level and then
// writes the passed in value.
func (f *Formatter) printWithIndent(value string) {
	for i := 0; i < f.indentLevel; i++ {
		io.WriteString(f.w, f.indentSeparator)
	}
	io.WriteString(f.w, value)
}

// indent increases the indent level.
func (f *Formatter) indent() {
	f.indentLevel++
}

// dedent decreases the indentation level.
func (f *Formatter) dedent() {
	f.indentLevel--
}

// newline writes a newline.
func (f *Formatter) newline() {
	io.WriteString(f.w, "\n")
}

// endWithComment writes a statement end (;) followed by inline comment if present.
func (f *Formatter) endWithComment(commentOrNil *proto.Comment) {
	io.WriteString(f.w, ";")
	if commentOrNil != nil {
		// trim left space
		inline := strings.TrimLeft(commentOrNil.Message(), " ")
		io.WriteString(f.w, fmt.Sprintf(" // %s", inline))
	}
}

// allowedLastStatements takes a list of stmts for which a new line won't be printed
// if lastStmt is among them. "embedded-comment" is checked by default so it can be omitted.
func (f *Formatter) allowedLastStatements(stmts ...string) {
	stmts = append(stmts, "embedded-comment")
	for _, each := range stmts {
		if f.lastStmt == each {
			return
		}
	}
	f.newline()
}

// VisitInlineOptions visits any options in opts and prints
// them as inline (w/o the 'option' prefix).
// Found in Normal, Oneof and Map fields.
func (f *Formatter) VisitInlineOptions(opts []*proto.Option) {
	optsLen := len(opts)
	if optsLen == 1 {
		io.WriteString(f.w, " [")
		f.VisitInlineOption(opts[0])
		io.WriteString(f.w, "]")
	} else if optsLen > 1 {
		// Format as an array.
		io.WriteString(f.w, " [")
		f.indent()
		for idx, opt := range opts {
			io.WriteString(f.w, "\n")
			f.printWithIndent("") // move f.w to right spot.
			f.VisitInlineOption(opt)
			// Actually doesn't parse if we add a comma at the end.
			if idx < optsLen-1 {
				io.WriteString(f.w, ",")
			}
		}
		io.WriteString(f.w, "\n")
		f.dedent()
		f.printWithIndent("]")
	}
}

// VisitInlineOption visits an option and prints it as inline. Used in
// VisitInlineOptions and while visiting enum field options (which aren't
// typed as Options, unfortunately)
func (f *Formatter) VisitInlineOption(o *proto.Option) {
	// never place a new line.
	// format in single line as optname = optval
	// whatever field this was will attach ';' and its inline comment.
	io.WriteString(f.w, fmt.Sprintf("%s = ", o.Name))
	f.formatLiteral(&o.Constant)
}

// currentNode gets the currentNode. Used mainly while visiting embedded comments so as to
// keep similar elements close.
func currentNode(n proto.Visitee) string {
	switch n.(type) {
	case *proto.Message:
		return "message"
	case *proto.NormalField:
		return "normalfield"
	case *proto.Enum:
		return "enum"
	case *proto.EnumField:
		return "enumfield"
	case *proto.Service:
		return "service"
	case *proto.RPC:
		return "rpc"
	case *proto.Oneof:
		return "oneof"
	case *proto.OneOfField:
		return "oneoffield"
	// MapField and Reserved have comments but dont
	// implement Documented :-(
	// case *proto.MapField:
	// 	return "mapfield"
	// case *proto.Reserved:
	// 	return "reserved"
	// case *proto.Extensions:
	// 	return "extensions"
	case *proto.Syntax:
		return "syntax"
	case *proto.Package:
		return "package"
	case *proto.Import:
		return "import"
	case *proto.Option:
		return "option"
	}
	panic("unreachable")
}
