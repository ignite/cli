package plugins

func getPluginId(plug Plugin) string {
	var plugId string
	if plug.Name != "" {
		plugId = plug.Name
	} else {
		repoSplit := strings.Split(plug.Repo, "/")
		repoName := repoSplit[len(repoSplit)-1]
		if plug.Subdir != "" {
			plugId = repoName + "_" + plug.Subdir
		} else {
			plugId = repoName
		}
	}

	return plugId
}

func formatPluginHome(chainId string, pluginId string) (string, error) {
	configDirPath, err := chaincfg.ConfigDirPath()
	if err != nil {
		return "", err
	}

	if pluginId != "" {
		return path.Join(configDirPath, "local-chains", chainId, PLUGINS_DIR, pluginId), nil
	}

	return path.Join(configDirPath, "local-chains", chainId, PLUGINS_DIR), nil
}

func listDirs(dir string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filteredFiles := []os.FileInfo{}
	for _, file := range files {
		if file.IsDir() {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles, nil
}

func listFiles(dir, pattern string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filteredFiles := []os.FileInfo{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		matched, err := regexp.MatchString(pattern, file.Name())
		if err != nil {
			return nil, err
		}

		if matched {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles, nil
}
