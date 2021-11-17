package plugin

import (
	"github.com/tendermint/starport/starport/chainconfig"
)

//
// Builder handles download build process for new plugin.
//

// TODO: How to divide GIT module?

const (
	pluginDir = "plugin"
)

// Builder provides interfaces to build plugins.
type Builder interface {
	Build(config chainconfig.Plugin) error
}

type builder struct {
}

func (b *builder) Build(config chainconfig.Plugin) error {
	// TODO:
	err := b.download(config.RepositoryURL)
	if err != nil {
		return err
	}

	err = b.build(config.Name)
	return err
}

func (b *builder) download(url string) error {
	_ = pluginDir

	// TODO:
	// Create `pluginDir` if not exist.
	// Clone repo from url on `pluginDir`.
	return nil
}

func (b *builder) build(name string) error {
	// TODO:
	// Build plugin.
	return nil
}

// NewBuilder creates new plugin builder.
// TODO: Need parameters?
func NewBuilder() (Builder, error) {
	// TODO:
	return &builder{}, nil
}
