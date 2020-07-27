package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/ilgooz/analytics-go"
	"github.com/tendermint/starport/pkg/analyticsutil"
)

const (
	analyticsEndpoint = "https://analytics.starport.cloud"
	analyticsKey      = "ib6mwzNSLW6qIFRTyftezJL8cX4jWkQY"
)

var (
	analyticsc           *analyticsutil.Client
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
	var (
		category string
		event    string
	)
	props := analytics.NewProperties()
	switch {
	case m.IsInstallation:
		category = "install"
	case m.Err == nil:
		category = "success"
	case m.Err != nil:
		category = "error"
		props.Set("value", m.Err.Error())
	}
	if m.IsInstallation {
		event = m.Login
	} else {
		event = rootCommand
		props.Set("label", strings.Join(fullCommand, " "))
	}
	props.Set("category", category)
	analyticsc.Track(analytics.Track{
		Event:      event,
		Properties: props,
	})
}

func prepLoginName() (name string, hadLogin bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "any", false
	}
	anonPath := filepath.Join(home, starportDir, starportAnonIdentity)
	data, err := ioutil.ReadFile(anonPath)
	if len(data) != 0 {
		return string(data), true
	}
	name = randomdata.SillyName()
	ioutil.WriteFile(anonPath, []byte(name), 0644)
	return name, false
}
