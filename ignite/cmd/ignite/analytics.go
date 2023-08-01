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
	gaid          = "<GA_TRACK_ID>"
	envDoNotTrack = "DO_NOT_TRACK"
)

var (
	gaclient           *gacli.Client
	igniteDir          = ".ignite"
	igniteAnonIdentity = "anon"
)

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
		// name represents the user name.
		Name string `json:"name" yaml:"name"`
		// doNotTrack represents the user track choice.
		DoNotTrack bool `json:"doNotTrack" yaml:"doNotTrack"`
	}
)

func init() {
	gaclient = gacli.New(gaid)
}

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

	var met gacli.Metric
	switch {
	case m.err == nil:
		met.Category = "success"
	case m.err != nil:
		met.Category = "error"
		met.Value = m.err.Error()
	}
	met.Action = m.command

	cmds := strings.Split(m.command, " ")
	met.Label = cmds[0]
	if len(cmds) > 0 {
		met.Label = cmds[1]
	}
	met.User = ident.Name
	met.Version = version.Version
	go func() {
		gaclient.Send(met)
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
		Label: "Now, Ignite collects metrics so we can constantly improve our tools. " +
			"Since you are running ignite for the first time, we should ask. " +
			"It would be great if we could collect your metrics. " +
			"Do you want to share them with us? We will only ask for it one time!",
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
