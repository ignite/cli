package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Pallinder/go-randomdata"

	"github.com/ignite/cli/ignite/pkg/gacli"
	"github.com/ignite/cli/ignite/version"
)

const (
	gaid     = "<GA_KEY>" // Google Analytics' tracking id.
	loginAny = "any"
)

var (
	gaclient           *gacli.Client
	igniteDir          = ".ignite"
	igniteAnonIdentity = "anon"
)

// Metric represents an analytics metric.
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
	if len(os.Args) > 1 { // first is ignite (binary name).
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
		return loginAny, false
	}
	if err := os.Mkdir(filepath.Join(home, igniteDir), 0o700); err != nil {
		return loginAny, false
	}
	anonPath := filepath.Join(home, igniteDir, igniteAnonIdentity)
	data, err := os.ReadFile(anonPath)
	if err != nil {
		return loginAny, false
	}
	if len(data) != 0 {
		return string(data), true
	}
	name = randomdata.SillyName()
	if err := os.WriteFile(anonPath, []byte(name), 0o700); err != nil {
		return loginAny, false
	}
	return name, false
}
