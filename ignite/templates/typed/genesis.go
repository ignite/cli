package typed

import (
	"context"
	"fmt"
	"strings"

	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/protoanalysis"
)

// ProtoGenesisStateMessage is the name of the proto message that represents the genesis state
const ProtoGenesisStateMessage = "GenesisState"

// PatchGenesisTypeImport patches types/genesis.go content from the issue:
// https://github.com/ignite/cli/issues/992
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

// GenesisStateHighestFieldNumber returns the highest field number in the genesis state proto message
// This allows to determine next the field numbers
func GenesisStateHighestFieldNumber(path string) (int, error) {
	pkgs, err := protoanalysis.Parse(context.Background(), nil, path)
	if err != nil {
		return 0, err
	}
	if len(pkgs) == 0 {
		return 0, fmt.Errorf("%s is not a proto file", path)
	}
	m, err := pkgs[0].MessageByName(ProtoGenesisStateMessage)
	if err != nil {
		return 0, err
	}

	return m.HighestFieldNumber, nil
}
