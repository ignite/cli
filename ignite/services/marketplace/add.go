package marketplace

import (
	"fmt"
	"os"
	"os/exec"
)

func AddPlugin(repo string) error {
	homeDir := os.Getenv("HOME")
	futurePluginPath := homeDir + "/.ignite/plugins/" + repo

	if _, err := os.Stat(futurePluginPath); err == nil {
		return fmt.Errorf("plugin %s already exists", repo)
	}

	command := "git"
	args := []string{"clone", "http://github.com/" + repo, futurePluginPath}

	cmd := exec.Command(command, args...)

	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return err
	}

	// Add plugin
	igniteCmd := exec.Command("ignite", "plugin", "add", "-g", futurePluginPath)
	igniteCmd.Stdout = os.Stdout
	err = igniteCmd.Run()
	if err != nil {
		return err
	}
	return nil
}
