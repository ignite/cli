package protoutil

import (
	"testing"

	"github.com/emicklei/proto"
	"github.com/stretchr/testify/require"
)

// Imports.
type testCreateImport struct {
	name, path, kind string
	out              *proto.Import
}

var testImport = []testCreateImport{
	{
		name: "simple import",
		path: "github.com/emicklei/proto.proto",
		kind: "weak",
		out: &proto.Import{
			Filename: "github.com/emicklei/proto.proto",
			Kind:     "weak",
		},
	},
	{
		name: "simple import",
		path: "github.com/emicklei/proto.proto",
		kind: "public",
		out: &proto.Import{
			Filename: "github.com/emicklei/proto.proto",
			Kind:     "public",
		},
	},
}

func TestCreateImport(t *testing.T) {
	for _, test := range testImport {
		spec := []ImportSpecOpts{}
		if test.kind == "weak" {
			spec = append(spec, Weak())
		} else if test.kind == "public" {
			spec = append(spec, Public())
		}
		imp := NewImport(test.path, spec...)
		require.Equal(t, test.out, imp, "expected %v, got %v", test.out, imp)
	}
}

// Packages.
type testCreatePackage struct {
	name string
	out  *proto.Package
}

var testPackage = []testCreatePackage{
	{
		name: "org.foo.hack",
		out: &proto.Package{
			Name: "org.foo.hack",
		},
	},
	{
		name: "simple.package",
		out: &proto.Package{
			Name: "simple.package",
		},
	},
}

func TestCreatePackage(t *testing.T) {
	for _, test := range testPackage {
		p := NewPackage(test.name)
		require.Equal(t, test.out, p, "expected %v, got %v", test.out, p)
	}
}

// Options
type testCreateOption struct {
	name, constant, setField string
	isCustom                 bool
	out                      *proto.Option
}

var testOption = []testCreateOption{
	{
		name:     "my_option",
		constant: "5",
		out: &proto.Option{
			Name:     "my_option",
			Constant: *NewLiteral("5"),
		},
	},
	{
		name:     "my_option",
		constant: "false",
		isCustom: true,
		out: &proto.Option{
			Name:     "(my_option)",
			Constant: *NewLiteral("false"),
		},
	},
	{
		name:     "my_option",
		constant: "2.341",
		setField: "my_field",
		isCustom: true,
		out: &proto.Option{
			Name:     "(my_option).my_field",
			Constant: *NewLiteral("2.341"),
		},
	},
}

func TestCreateOption(t *testing.T) {
	for _, test := range testOption {
		opts := []OptionSpecOpts{}
		if test.isCustom {
			opts = []OptionSpecOpts{Custom()}
		}
		if test.setField != "" {
			opts = append(opts, SetField(test.setField))
		}
		opt := NewOption(test.name, test.constant, opts...)
		require.Equal(t, test.out, opt, "expected %v, got %v", test.out, opt)
	}
}

// RPCs.
type testCreateRPC struct {
	name, inputType, outputType string
	streamsReq, streamsResp     bool
	options                     []*proto.Option
}

var rpcTests = []testCreateRPC{
	{
		name:       "my_rpc",
		inputType:  "my_input_type",
		outputType: "my_output_type",
	},
	{
		name:        "my_rpc",
		inputType:   "my_input_type",
		outputType:  "my_output_type",
		streamsReq:  true,
		streamsResp: true,
	},
	{
		name:       "my_rpc",
		inputType:  "my_input_type",
		outputType: "my_output_type",
		options: []*proto.Option{
			NewOption("my_option", "5"),
			NewOption("gogoproto.nullable", "false", Custom(), SetField("set")),
		},
	},
}

