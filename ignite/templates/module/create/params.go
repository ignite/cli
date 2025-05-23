package modulecreate

import (
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
)

// NewModuleParam returns the generator to scaffold a new parameter inside a module.
func NewModuleParam(opts ParamsOptions) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(paramsProtoModify(opts))
	g.RunFn(paramsTypesModify(opts))
	return g, nil
}

func paramsProtoModify(opts ParamsOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := opts.ProtoFile("params.proto")
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
			_, err := protoutil.GetFieldByName(params, paramField.ProtoFieldName())
			if err == nil {
				return errors.Errorf("duplicate field %s in %s", paramField.ProtoFieldName(), params.Name)
			}

			param := protoutil.NewField(
				paramField.ProtoFieldName(),
				paramField.DataType(),
				protoutil.NextUniqueID(params),
			)
			protoutil.Append(params, param)
		}

		newFile := genny.NewFileS(path, protoutil.Print(protoFile))
		return r.File(newFile)
	}
}

func paramsTypesModify(opts ParamsOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", opts.ModuleName, "types/params.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		var (
			content               = f.String()
			globalOpts            = make([]xast.GlobalOptions, len(opts.Params))
			newParamsModifier     = make([]xast.FunctionOptions, 0)
			defaultParamsModifier = make([]xast.FunctionOptions, len(opts.Params))
			validateModifier      = make([]xast.FunctionOptions, len(opts.Params))
		)
		for i, param := range opts.Params {
			// param key and default value.
			globalOpts[i] = xast.WithGlobal(
				fmt.Sprintf("Default%s", param.Name.UpperCamel),
				param.DataType(),
				param.Value(),
			)

			// add parameter to the struct into the new method.
			newParamsModifier = append(
				newParamsModifier,
				xast.AppendFuncParams(param.ProtoFieldName(), param.DataType(), -1),
				xast.AppendFuncStruct(
					"Params",
					param.Name.UpperCamel,
					param.ProtoFieldName(),
				),
			)

			// add default parameter.
			defaultParamsModifier[i] = xast.AppendInsideFuncCall(
				"NewParams",
				fmt.Sprintf("Default%s", param.Name.UpperCamel),
				-1,
			)

			// add param field to the validate method.
			replacementValidate := fmt.Sprintf(
				`if err := validate%[1]v(p.%[1]v); err != nil { return err }`,
				param.Name.UpperCamel,
			)
			validateModifier[i] = xast.AppendFuncCode(replacementValidate)

			// add param field to the validate method.
			templateValidation := `// validate%[1]v validates the %[1]v parameter.
func validate%[1]v(v %[2]v) error {
	// TODO implement validation
	return nil
}`
			validationFunc := fmt.Sprintf(
				templateValidation,
				param.Name.UpperCamel,
				param.DataType(),
			)
			content, err = xast.AppendFunction(content, validationFunc)
			if err != nil {
				return err
			}
		}

		content, err = xast.InsertGlobal(content, xast.GlobalTypeConst, globalOpts...)
		if err != nil {
			return err
		}

		content, err = xast.ModifyFunction(content, "NewParams", newParamsModifier...)
		if err != nil {
			return err
		}

		content, err = xast.ModifyFunction(content, "DefaultParams", defaultParamsModifier...)
		if err != nil {
			return err
		}

		content, err = xast.ModifyFunction(content, "Validate", validateModifier...)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
