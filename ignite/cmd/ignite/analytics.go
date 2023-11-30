package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/manifoldco/promptui"

	"github.com/ignite/cli/ignite/pkg/gacli"
	"github.com/ignite/cli/ignite/version"
)

const (
	gaID           = "G-<ID>"
	gaSecret       = "<API_SECRET>"
	envDoNotTrack  = "DO_NOT_TRACK"
	igniteDir      = ".ignite"
	igniteIdentity = "identity.json"
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

	dntInfo, err := checkDNT()
	if err != nil || dntInfo.DoNotTrack {
		return
	}

	met := gacli.Metric{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		FullCmd:   m.command,
		SessionId: dntInfo.Name,
		Version:   version.Version,
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

	go gaclient.SendMetric(met)
}

func checkDNT() (identity, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return identity{}, err
	}
	if err := os.Mkdir(filepath.Join(home, igniteDir), 0o700); err != nil && !os.IsExist(err) {
		return identity{}, err
	}
	identityPath := filepath.Join(home, igniteDir, igniteIdentity)
	data, err := os.ReadFile(identityPath)
	if err != nil && !os.IsNotExist(err) {
		return identity{}, err
	}

	var i identity
	if err := json.Unmarshal(data, &i); err == nil {
		return i, nil
	}

	i.Name = randomdata.SillyName()
	i.DoNotTrack = false

	prompt := promptui.Select{
		Label: "Ignite collects metrics about command usage. " +
			"All data is anonymous and helps to improve Ignite. " +
			"Ignite respect the DNT rules (consoledonottrack.com). " +
			"Would you agree to share these metrics with us?",
		Items: []string{"Yes", "No"},
	}
	resultID, _, err := prompt.Run()
	if err != nil {
		return identity{}, err
	}

	if resultID != 0 {
		i.DoNotTrack = true
	}

	data, err = json.Marshal(&i)
	if err != nil {
		return i, err
	}

	return i, os.WriteFile(identityPath, data, 0o700)
}
