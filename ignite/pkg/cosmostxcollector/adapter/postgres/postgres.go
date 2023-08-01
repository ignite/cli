package postgres

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/lib/pq"

	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/query"
)

const (
	DefaultPort = 5432
	DefaultHost = "127.0.0.1"
)

const (
	adapterType = "postgres"

	sqlSelectBlockHeight = `
		SELECT COALESCE(MAX(height), 0)
		FROM tx
	`
	sqlSelectEventAttrs = `
		SELECT event_id, name, value FROM attribute
		WHERE event_id = ANY($1)
		ORDER BY event_id
	`
	sqlInsertTX = `
		INSERT INTO tx (hash, index, height, block_time)
		VALUES ($1, $2, $3, $4)
	`
	sqlInsertEvent = `
		INSERT INTO event (tx_hash, type, index)
		VALUES ($1, $2, $3) RETURNING id
	`
	sqlInsertEventAttr = `
		INSERT INTO attribute (event_id, name, value)
		VALUES ($1, $2, $3)
	`
	sqlInsertRawTX = `
		INSERT INTO raw_tx (hash, data)
		VALUES ($1, $2)
	`
)

//go:embed schemas/*
var fsSchemas embed.FS

// ErrClosed is returned when database connection is not open.
var ErrClosed = errors.New("no database connection")

// Option defines an option for the adapter.
type Option func(*Adapter)

// WithHost configures a database host name or IP.
func WithHost(host string) Option {
	return func(a *Adapter) {
		a.host = host
	}
}

// WithPort configures a database port.
func WithPort(port uint) Option {
	return func(a *Adapter) {
		a.port = port
	}
}

// WithUser configures a database user.
func WithUser(user string) Option {
	return func(a *Adapter) {
		a.user = user
	}
}

// WithPassword configures a database password.
func WithPassword(password string) Option {
	return func(a *Adapter) {
		a.password = password
	}
}

// WithParams configures extra database parameters.
func WithParams(params map[string]string) Option {
	return func(a *Adapter) {
		a.params = params
	}
}

// NewAdapter creates a new PostgreSQL adapter.
func NewAdapter(database string, options ...Option) (Adapter, error) {
	adapter := Adapter{
		host:     DefaultHost,
		port:     DefaultPort,
		database: database,
		schemas:  NewSchemas(fsSchemas, ""),
	}

	for _, o := range options {
		o(&adapter)
	}

	db, err := sql.Open("postgres", createPostgresURI(adapter))
	if err != nil {
		return Adapter{}, err
	}

	adapter.db = db

	return adapter, nil
}

// Adapter implements a data backend adapter for PostgreSQL.
type Adapter struct {
	host, user, password, database string
	port                           uint
	params                         map[string]string
	db                             *sql.DB
	schemas                        Schemas
}

// UpdateSchema updates the database schema to the latest version available.
// It applies all available schemas that were not applied already.
func (a Adapter) UpdateSchema(ctx context.Context, s Schemas) error {
	db, err := a.getDB()
	if err != nil {
		return err
	}

	// Create the schema table if it doesn't exist
	if _, err := db.ExecContext(ctx, s.GetTableDDL()); err != nil {
		return fmt.Errorf("failed to check schema table: %w", err)
	}

	// Get the current schema version
	var v uint64
	if err := db.QueryRowContext(ctx, s.GetSchemaVersionSQL()).Scan(&v); err != nil {
		return fmt.Errorf("failed to read current schema version: %w", err)
	}

	return s.WalkFrom(v+1, func(version uint64, script []byte) error {
		if _, err := db.ExecContext(ctx, string(script)); err != nil {
			return fmt.Errorf("error applying schema version %d: %w", version, err)
		}

		return nil
	})
}

func (a Adapter) GetType() string {
	return adapterType
}

func (a Adapter) Init(ctx context.Context) error {
	return a.UpdateSchema(ctx, a.schemas)
}

func (a Adapter) Save(ctx context.Context, txs []cosmosclient.TX) error {
	db, err := a.getDB()
	if err != nil {
		return err
	}

	// Start a transaction
	sqlTx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Rollback won't have any effect if the transaction is committed before
	defer sqlTx.Rollback()

	// Prepare insert statements to speed up "bulk" saving times
	txStmt, err := sqlTx.PrepareContext(ctx, sqlInsertTX)
	if err != nil {
		return err
	}

	defer txStmt.Close()

	evtStmt, err := sqlTx.PrepareContext(ctx, sqlInsertEvent)
	if err != nil {
		return err
	}

	defer evtStmt.Close()

	attrStmt, err := sqlTx.PrepareContext(ctx, sqlInsertEventAttr)
	if err != nil {
		return err
	}

	defer attrStmt.Close()

	// All the transactions are saved within the context of the same database
	// transactions and because of that either all block transactions are
	// saved or none of them.
	for _, tx := range txs {
		if err := saveRawTX(ctx, sqlTx, tx.Raw); err != nil {
			return err
		}

		if err := saveTX(ctx, txStmt, evtStmt, attrStmt, tx); err != nil {
			return err
		}
	}

	return sqlTx.Commit()
}

