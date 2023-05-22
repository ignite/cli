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
	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/templates/app"
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
	if version != chainconfig.LatestVersion {
		// migrate config file
		// Convert the current config to the latest version and update the YAML file
		var buf bytes.Buffer
		f.Seek(0, 0)
		if err := chainconfig.MigrateLatest(f, &buf); err != nil {
			return errf(err)
		}
		if err := os.WriteFile(configPath, buf.Bytes(), 0o755); err != nil {
			return errf(fmt.Errorf("config file migration failed: %w", err))
		}
		d.ev.Send(fmt.Sprintf("config file %s", colors.Success("migrated")),
			events.Icon(icons.OK), events.ProgressFinish())
	}
	d.ev.Send("config file OK", events.Icon(icons.OK), events.ProgressFinish())

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

	const toolsGoFile = "tools/tools.go"
	_, err := os.Stat(toolsGoFile)

	switch {
	case err == nil:
		// tools.go exists
		d.ev.Send(fmt.Sprintf("%s exists", toolsGoFile), events.Icon(icons.OK),
			events.ProgressFinish())
		// TODO ensure tools.go has the required dependencies

	case os.IsNotExist(err):
		// create tools.go
		pathInfo, err := gomodulepath.ParseAt(".")
		if err != nil {
			return errf(err)
		}
		g, err := app.NewGenerator(&app.Options{
			ModulePath:       pathInfo.RawPath,
			AppName:          pathInfo.Package,
			BinaryNamePrefix: pathInfo.Root,
			IncludePrefixes:  []string{toolsGoFile},
		})
		if err != nil {
			return errf(err)
		}
		// run generator
		runner := genny.WetRunner(ctx)
		if err := runner.With(g); err != nil {
			return errf(err)
		}
		if err := runner.Run(); err != nil {
			return errf(err)
		}
		d.ev.Send(fmt.Sprintf("%s %s", toolsGoFile, colors.Success("created")),
			events.Icon(icons.OK), events.ProgressFinish())

		d.ev.Send("Installing dependency tools", events.ProgressStart())
		if err := cosmosgen.InstallDepTools(ctx, "."); err != nil {
			return errf(err)
		}
		for _, dep := range cosmosgen.DepTools() {
			d.ev.Send(fmt.Sprintf("%s %s", path.Base(dep), colors.Success("installed")),
				events.Icon(icons.OK), events.ProgressFinish())
		}

	default:
		return errf(err)
	}
	return nil
}
