package postgres

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosmetric/query"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"golang.org/x/sync/errgroup"

	_ "github.com/lib/pq" // required to register postgres sql driver
)

const (
	DefaultPort            = 5432
	DefaultHost            = "127.0.0.1"
	DefaultSaveConcurrency = 0 // no limit

	adapterType = "postgres"

	queryBlockHeight = `
		SELECT COALESCE(MAX(height), 0)
		FROM tx
	`
	queryInsertTX = `
		INSERT INTO tx (hash, index, height, block_time)
		VALUES ($1, $2, $3, $4)
	`
	queryInsertAttr = `
		INSERT INTO attribute (tx_hash, event_type, event_index, name, value)
		VALUES ($1, $2, $3, $4, $5)
	`
	querySchemaExists = `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = 'schema'
		)
	`
	querySchemaVersion = `
		SELECT COALESCE(MAX(version), 0)
		FROM schema
	`
	queryInsertRawTX = `
		INSERT INTO raw_tx (hash, data)
		VALUES ($1, $2)
	`

	// Latest schema version that the adapter should apply.
	// This version should be updated when new schema/*.sql files are added
	// to match the name of the latest file, otherwise the new schemas won't
	// be applied. All schema file names MUST be numeric.
	schemaVersion = 1
)

//go:embed schemas/*
var fsSchemas embed.FS

var (
	// ErrClosed is returned when database connection is not open.
	ErrClosed = errors.New("no database connection")
)

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

// WithSaveConcurrency configures the max number of parallel
// INSERT statements allowed during save.
func WithSaveConcurrency(concurrency int) Option {
	return func(a *Adapter) {
		a.saveConcurrency = concurrency
	}
}

// NewAdapter creates a new PostgreSQL adapter.
func NewAdapter(database string, options ...Option) (Adapter, error) {
	adapter := Adapter{
		host:            DefaultHost,
		port:            DefaultPort,
		saveConcurrency: DefaultSaveConcurrency,
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
	saveConcurrency                int
	db                             *sql.DB
}

func (a Adapter) GetType() string {
	return adapterType
}

func (a Adapter) Init(ctx context.Context) error {
	current, err := a.getSchemaVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to read schema version: %w", err)
	}

	if current == schemaVersion {
		return nil
	} else if current > schemaVersion {
		return fmt.Errorf("latest schema version is v%d, found v%d", schemaVersion, current)
	}

	for i := current + 1; i <= schemaVersion; i++ {
		name := fmt.Sprintf("%d.sql", i)
		if err := a.applySchema(ctx, name); err != nil {
			return fmt.Errorf("error applying schema %s: %w", name, err)
		}
	}

	return nil
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
	txStmt, err := sqlTx.PrepareContext(ctx, queryInsertTX)
	if err != nil {
		return err
	}

	defer txStmt.Close()

	attrStmt, err := sqlTx.PrepareContext(ctx, queryInsertAttr)
	if err != nil {
		return err
	}

	defer attrStmt.Close()

	wg, groupCtx := errgroup.WithContext(ctx)
	if a.saveConcurrency > 0 {
		wg.SetLimit(a.saveConcurrency)
	}

	for _, tx := range txs {
		tx := tx

		wg.Go(func() error {
			return saveRawTX(groupCtx, sqlTx, tx.Raw)
		})
		wg.Go(func() error {
			return saveTX(groupCtx, txStmt, attrStmt, tx)
		})
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	return sqlTx.Commit()
}

func (a Adapter) GetLatestHeight(ctx context.Context) (height int64, err error) {
	db, err := a.getDB()
	if err != nil {
		return 0, err
	}

	row := db.QueryRowContext(ctx, queryBlockHeight)
	if err = row.Scan(&height); err != nil {
		return 0, err
	}

	return height, nil
}

func (a Adapter) Query(ctx context.Context, q query.Query) (query.Cursor, error) {
	db, err := a.getDB()
	if err != nil {
		return nil, err
	}

	query, err := parseQuery(q)
	if err != nil {
		return nil, err
	}

	args := extractQueryArgs(q)
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &cursor{rows}, nil
}

func (a Adapter) getDB() (*sql.DB, error) {
	if a.db == nil {
		return nil, ErrClosed
	}

	return a.db, nil
}

func (a Adapter) getSchemaVersion(ctx context.Context) (version uint, err error) {
	db, err := a.getDB()
	if err != nil {
		return 0, err
	}

	exists := false
	row := db.QueryRowContext(ctx, querySchemaExists)
	if err = row.Scan(&exists); err != nil {
		return 0, err
	}

	if !exists {
		return 0, nil
	}

	row = db.QueryRowContext(ctx, querySchemaVersion)
	if err = row.Scan(&version); err != nil {
		return 0, err
	}

	return version, nil
}

func (a Adapter) applySchema(ctx context.Context, filename string) error {
	script, err := fsSchemas.ReadFile(fmt.Sprintf("schemas/%s", filename))
	if err != nil {
		return err
	}

	db, err := a.getDB()
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, string(script))

	return err
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
		query := url.Values{}
		for k, v := range a.params {
			query.Set(k, v)
		}

		uri.RawQuery = query.Encode()
	}

	return uri.String()
}

func saveRawTX(ctx context.Context, sqlTx *sql.Tx, rtx *ctypes.ResultTx) error {
	hash := rtx.Hash.String()
	raw, err := json.Marshal(rtx)
	if err != nil {
		return fmt.Errorf("failed to encode raw TX %s: %w", hash, err)
	}

	if _, err := sqlTx.ExecContext(ctx, queryInsertRawTX, hash, raw); err != nil {
		return fmt.Errorf("error saving raw TX %s: %w", hash, err)
	}

	return nil
}

func saveTX(ctx context.Context, txStmt, attrStmt *sql.Stmt, tx cosmosclient.TX) error {
	hash := tx.Raw.Hash.String()
	if _, err := txStmt.ExecContext(ctx, hash, tx.Raw.Index, tx.Raw.Height, tx.BlockTime); err != nil {
		return fmt.Errorf("error saving TX %s: %w", hash, err)
	}

	events, err := tx.GetEvents()
	if err != nil {
		return err
	}

	for i, evt := range events {
		for _, attr := range evt.Attributes {
			if _, err := attrStmt.ExecContext(ctx, hash, evt.Type, i, attr.Key, attr.Value); err != nil {
				return fmt.Errorf("error saving event attr '%s.%s': %w", evt.Type, attr.Key, err)
			}
		}
	}

	return nil
}

func extractQueryArgs(q query.Query) (args []any) {
	// When the query is a call to a postgres function
	// add the arguments before the filter values
	if call := q.GetCall(); len(call.Args) > 0 {
		args = append(args, call.Args...)
	}

	// Add the values from the filters
	for _, f := range q.GetFilters() {
		if a := f.GetValue(); a != nil {
			args = append(args, a)
		}
	}

	return args
}
