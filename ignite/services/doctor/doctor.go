package doctor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"

	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosbuf"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosgen"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/events"
	"github.com/ignite/cli/v29/ignite/pkg/goanalysis"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
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

	d.ev.Send("Checking buf config file version")

	// Check if the appPath contains the buf.work.yaml file in the root folder.
	bufWorkFile := path.Join(appPath, "buf.work.yaml")
	if _, err := os.Stat(bufWorkFile); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return errf(errors.Errorf("unable to check if buf.work.yaml exists: %w", err))
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

	d.ev.Send(
		"buf config files migrated",
		events.Icon(icons.OK),
		events.Indent(1),
		events.ProgressFinish(),
	)

	return nil
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

// FixDependencyTools ensures that:
// - tools/tools.go is present and populated properly
// - dependency tools are installed.
// Deprecated: This isn't required since Go 1.14.
func (d *Doctor) FixDependencyTools(ctx context.Context) error {
	const toolsFile = "tools/tools.go"

	errf := func(err error) error {
		return errors.Errorf("doctor fix dependency tools: %w", err)
	}

	d.ev.Send("Checking dependency tools:")

	_, err := os.Stat(toolsFile)

	switch {
	case err == nil:
		d.ev.Send(
			fmt.Sprintf("%s %s", toolsFile, colors.Success("exists")),
			events.Icon(icons.OK),
			events.Indent(1),
		)

		updated, err := d.ensureDependencyImports(toolsFile)
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
			events.Indent(1),
			events.ProgressFinish(),
		)

	case os.IsNotExist(err):
		if err := d.createToolsFile(ctx, toolsFile); err != nil {
			return errf(err)
		}

		d.ev.Send(
			fmt.Sprintf("tools file %s", colors.Success("created")),
			events.Icon(icons.OK),
			events.Indent(1),
		)

	default:
		return errf(err)
	}

	return nil
}

func (d Doctor) createToolsFile(ctx context.Context, toolsFilename string) error {
	absPath, err := os.Getwd()
	if err != nil {
		return err
	}

	pathInfo, err := gomodulepath.ParseAt(absPath)
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

	runner := xgenny.NewRunner(ctx, absPath)
	if _, err := runner.RunAndApply(g); err != nil {
		return err
	}

	d.ev.Send(
		fmt.Sprintf("%s %s", toolsFilename, colors.Success("created")),
		events.Icon(icons.OK),
		events.Indent(1),
	)

	d.ev.Send("Installing dependency tools", events.ProgressUpdate())
	if err := cosmosgen.InstallDepTools(ctx, "."); err != nil {
		return err
	}

	for _, dep := range cosmosgen.DepTools() {
		d.ev.Send(
			fmt.Sprintf("%s %s", path.Base(dep), colors.Success("installed")),
			events.Icon(icons.OK),
			events.Indent(1),
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

	err = os.WriteFile(toolsFilename, buf.Bytes(), 0o600)
	if err != nil {
		return false, err
	}

	d.ev.Send(
		fmt.Sprintf("tools dependencies  %s", colors.Success("OK")),
		events.Icon(icons.OK),
		events.Indent(1),
		events.ProgressFinish(),
	)

	return true, nil
}
