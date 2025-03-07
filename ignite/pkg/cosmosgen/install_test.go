package cosmosgen_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/mod/modfile"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosgen"
)

func TestMissingTools(t *testing.T) {
	var (
		tools        = cosmosgen.DepTools()
		someTools    = tools[:2]
		missingTools = tools[2:]
	)
	tests := []struct {
		name    string
		modFile *modfile.File
		want    []string
	}{
		{
			name:    "no missing tools",
			modFile: createModFileWithTools(t, tools...),
			want:    nil,
		},
		{
			name:    "some missing tools",
			modFile: createModFileWithTools(t, someTools...),
			want:    missingTools,
		},
		{
			name:    "all tools missing",
			modFile: createModFileWithTools(t),
			want:    tools,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cosmosgen.MissingTools(tt.modFile)
			require.EqualValues(t, tt.want, got)
		})
	}
}

func TestUnusedTools(t *testing.T) {
	tests := []struct {
		name    string
		modFile *modfile.File
		want    []string
	}{
		{
			name: "all unused tools",
			modFile: createModFileWithTools(t,
				"github.com/regen-network/cosmos-proto/protoc-gen-gocosmos",
				"github.com/ignite-hq/cli/ignite/pkg/cmdrunner",
				"github.com/ignite-hq/cli/ignite/pkg/cmdrunner/step",
			),
			want: []string{
				"github.com/regen-network/cosmos-proto/protoc-gen-gocosmos",
				"github.com/ignite-hq/cli/ignite/pkg/cmdrunner",
				"github.com/ignite-hq/cli/ignite/pkg/cmdrunner/step",
			},
		},
		{
			name: "some unused tools",
			modFile: createModFileWithTools(t,
				"github.com/ignite-hq/cli/ignite/pkg/cmdrunner",
			),
			want: []string{"github.com/ignite-hq/cli/ignite/pkg/cmdrunner"},
		},
		{
			name:    "no tools unused",
			modFile: createModFileWithTools(t, ""),
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cosmosgen.UnusedTools(tt.modFile)
			require.EqualValues(t, tt.want, got)
		})
	}
}

// createModFileWithTools helper function to create a modfile.File with given tool paths.
// This simulates the Tool entries in a go.mod file.
func createModFileWithTools(t *testing.T, toolPaths ...string) *modfile.File {
	// create a minimal go.mod content
	content := "module test\n\ngo 1.24\n\n"

	// parse the basic module
	f, err := modfile.Parse("go.mod", []byte(content), nil)
	if err != nil {
		t.Logf("failed to parse test go.mod content: %v", err)
		t.FailNow()
	}

	// add the tools
	for _, path := range toolPaths {
		if err := f.AddTool(path); err != nil {
			t.Logf("failed to add tool %s to go.mod: %v", path, err)
			t.FailNow()
		}
	}

	return f
}
