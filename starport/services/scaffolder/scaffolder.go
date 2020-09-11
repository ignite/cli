// Package scaffolder initializes Starport apps and modifies existing ones
// to add more features in a later time.
package scaffolder

// Scaffolder is Starport app scaffolder.
type Scaffolder struct {
	// path is app's path on the filesystem.
	path string
}

// New initializes a new Scaffolder for app at path.
func New(path string) *Scaffolder {
	return &Scaffolder{
		path: path,
	}
}
