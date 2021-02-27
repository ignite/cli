package tsc

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/imdario/mergo"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
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

var placeOnce sync.Once

// Generate transpiles TS into JS by given TS config.
func Generate(ctx context.Context, config Config) error {
	var err error

	placeOnce.Do(func() { err = nodetime.PlaceBinary() })

	if err != nil {
		return err
	}

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
	command := []string{
		nodetime.BinaryPath,
		nodetime.CommandTSC,
		"-b",
		path,
	}

	// execute the command.
	return cmdrunner.Exec(ctx, command[0], command[1:]...)
}
