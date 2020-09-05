// Package scaffolder, initializes Starports and modifies to add
// more features in a later time.
package scaffolder

type Scaffolder struct {
	// name is app's name
	name string
}

// New creates a new Scaffolder for given app name.
func New(name string) *Scaffolder {
	return &Scaffolder{
		name: name,
	}
}
