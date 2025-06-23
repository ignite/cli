package ignitecmd

import (
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis/app"
	"github.com/ignite/cli/v29/ignite/pkg/gomodule"
	"github.com/ignite/cli/v29/ignite/services/chain"
)

func NewChainModulesList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List all Cosmos SDK modules in the app",
		Long:  "The list command lists all modules in the app.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			session := cliui.New(cliui.StartSpinner())
			defer session.End()

			chainOption := []chain.Option{
				chain.WithOutputer(session),
				chain.CollectEvents(session.EventBus()),
			}

			c, err := chain.NewWithHomeFlags(cmd, chainOption...)
			if err != nil {
				return err
			}

			modules, err := app.FindRegisteredModules(c.AppPath())
			if err != nil {
				return err
			}

			if len(modules) == 0 {
				session.Println("no modules found")
				return nil
			}

			modFile, err := gomodule.ParseAt(c.AppPath())
			if err != nil {
				return err
			}

			deps, err := gomodule.ResolveDependencies(modFile, false)
			if err != nil {
				return err
			}

			depMap := make(map[string]string)
			for _, dep := range deps {
				depMap[dep.Path] = dep.Version
			}

			// create a map of replaced modules for easy lookup
			// check the original required modules, not the resolved ones
			replacedMap := make(map[string]bool)
			for _, replace := range modFile.Replace {
				replacedMap[replace.Old.Path] = true
			}

			// get the app's module path to identify app modules
			appModulePath := modFile.Module.Mod.Path

			var entries [][]string
			for _, m := range modules {
				ver := depMap[m]
				modName := m

				switch {
				case strings.HasPrefix(m, appModulePath+"/"):
					ver = "main"
				case strings.HasPrefix(m, cosmosSDKModulePrefix+"/"):
					ver = depMap[cosmosSDKModulePrefix]
					modName = strings.TrimPrefix(m, cosmosSDKModulePrefix+"/")
				case strings.Contains(m, ibcModulePrefix+"/v"):
					modName, ver = getIBCVersion(m, depMap)
				case isModuleReplaced(m, replacedMap):
					ver = "locally replaced"
				}

				if ver == "" {
					ver = findBestMatchingVersion(m, depMap)
					if ver == "" {
						ver = "-"
					}
				}

				entries = append(entries, []string{modName, ver})
			}

			session.StopSpinner()

			// Sort entries by module name
			sort.SliceStable(entries, func(i, j int) bool {
				return entries[i][0] < entries[j][0]
			})

			header := []string{"module", "version"}
			return session.PrintTable(header, entries...)
		},
	}

	return c
}

const (
	cosmosSDKModulePrefix = "github.com/cosmos/cosmos-sdk"
	ibcModulePrefix       = "github.com/cosmos/ibc-go"
)

// isModuleReplaced checks if a module path (or its parent paths) is in the replaced map.
func isModuleReplaced(modulePath string, replacedMap map[string]bool) bool {
	checkPath := modulePath
	for checkPath != "" && checkPath != "." {
		if replacedMap[checkPath] {
			return true
		}
		// check parent path
		if idx := strings.LastIndex(checkPath, "/"); idx > 0 {
			checkPath = checkPath[:idx]
		} else {
			break
		}
	}
	return false
}

// for a given module path by checking progressively shorter paths.
func findBestMatchingVersion(modulePath string, depMap map[string]string) string {
	checkPath := modulePath
	for checkPath != "" && checkPath != "." {
		if version, exists := depMap[checkPath]; exists {
			return version
		}
		// check parent path
		if idx := strings.LastIndex(checkPath, "/"); idx > 0 {
			checkPath = checkPath[:idx]
		} else {
			break
		}
	}
	return ""
}

// getIBCVersion tries to extract the ibc-go version from the module path or dependencies.
func getIBCVersion(modulePath string, depMap map[string]string) (string, string) {
	// find the root ibc-go module path (with major version)
	parts := strings.Split(modulePath, "/")
	for i := range parts {
		if parts[i] == "ibc-go" && i+1 < len(parts) && strings.HasPrefix(parts[i+1], "v") {
			root := strings.Join(parts[:i+2], "/")
			ver := depMap[root]
			// clean module name after root
			mod := strings.TrimPrefix(modulePath, root+"/")
			return mod, ver
		}
	}
	return modulePath, ""
}
