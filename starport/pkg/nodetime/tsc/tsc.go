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

var (
	defaultConfig = func() Config {
		return Config{
			CompilerOptions: CompilerOptions{
				Target:    "es2020",
				Module:    "es6",
				TypeRoots: []string{"/snapshot/gen-nodetime/node_modules/@types"},
				Types:     []string{"node"},
				Paths: map[string][]string{
					"*":    []string{"/snapshot/gen-nodetime/node_modules/*"},
					"long": []string{"/snapshot/gen-nodetime/node_modules/long/index.js"},
				},
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
	Declaration bool                `json:"declaration"`
	Paths       map[string][]string `json:"paths"`
	Target      string              `json:"target"`
	Module      string              `json:"module"`
	TypeRoots   []string            `json:"typeRoots"`
	Types       []string            `json:"types"`
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
