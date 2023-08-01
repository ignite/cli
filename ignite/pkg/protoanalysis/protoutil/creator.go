// Package protoutil wraps proto structs to allow easier creation, protobuf lang is small enough
// to easily allow this.
package protoutil

import (
	"fmt"
	"strconv"

	"github.com/emicklei/proto"
)

// TODO: Can also support comments/inline comments? -- Probably, formatting is currently
// flaky with how it prints them, though.

// Values for the kind of import.
const (
	KindWeak   = "weak"
	KindPublic = "public"
)

// NewLiteral creates a new Literal:
//
// // true
// l := NewLiteral("true")
//
// // 1
// l := NewLiteral("1")
//
// // "foo"
// l := NewLiteral("foo")
//
// Currently doesn't support creating compound literals (arrays/maps).
func NewLiteral(lit string) *proto.Literal {
	return &proto.Literal{
		Source:   lit,
		IsString: isString(lit),
	}
}

// ImportSpec holds information relevant to the import statement.
type ImportSpec struct {
	path string
	kind string
}

// ImportSpecOptions is a type alias for a callable accepting an ImportSpec.
type ImportSpecOptions func(i *ImportSpec)

// Weak allows you to set the kind of the import statement to 'weak'.
func Weak() ImportSpecOptions {
	return func(i *ImportSpec) {
		i.kind = KindWeak
	}
}

// Public allows you to set the kind of the import statement to 'public'.
func Public() ImportSpecOptions {
	return func(i *ImportSpec) {
		i.kind = KindPublic
	}
}

// NewImport creates a new import statement node:
//
//		// import "myproto.proto";
//	 imp := NewImport("myproto.proto")
//
// By default, no kind is assigned to it, by using Weak or Public, this can be specified:
//
//	// import weak "myproto.proto";
//	imp := NewImport("myproto.proto", Weak())
func NewImport(path string, opts ...ImportSpecOptions) *proto.Import {
	i := ImportSpec{path: path}
	for _, opt := range opts {
		opt(&i)
	}

	return &proto.Import{
		Filename: i.path,
		Kind:     i.kind,
	}
}

// NewPackage creates a new package statement node:
//
//	// package foo.bar;
//	pkg := NewPackage("foo.bar")
func NewPackage(path string) *proto.Package {
	return &proto.Package{
		Name: path,
	}
}

// OptionSpec holds information relevant to the option statement.
type OptionSpec struct {
	name     string
	setter   string
	constant string
	custom   bool
}

// OptionSpecOptions is a function that accepts an OptionSpec.
type OptionSpecOptions func(o *OptionSpec)

// Custom denotes the option as being a custom option.
func Custom() OptionSpecOptions {
	return func(f *OptionSpec) {
		f.custom = true
	}
}

// SetField allows setting specific fields for a given option
// that denotes a type with fields.
//
//	// option (my_opt).field = "Value";
//	opt := NewOption("my_opt", "Value", Custom(), Setter("field"))
func SetField(name string) OptionSpecOptions {
	return func(f *OptionSpec) {
		f.setter = name
	}
}

// NewOption creates a new option statement node:
//
//	// option foo = 1;
//	opt := NewOption("foo", "1")
//
// Custom options can be marked as such by using Custom, this wraps the option name
// in parenthesis:
//
//	// option (foo) = 1;
//	opt := NewOption("foo", "1", Custom())
//
// Since option constants can accept a number of types, strings that require quotation
// should be passed as raw strings:
//
//	// option foo = "bar";
//	opt := NewOption("foo", `bar`)
func NewOption(name, constant string, opts ...OptionSpecOptions) *proto.Option {
	o := OptionSpec{name: name, constant: constant}
	for _, opt := range opts {
		opt(&o)
	}
	if o.custom {
		o.name = fmt.Sprintf("(%s)", o.name)
	}
	// add the field we are setting outside the parentheses.
	if o.setter != "" {
		o.name = fmt.Sprintf("%s.%s", o.name, o.setter)
	}
	return &proto.Option{
		Name:     o.name,
		Constant: *NewLiteral(o.constant),
	}
}

/// Service + PRC

// RPCSpec holds information relevant to the rpc statement.
type RPCSpec struct {
	name, inputType, outputType string
	streamsReq, streamsResp     bool
	options                     []*proto.Option
}

// RPCSpecOptions is a type alias for a callable accepting an RPCSpec.
type RPCSpecOptions func(i *RPCSpec)

// StreamRequest marks request as streaming.
func StreamRequest() RPCSpecOptions {
	return func(r *RPCSpec) {
		r.streamsReq = true
	}
}

// StreamResponse marks response as streaming.
func StreamResponse() RPCSpecOptions {
	return func(r *RPCSpec) {
		r.streamsResp = true
	}
}

