package postgres

import "database/sql"

type cursor struct {
	rows *sql.Rows
}

func (c cursor) Err() error {
	return c.rows.Err()
}

func (c cursor) Next() bool {
	return c.rows.Next()
}

func (c cursor) Scan(values ...any) error {
	return c.rows.Scan(values...)
}

func (c cursor) Close() error {
	return c.rows.Close()
}
