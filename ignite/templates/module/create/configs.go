package modulecreate

import (
	"path/filepath"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
)

// NewModuleConfigs returns the generator to scaffold a new configs inside a module.
func NewModuleConfigs(opts ConfigsOptions) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(configsProtoModify(opts))
	return g, nil
}

func configsProtoModify(opts ConfigsOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		// here we do not use opts.ProtoFile as it will append an extra opts.ProtoVer in the path
		path := filepath.Join(opts.ProtoDir, opts.AppName, opts.ModuleName, "module", opts.ProtoVer, "module.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		protoFile, err := protoutil.ParseProtoFile(f)
		if err != nil {
			return err
		}

		params, err := protoutil.GetMessageByName(protoFile, "Module")
		if err != nil {
			return errors.Errorf("couldn't find message 'Module' in %s: %w", path, err)
		}

		for _, paramField := range opts.Configs {
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
