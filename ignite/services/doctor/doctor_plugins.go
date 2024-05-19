package doctor

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/ignite/cli/v29/ignite/config"
	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/events"
)

// MigratePluginsConfig migrates plugins config to Ignite App config if required.
func (d Doctor) MigratePluginsConfig() error {
	errf := func(err error) error {
		return errors.Errorf("doctor migrate plugins config: %w", err)
	}

	d.ev.Send("Checking for legacy plugin config files:")
	d.ev.Send("Searching global plugins config file", events.ProgressStart())

	if err := d.migrateGlobalPluginConfig(); err != nil {
		return errf(err)
	}

	d.ev.Send("Searching local plugins config file", events.ProgressUpdate())

	if err := d.migrateLocalPluginsConfig(); err != nil {
		return errf(err)
	}

	d.ev.Send(
		fmt.Sprintf("plugin config files %s", colors.Success("OK")),
		events.Icon(icons.OK),
		events.Indent(1),
		events.ProgressFinish(),
	)

	return nil
}

func (d Doctor) migrateGlobalPluginConfig() error {
	globalPath, err := config.DirPath()
	if err != nil {
		return err
	}

	// Global apps directory is always available because it is
	// created if it doesn't exists when any command is executed.
	appsPath := filepath.Join(globalPath, "apps", "igniteapps.yml")
	if _, err := os.Stat(appsPath); err == nil {
		d.ev.Send(
			fmt.Sprintf("%s %s", appsPath, colors.Success("exists")),
			events.Icon(icons.OK),
			events.Indent(1),
		)

		// Ignite apps config file exists in global directory
		return nil
	}

	legacyPath, err := findPluginsConfigPath(filepath.Join(globalPath, "plugins"))
	if err != nil {
		return err
	} else if legacyPath == "" {
		// Nothing to migrate when the legacy plugins config path doesn't exist
		return nil
	}

	if err := d.migratePluginsConfigFiles(legacyPath, appsPath); err != nil {
		return err
	}

	d.ev.SendInfo(
		fmt.Sprintf("directory %s can safely be removed", filepath.Dir(legacyPath)),
		events.Icon(icons.Info),
		events.Indent(1),
	)

	return nil
}

func (d Doctor) migrateLocalPluginsConfig() error {
	localPath, err := chainconfig.LocateDefault(".")
	if err != nil {
		if errors.Is(err, chainconfig.ErrConfigNotFound) {
			// When app config is not found it means the doctor
			// command is not being run within a blockchain app,
			// so there is not local config to migrate
			return nil
		}

		return err
	}

	localPath, err = filepath.Abs(filepath.Dir(localPath))
	if err != nil {
		return err
	}

	appsPath := filepath.Join(localPath, "igniteapps.yml")
	if _, err := os.Stat(appsPath); err == nil {
		d.ev.Send(
			fmt.Sprintf("%s %s", appsPath, colors.Success("exists")),
			events.Icon(icons.OK),
			events.Indent(1),
		)

		// Ignite apps config file exists in current directory
		return nil
	}

	legacyPath, err := findPluginsConfigPath(localPath)
	if err != nil {
		return err
	} else if legacyPath == "" {
		// Nothing to migrate when plugins config file is not found in current directory
		return nil
	}

	return d.migratePluginsConfigFiles(legacyPath, appsPath)
}

func (d Doctor) migratePluginsConfigFiles(pluginsPath, appsPath string) error {
	pluginsFile, err := os.Open(pluginsPath)
	if err != nil {
		return err
	}

	defer pluginsFile.Close()

	appsFile, err := os.OpenFile(appsPath, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}

	defer appsFile.Close()

	if err = migratePluginsConfig(pluginsFile, appsFile); err != nil {
		return err
	}

	d.ev.Send(
		fmt.Sprintf("migrated config file %s to %s", colors.Faint(pluginsPath), colors.Faint(appsPath)),
		events.Icon(icons.OK),
		events.Indent(1),
	)
	d.ev.SendInfo(
		fmt.Sprintf("file %s can safely be removed", pluginsPath),
		events.Icon(icons.Info),
		events.Indent(1),
	)

	return nil
}

func migratePluginsConfig(r io.Reader, w io.Writer) error {
	bz, err := updatePluginsConfig(r)
	if err != nil {
		return err
	}

	_, err = w.Write(bz)
	if err != nil {
		return err
	}
	return nil
}

func updatePluginsConfig(r io.Reader) ([]byte, error) {
	var cfg map[string]any
	err := yaml.NewDecoder(r).Decode(&cfg)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	if apps, ok := cfg["plugins"]; ok {
		cfg["apps"] = apps
		delete(cfg, "plugins")
	}

	var buf bytes.Buffer
	if err = yaml.NewEncoder(&buf).Encode(cfg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func findPluginsConfigPath(dir string) (string, error) {
	for _, ext := range []string{"yml", "yaml"} {
		path := filepath.Join(dir, fmt.Sprintf("plugins.%s", ext))
		_, err := os.Stat(path)
		if err == nil {
			// File found
			return path, nil
		}

		if !os.IsNotExist(err) {
			return "", err
		}
	}
	return "", nil
}
