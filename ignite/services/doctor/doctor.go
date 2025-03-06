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
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/events"
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