// WithRPCOptions adds options to the RPC.
func WithRPCOptions(option ...*proto.Option) RPCSpecOptions {
	return func(o *RPCSpec) {
		o.options = append(o.options, option...)
	}
}

// NewRPC creates a new RPC statement node:
//
//	// rpc Foo(Bar) returns(Bar) {}
//	rpc := NewRPC("Foo", "Bar", "Bar")
//
// No options are attached by default, use WithRPCOptions to add options as required:
//
//	// rpc Foo(Bar) returns(Bar) {
//	//  option (foo) = 1;
//	// }
//	rpc := NewRPC("Foo", "Bar", "Bar", WithRPCOptions(NewOption("foo", "1")))
func NewRPC(name, inputType, outputType string, opts ...RPCSpecOptions) *proto.RPC {
	r := RPCSpec{name: name, inputType: inputType, outputType: outputType}
	for _, opt := range opts {
		opt(&r)
	}

	rpc := &proto.RPC{
		Name:           r.name,
		RequestType:    r.inputType,
		ReturnsType:    r.outputType,
		StreamsRequest: r.streamsReq,
		StreamsReturns: r.streamsResp,
	}
	if len(r.options) > 0 {
		for _, opt := range r.options {
			rpc.Elements = append(rpc.Elements, opt)
		}
	}
	return rpc
}

// ServiceSpec holds information relevant to the service statement.
type ServiceSpec struct {
	name string
	rpcs []*proto.RPC
	opts []*proto.Option
}

// ServiceSpecOptions is a type alias for a callable accepting a ServiceSpec.
type ServiceSpecOptions func(i *ServiceSpec)

// WithRPCs adds rpcs to the service.
func WithRPCs(rpcs ...*proto.RPC) ServiceSpecOptions {
	return func(s *ServiceSpec) {
		s.rpcs = append(s.rpcs, rpcs...)
	}
}

// WithServiceOptions adds options to the service.
func WithServiceOptions(options ...*proto.Option) ServiceSpecOptions {
	return func(s *ServiceSpec) {
		s.opts = append(s.opts, options...)
	}
}

// NewService creates a new service statement node:
//
//	// service Foo {}
//	service := NewService("Foo")
//
// No rpcs/options are attached by default, use WithRPCs and
// WithServiceOptions to add them as required:
//
//	 // service Foo {
//	 //  option (foo) = 1;
//	 //  rpc Bar(Bar) returns (Bar) {}
//	 // }
//		opt := NewOption("foo", "1")
//	 rpc := NewRPC("Bar", "Bar", "Bar")
//	 service := NewService("Foo", WithServiceOptions(opt), WithRPCs(rpc))
//
// By default, options are added first and then the rpcs.
func NewService(name string, opts ...ServiceSpecOptions) *proto.Service {
	s := ServiceSpec{name: name}
	for _, opt := range opts {
		opt(&s)
	}
	service := &proto.Service{
		Name: s.name,
	}
	for _, opt := range s.opts {
		service.Elements = append(service.Elements, opt)
	}
	for _, rpc := range s.rpcs {
		service.Elements = append(service.Elements, rpc)
	}
	return service
}

/// Message + NormalField

// FieldSpec holds information relevant to the field statement.
type FieldSpec struct {
	name, typename               string
	sequence                     int
	repeated, optional, required bool
	options                      []*proto.Option
}

// FieldSpecOptions is a type alias for a callable accepting a FieldSpec.
type FieldSpecOptions func(f *FieldSpec)

// Repeated marks the field as repeated.
func Repeated() FieldSpecOptions {
	return func(f *FieldSpec) {
		f.repeated = true
	}
}

// Optional marks the field as optional.
func Optional() FieldSpecOptions {
	return func(f *FieldSpec) {
		f.optional = true
	}
}

// Required marks the field as required.
func Required() FieldSpecOptions {
	return func(f *FieldSpec) {
		f.required = true
	}
}

// WithFieldOptions adds options to the field.
func WithFieldOptions(options ...*proto.Option) FieldSpecOptions {
	return func(f *FieldSpec) {
		f.options = append(f.options, options...)
	}
}

// NewField creates a new field statement node:
//
//	// int32 Foo = 1;
//	field := NewField("Foo", "int32", 1)
//
// Fields aren't marked as repeated, required or optional. Use Repeated, Optional
// and Required to mark the field as such.
//
//	// repeated int32 Foo = 1;
//	field := NewField("Foo", "int32", 1, Repeated())
func NewField(name, typename string, sequence int, opts ...FieldSpecOptions) *proto.NormalField {
	f := FieldSpec{name: name, typename: typename, sequence: sequence}
	for _, opt := range opts {
		opt(&f)
	}

	// Check qualifiers? Though protoc will shout if we do stupid things.
	field := &proto.NormalField{
		Field: &proto.Field{
			Name:     f.name,
			Sequence: f.sequence,
			Type:     f.typename,
			Options:  []*proto.Option{},
		},
		Repeated: f.repeated,
		Required: f.required,
		Optional: f.optional,
	}
	if len(f.options) > 0 {
		field.Options = append(field.Options, f.options...)
	}
	return field
}

