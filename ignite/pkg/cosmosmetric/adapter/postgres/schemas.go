package postgres

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	// SchemasDir defines the name for the embedded schema directory
	SchemasDir = "schemas"

	defaultSchemasTableName = "schema"

	tplSchemaTableDDL = `
		CREATE TABLE IF NOT EXISTS %[1]v (
			version     SMALLINT NOT NULL,
			created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

			CONSTRAINT %[1]v_pk PRIMARY KEY (version)
		)
	`
	tplSchemaVersionSQL = `
		SELECT COALESCE(MAX(version), 0)
		FROM %s
	`
)

// SchemasWalkFunc is the type of the function called by WalkFrom.
type SchemasWalkFunc func(version uint64, script []byte) error

// NewSchemas creates a new embedded SQL schema manager.
// The embedded FS is used to iterate the schema files.
// By default the applied schema versions are stored in the "schema"
// table but the name can have a prefix namespace when different
// packages are storing the schemas in the same database.
func NewSchemas(fs embed.FS, namespace string) Schemas {
	tableName := defaultSchemasTableName
	if namespace != "" {
		tableName = fmt.Sprintf("%s_%s", namespace, tableName)
	}

	return Schemas{tableName, fs}
}

// Schemas defines a type to manage versioning of embedded SQL schemas.
// Each schema file must live inside the embedded schemas directory and the name
// of each schema file must be numeric, where the number represents the version.
type Schemas struct {
	tableName string
	fs        embed.FS
}

// GetTableDDL returns the DDL to create the schemas table.
func (s Schemas) GetTableDDL() string {
	return fmt.Sprintf(tplSchemaTableDDL, s.tableName)
}

// GetSchemaVersionSQL returns the SQL query to get the current schema version.
func (s Schemas) GetSchemaVersionSQL() string {
	return fmt.Sprintf(tplSchemaVersionSQL, s.tableName)
}

// WalkFrom calls a function for SQL schemas starting from a specific version.
// This is useful to apply newer schemas that are not yet applied.
func (s Schemas) WalkFrom(fromVersion uint64, fn SchemasWalkFunc) error {
	return fs.WalkDir(s.fs, SchemasDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to read schema %s: %w", path, err)
		}

		if path == SchemasDir {
			return nil
		}

		// Extract the schema file version from the file name
		version := extractSchemaVersion(path)
		if version == 0 {
			return fmt.Errorf("invalid schema file name '%s'", path)
		}

		// Skip current schema file
		if version < fromVersion {
			return nil
		}

		// Read the SQL script from the schema file
		script, err := s.fs.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read schema '%s': %w", path, err)
		}

		return fn(version, script)
	})
}

func extractSchemaVersion(fileName string) uint64 {
	name := strings.TrimSuffix(
		filepath.Base(fileName),
		filepath.Ext(fileName),
	)

	// The names of the schema files MUST be numeric
	version, err := strconv.ParseUint(name, 10, 0)
	if err != nil {
		return 0
	}

	return version
}