func TestCreateRPC(t *testing.T) {
	for _, test := range rpcTests {
		opts := []RPCSpecOpts{}
		if test.streamsReq {
			opts = append(opts, StreamRequest())
		}
		if test.streamsResp {
			opts = append(opts, StreamResponse())
		}
		if len(test.options) > 0 {
			opts = append(opts, WithRPCOptions(test.options...))
		}
		rpc := NewRPC(test.name, test.inputType, test.outputType, opts...)

		require.Equal(t, test.name, rpc.Name, "expected %v, got %v", test.name, rpc.Name)
		require.Equal(t, test.inputType, rpc.RequestType, "expected %v, got %v", test.inputType, rpc.ReturnsType)
		require.Equal(t, test.outputType, rpc.ReturnsType, "expected %v, got %v", test.outputType, rpc.ReturnsType)
		require.Equal(t, test.streamsReq, rpc.StreamsRequest, "expected %v, got %v", test.streamsReq, rpc.StreamsRequest)
		require.Equal(t, test.streamsResp, rpc.StreamsReturns, "expected %v, got %v", test.streamsResp, rpc.StreamsReturns)
		for i, opt := range rpc.Elements {
			opt, ok := opt.(*proto.Option)
			require.True(t, ok, "expected option, got %T", opt)
			require.Equal(t, test.options[i], opt, "expected %v, got %v", test.options[i], opt)
		}
		require.Equal(t, len(test.options), len(rpc.Elements), "expected %v, got %v", len(test.options), len(rpc.Elements))
	}
}

// Service Creation.
type testCreateService struct {
	name    string
	rpcs    []*proto.RPC
	options []*proto.Option
}

var serviceTests = []testCreateService{
	{
		name: "my_service",
		rpcs: []*proto.RPC{
			NewRPC("my_rpc", "my_input_type", "my_output_type"),
			NewRPC("my_other_rpc", "my_other_input_type", "my_other_output_type", StreamRequest(), StreamResponse()),
		},
		options: []*proto.Option{NewOption("my_option", "with a great value")},
	},
}

func TestCreateService(t *testing.T) {
	for _, test := range serviceTests {
		opts := []ServiceSpecOpts{}
		opts = append(opts, WithRPCs(test.rpcs...))
		opts = append(opts, WithServiceOptions(test.options...))
		rpc := NewService(test.name, opts...)

		require.Equal(t, test.name, rpc.Name, "expected %v, got %v", test.name, rpc.Name)
		// careful, options come first, then rpcs.
		lenOpts, lenRPCs := len(test.options), len(test.rpcs)
		require.True(t, len(rpc.Elements) == lenOpts+lenRPCs, "expected %v, got %v", lenOpts+lenRPCs, len(rpc.Elements))
		for i, opt := range rpc.Elements {
			if i < lenOpts {
				opt, ok := opt.(*proto.Option)
				require.True(t, ok, "expected option, got %T", opt)
				require.Equal(t, test.options[i], opt, "expected %v, got %v", test.options[i], opt)
			} else {
				rpc, ok := opt.(*proto.RPC)
				require.True(t, ok, "expected rpc, got %T", opt)
				require.Equal(t, test.rpcs[i-lenOpts], rpc, "expected %v, got %v", test.rpcs[i-lenOpts], rpc)
			}
		}
	}
}

// Fields
type testCreateField struct {
	name, typeName               string
	sequence                     int
	repeated, optional, required bool
	options                      []*proto.Option
}

var fieldTests = []testCreateField{
	{
		name:     "my_field",
		typeName: "my_type",
		sequence: 1,
		repeated: true,
	},
	{
		name:     "my_field",
		typeName: "my_type",
		sequence: 2,
		optional: true,
	},
	{
		name:     "my_field",
		typeName: "my_type",
		sequence: 3,
		required: true,
	},
	{
		name:     "my_field",
		typeName: "my_type",
		sequence: 4,
		options: []*proto.Option{
			NewOption("my_option", "5"),
			NewOption("gogoproto.nullable", "false", Custom(), SetField("set")),
		},
	},
}

