package plugins

func (m *Manager) inject() {
	for _, cmdPlugin := range m.cmdPlugins {
		c := &cobra.Command{
			Use:   cmdPlugin.Usage(),
			Short: cmdPlugin.ShortDesc(),
			Long:  cmdPlugin.LongDesc(),
			Args: cobra.ExactArgs(hookPlugin.ExactArgs())
		}
	}

	for _, hookPlugin := range m.hookPlugins {
		
	}
}
