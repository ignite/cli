package doctor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"

	"github.com/gobuffalo/genny/v2"

	chainconfig "github.com/ignite/cli/ignite/config/chain"
	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/cosmosgen"
	"github.com/ignite/cli/ignite/pkg/events"
	"github.com/ignite/cli/ignite/pkg/goanalysis"
	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/xast"
	"github.com/ignite/cli/ignite/templates/app"
)

const (
	// ToolsFile defines the app relative path to the Go tools file.
	ToolsFile = "tools/tools.go"
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

// MigrateConfig migrates the chain config if required.
func (d *Doctor) MigrateConfig(_ context.Context) error {
	errf := func(err error) error {
		return fmt.Errorf("doctor migrate config: %w", err)
	}

	d.ev.Send("Checking chain config file:", events.ProgressFinish())

	configPath, err := chainconfig.LocateDefault(".")
	if err != nil {
		return errf(err)
	}

	f, err := os.Open(configPath)
	if err != nil {
		return errf(err)
	}
	defer f.Close()

	version, err := chainconfig.ReadConfigVersion(f)
	if err != nil {
		return errf(err)
	}

	status := "OK"

	if version != chainconfig.LatestVersion {
		f.Seek(0, 0)

		// migrate config file
		// Convert the current config to the latest version and update the YAML file
		var buf bytes.Buffer
		if err := chainconfig.MigrateLatest(f, &buf); err != nil {
			return errf(err)
		}

		if err := os.WriteFile(configPath, buf.Bytes(), 0o755); err != nil {
			return errf(fmt.Errorf("config file migration failed: %w", err))
		}

		status = "migrated"
	}

	d.ev.Send(
		fmt.Sprintf("config file %s", colors.Success(status)),
		events.Icon(icons.OK),
		events.ProgressFinish(),
	)

	return nil
}

// FixDependencyTools ensures that:
// - tools/tools.go is present and populated properly
// - dependency tools are installed.
func (d *Doctor) FixDependencyTools(ctx context.Context) error {
	errf := func(err error) error {
		return fmt.Errorf("doctor fix dependency tools: %w", err)
	}

	d.ev.Send("Checking dependency tools:", events.ProgressFinish())

	_, err := os.Stat(ToolsFile)

	switch {
	case err == nil:
		d.ev.Send(
			fmt.Sprintf("%s %s", ToolsFile, colors.Success("exists")),
			events.Icon(icons.OK),
			events.ProgressUpdate(),
		)

		updated, err := d.ensureDependencyImports(ToolsFile)
		if err != nil {
			return errf(err)
		}

		status := "OK"
		if updated {
			status = "updated"
		}

		d.ev.Send(
			fmt.Sprintf("tools file %s", colors.Success(status)),
			events.Icon(icons.OK),
			events.ProgressFinish(),
		)

	case os.IsNotExist(err):
		if err := d.createToolsFile(ctx, ToolsFile); err != nil {
			return errf(err)
		}

	default:
		return errf(err)
	}

	return nil
}

func (d Doctor) createToolsFile(ctx context.Context, toolsFilename string) error {
	pathInfo, err := gomodulepath.ParseAt(".")
	if err != nil {
		return err
	}

	g, err := app.NewGenerator(&app.Options{
		ModulePath:       pathInfo.RawPath,
		AppName:          pathInfo.Package,
		BinaryNamePrefix: pathInfo.Root,
		IncludePrefixes:  []string{toolsFilename},
	})
	if err != nil {
		return err
	}

	runner := genny.WetRunner(ctx)
	if err := runner.With(g); err != nil {
		return err
	}

	if err := runner.Run(); err != nil {
		return err
	}

	d.ev.Send(
		fmt.Sprintf("%s %s", toolsFilename, colors.Success("created")),
		events.Icon(icons.OK),
		events.ProgressFinish(),
	)

	d.ev.Send("Installing dependency tools", events.ProgressStart())
	if err := cosmosgen.InstallDepTools(ctx, "."); err != nil {
		return err
	}

	for _, dep := range cosmosgen.DepTools() {
		d.ev.Send(
			fmt.Sprintf("%s %s", path.Base(dep), colors.Success("installed")),
			events.Icon(icons.OK),
			events.ProgressFinish(),
		)
	}

	return nil
}

func (d Doctor) ensureDependencyImports(toolsFilename string) (bool, error) {
	d.ev.Send("Ensuring required tools imports", events.ProgressStart())

	f, _, err := xast.ParseFile(toolsFilename)
	if err != nil {
		return false, err
	}

	var (
		buf     bytes.Buffer
		missing = cosmosgen.MissingTools(f)
		unused  = cosmosgen.UnusedTools(f)
	)

	// Check if the tools file should be fixed
	if len(missing) == 0 && len(unused) == 0 {
		return false, nil
	}

	err = goanalysis.UpdateInitImports(f, &buf, missing, unused)
	if err != nil {
		return false, err
	}

	err = os.WriteFile(toolsFilename, buf.Bytes(), 0o644)
	if err != nil {
		return false, err
	}

	return true, nil
}
