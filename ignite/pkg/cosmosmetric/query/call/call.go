package call

// Option defines options for the call.
type Option func(*Call)

// WithArgs assigns one or more arguments for the call.
func WithArgs(args ...any) Option {
	return func(c *Call) {
		c.args = args
	}
}

// WithFields assigns one or more custom field names to select from the results.
func WithFields(fields ...string) Option {
	return func(c *Call) {
		c.fields = fields
	}
}

// New creates a new query call.
func New(name string, options ...Option) Call {
	c := Call{name: name}

	for _, o := range options {
		o(&c)
	}

	return c
}

// Call defines a data backend function or view call to get the query results.
type Call struct {
	name   string
	args   []any
	fields []string
}

// Name of the function or view to call.
func (c Call) Name() string {
	return c.name
}

// Args is a list of arguments for the call.
func (c Call) Args() []any {
	return c.args
}

// Fields is a list of custom field names to select from the call result.
func (c Call) Fields() []string {
	return c.fields
}
