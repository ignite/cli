package ignitecmd

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
	gaid          = "<GA_KEY>" // Google Analytics' tracking id.
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

func addCmdMetric(m metric) error {
	if os.Getenv(envDoNotTrack) == "1" {
		return nil
	}

	ident, err := prepareMetrics()
	if err != nil {
		return err
	}

	fullCommand := os.Args
	var rootCommand string
	if len(os.Args) > 1 { // first is ignite (binary name).
		rootCommand = os.Args[1]
	}

	var met gacli.Metric
	switch {
	case m.err == nil:
		met.Category = "success"
	case m.err != nil:
		met.Category = "error"
		met.Value = m.err.Error()
	}
	met.Action = rootCommand
	met.Label = strings.Join(fullCommand, " ")
	met.User = ident.Name
	met.Version = version.Version
	return gaclient.Send(met)
}

func prepareMetrics() (identity, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return identity{}, err
	}
	if err := os.Mkdir(filepath.Join(home, igniteDir), 0o700); err != nil {
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

	return i, os.WriteFile(anonPath, data, 0o700)
}
