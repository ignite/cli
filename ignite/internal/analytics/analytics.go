package analytics

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/manifoldco/promptui"

	"github.com/ignite/cli/ignite/pkg/gacli"
	"github.com/ignite/cli/ignite/pkg/randstr"
	"github.com/ignite/cli/ignite/version"
)

const (
	telemetryEndpoint  = "https://telemetry-cli.ignite.com"
	envDoNotTrack      = "DO_NOT_TRACK"
	igniteDir          = ".ignite"
	igniteAnonIdentity = "anon_identity.json"
)

var gaclient gacli.Client

type (
	// metric represents an analytics metric.
	options struct {
		// err sets metrics type as an error metric.
		err error
	}

	// anonIdentity represents an analytics identity file.
	anonIdentity struct {
		// name represents the username.
		Name string `json:"name" yaml:"name"`
		// doNotTrack represents the user track choice.
		DoNotTrack bool `json:"doNotTrack" yaml:"doNotTrack"`
	}
)

func init() {
	gaclient = gacli.New(telemetryEndpoint)
}

// Option configures ChainCmd.
type Option func(*options)

// WithError with application command error.
func WithError(error error) Option {
	return func(m *options) {
		m.err = error
	}
}

func SendMetric(wg *sync.WaitGroup, args []string, opts ...Option) {
	// only the app name
	if len(args) <= 1 {
		return
	}

	// apply analytics options.
	var opt options
	for _, o := range opts {
		o(&opt)
	}

	envDoNotTrackVar := os.Getenv(envDoNotTrack)
	if envDoNotTrackVar == "1" || strings.ToLower(envDoNotTrackVar) == "true" {
		return
	}

	if args[1] == "version" {
		return
	}

	fullCmd := strings.Join(args[1:], " ")

	dntInfo, err := checkDNT()
	if err != nil || dntInfo.DoNotTrack {
		return
	}

	met := gacli.Metric{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		FullCmd:   fullCmd,
		SessionID: dntInfo.Name,
		Version:   version.Version,
	}

	switch {
	case opt.err == nil:
		met.Status = "success"
	case opt.err != nil:
		met.Status = "error"
		met.Error = opt.err.Error()
	}
	met.Cmd = args[1]

	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = gaclient.SendMetric(met)
	}()
}

func checkDNT() (anonIdentity, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return anonIdentity{}, err
	}
	if err := os.Mkdir(filepath.Join(home, igniteDir), 0o700); err != nil && !os.IsExist(err) {
		return anonIdentity{}, err
	}
	identityPath := filepath.Join(home, igniteDir, igniteAnonIdentity)
	data, err := os.ReadFile(identityPath)
	if err != nil && !os.IsNotExist(err) {
		return anonIdentity{}, err
	}

	var i anonIdentity
	if err := json.Unmarshal(data, &i); err == nil {
		return i, nil
	}

	i.Name = randstr.Runes(10)
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
		return anonIdentity{}, err
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
