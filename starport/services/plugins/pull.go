package plugins

// MUST BE RAN BEFORE BUILD
func (m *Manager) pull(ctx context.Context, cfg chaincfg.Config) error {
	for _, plug := range cfg.Plugins {
		// Seperate individual plugins by ID
		plugId := getPluginId(plug)

		// Check GOPATH for plugin

		// Get plugin home
		dst, err := formatPluginHome(m.ChainId, plugId)
		if err != nil {
			return err
		}

		if err := download(plug.Repo, plug.Subdir, dst); err != nil {
			return err
		}
	}

	return nil
}

func download(repo string, subdir string, dst string) error {
	url := "https://" + repo + ".git"
	if subdir != "" {
		url += "//" + subdir
	}

	if err := gogetter.Get(dst, url); err != nil {
		return err
	}

	return nil
}
