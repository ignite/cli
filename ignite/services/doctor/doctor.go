package doctor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"

	"golang.org/x/mod/modfile"

	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosbuf"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/events"
	"github.com/ignite/cli/v29/ignite/pkg/goanalysis"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/templates/app"
)

// DONTCOVER: Doctor read and write the filesystem intensively, so it's better
// to rely on integration tests only. See integration/doctor package.
type Doctor struct {
	ev events.Bus
}

// New returns a new doctor.
func New(opts ...Option) *Doctor {
	d := &Doctor{}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

type Option func(*Doctor)

// CollectEvents sets doctor event bus.
func CollectEvents(ev events.Bus) Option {
	return func(d *Doctor) {
		d.ev = ev
	}
}

// MigrateBufConfig migrates the buf chain config if required.
func (d *Doctor) MigrateBufConfig(ctx context.Context, cacheStorage cache.Storage, appPath, configPath string) error {
	errf := func(err error) error {
		return errors.Errorf("doctor migrate buf config: %w", err)
	}

	d.ev.Send("Checking buf config file version:")

	// Check if the appPath contains the buf.work.yaml file in the root folder.
	// The buf.work.yaml file does not exist in buf v2 config, so it is a good
	// indicator that the buf config is already migrated.
	bufWorkFile := path.Join(appPath, "buf.work.yaml")
	if _, err := os.Stat(bufWorkFile); os.IsNotExist(err) {
		d.ev.Send(
			fmt.Sprintf("buf files %s", colors.Success("OK")),
			events.Icon(icons.OK),
			events.Indent(1),
		)
		return nil
	} else if err != nil {
		return errf(errors.Errorf("unable to check buf files have been migrated: %w", err))
	}

	d.ev.Send("Migrating buf config file to v2")

	configFile, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer configFile.Close()

	protoPath, err := chainconfig.ReadProtoPath(configFile)
	if err != nil {
		return errf(err)
	}

	b, err := cosmosbuf.New(cacheStorage, appPath)
	if err != nil {
		return errf(err)
	}

	if err := b.Migrate(ctx, protoPath); err != nil {
		return errf(err)
	}

	runner := xgenny.NewRunner(ctx, appPath)
	_, err = boxBufFiles(runner, protoPath)
	if err != nil {
		return err
	}

	d.ev.Send(
		"buf config files migrated",
		events.Icon(icons.OK),
		events.Indent(1),
		events.ProgressFinish(),
	)

	return nil
}

// BoxBufFiles box all buf files.
func boxBufFiles(runner *xgenny.Runner, protoDir string) (xgenny.SourceModification, error) {
	g, err := app.NewBufGenerator(protoDir)
	if err != nil {
		return xgenny.SourceModification{}, err
	}
	return runner.RunAndApply(g)
}

// MigrateChainConfig migrates the chain config if required.
func (d *Doctor) MigrateChainConfig(configPath string) error {
	errf := func(err error) error {
		return errors.Errorf("doctor migrate config: %w", err)
	}

	d.ev.Send("Checking chain config file:")
	configFile, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer configFile.Close()

	version, err := chainconfig.ReadConfigVersion(configFile)
	if err != nil {
		return errf(err)
	}

	status := "OK"
	if version != chainconfig.LatestVersion {
		_, err := configFile.Seek(0, 0)
		if err != nil {
			return errf(errors.Errorf("failed to reset the file: %w", err))
		}
		// migrate config file
		// Convert the current config to the latest version and update the YAML file
		var buf bytes.Buffer
		if err := chainconfig.MigrateLatest(configFile, &buf); err != nil {
			return errf(err)
		}

		if err := os.WriteFile(configPath, buf.Bytes(), 0o600); err != nil {
			return errf(errors.Errorf("config file migration failed: %w", err))
		}

		status = "migrated"
	}

	d.ev.Send(
		fmt.Sprintf("config file %s", colors.Success(status)),
		events.Icon(icons.OK),
		events.Indent(1),
		events.ProgressFinish(),
	)

	return nil
}

// MigrateToolsGo ensures that.
// - go.mod is bumped to go 1.25.
// - removes tools.go file from chain.
// - add all tools to go.mod.
func (d *Doctor) MigrateToolsGo(appPath string) error {
	errf := func(err error) error {
		return errors.Errorf("doctor migrate tools.go: %w", err)
	}

	const (
		// toolsFile defines the app relative path to the Go tools file.
		toolsFile = "tools/tools.go"
		// goModFile defines the app relative path to the Go module file.
		goModFile = "go.mod"
	)

	_, err := os.Stat(toolsFile)
	if os.IsNotExist(err) { // file doesn't exist, nothing to do
		return nil
	}

	d.ev.Send("Migrating dependency tools:")

	toolsAst, _, err := xast.ParseFile(toolsFile)
	if err != nil {
		return errf(errors.Errorf("failed to parse tools.go file: %w", err))
	}

	imports := goanalysis.FormatImports(toolsAst)
	if len(imports) == 0 {
		d.ev.Send(
			"no tools to migrate",
			events.Icon(icons.OK),
			events.Indent(1),
			events.ProgressFinish(),
		)
		return nil
	}

	goModPath := path.Join(appPath, goModFile)
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return errf(errors.Errorf("failed to read go.mod file: %w", err))
	}

	goModAst, err := modfile.Parse(goModPath, data, nil)
	if err != nil {
		return errf(errors.Errorf("failed to parse go.mod file: %w", err))
	}

	// bump to go 1.25
	if goModAst.Go.Version < "1.24" {
		goModAst.Go.Version = "1.25"
	}

	for _, imp := range imports {
		_ = goModAst.AddTool(imp)
	}

	// remove the tools.go file
	if err := os.Remove(toolsFile); err != nil {
		return errf(errors.Errorf("failed to remove tools.go file: %w", err))
	}

	// write the updated go.mod file
	data, err = goModAst.Format()
	if err != nil {
		return errf(errors.Errorf("failed to format go.mod file: %w", err))
	}

	if err := os.WriteFile(goModPath, data, 0o600); err != nil {
		return errf(errors.Errorf("failed to write go.mod file: %w", err))
	}
	d.ev.Send(
		fmt.Sprintf("tools migrated to %s", colors.Success(goModFile)),
		events.Icon(icons.OK),
		events.Indent(1),
		events.ProgressFinish(),
	)

	return nil
}
