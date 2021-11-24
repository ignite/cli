package plugins

import (
	"context"
	"os"
	"time"

	gogetter "github.com/hashicorp/go-getter"
	chaincfg "github.com/tendermint/starport/starport/chainconfig"
)

// MUST BE RAN BEFORE BUILD
func (m *Manager) pull(ctx context.Context, cfg chaincfg.Config) error {
	for _, cfgPlugin := range cfg.Plugins {
		// Seperate individual plugins by ID
		plugId := getPluginId(cfgPlugin)

		// Check GOPATH for plugin?

		// Get plugin home
		dst, err := formatPluginHome(m.ChainId, plugId)
		if err != nil {
			return err
		}

		_, err = os.Stat(dst)
		if err == nil {
			err = os.RemoveAll(dst)
			if err != nil {
				return err
			}
		}

		if err := download(cfgPlugin.Repo, cfgPlugin.Subdir, dst); err != nil {
			return err
		}
	}

	return nil
}

func download(repo string, subdir string, dst string) error {
	url := "git::https://" + repo + ".git"
	// url = repo
	if subdir != "" {
		url += ("//" + subdir)
	}

	// Not cancelling for some noob reason
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	client := gogetter.Client{
		Ctx:  ctx,
		Src:  url,
		Dst:  dst,
		Mode: gogetter.ClientModeAny,
	}

	if err := client.Get(); err != nil {
		return err
	}

	return nil
}
