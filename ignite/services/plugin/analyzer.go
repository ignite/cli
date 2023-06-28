package plugin

import "context"

// NewAnalyzer creates a new app analyzer.
func NewAnalyzer() Analyzer {
	return analyzer{}
}

type analyzer struct{}

// TODO: Implement dependency analyzer.

// Deoendencies returns chain app dependencies.
func (a analyzer) Dependencies(_ context.Context) ([]*Dependency, error) {
	return []*Dependency{
		{Path: "Foo"},
	}, nil
}
