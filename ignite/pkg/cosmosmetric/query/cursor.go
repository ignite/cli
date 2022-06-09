package query

// Cursor defines a cursor to iterate query results.
type Cursor interface {
	// Err returns the last error seen by the Cursor, or nil if no error has occurred.
	Err() error

	// Next prepares the next result to be read with the Scan method.
	// It returns true on success, or false if there is no next result row or an error
	// happened while preparing it. Err should be consulted to distinguish between the
	// two cases.
	//
	// Every call to Scan, even the first one, must be preceded by a call to Next.
	Next() bool

	// Scan copies the query row into the pointed values.
	Scan(values ...any) error

	// Close closes the cursor preventing further iterations.
	// If Next is called and returns false and there are no further result sets,
	// the cursor is closed automatically and it will suffice to check the result of Err.
	// This method is idempotent meaning that after the first call, any subsequent calls
	// will not change the state.
	Close() error
}