func (a Adapter) GetLatestHeight(ctx context.Context) (height int64, err error) {
	db, err := a.getDB()
	if err != nil {
		return 0, err
	}

	row := db.QueryRowContext(ctx, sqlSelectBlockHeight)
	if err = row.Scan(&height); err != nil {
		return 0, err
	}

	return height, nil
}

func (a Adapter) QueryEvents(ctx context.Context, q query.EventQuery) ([]query.Event, error) {
	db, err := a.getDB()
	if err != nil {
		return nil, err
	}

	sql := parseEventQuery(q)
	args := extractEventQueryArgs(q)
	rows, err := db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	var (
		events   []query.Event
		eventIDs []int64

		// Keep an index of the event position within the events slice
		// to find them later when updating their attributes.
		eventIndexes = make(map[int64]int)
	)

	for i := 0; rows.Next(); i++ {
		e := query.Event{}
		if err := rows.Scan(&e.ID, &e.Index, &e.TXHash, &e.Type, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to read event: %w", err)
		}

		events = append(events, e)
		eventIDs = append(eventIDs, e.ID)

		eventIndexes[e.ID] = i
	}

	// Don't query attributes when there are no events
	if len(events) == 0 {
		return events, nil
	}

	// Select the attributes for the events that matched the query
	rows, err = db.QueryContext(ctx, sqlSelectEventAttrs, pq.Array(eventIDs))
	if err != nil {
		return nil, err
	}

	// Update the attributes of the selected events
	for rows.Next() {
		var (
			eventID int64
			name    string
			value   []byte
		)

		if err := rows.Scan(&eventID, &name, &value); err != nil {
			return nil, fmt.Errorf("failed to read event attribute: %w", err)
		}

		i := eventIndexes[eventID]
		events[i].Attributes = append(events[i].Attributes, query.NewAttribute(name, value))
	}

	return events, nil
}

func (a Adapter) Query(ctx context.Context, q query.Query) (query.Cursor, error) {
	db, err := a.getDB()
	if err != nil {
		return nil, err
	}

	sql, err := parseQuery(q)
	if err != nil {
		return nil, err
	}

	args := extractQueryArgs(q)
	rows, err := db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (a Adapter) getDB() (*sql.DB, error) {
	if a.db == nil {
		return nil, ErrClosed
	}

	return a.db, nil
}

func createPostgresURI(a Adapter) string {
	uri := url.URL{
		Scheme: adapterType,
		Host:   fmt.Sprintf("%s:%d", a.host, a.port),
		Path:   a.database,
	}

	if a.user != "" {
		if a.password != "" {
			uri.User = url.UserPassword(a.user, a.password)
		} else {
			uri.User = url.User(a.user)
		}
	}

	// Add extra params as query arguments
	if a.params != nil {
		val := url.Values{}
		for k, v := range a.params {
			val.Set(k, v)
		}

		uri.RawQuery = val.Encode()
	}

	return uri.String()
}

func saveRawTX(ctx context.Context, sqlTx *sql.Tx, rtx *ctypes.ResultTx) error {
	hash := rtx.Hash.String()
	raw, err := json.Marshal(rtx)
	if err != nil {
		return fmt.Errorf("failed to encode raw TX %s: %w", hash, err)
	}

	if _, err := sqlTx.ExecContext(ctx, sqlInsertRawTX, hash, raw); err != nil {
		return fmt.Errorf("error saving raw TX %s: %w", hash, err)
	}

	return nil
}

func saveTX(ctx context.Context, txStmt, evtStmt, attrStmt *sql.Stmt, tx cosmosclient.TX) error {
	hash := tx.Raw.Hash.String()
	if _, err := txStmt.ExecContext(ctx, hash, tx.Raw.Index, tx.Raw.Height, tx.BlockTime); err != nil {
		return fmt.Errorf("error saving TX %s: %w", hash, err)
	}

	events, err := tx.GetEvents()
	if err != nil {
		return err
	}

	for i, evt := range events {
		var evtID int

		row := evtStmt.QueryRowContext(ctx, hash, evt.Type, i)
		if err := row.Err(); err != nil {
			return fmt.Errorf("error saving event '%s': %w", evt.Type, err)
		}

		if err := row.Scan(&evtID); err != nil {
			return fmt.Errorf("error reading event ID: %w", err)
		}

		for _, attr := range evt.Attributes {
			if _, err := attrStmt.ExecContext(ctx, evtID, attr.Key, attr.Value); err != nil {
				return fmt.Errorf("error saving event attr '%s.%s': %w", evt.Type, attr.Key, err)
			}
		}
	}

	return nil
}

func extractQueryArgs(q query.Query) []any {
	// When the query is a call to a postgres function
	// add the arguments before the filter values
	args := q.Args()

	// Add the values from the filters
	for _, f := range q.Filters() {
		if a := f.Value(); a != nil {
			args = append(args, a)
		}
	}

	return args
}

func extractEventQueryArgs(q query.EventQuery) (args []any) {
	for _, f := range q.Filters() {
		if a := f.Value(); a != nil {
			args = append(args, a)
		}
	}

	return args
}
