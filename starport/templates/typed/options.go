package typed

// Field ...
type Field struct {
	Name     string
	Datatype string
}

// Options ...
type Options struct {
	AppName    string
	ModulePath string
	TypeName   string
	Fields     []Field
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	return nil
}
