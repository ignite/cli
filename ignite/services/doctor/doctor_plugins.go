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

// MigratePluginsConfig migrates plugins config to Ignite Extension config if required.
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

	// Global extensions directory is always available because it is
	// created if it doesn't exists when any command is executed.
	extensionsPath := filepath.Join(globalPath, "extensions", "extensions.yml")
	if _, err := os.Stat(extensionsPath); err == nil {
		d.ev.Send(
			fmt.Sprintf("%s %s", extensionsPath, colors.Success("exists")),
			events.Icon(icons.OK),
			events.Indent(1),
		)

		// Ignite extensions config file exists in global directory
		return nil
	}

	// old plugins config file
	legacyPath, err := findConfigPath(filepath.Join(globalPath, "plugins"), "plugins")
	if err != nil {
		return err
	}

	// old plugin file not found, check for app config file
	if legacyPath == "" {
		legacyPath, err = findConfigPath(filepath.Join(globalPath, "apps"), "igniteapps")
		if err != nil {
			return err
		}

		// no legacy config file found
		if legacyPath == "" {
			return nil
		}
	}

	if err := d.migratePluginsConfigFiles(legacyPath, extensionsPath); err != nil {
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

	extensionsPath := filepath.Join(localPath, "extensions.yml")
	if _, err := os.Stat(extensionsPath); err == nil {
		d.ev.Send(
			fmt.Sprintf("%s %s", extensionsPath, colors.Success("exists")),
			events.Icon(icons.OK),
			events.Indent(1),
		)

		// Ignite extensions config file exists in current directory
		return nil
	}

	legacyPath, err := findConfigPath(localPath, "plugins")
	if err != nil {
		return err
	}

	// old plugin file not found, check for app config file
	if legacyPath == "" {
		legacyPath, err = findConfigPath(localPath, "igniteapps")
		if err != nil {
			return err
		}

		// no legacy config file found
		if legacyPath == "" {
			return nil
		}
	}

	return d.migratePluginsConfigFiles(legacyPath, extensionsPath)
}

func (d Doctor) migratePluginsConfigFiles(legacyPath, extensionsPath string) error {
	legacyFile, err := os.Open(legacyPath)
	if err != nil {
		return err
	}

	defer legacyFile.Close()

	extensionFile, err := os.OpenFile(extensionsPath, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}

	defer extensionFile.Close()

	if err = migratePluginsConfig(legacyFile, extensionFile); err != nil {
		return err
	}

	d.ev.Send(
		fmt.Sprintf("migrated config file %s to %s", colors.Faint(legacyPath), colors.Faint(extensionsPath)),
		events.Icon(icons.OK),
		events.Indent(1),
	)
	d.ev.SendInfo(
		fmt.Sprintf("file (and folder) %s can safely be removed", legacyPath),
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

	if extensions, ok := cfg["plugins"]; ok {
		cfg["extensions"] = extensions
		delete(cfg, "plugins")
	}

	if extensions, ok := cfg["apps"]; ok {
		cfg["extensions"] = extensions
		delete(cfg, "apps")
	}

	var buf bytes.Buffer
	if err = yaml.NewEncoder(&buf).Encode(cfg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func findConfigPath(dir, fileNameWithoutExtension string) (string, error) {
	for _, ext := range []string{"yml", "yaml"} {
		path := filepath.Join(dir, fmt.Sprintf("%s.%s", fileNameWithoutExtension, ext))
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
