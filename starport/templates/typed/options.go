package typed

// Field ...
type Field struct {
	Name         string
	Datatype     string
	DatatypeName string
}

// Options ...
type Options struct {
	AppName    string
	ModuleName string
	ModulePath string
	TypeName   string
	Fields     []Field
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	return nil
}
