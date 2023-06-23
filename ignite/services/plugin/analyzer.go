package plugin

import "context"

// NewClientAPI creates a new app ClientAPI.
func NewClientAPI() clientAPI {
	return clientAPI{}
}

type clientAPI struct{}

// TODO: Implement dependency ClientAPI.

// Deoendencies returns chain app dependencies.
func (c clientAPI) Dependencies(_ context.Context) ([]*Dependency, error) {
	return []*Dependency{
		{Path: "Foo"},
	}, nil
}