// MessageSpec holds information relevant to the message statement.
type MessageSpec struct {
	name     string
	fields   []*proto.NormalField
	enums    []*proto.Enum
	options  []*proto.Option
	isExtend bool
}

// MessageSpecOptions is a type alias for a callable accepting a MessageSpec.
type MessageSpecOptions func(i *MessageSpec)

// WithMessageOptions adds options to the message.
func WithMessageOptions(options ...*proto.Option) MessageSpecOptions {
	return func(m *MessageSpec) {
		m.options = append(m.options, options...)
	}
}

// WithFields adds fields to the message.
func WithFields(fields ...*proto.NormalField) MessageSpecOptions {
	return func(m *MessageSpec) {
		m.fields = append(m.fields, fields...)
	}
}

// WithEnums adds enums to the message.
func WithEnums(enum ...*proto.Enum) MessageSpecOptions {
	return func(m *MessageSpec) {
		m.enums = append(m.enums, enum...)
	}
}

func Extend() MessageSpecOptions {
	return func(m *MessageSpec) {
		m.isExtend = true
	}
}

// NewMessage creates a new message statement node:
//
//	// message Foo {}
//	message := NewMessage("Foo")
//
// No fields/enums/options are attached by default, use WithMessageFields, WithEnums,
// and WithMessageOptions to add them as required:
//
//	 // message Foo {
//	 //  option (foo) = 1;
//	 //  int32 Bar = 1;
//	 // }
//		opt := NewOption("foo", "1")
//	 field := NewField("int32", "Bar", 1)
//	 message := NewMessage("Foo", WithMessageOptions(opt), WithFields(field))
//
// By default, options are added first, then fields and then enums.
func NewMessage(name string, opts ...MessageSpecOptions) *proto.Message {
	m := MessageSpec{name: name}
	for _, opt := range opts {
		opt(&m)
	}
	message := &proto.Message{
		Name:     m.name,
		IsExtend: m.isExtend,
	}
	for _, opt := range m.options {
		message.Elements = append(message.Elements, opt)
	}

	// Verify that fields have unique sequence? Though, again, protoc will shout if
	// it isn't the case.
	for _, field := range m.fields {
		message.Elements = append(message.Elements, field)
	}
	for _, enum := range m.enums {
		message.Elements = append(message.Elements, enum)
	}
	return message
}

// EnumFieldSpec holds information relevant to the enum field statement.
type EnumFieldSpec struct {
	name    string
	value   int
	options []*proto.Option
}

// EnumFieldSpecOptions is a type alias for a callable accepting an EnumFieldSpec.
type EnumFieldSpecOptions func(f *EnumFieldSpec)

// WithEnumFieldOptions adds options to the enum field.
func WithEnumFieldOptions(options ...*proto.Option) EnumFieldSpecOptions {
	return func(f *EnumFieldSpec) {
		f.options = append(f.options, options...)
	}
}

// NewEnumField creates a new enum field statement node:
//
//	// BAR = 1;
//	field := NewEnumField("BAR", 1)
//
// No options are attached by default, use WithEnumFieldOptions to add them as
// required:
//
//	// BAR = 1 [option (foo) = 1];
//	field := NewEnumField("BAR", 1, WithEnumFieldOptions(NewOption("foo", "1")))
func NewEnumField(name string, value int, opts ...EnumFieldSpecOptions) *proto.EnumField {
	f := EnumFieldSpec{name: name, value: value}
	for _, opt := range opts {
		opt(&f)
	}

	field := &proto.EnumField{
		Name:    f.name,
		Integer: f.value,
	}
	for _, opt := range f.options {
		field.Elements = append(field.Elements, opt)
	}
	return field
}

// EnumSpec holds information relevant to the enum statement.
type EnumSpec struct {
	name    string
	fields  []*proto.EnumField
	options []*proto.Option
}

// EnumSpecOpts is a type alias for a callable accepting an EnumSpec.
type EnumSpecOpts func(i *EnumSpec)

// WithEnumOptions adds options to the enum.
func WithEnumOptions(options ...*proto.Option) EnumSpecOpts {
	return func(e *EnumSpec) {
		e.options = append(e.options, options...)
	}
}

// WithEnumFields adds fields to the enum.
func WithEnumFields(fields ...*proto.EnumField) EnumSpecOpts {
	return func(e *EnumSpec) {
		e.fields = append(e.fields, fields...)
	}
}