func TestCreateField(t *testing.T) {
	for _, test := range fieldTests {
		opts := []FieldSpecOpts{}
		if test.repeated {
			opts = append(opts, Repeated())
		}
		if test.optional {
			opts = append(opts, Optional())
		}
		if test.required {
			opts = append(opts, Required())
		}
		opts = append(opts, WithFieldOptions(test.options...))
		field := NewField(test.typeName, test.name, test.sequence, opts...)

		require.Equal(t, test.name, field.Name, "expected %v, got %v", test.name, field.Name)
		require.Equal(t, test.typeName, field.Type, "expected %v, got %v", test.typeName, field.Type)
		require.Equal(t, test.sequence, field.Sequence, "expected %v, got %v", test.sequence, field.Sequence)
		require.Equal(t, test.repeated, field.Repeated, "expected %v, got %v", test.repeated, field.Repeated)
		require.Equal(t, test.optional, field.Optional, "expected %v, got %v", test.optional, field.Optional)
		require.Equal(t, test.required, field.Required, "expected %v, got %v", test.required, field.Required)
		for i, opt := range field.Options {
			require.Equal(t, test.options[i], opt, "expected %v, got %v", test.options[i], opt)
		}
		require.Equal(t, len(test.options), len(field.Options), "expected %v, got %v", len(test.options), len(field.Options))
	}
}

// Messages.
type testCreateMessage struct {
	name     string
	fields   []*proto.NormalField
	enums    []*proto.Enum
	options  []*proto.Option
	isExtend bool
}

var messageTests = []testCreateMessage{
	{
		name: "my_message",
		fields: []*proto.NormalField{
			NewField("my_field", "my_type", 1),
			NewField("my_other_field", "my_other_type", 2),
		},
	},
	{
		name: "my_message",
		fields: []*proto.NormalField{
			NewField("my_field", "my_type", 1),
			NewField("my_other_field", "my_other_type", 2),
		},
		enums: []*proto.Enum{NewEnum("my_enum")},
		options: []*proto.Option{
			NewOption("my_option", "with a great value"),
			NewOption("gogoproto.nullable", "false", Custom(), SetField("set")),
		},
		isExtend: true,
	},
}

func TestCreateMessage(t *testing.T) {
	for _, test := range messageTests {
		opts := []MessageSpecOpts{}
		opts = append(opts, WithFields(test.fields...))
		opts = append(opts, WithEnums(test.enums...))
		opts = append(opts, WithMessageOptions(test.options...))
		if test.isExtend {
			opts = append(opts, Extend())
		}
		message := NewMessage(test.name, opts...)

		require.Equal(t, test.name, message.Name, "expected %v, got %v", test.name, message.Name)
		require.Equal(t, test.isExtend, message.IsExtend, "expected %v, got %v", test.isExtend, message.IsExtend)

		// options added first, then fields and then enums.
		lenOpts, lenFields, lenEnums := len(test.options), len(test.fields), len(test.enums)
		for i, field := range message.Elements {
			if i < lenOpts {
				opt, ok := field.(*proto.Option)
				require.True(t, ok, "expected option, got %T", field)
				require.Equal(t, test.options[i], opt, "expected %v, got %v", test.options[i], opt)
			} else if i < lenOpts+lenFields {
				field, ok := field.(*proto.NormalField)
				require.True(t, ok, "expected field, got %T", field)
				require.Equal(t, test.fields[i-lenOpts], field, "expected %v, got %v", test.fields[i-lenOpts], field)
			} else {
				enum, ok := field.(*proto.Enum)
				require.True(t, ok, "expected enum, got %T", field)
				require.Equal(t, test.enums[i-lenOpts-lenFields], enum, "expected %v, got %v", test.enums[i-lenOpts-lenFields], enum)
			}
		}
		require.True(t, lenOpts+lenFields+lenEnums == len(message.Elements), "expected %v, got %v", lenOpts+lenFields+lenEnums, len(message.Elements))
	}
}

