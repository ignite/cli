package modulecreate

import (
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/v28/ignite/pkg/xast"
	"github.com/ignite/cli/v28/ignite/templates/module"
)

// NewModuleParam returns the generator to scaffold a new parameter inside a module.
func NewModuleParam(replacer placeholder.Replacer, opts ParamsOptions) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(paramsProtoModify(opts))
	g.RunFn(paramsTypesModify(replacer, opts))
	return g, nil
}

func paramsProtoModify(opts ParamsOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "proto", opts.AppName, opts.ModuleName, "params.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		protoFile, err := protoutil.ParseProtoFile(f)
		if err != nil {
			return err
		}

		params, err := protoutil.GetMessageByName(protoFile, "Params")
		if err != nil {
			return errors.Errorf("couldn't find message 'Params' in %s: %w", path, err)
		}
		for _, paramField := range opts.Params {
			param := protoutil.NewField(
				paramField.Name.LowerCamel,
				paramField.DataType(),
				protoutil.NextUniqueID(params),
			)
			protoutil.Append(params, param)
		}

		newFile := genny.NewFileS(path, protoutil.Print(protoFile))
		return r.File(newFile)
	}
}

func paramsTypesModify(replacer placeholder.Replacer, opts ParamsOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/params.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content := f.String()

		globalOpts := make([]xast.GlobalOptions, len(opts.Params))
		modifyNewParams := make([]xast.FunctionOptions, len(opts.Params))
		modifyDefaultParams := make([]xast.FunctionOptions, len(opts.Params))
		modifyValidate := make([]xast.FunctionOptions, len(opts.Params))
		//xast.AppendCode()
		for _, param := range opts.Params {
			// param key and default value.
			globalOpts = append(
				globalOpts,
				xast.WithGlobal(fmt.Sprintf("Default%s", param.Name.UpperCamel), param.DataType(), param.Value()),
			)

			// param key and default value.
			content, err = xast.ModifyFunction(
				content,
				"",
				xast.GlobalTypeConst,
				xast.WithGlobal(fmt.Sprintf("Default%s", param.Name.UpperCamel), param.DataType(), param.Value()),
			)
			if err != nil {
				return err
			}

			// add parameter to the new method.
			templateNewParam := "%[2]v %[3]v,\n%[1]v"
			replacementNewParam := fmt.Sprintf(
				templateNewParam,
				module.PlaceholderParamsNewParam,
				param.Name.LowerCamel,
				param.DataType(),
			)
			content = replacer.Replace(content, module.PlaceholderParamsNewParam, replacementNewParam)

			// add parameter to the struct into the new method.
			templateNewStruct := "%[2]v: %[3]v,\n%[1]v"
			replacementNewStruct := fmt.Sprintf(
				templateNewStruct,
				module.PlaceholderParamsNewStruct,
				param.Name.UpperCamel,
				param.Name.LowerCamel,
			)
			content = replacer.Replace(content, module.PlaceholderParamsNewStruct, replacementNewStruct)

			// add default parameter.
			templateDefault := `Default%[2]v,
%[1]v`
			replacementDefault := fmt.Sprintf(
				templateDefault,
				module.PlaceholderParamsDefault,
				param.Name.UpperCamel,
			)
			content = replacer.Replace(content, module.PlaceholderParamsDefault, replacementDefault)

			// add param field to the validate method.
			templateValidate := `if err := validate%[2]v(p.%[2]v); err != nil {
   		return err
   	}
	%[1]v`
			replacementValidate := fmt.Sprintf(
				templateValidate,
				module.PlaceholderParamsValidate,
				param.Name.UpperCamel,
			)
			content = replacer.Replace(content, module.PlaceholderParamsValidate, replacementValidate)

			// add param field to the validate method.
			templateValidation := `// validate%[2]v validates the %[2]v parameter.
func validate%[2]v(v %[3]v) error {
	// TODO implement validation
	return nil
}

%[1]v`
			replacementValidation := fmt.Sprintf(
				templateValidation,
				module.PlaceholderParamsValidation,
				param.Name.UpperCamel,
				param.DataType(),
			)
			content = replacer.Replace(content, module.PlaceholderParamsValidation, replacementValidation)

		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
