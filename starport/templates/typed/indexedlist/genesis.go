package indexedlist

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/templates/typed"
)

func genesisModify(replacer placeholder.Replacer, opts *typed.Options, g *genny.Generator) {
	g.RunFn(genesisProtoModify(replacer, opts))
	g.RunFn(genesisTypesModify(replacer, opts))
	g.RunFn(genesisModuleModify(replacer, opts))
}

func genesisProtoModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "proto", opts.ModuleName, "genesis.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateProtoImport := `import "%[2]v/%[3]v.proto";
%[1]v`
		replacementProtoImport := fmt.Sprintf(
			templateProtoImport,
			typed.PlaceholderGenesisProtoImport,
			opts.ModuleName,
			opts.TypeName.Snake,
		)
		content := replacer.Replace(f.String(), typed.PlaceholderGenesisProtoImport, replacementProtoImport)

		// Add gogo.proto
		replacementGogoImport := typed.EnsureGogoProtoImported(path, typed.PlaceholderGenesisProtoImport)
		content = replacer.Replace(content, typed.PlaceholderGenesisProtoImport, replacementGogoImport)

		// Parse proto file to determine the field numbers
		highestNumber, err := typed.GenesisStateHighestFieldNumber(path)
		if err != nil {
			return err
		}

		templateProtoState := `repeated %[2]v %[3]vList = %[4]v [(gogoproto.nullable) = false];
  repeated %[2]vCount %[3]vCountList = %[5]v [(gogoproto.nullable) = false];
  %[1]v`
		replacementProtoState := fmt.Sprintf(
			templateProtoState,
			typed.PlaceholderGenesisProtoState,
			opts.TypeName.UpperCamel,
			opts.TypeName.LowerCamel,
			highestNumber+1,
			highestNumber+2,
		)
		content = replacer.Replace(content, typed.PlaceholderGenesisProtoState, replacementProtoState)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisTypesModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/genesis.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content := typed.PatchGenesisTypeImport(replacer, f.String())

		templateTypesImport := `"fmt"`
		content = replacer.ReplaceOnce(content, typed.PlaceholderGenesisTypesImport, templateTypesImport)

		templateTypesDefault := `%[2]vList: []%[2]v{},
%[2]vCountList: []%[2]vCount{},
%[1]v`
		replacementTypesDefault := fmt.Sprintf(
			templateTypesDefault,
			typed.PlaceholderGenesisTypesDefault,
			opts.TypeName.UpperCamel,
		)
		content = replacer.Replace(content, typed.PlaceholderGenesisTypesDefault, replacementTypesDefault)

		// Array of the indexes of the element
		var indexArgs []string
		for _, index := range opts.Indexes {
			indexArgs = append(indexArgs, "elem."+index.Name.UpperCamel)
		}
		indexCall := strings.Join(indexArgs, ",\n")
		indexCall += ",\n"

		templateTypesValidate := `
// Checkout %[2]vCounts to perform verification
%[2]vCountMap := make(map[string]uint64)
for _, elem := range gs.%[3]vCountList {
	countKey := string(All%[3]vKeyPath(
		%[4]v))
	if _, ok := %[2]vCountMap[countKey]; ok {
		return fmt.Errorf("duplicated %[2]v count")
	}
	%[2]vCountMap[countKey] = elem.Count
}

// Check for duplicated ID in %[2]v
%[2]vIdMap := make(map[string]struct{})
for _, elem := range gs.%[3]vList {
	elemKey := string(%[3]vKeyPath(
		%[4]v elem.Id,
	))
	if _, ok := %[2]vIdMap[elemKey]; ok {
		return fmt.Errorf("duplicated id for %[2]v")
	}
	countKey := string(All%[3]vKeyPath(
		%[4]v))
	if elem.Id >= %[2]vCountMap[countKey] {
		return fmt.Errorf("%[2]v id should be lower or equal than the biggest reported id")
	}
	%[2]vIdMap[elemKey] = struct{}{}
}
%[1]v`
		replacementTypesValidate := fmt.Sprintf(
			templateTypesValidate,
			typed.PlaceholderGenesisTypesValidate,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
			indexCall,
		)
		content = replacer.Replace(content, typed.PlaceholderGenesisTypesValidate, replacementTypesValidate)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisModuleModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "genesis.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Array of the indexes as field
		var indexField []string
		for _, index := range opts.Indexes {
			indexField = append(indexField, fmt.Sprintf("%[1]v: elem.%[1]v", index.Name.UpperCamel))
		}
		indexFieldStr := strings.Join(indexField, ",\n")
		indexFieldStr += ",\n"

		templateModuleInit := `// Set all the %[2]v
for _, elem := range genState.%[3]vList {
	k.Set%[3]v(ctx, elem)
}

// Set %[2]v count
for _, elem := range genState.%[3]vCountList {
	count := types.%[3]vCount{
		%[4]v  Count: elem.Count,
	}
	k.Set%[3]vCount(ctx, count)
}
%[1]v`
		replacementModuleInit := fmt.Sprintf(
			templateModuleInit,
			typed.PlaceholderGenesisModuleInit,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
			indexFieldStr,
		)
		content := replacer.Replace(f.String(), typed.PlaceholderGenesisModuleInit, replacementModuleInit)

		templateModuleExport := `genesis.%[2]vList = k.GetAll%[2]v(ctx)
genesis.%[2]vCountList = k.GetAll%[2]vCount(ctx)
%[1]v`
		replacementModuleExport := fmt.Sprintf(
			templateModuleExport,
			typed.PlaceholderGenesisModuleExport,
			opts.TypeName.UpperCamel,
			opts.TypeName.LowerCamel,
		)
		content = replacer.Replace(content, typed.PlaceholderGenesisModuleExport, replacementModuleExport)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