// Enum fields.
type testCreateEnumField struct {
	name    string
	value   int
	options []*proto.Option
}

var enumFieldTests = []testCreateEnumField{
	{
		name:  "my_field",
		value: 1,
	},
	{
		name:  "my_field",
		value: 2,
		options: []*proto.Option{
			NewOption("my_option", "with a great value"),
			NewOption("gogoproto.nullable", "false", Custom(), SetField("set")),
		},
	},
}

func TestCreateEnumField(t *testing.T) {
	for _, test := range enumFieldTests {
		opts := []EnumFieldSpecOpts{}
		opts = append(opts, WithEnumFieldOptions(test.options...))
		field := NewEnumField(test.name, test.value, opts...)

		require.Equal(t, test.name, field.Name, "expected %v, got %v", test.name, field.Name)
		require.Equal(t, test.value, field.Integer, "expected %v, got %v", test.value, field.Integer)
		for i, opt := range field.Elements {
			opt, ok := opt.(*proto.Option)
			require.True(t, ok, "expected option, got %T", opt)
			require.Equal(t, test.options[i], opt, "expected %v, got %v", test.options[i], opt)
		}
		require.Equal(t, len(test.options), len(field.Elements), "expected %v, got %v", len(test.options), len(field.Elements))
	}
}

// Enums:
type testCreateEnum struct {
	name    string
	options []*proto.Option
	values  []*proto.EnumField
}

var enumTests = []testCreateEnum{
	{
		name: "my_enum",
		values: []*proto.EnumField{
			NewEnumField("my_value", 1),
			NewEnumField("my_other_value", 2),
		},
	},
	{
		name: "my_enum",
		values: []*proto.EnumField{
			NewEnumField("my_value", 1),
			NewEnumField("my_other_value", 2),
		},
		options: []*proto.Option{
			NewOption("my_option", "with a great value"),
		},
	},
}

func TestCreateEnum(t *testing.T) {
	for _, test := range enumTests {
		opts := []EnumSpecOpts{}
		opts = append(opts, WithEnumFields(test.values...))
		opts = append(opts, WithEnumOptions(test.options...))
		enum := NewEnum(test.name, opts...)

		require.Equal(t, test.name, enum.Name, "expected %v, got %v", test.name, enum.Name)
		lenFields, lenOptions := len(test.values), len(test.options)
		for i, opt := range enum.Elements[:lenOptions] {
			opt, ok := opt.(*proto.Option)
			require.True(t, ok, "expected option, got %T", opt)
			require.Equal(t, test.options[i], opt, "expected %v, got %v", test.options[i], opt)
		}
		for i, value := range enum.Elements[lenOptions:] {
			value, ok := value.(*proto.EnumField)
			require.True(t, ok, "expected enum field, got %T", value)
			require.Equal(t, test.values[i], value, "expected %v, got %v", test.values[i], value)
		}
		require.Equal(t, lenOptions+lenFields, len(enum.Elements), "expected %v, got %v", lenOptions+lenFields, len(enum.Elements))
	}
}

// OneOf fields:
type testCreateOneofField struct {
	name, typeName string
	sequence       int
	options        []*proto.Option
}

var oneoffieldTests = []testCreateOneofField{
	{
		name:     "my_field",
		typeName: "my_type",
		sequence: 1,
	},
	{
		name:     "my_field",
		typeName: "my_type",
		sequence: 4,
		options: []*proto.Option{
			NewOption("my_option", "5"),
			NewOption("gogoproto.nullable", "false", Custom(), SetField("set")),
		},
	},
}

