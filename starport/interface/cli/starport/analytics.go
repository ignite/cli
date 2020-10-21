package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/tendermint/starport/starport/internal/version"
	"github.com/tendermint/starport/starport/pkg/gacli"
)

// Google Analytics' tracking id.
const gaid = "UA-51029217-18"

var (
	gaclient             *gacli.Client
	starportDir          = ".starport"
	starportAnonIdentity = "anon"
)

type Metric struct {
	// IsInstallation sets metrics type as an installation metric.
	IsInstallation bool

	// Err sets metrics type as an error metric.
	Err error

	// Login is the name of anon user.
	Login string
}

func addMetric(m Metric) {
	fullCommand := os.Args
	var rootCommand string
	if len(os.Args) > 1 { // first is starport (binary name).
		rootCommand = os.Args[1]
	}

	var met gacli.Metric
	switch {
	case m.IsInstallation:
		met.Category = "install"
	case m.Err == nil:
		met.Category = "success"
	case m.Err != nil:
		met.Category = "error"
		met.Value = m.Err.Error()
	}
	if m.IsInstallation {
		met.Action = m.Login
	} else {
		met.Action = rootCommand
		met.Label = strings.Join(fullCommand, " ")
	}
	user, _ := prepLoginName()
	met.User = user
	met.Version = version.Version
	gaclient.Send(met)
}

func prepLoginName() (name string, hadLogin bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "any", false
	}
	os.Mkdir(filepath.Join(home, starportDir), 0700)
	anonPath := filepath.Join(home, starportDir, starportAnonIdentity)
	data, err := ioutil.ReadFile(anonPath)
	if len(data) != 0 {
		return string(data), true
	}
	name = randomdata.SillyName()
	ioutil.WriteFile(anonPath, []byte(name), 0700)
	return name, false
}
