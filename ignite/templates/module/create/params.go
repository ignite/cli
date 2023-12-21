package modulecreate

import (
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/pkg/protoanalysis/protoutil"
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
		for _, param := range opts.Params {
			// param key and default value.
			templateVars := `var (
	// Key%[2]v represents the %[2]v parameter.
	Key%[2]v = []byte("%[2]v")
	// Default%[2]v represents the %[2]v default value.
	// TODO: Determine the default value.
	Default%[2]v %[3]v = %[4]v
)

%[1]v`
			replacementVars := fmt.Sprintf(
				templateVars,
				module.PlaceholderParamsVars,
				param.Name.UpperCamel,
				param.DataType(),
				param.Value(),
			)
			content = replacer.Replace(content, module.PlaceholderParamsVars, replacementVars)

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
func validate%[2]v(v interface{}) error {
	%[3]v, ok := v.(%[4]v)
	if !ok {
		return fmt.Errorf("invalid parameter type: %%T", v)
	}

	// TODO implement validation
	_ = %[3]v

	return nil
}

%[1]v`
			replacementValidation := fmt.Sprintf(
				templateValidation,
				module.PlaceholderParamsValidation,
				param.Name.UpperCamel,
				param.Name.LowerCamel,
				param.DataType(),
			)
			content = replacer.Replace(content, module.PlaceholderParamsValidation, replacementValidation)

		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
