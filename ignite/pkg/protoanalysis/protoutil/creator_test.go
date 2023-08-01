package protoutil_test

import (
	"testing"

	"github.com/emicklei/proto"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/protoanalysis/protoutil"
)

// Imports.
func TestCreateImport(t *testing.T) {
	cases := []struct {
		name, path, kind string
		out              *proto.Import
	}{
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

	for _, test := range cases {
		var spec []protoutil.ImportSpecOptions
		switch test.kind {
		case "weak":
			spec = append(spec, protoutil.Weak())
		case "public":
			spec = append(spec, protoutil.Public())
		}
		imp := protoutil.NewImport(test.path, spec...)
		require.Equal(t, test.out, imp, "expected %v, got %v", test.out, imp)
	}
}

// Packages.
func TestCreatePackage(t *testing.T) {
	cases := []struct {
		name string
		out  *proto.Package
	}{
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

	for _, test := range cases {
		p := protoutil.NewPackage(test.name)
		require.Equal(t, test.out, p, "expected %v, got %v", test.out, p)
	}
}

// Options
func TestCreateOption(t *testing.T) {
	cases := []struct {
		name, constant, setField string
		isCustom                 bool
		out                      *proto.Option
	}{
		{
			name:     "my_option",
			constant: "5",
			out: &proto.Option{
				Name:     "my_option",
				Constant: *protoutil.NewLiteral("5"),
			},
		},
		{
			name:     "my_option",
			constant: "false",
			isCustom: true,
			out: &proto.Option{
				Name:     "(my_option)",
				Constant: *protoutil.NewLiteral("false"),
			},
		},
		{
			name:     "my_option",
			constant: "2.341",
			setField: "my_field",
			isCustom: true,
			out: &proto.Option{
				Name:     "(my_option).my_field",
				Constant: *protoutil.NewLiteral("2.341"),
			},
		},
	}

	for _, test := range cases {
		var opts []protoutil.OptionSpecOptions
		if test.isCustom {
			opts = []protoutil.OptionSpecOptions{protoutil.Custom()}
		}
		if test.setField != "" {
			opts = append(opts, protoutil.SetField(test.setField))
		}
		opt := protoutil.NewOption(test.name, test.constant, opts...)
		require.Equal(t, test.out, opt, "expected %v, got %v", test.out, opt)
	}
}

// RPCs.
func TestCreateRPC(t *testing.T) {
	cases := []struct {
		name, inputType, outputType string
		streamsReq, streamsResp     bool
		options                     []*proto.Option
	}{
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
				protoutil.NewOption("my_option", "5"),
				protoutil.NewOption("gogoproto.nullable", "false", protoutil.Custom(), protoutil.SetField("set")),
			},
		},
	}

	for _, test := range cases {
		var opts []protoutil.RPCSpecOptions
		if test.streamsReq {
			opts = append(opts, protoutil.StreamRequest())
		}
		if test.streamsResp {
			opts = append(opts, protoutil.StreamResponse())
		}
		if len(test.options) > 0 {
			opts = append(opts, protoutil.WithRPCOptions(test.options...))
		}
		rpc := protoutil.NewRPC(test.name, test.inputType, test.outputType, opts...)

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

// Services.
func TestCreateService(t *testing.T) {
	cases := []struct {
		name    string
		rpcs    []*proto.RPC
		options []*proto.Option
	}{
		{
			name: "my_service",
			rpcs: []*proto.RPC{
				protoutil.NewRPC("my_rpc", "my_input_type", "my_output_type"),
				protoutil.NewRPC("my_other_rpc", "my_other_input_type", "my_other_output_type", protoutil.StreamRequest(), protoutil.StreamResponse()),
			},
			options: []*proto.Option{protoutil.NewOption("my_option", "with a great value")},
		},
	}

	for _, test := range cases {
		var opts []protoutil.ServiceSpecOptions
		opts = append(opts, protoutil.WithRPCs(test.rpcs...))
		opts = append(opts, protoutil.WithServiceOptions(test.options...))
		rpc := protoutil.NewService(test.name, opts...)

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
func TestCreateField(t *testing.T) {
	cases := []struct {
		name, typeName               string
		sequence                     int
		repeated, optional, required bool
		options                      []*proto.Option
	}{
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
				protoutil.NewOption("my_option", "5"),
				protoutil.NewOption("gogoproto.nullable", "false", protoutil.Custom(), protoutil.SetField("set")),
			},
		},
	}

	for _, test := range cases {
		var opts []protoutil.FieldSpecOptions
		if test.repeated {
			opts = append(opts, protoutil.Repeated())
		}
		if test.optional {
			opts = append(opts, protoutil.Optional())
		}
		if test.required {
			opts = append(opts, protoutil.Required())
		}
		opts = append(opts, protoutil.WithFieldOptions(test.options...))
		field := protoutil.NewField(test.name, test.typeName, test.sequence, opts...)

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
func TestCreateMessage(t *testing.T) {
	cases := []struct {
		name     string
		fields   []*proto.NormalField
		enums    []*proto.Enum
		options  []*proto.Option
		isExtend bool
	}{
		{
			name: "my_message",
			fields: []*proto.NormalField{
				protoutil.NewField("my_field", "my_type", 1),
				protoutil.NewField("my_other_field", "my_other_type", 2),
			},
		},
		{
			name: "my_message",
			fields: []*proto.NormalField{
				protoutil.NewField("my_field", "my_type", 1),
				protoutil.NewField("my_other_field", "my_other_type", 2),
			},
			enums: []*proto.Enum{protoutil.NewEnum("my_enum")},
			options: []*proto.Option{
				protoutil.NewOption("my_option", "with a great value"),
				protoutil.NewOption("gogoproto.nullable", "false", protoutil.Custom(), protoutil.SetField("set")),
			},
			isExtend: true,
		},
	}

	for _, test := range cases {
		var opts []protoutil.MessageSpecOptions
		opts = append(opts, protoutil.WithFields(test.fields...))
		opts = append(opts, protoutil.WithEnums(test.enums...))
		opts = append(opts, protoutil.WithMessageOptions(test.options...))
		if test.isExtend {
			opts = append(opts, protoutil.Extend())
		}
		message := protoutil.NewMessage(test.name, opts...)

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
func TestCreateEnumField(t *testing.T) {
	cases := []struct {
		name    string
		value   int
		options []*proto.Option
	}{
		{
			name:  "my_field",
			value: 1,
		},
		{
			name:  "my_field",
			value: 2,
			options: []*proto.Option{
				protoutil.NewOption("my_option", "with a great value"),
				protoutil.NewOption("gogoproto.nullable", "false", protoutil.Custom(), protoutil.SetField("set")),
			},
		},
	}

	for _, test := range cases {
		var opts []protoutil.EnumFieldSpecOptions
		opts = append(opts, protoutil.WithEnumFieldOptions(test.options...))
		field := protoutil.NewEnumField(test.name, test.value, opts...)

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
func TestCreateEnum(t *testing.T) {
	cases := []struct {
		name    string
		options []*proto.Option
		values  []*proto.EnumField
	}{
		{
			name: "my_enum",
			values: []*proto.EnumField{
				protoutil.NewEnumField("my_value", 1),
				protoutil.NewEnumField("my_other_value", 2),
			},
		},
		{
			name: "my_enum",
			values: []*proto.EnumField{
				protoutil.NewEnumField("my_value", 1),
				protoutil.NewEnumField("my_other_value", 2),
			},
			options: []*proto.Option{
				protoutil.NewOption("my_option", "with a great value"),
			},
		},
	}

	for _, test := range cases {
		var opts []protoutil.EnumSpecOpts
		opts = append(opts, protoutil.WithEnumFields(test.values...))
		opts = append(opts, protoutil.WithEnumOptions(test.options...))
		enum := protoutil.NewEnum(test.name, opts...)

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
func TestCreateOneofField(t *testing.T) {
	cases := []struct {
		name, typeName string
		sequence       int
		options        []*proto.Option
	}{
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
				protoutil.NewOption("my_option", "5"),
				protoutil.NewOption("gogoproto.nullable", "false", protoutil.Custom(), protoutil.SetField("set")),
			},
		},
	}

	for _, test := range cases {
		opts := []protoutil.OneofFieldOptions{protoutil.WithOneofFieldOptions(test.options...)}
		field := protoutil.NewOneofField(test.name, test.typeName, test.sequence, opts...)

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
func TestCreateOneof(t *testing.T) {
	cases := []struct {
		name    string
		options []*proto.Option
		values  []*proto.OneOfField
	}{
		{
			name: "oneof_this",
			values: []*proto.OneOfField{
				protoutil.NewOneofField("my_value", "my_type", 1),
				protoutil.NewOneofField("my_other_value", "my_type", 2),
			},
		},
		{
			name: "oneof_that",
			values: []*proto.OneOfField{
				protoutil.NewOneofField("my_value", "my_type", 1),
			},
			options: []*proto.Option{
				protoutil.NewOption("my_option", "with a great value"),
			},
		},
	}

	for _, test := range cases {
		var opts []protoutil.OneofSpecOptions
		opts = append(opts, protoutil.WithOneofFields(test.values...))
		opts = append(opts, protoutil.WithOneofOptions(test.options...))
		oneof := protoutil.NewOneof(test.name, opts...)

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
	msg := protoutil.NewMessage("my_message")
	protoutil.AttachComment(msg, "my comment")
	require.Equal(t, " my comment", msg.Comment.Lines[0], "expected %v, got %v", "my comment", msg.Comment.Lines[0])

	// Attach comment to rpc call
	rpc := protoutil.NewRPC("my_rpc", "my_request", "my_response")
	protoutil.AttachComment(rpc, "my comment")
	require.Equal(t, " my comment", rpc.Comment.Lines[0], "expected %v, got %v", "my comment", rpc.Comment.Lines[0])

	// Attach comment to a service
	svc := protoutil.NewService("my_service")
	protoutil.AttachComment(svc, "my comment")
	require.Equal(t, " my comment", svc.Comment.Lines[0], "expected %v, got %v", "my comment", svc.Comment.Lines[0])
}

// Test literal creation (indirectly tests the isString function.)
func TestIsString(t *testing.T) {
	require.True(t, protoutil.NewLiteral("string").IsString)
	require.True(t, protoutil.NewLiteral("THIS/PATH/IS/STRING").IsString)

	// Don't report "true" and "false" as strings
	require.False(t, protoutil.NewLiteral("true").IsString)
	require.False(t, protoutil.NewLiteral("false").IsString)

	// Don't report numbers as strings
	require.False(t, protoutil.NewLiteral("1").IsString)
	require.False(t, protoutil.NewLiteral("1.0").IsString)
	require.False(t, protoutil.NewLiteral("1.0e-10").IsString)
	require.False(t, protoutil.NewLiteral("1.0e+10").IsString)
	require.False(t, protoutil.NewLiteral("1.0e10").IsString)
	require.False(t, protoutil.NewLiteral("3.1929348317293483e-10").IsString)

	// A single numbers means not a string, parser would fail with that either way.
	require.True(t, protoutil.NewLiteral("isthisastringohnoitactuallyisn't1.0").IsString)
}
