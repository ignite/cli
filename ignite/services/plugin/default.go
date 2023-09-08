package plugin

import (
	"errors"
	"fmt"
	"strings"

	"github.com/russross/blackfriday/v2"

	"github.com/ignite/cli"
)

// ErrNetworkPluginNotFound indicates that the network plugin is not found within the default plugins definition.
var ErrNetworkPluginNotFound = errors.New("default network plugin not found")

// DefaultPlugin defines a default Ignite plugin.
type DefaultPlugin struct {
	Use     string
	Short   string
	Aliases []string
	Path    string
}

// GetDefaultPlugins returns the list of default Ignite plugins.
func GetDefaultPlugins() ([]DefaultPlugin, error) {
	return parsePluginsMarkdown(cli.IgniteAppsDoc())
}

// GetDefaultNetworkPlugin returns the default network plugin.
func GetDefaultNetworkPlugin() (DefaultPlugin, error) {
	defaultPlugins, err := GetDefaultPlugins()
	if err != nil {
		return DefaultPlugin{}, err
	}

	for _, p := range defaultPlugins {
		if p.Use == "network" {
			return p, nil
		}
	}

	return DefaultPlugin{}, ErrNetworkPluginNotFound
}

func parsePluginsMarkdown(md []byte) ([]DefaultPlugin, error) {
	var (
		err      error
		plugins  []DefaultPlugin
		listNode *blackfriday.Node

		parser = blackfriday.New(blackfriday.WithNoExtensions())
		node   = parser.Parse(md)
	)

	// Locate the first list definition within the document
	node.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if !entering {
			return blackfriday.SkipChildren
		}

		if node.Type == blackfriday.List {
			listNode = node
			return blackfriday.Terminate
		}

		return blackfriday.GoToNext
	})

	if listNode == nil {
		return nil, errors.New("official Ignite Apps list not found")
	}

	// Extract ignite app list data from each list item
	listNode.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if !entering {
			return blackfriday.SkipChildren
		}

		if node.Type == blackfriday.Link {
			alias := string(node.FirstChild.Literal[0])
			url := string(node.LinkData.Destination)
			path := parsePluginRepoURL(url)
			if path == "" {
				err = fmt.Errorf("invalid Ignite App repository URL: %s", url)
				return blackfriday.Terminate
			}

			plugins = append(plugins, DefaultPlugin{
				Use:     string(node.FirstChild.Literal),
				Short:   string(node.Next.Literal[2:]),
				Aliases: []string{alias},
				Path:    path,
			})
		}

		return blackfriday.GoToNext
	})

	if err != nil {
		return nil, err
	}
	return plugins, nil
}

func parsePluginRepoURL(url string) string {
	elems := strings.Split(url, "/tree/")
	if len(elems) == 2 {
		return strings.Join(elems, "@")
	}
	return ""
}