// NewEnum creates a new enum statement node:
//
//	// enum Foo {
//	//  BAR = 1;
//	// }
//	enum := NewEnum("Foo", WithEnumFields(NewEnumField("BAR", 1)))
//
// No options are attached by default, use WithEnumOptions to add them as
// required:
//
//	// enum Foo {
//	//  BAR = 1 [option (foo) = 1];
//	// }
//	enum := NewEnum("Foo", WithEnumOptions(NewOption("foo", "1")), WithEnumFields(NewEnumField("BAR", 1)))
//
// By default, options are added first, then fields.
func NewEnum(name string, opts ...EnumSpecOpts) *proto.Enum {
	e := EnumSpec{name: name}
	for _, opt := range opts {
		opt(&e)
	}
	enum := &proto.Enum{
		Name: e.name,
	}
	for _, opt := range e.options {
		enum.Elements = append(enum.Elements, opt)
	}
	for _, field := range e.fields {
		enum.Elements = append(enum.Elements, field)
	}
	return enum
}

// OneofFieldSpec holds information relevant to the oneof field statement.
type OneofFieldSpec struct {
	name, typename string
	sequence       int
	options        []*proto.Option
}

// OneofFieldOptions is a type alias for a callable accepting a OneOfField.
type OneofFieldOptions func(f *OneofFieldSpec)

// WithOneofFieldOptions adds options to the oneof field.
func WithOneofFieldOptions(options ...*proto.Option) OneofFieldOptions {
	return func(f *OneofFieldSpec) {
		f.options = append(f.options, options...)
	}
}

// NewOneofField creates a new oneof field statement node:
//
//		// Needs to placed in oneof block.
//	 // int32 Foo = 1;
//	 field := NewOneofField("Foo", "int32", 1)
//
// Additional options can be created and attached to the field to the field via
// WithOneOfFieldOptions:
//
//	// int32 Foo = 1 [option (foo) = 1];
//	field := NewOneofField("Foo", "int32", 1, WithOneOfFieldOptions(NewOption("foo", "1")))
func NewOneofField(name, typename string, sequence int, opts ...OneofFieldOptions) *proto.OneOfField {
	f := OneofFieldSpec{name: name, typename: typename, sequence: sequence}
	for _, opt := range opts {
		opt(&f)
	}
	field := &proto.OneOfField{
		Field: &proto.Field{
			Name:     f.name,
			Sequence: f.sequence,
			Type:     f.typename,
			Options:  []*proto.Option{},
		},
	}
	field.Options = append(field.Options, f.options...)
	return field
}

// OneofSpec holds information relevant to the enum statement.
type OneofSpec struct {
	name    string
	options []*proto.Option
	fields  []*proto.OneOfField
}

// OneofSpecOptions is a type alias for a callable accepting a OneOfSpec.
type OneofSpecOptions func(o *OneofSpec)

// WithOneofOptions adds options to the oneof.
func WithOneofOptions(options ...*proto.Option) OneofSpecOptions {
	return func(o *OneofSpec) {
		o.options = append(o.options, options...)
	}
}

// WithOneofFields adds fields to the oneof.
func WithOneofFields(fields ...*proto.OneOfField) OneofSpecOptions {
	return func(o *OneofSpec) {
		o.fields = append(o.fields, fields...)
	}
}

// NewOneof creates a new oneof statement node:
//
//	// oneof Foo {
//	//  int32 Foo = 1;
//	// }
//	oneof := NewOneof("Foo", WithOneOfFields(NewOneOfField("Foo", "int32", 1)))
//
// No options are attached by default, use WithOneOfOptions to add them as required.
func NewOneof(name string, opts ...OneofSpecOptions) *proto.Oneof {
	o := OneofSpec{name: name}
	for _, opt := range opts {
		opt(&o)
	}
	oneof := &proto.Oneof{
		Name: o.name,
	}
	for _, opt := range o.options {
		oneof.Elements = append(oneof.Elements, opt)
	}
	for _, field := range o.fields {
		oneof.Elements = append(oneof.Elements, field)
	}
	return oneof
}

// AttachComment attaches a comment top level nodes. Currently only supports Messages, RPC's
// and Services. Silently ignores other nodes though they can easily be added by just appending
// a new case to the switch statement.
func AttachComment(n proto.Visitee, comment string) {
	c := &proto.Comment{
		// Attach a starting space here, i.e // text and not //text
		Lines: []string{" " + comment},
	}
	switch n := n.(type) {
	case *proto.Message:
		n.Comment = c
	case *proto.RPC:
		n.Comment = c
	case *proto.Service:
		n.Comment = c
	}
}

// Check if s is a string, exclude special cases of "false" and "true".
func isString(s string) bool {
	if s == "true" || s == "false" {
		return false
	}
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return false
	}
	return true
}