func TestCreateOneofField(t *testing.T) {
	for _, test := range oneoffieldTests {
		opts := []OneOfFieldOpts{WithOneOfFieldOptions(test.options...)}
		field := NewOneOfField(test.typeName, test.name, test.sequence, opts...)

		require.Equal(t, test.name, field.Name, "expected %v, got %v", test.name, field.Name)
		require.Equal(t, test.typeName, field.Type, "expected %v, got %v", test.typeName, field.Type)
		require.Equal(t, test.sequence, field.Sequence, "expected %v, got %v", test.sequence, field.Sequence)

		for i, opt := range field.Options {
			require.Equal(t, test.options[i], opt, "expected %v, got %v", test.options[i], opt)
		}
		require.Equal(t, len(test.options), len(field.Options), "expected %v, got %v", len(test.options), len(field.Options))
	}
}

// Oneof:
type testCreateOneof struct {
	name    string
	options []*proto.Option
	values  []*proto.OneOfField
}

var oneofTests = []testCreateOneof{
	{
		name: "oneof_this",
		values: []*proto.OneOfField{
			NewOneOfField("my_value", "my_type", 1),
			NewOneOfField("my_other_value", "my_type", 2),
		},
	},
	{
		name: "oneof_that",
		values: []*proto.OneOfField{
			NewOneOfField("my_value", "my_type", 1),
		},
		options: []*proto.Option{
			NewOption("my_option", "with a great value"),
		},
	},
}

func TestCreateOneof(t *testing.T) {
	for _, test := range oneofTests {
		opts := []OneOfSpecOpts{}
		opts = append(opts, WithOneOfFields(test.values...))
		opts = append(opts, WithOneOfOptions(test.options...))
		oneof := NewOneOf(test.name, opts...)

		require.Equal(t, test.name, oneof.Name, "expected %v, got %v", test.name, oneof.Name)
		lenFields, lenOptions := len(test.values), len(test.options)
		for i, opt := range oneof.Elements[:lenOptions] {
			opt, ok := opt.(*proto.Option)
			require.True(t, ok, "expected option, got %T", opt)
			require.Equal(t, test.options[i], opt, "expected %v, got %v", test.options[i], opt)
		}
		for i, value := range oneof.Elements[lenOptions:] {
			value, ok := value.(*proto.OneOfField)
			require.True(t, ok, "expected oneof field, got %T", value)
			require.Equal(t, test.values[i], value, "expected %v, got %v", test.values[i], value)
		}
		require.Equal(t, lenOptions+lenFields, len(oneof.Elements), "expected %v, got %v", lenOptions+lenFields, len(oneof.Elements))
	}
}

func TestAttachComment(t *testing.T) {
	// Attach comment to message
	msg := NewMessage("my_message")
	AttachComment(msg, "my comment")
	require.Equal(t, " my comment", msg.Comment.Lines[0], "expected %v, got %v", "my comment", msg.Comment.Lines[0])

	// Attach comment to an rpc call
	rpc := NewRPC("my_rpc", "my_request", "my_response")
	AttachComment(rpc, "my comment")
	require.Equal(t, " my comment", rpc.Comment.Lines[0], "expected %v, got %v", "my comment", rpc.Comment.Lines[0])

	// Attach comment to a service
	svc := NewService("my_service")
	AttachComment(svc, "my comment")
	require.Equal(t, " my comment", svc.Comment.Lines[0], "expected %v, got %v", "my comment", svc.Comment.Lines[0])
}

func TestIsString(t *testing.T) {
	require.True(t, isString("string"))
	require.True(t, isString("THIS/PATH/IS/STRING"))

	// Don't report "true" and "false" as strings
	require.False(t, isString("true"))
	require.False(t, isString("false"))

	// Dont report numbers as strings
	require.False(t, isString("1"))
	require.False(t, isString("1.0"))
	require.False(t, isString("1.0e-10"))
	require.False(t, isString("1.0e+10"))
	require.False(t, isString("1.0e10"))
	require.False(t, isString("3.1929348317293483e-10"))

	// A single numbers means not a string, parser would fail with that either way.
	require.True(t, isString("isthisastringohnoitactuallyisn't1.0"))
}
