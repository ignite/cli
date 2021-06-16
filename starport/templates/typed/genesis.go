package typed

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/placeholder"
)

func (t *typedStargate) genesisModify(replacer placeholder.Replacer, opts *Options, g *genny.Generator) {
	g.RunFn(t.genesisProtoModify(replacer, opts))
	g.RunFn(t.genesisTypesModify(replacer, opts))
	g.RunFn(t.genesisModuleModify(replacer, opts))
}

func (t *typedStargate) genesisProtoModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/genesis.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateProtoImport := `%[1]v
import "%[2]v/%[3]v.proto";`
		replacementProtoImport := fmt.Sprintf(
			templateProtoImport,
			PlaceholderGenesisProtoImport,
			opts.ModuleName,
			opts.TypeName.LowerCamel,
		)
		content := replacer.Replace(f.String(), PlaceholderGenesisProtoImport, replacementProtoImport)

		// Determine the new field number
		fieldNumber := strings.Count(content, PlaceholderGenesisProtoStateField) + 1

		templateProtoState := `%[1]v
		repeated %[2]v %[3]vList = %[4]v; %[5]v
		uint64 %[3]vCount = %[6]v; %[5]v`
		replacementProtoState := fmt.Sprintf(
			templateProtoState,
			PlaceholderGenesisProtoState,
			opts.TypeName.UpperCamel,
			opts.TypeName.LowerCamel,
			fieldNumber,
			PlaceholderGenesisProtoStateField,
			fieldNumber+1,
		)
		content = replacer.Replace(content, PlaceholderGenesisProtoState, replacementProtoState)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) genesisTypesModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/genesis.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content := PatchGenesisTypeImport(replacer, f.String())

		templateTypesImport := `"fmt"`
		content = replacer.ReplaceOnce(content, PlaceholderGenesisTypesImport, templateTypesImport)

		templateTypesDefault := `%[1]v
%[2]vList: []*%[2]v{},`
		replacementTypesDefault := fmt.Sprintf(
			templateTypesDefault,
			PlaceholderGenesisTypesDefault,
			opts.TypeName.UpperCamel,
		)
		content = replacer.Replace(content, PlaceholderGenesisTypesDefault, replacementTypesDefault)

		templateTypesValidate := `%[1]v
// Check for duplicated ID in %[2]v
%[2]vIdMap := make(map[uint64]bool)

for _, elem := range gs.%[3]vList {
	if _, ok := %[2]vIdMap[elem.Id]; ok {
		return fmt.Errorf("duplicated id for %[2]v")
	}
	%[2]vIdMap[elem.Id] = true
}`
		replacementTypesValidate := fmt.Sprintf(
			templateTypesValidate,
			PlaceholderGenesisTypesValidate,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
		)
		content = replacer.Replace(content, PlaceholderGenesisTypesValidate, replacementTypesValidate)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) genesisModuleModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/genesis.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateModuleInit := `%[1]v
// Set all the %[2]v
for _, elem := range genState.%[3]vList {
	k.Set%[3]v(ctx, *elem)
}

// Set %[2]v count
k.Set%[3]vCount(ctx, genState.%[3]vCount)
`
		replacementModuleInit := fmt.Sprintf(
			templateModuleInit,
			PlaceholderGenesisModuleInit,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
		)
		content := replacer.Replace(f.String(), PlaceholderGenesisModuleInit, replacementModuleInit)

		templateModuleExport := `%[1]v
// Get all %[2]v
%[2]vList := k.GetAll%[3]v(ctx)
for _, elem := range %[2]vList {
	elem := elem
	genesis.%[3]vList = append(genesis.%[3]vList, &elem)
}

// Set the current count
genesis.%[3]vCount = k.Get%[3]vCount(ctx)
`
		replacementModuleExport := fmt.Sprintf(
			templateModuleExport,
			PlaceholderGenesisModuleExport,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
		)
		content = replacer.Replace(content, PlaceholderGenesisModuleExport, replacementModuleExport)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// PatchGenesisTypeImport patches types/genesis.go content from the issue:
// https://github.com/tendermint/starport/issues/992
func PatchGenesisTypeImport(replacer placeholder.Replacer, content string) string {
	patternToCheck := "import ("
	replacement := fmt.Sprintf(`import (
%[1]v
)`, PlaceholderGenesisTypesImport)

	if !strings.Contains(content, patternToCheck) {
		content = replacer.Replace(content, PlaceholderGenesisTypesImport, replacement)
	}

	return content
}
