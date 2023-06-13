package plugin

import "context"

// NewAnalizer creates a new app analizer.
func NewAnalizer() Analizer {
	return analizer{}
}

type analizer struct{}

// TODO: Implement dependency analizer.

// Deoendencies returns chain app dependencies.
func (a analizer) Dependencies(_ context.Context) ([]*Dependency, error) {
	return []*Dependency{
		{Path: "Foo"},
	}, nil
}
