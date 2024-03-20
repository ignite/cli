package analytics

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/gacli"
	"github.com/ignite/cli/v29/ignite/pkg/gitpod"
	"github.com/ignite/cli/v29/ignite/pkg/randstr"
	"github.com/ignite/cli/v29/ignite/version"
)

const (
	telemetryEndpoint  = "https://telemetry-cli.ignite.com"
	envDoNotTrack      = "DO_NOT_TRACK"
	envCI              = "CI"
	envGitHubActions   = "GITHUB_ACTIONS"
	igniteDir          = ".ignite"
	igniteAnonIdentity = "anon_identity.json"
)

var gaclient gacli.Client

// anonIdentity represents an analytics identity file.
type anonIdentity struct {
	// Name represents the username.
	Name string `json:"name" yaml:"name"`
	// DoNotTrack represents the user track choice.
	DoNotTrack bool `json:"doNotTrack" yaml:"doNotTrack"`
}

func init() {
	gaclient = gacli.New(telemetryEndpoint)
}

// SendMetric send command metrics to analytics.
func SendMetric(wg *sync.WaitGroup, cmd *cobra.Command) {
	if cmd.Name() == "version" {
		return
	}

	dntInfo, err := checkDNT()
	if err != nil || dntInfo.DoNotTrack {
		return
	}

	path := cmd.CommandPath()
	met := gacli.Metric{
		Name:      cmd.Name(),
		Cmd:       path,
		Tag:       strings.ReplaceAll(path, " ", "+"),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		SessionID: dntInfo.Name,
		Version:   version.Version,
		IsGitPod:  gitpod.IsOnGitpod(),
		IsCI:      getIsCI(),
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = gaclient.SendMetric(met)
	}()
}

// checkDNT check if the user allow to track data or if the DO_NOT_TRACK
// env var is set https://consoledonottrack.com/
func checkDNT() (anonIdentity, error) {
	if dnt, err := strconv.ParseBool(os.Getenv(envDoNotTrack)); err != nil || dnt {
		return anonIdentity{DoNotTrack: true}, nil
	}

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
		Label: "Ignite uses anonymized metrics to enhance the application, " +
			"focusing on features such as command usage. We do not collect " +
			"identifiable personal information. Your privacy is important to us. " +
			"For more details, please visit our Privacy Policy at https://ignite.com/privacy " +
			"and our Terms of Use at https://ignite.com/terms-of-use. " +
			"Do you consent to the collection of these usage metrics for analytics purposes?",
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

func getIsCI() bool {
	ci, err := strconv.ParseBool(os.Getenv(envCI))
	if err != nil {
		return false
	}

	if ci {
		return true
	}

	ci, err = strconv.ParseBool(os.Getenv(envGitHubActions))
	if err != nil {
		return false
	}

	return ci
}
