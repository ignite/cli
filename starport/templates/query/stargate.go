package query

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny"
)

// NewStargate returns the generator to scaffold a empty query in a Stargate module
func NewStargate(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(protoQueryModify(opts))
	g.RunFn(cliQueryModify(opts))

	return g, Box(stargateTemplate, opts, g)
}

func protoQueryModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/query.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// RPC service
		templateRPC := `%[1]v
	rpc %[2]v(Query%[2]vRequest) returns (Query%[2]vResponse) {
		option (google.api.http).get = "/%[4]v/%[5]v/%[6]v/%[3]v";
	}
`
		replacementRPC := fmt.Sprintf(
			templateRPC,
			Placeholder2,
			strings.Title(opts.QueryName),
			opts.QueryName,
			opts.OwnerName,
			opts.AppName,
			opts.ModuleName,
		)
		content := strings.Replace(f.String(), Placeholder2, replacementRPC, 1)

		// Fields for request and response
		var reqFields string
		for i, field := range opts.ReqFields {
			reqFields += fmt.Sprintf("  %s %s = %d;\n", field.Datatype, field.Name, i+1)
		}
		var resFields string
		for i, field := range opts.ResFields {
			resFields += fmt.Sprintf("  %s %s = %d;\n", field.Datatype, field.Name, i+1)
		}

		// Messages
		templateMessages := `%[1]v
message Query%[2]vRequest {
%[3]v}

message Query%[2]vResponse {
%[4]v}
`
		replacementMessages := fmt.Sprintf(
			templateMessages,
			Placeholder3,
			strings.Title(opts.QueryName),
			reqFields,
			resFields,
		)
		content = strings.Replace(content, Placeholder3, replacementMessages, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func cliQueryModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/cli/query.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		template := `%[1]v

	cmd.AddCommand(Cmd%[2]v())
`
		replacement := fmt.Sprintf(
			template,
			Placeholder,
			strings.Title(opts.QueryName),
		)
		content := strings.Replace(f.String(), Placeholder, replacement, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
