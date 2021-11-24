package plugins

import (
	"context"
	"log"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	chaincfg "github.com/tendermint/starport/starport/chainconfig"
	"github.com/tendermint/starport/starport/pkg/goenv"
)

// Set PersistentPreRunE here (https://pkg.go.dev/github.com/spf13/cobra#Command)
// Find a way to have this run on setup, not on every usage of the command

// This is the solution to the nested command additions,
// but is extremely tolling to run before every command.
// For this reason, you must figure out how to run this
// on limited intervals.
func PersistentPreRunE(cmd *cobra.Command, args []string) error {
	// Check for existence of plugins/hooks
	// If plugins exist in the config somewhere, extract them
	// Finally, inject it in the command

	// from old persistentprerun
	if err := goenv.ConfigurePath(); err != nil {
		return err
	}

	// Get chain config
	var cfg chaincfg.Config
	var cmdPlugins []CmdPlugin
	var hookPlugins []HookPlugin
	var chainId string
	ctx := context.Background()
	var configPlugins []chaincfg.Plugin

	if len(configPlugins) == 0 {
		return nil
	}

	for _, configPlugin := range configPlugins {
		log.Println("Config: ", configPlugin.Name)
		// Make sure plugin is downloaded
		// If not, download it
		// If it is, extract it

		pluginId := getPluginId(configPlugin)
		pluginHome, err := formatPluginHome(chainId, pluginId)
		if err != nil {
			return err
		}

		cmdPlugins, err = extractCommandPlugins(ctx, pluginHome, cmd, cfg)
		if err != nil {
			return err
		}

		hookPlugins, err = extractHookPlugins(ctx, pluginHome, cmd, cfg)
		if err != nil {
			return err
		}
	}

	if len(cmdPlugins) > 0 {
		for _, cmdPlugin := range cmdPlugins {
			for _, comd := range cmdPlugin.Registry() {
				cmdPath := strings.Split(cmd.CommandPath(), " ")
				cmdPluginCommandPath := append(comd.ParentCommand(), comd.Name())
				if reflect.DeepEqual(cmdPath, cmdPluginCommandPath) {
					c := &cobra.Command{
						Use:   comd.Usage(),
						Short: comd.ShortDesc(),
						Long:  comd.LongDesc(),
						Args:  cobra.ExactArgs(comd.NumArgs()),
						RunE:  comd.Exec,
					}

					cmd.AddCommand(c)
				}
			}
		}
	}

	if len(hookPlugins) > 0 {
		for _, hookPlugin := range hookPlugins {
			for _, hook := range hookPlugin.Registry() {
				hookPath := strings.Split(cmd.CommandPath(), " ")
				hookPluginCommandPath := append(hook.ParentCommand(), hook.Name())
				if reflect.DeepEqual(hookPath, hookPluginCommandPath) {
					switch hook.Type() {
					case "pre":
						cmd.PreRunE = hook.PreRun
					case "post":
						cmd.PostRunE = hook.PostRun
					default:
						cmd.PreRunE = hook.PreRun
						cmd.PostRunE = hook.PostRun
					}
				}
			}
		}
	}

	return nil
}
