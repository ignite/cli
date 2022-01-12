package tsc

import (
	"context"
	"os"
	"path/filepath"

	"github.com/imdario/mergo"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/exec"
	"github.com/tendermint/starport/starport/pkg/confile"
	"github.com/tendermint/starport/starport/pkg/nodetime"
)

const nodeModulesPath = "/snapshot/gen-nodetime/node_modules"

var (
	defaultConfig = func() Config {
		return Config{
			CompilerOptions: CompilerOptions{
				BaseURL:          nodeModulesPath,
				ModuleResolution: "node",
				Target:           "es2020",
				Module:           "es2020",
				TypeRoots:        []string{filepath.Join(nodeModulesPath, "@types")},
				SkipLibCheck:     true,
			},
		}
	}
	tsconfig = func(dir string) string { return filepath.Join(dir, "tsconfig.json") }
)

// Config represents tsconfig.json.
type Config struct {
	Include         []string        `json:"include"`
	CompilerOptions CompilerOptions `json:"compilerOptions"`
}

// CompilerOptions section of tsconfig.json.
type CompilerOptions struct {
	BaseURL          string   `json:"baseUrl"`
	ModuleResolution string   `json:"moduleResolution"`
	Target           string   `json:"target"`
	Module           string   `json:"module"`
	TypeRoots        []string `json:"typeRoots"`
	Declaration      bool     `json:"declaration"`
	SkipLibCheck     bool     `json:"skipLibCheck"`
}

// Generate transpiles TS into JS by given TS config.
func Generate(ctx context.Context, config Config) error {
	command, cleanup, err := nodetime.Command(nodetime.CommandTSC)
	if err != nil {
		return err
	}
	defer cleanup()

	dconfig := defaultConfig()

	if err := mergo.Merge(&dconfig, config, mergo.WithOverride); err != nil {
		return err
	}

	// save the config into a temp file in the fs.
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	path := tsconfig(dir)

	if err := confile.
		New(confile.DefaultJSONEncodingCreator, path).
		Save(dconfig); err != nil {
		return err
	}

	// command constructs the tsc command.
	command = append(command, []string{
		"-b",
		path,
	}...)

	// execute the command.
	return exec.Exec(ctx, command, exec.IncludeStdLogsToError())
}
