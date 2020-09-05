// Package scaffolder initializes Starport apps and modifies existing ones
// to add more features in a later time.
package scaffolder

// Scaffolder is Starport app scaffolder.
type Scaffolder struct {
	// name is app's name
	name string
}

// New creates a new Scaffolder for given app name.
// app name can be belong to an existent app or the one wished
// to be initialized.
func New(name string) *Scaffolder {
	return &Scaffolder{
		name: name,
	}
}
