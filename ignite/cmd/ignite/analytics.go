package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/manifoldco/promptui"

	"github.com/ignite/cli/ignite/pkg/gacli"
	"github.com/ignite/cli/ignite/version"
)

const (
	gaID               = "G-<ID>"
	gaSecret           = "<API_SECRET>"
	envDoNotTrack      = "DO_NOT_TRACK"
	igniteDir          = ".ignite"
	igniteAnonIdentity = "anon"
)

var gaclient gacli.Client

type (
	// metric represents an analytics metric.
	metric struct {
		// err sets metrics type as an error metric.
		err error
		// command is the command name.
		command string
	}

	// identity represents an analytics identity file.
	identity struct {
		// name represents the username.
		Name string `json:"name" yaml:"name"`
		// doNotTrack represents the user track choice.
		DoNotTrack bool `json:"doNotTrack" yaml:"doNotTrack"`
	}
)

func addCmdMetric(m metric) {
	envDoNotTrackVar := os.Getenv(envDoNotTrack)
	if envDoNotTrackVar == "1" || strings.ToLower(envDoNotTrackVar) == "true" {
		return
	}

	if m.command == "ignite version" {
		return
	}

	ident, err := prepareMetrics()
	if err != nil {
		return
	}

	met := gacli.Metric{
		FullCmd: m.command,
		User:    ident.Name,
		Version: version.Version,
	}

	switch {
	case m.err == nil:
		met.Status = "success"
	case m.err != nil:
		met.Status = "error"
		met.Error = m.err.Error()
	}

	cmds := strings.Split(m.command, " ")
	met.Cmd = cmds[0]
	if len(cmds) > 0 {
		met.Cmd = cmds[1]
	}
	go func() {
		gaclient.SendMetric(met)
	}()
}

func prepareMetrics() (identity, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return identity{}, err
	}
	if err := os.Mkdir(filepath.Join(home, igniteDir), 0o700); err != nil && !os.IsExist(err) {
		return identity{}, err
	}
	anonPath := filepath.Join(home, igniteDir, igniteAnonIdentity)
	data, err := os.ReadFile(anonPath)
	if err != nil && !os.IsNotExist(err) {
		return identity{}, err
	}

	i := identity{
		Name:       randomdata.SillyName(),
		DoNotTrack: false,
	}
	if len(data) > 0 {
		return i, json.Unmarshal(data, &i)
	}

	prompt := promptui.Prompt{
		Label: "Ignite would like to collect metrics about command usage. " +
			"All data will be anonymous and helps to improve Ignite. " +
			"Ignite respect the DNT rules (consoledonottrack.com). " +
			"Would you agree to share these metrics with us?",
		IsConfirm: true,
	}
	if _, err := prompt.Run(); err != nil {
		return identity{}, err
	}

	data, err = json.Marshal(&i)
	if err != nil {
		return i, err
	}

	return i, os.WriteFile(anonPath, data, 0o700)
}
