package postgres

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"testing/fstest"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/query"
)

var (
	eventFields     = []string{"id", "index", "tx_hash", "type", "created_at"}
	eventAttrFields = []string{"event_id", "name", "value"}
)

func TestUpdateSchema(t *testing.T) {
	// Arrange
	db, mock := createMatchEqualSQLMock(t)
	defer db.Close()

	ctx := context.Background()
	tplSchemaScript := `BEGIN;
		INSERT INTO schema(version)
		VALUES(%d)
	;%sCOMMIT;`

	// Arrange: Schema files
	schemasData := []string{"/* FOO */", "/* BAR */"}
	fs := fstest.MapFS{
		"schemas/1.sql": &fstest.MapFile{Data: []byte(schemasData[0])},
		"schemas/2.sql": &fstest.MapFile{Data: []byte(schemasData[1])},
	}
	s := NewSchemas(fs, "")

	// Arrange: Prepare database adapter
	adapter := Adapter{
		db:      db,
		schemas: s,
	}

	// Arrange: Database mock and expectations
	mock.
		ExpectExec(s.GetTableDDL()).
		WillReturnResult(
			// DDL execution won't affect any rows or IDs
			sqlmock.NewResult(0, 0),
		)
	mock.
		ExpectQuery(s.GetSchemaVersionSQL()).
		WillReturnRows(
			// Zero is returned to signal that there are no versions applied.
			// When no versions are applied schema walk will start from version 1.
			sqlmock.NewRows([]string{"version"}).AddRow(uint64(0)),
		)

	for i, data := range schemasData {
		version := i + 1
		script := fmt.Sprintf(tplSchemaScript, version, data)

		// Add database mock and expectation for the current schema version
		mock.
			ExpectExec(script).
			WillReturnResult(sqlmock.NewResult(0, 0))
	}

	// Act
	err := adapter.UpdateSchema(ctx, s)

	// Assert
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSave(t *testing.T) {
	// Arrange
	db, mock := createMatchEqualSQLMock(t)
	defer db.Close()

	adapter := Adapter{db: db}
	ctx := context.Background()
	hash := "F2564C78071E26643AE9B3E2A19FA0DC10D4D9E873AA0BE808660123F11A1E78"

	// Arrange: A Cosmos client TX to save
	evtAttr := abci.EventAttribute{
		Key:   []byte("recipient"),
		Value: []byte("cosmos1crje20aj4gxdtyct7z3knxqry2jqt2fuaey6u5"),
	}
	evt := abci.Event{
		Type:       "transfer",
		Attributes: []abci.EventAttribute{evtAttr},
	}

	h, _ := hex.DecodeString(hash) // TODO: How to properly generate TX hash for the result?
	tx := cosmosclient.TX{
		// Tendermint API search result
		Raw: &ctypes.ResultTx{
			Hash:   h,
			Height: 1,
			Index:  0,
			TxResult: abci.ResponseDeliverTx{
				Events: []abci.Event{evt},
			},
		},
	}

	// Arrange: JSON of the raw transaction result
	jsonResTX, err := json.Marshal(tx.Raw)
	require.NoError(t, err)

	// Arrange: Database mock and expectations for prepared SQL statements
	mock.ExpectBegin()

	txStmt := mock.ExpectPrepare(`
		INSERT INTO tx (hash, index, height, block_time)
		VALUES ($1, $2, $3, $4)
	`)
	evtStmt := mock.ExpectPrepare(`
		INSERT INTO event (tx_hash, type, index)
		VALUES ($1, $2, $3) RETURNING id
	`)
	attrStmt := mock.ExpectPrepare(`
		INSERT INTO attribute (event_id, name, value)
		VALUES ($1, $2, $3)
	`)

	// Arrange: Database mock and expectations for INSERT statement executions
	insertResult := sqlmock.NewResult(0, 1)
	evtIndex := 0
	evtID := int64(1)
	jsonEvtAttrValue := []byte(fmt.Sprintf(`"%s"`, evtAttr.Value))

	mock.
		ExpectExec(`
			INSERT INTO raw_tx (hash, data)
			VALUES ($1, $2)
		`).
		WithArgs(tx.Raw.Hash.String(), jsonResTX).
		WillReturnResult(insertResult)

	txStmt.
		ExpectExec().
		WithArgs(hash, tx.Raw.Index, tx.Raw.Height, tx.BlockTime).
		WillReturnResult(insertResult)
	evtStmt.
		ExpectQuery().
		WithArgs(hash, evt.Type, evtIndex).
		WillReturnRows(
			sqlmock.NewRows([]string{"event_id"}).AddRow(evtID),
		)
	attrStmt.
		ExpectExec().
		WithArgs(evtID, string(evtAttr.Key), jsonEvtAttrValue).
		WillReturnResult(insertResult)

	mock.ExpectCommit()

	// Act
	err = adapter.Save(ctx, []cosmosclient.TX{tx})

	// Assert
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLatestHeight(t *testing.T) {
	// Arrange
	db, mock := createMatchEqualSQLMock(t)
	defer db.Close()

	adapter := Adapter{db: db}
	ctx := context.Background()

	// Arrange: Database mock and expectations
	wantHeight := int64(42)

	mock.
		ExpectQuery(sqlSelectBlockHeight).
		WillReturnRows(
			sqlmock.NewRows([]string{"height"}).AddRow(wantHeight),
		)

	// Act
	height, err := adapter.GetLatestHeight(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, wantHeight, height)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestQuery(t *testing.T) {
	// Arrange
	var rowValue string

	db, mock := createMatchEqualSQLMock(t)
	defer db.Close()

	adapter := Adapter{db: db}
	ctx := context.Background()

	// Arrange: Query
	qry := query.New("baz", query.Fields("foo"))

	// Arrange: Database mock and expectations
	wantRowValue := "expected"
	fields := []string{"foo"}
	rows := sqlmock.NewRows(fields).AddRow(wantRowValue)

	mock.
		ExpectQuery(`
			SELECT DISTINCT foo
			FROM baz
			WHERE true
			LIMIT 30 OFFSET 0
		`).
		WillReturnRows(rows)

	// Act
	cr, err := adapter.Query(ctx, qry)
	if cr.Next() {
		cr.Scan(&rowValue)
	}

	// Assert
	require.NoError(t, err, "expected no query errors on execution")
	require.Equal(t, wantRowValue, rowValue)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestQueryCursor(t *testing.T) {
	// Arrange
	var (
		rowValue            string
		cursorNextSucceeded bool
		err                 error
	)

	db, mock := createMatchEqualSQLMock(t)
	defer db.Close()

	adapter := Adapter{db: db}
	ctx := context.Background()

	// Arrange: Query
	qry := query.New("baz", query.Fields("foo"))

	// Arrange: Database mock and expectations
	wantRowValue := "expected"
	fields := []string{"foo"}
	rows := sqlmock.NewRows(fields).AddRow(wantRowValue)

	mock.
		ExpectQuery(`
			SELECT DISTINCT foo
			FROM baz
			WHERE true
			LIMIT 30 OFFSET 0
		`).
		WillReturnRows(rows)

	// Act
	cr, _ := adapter.Query(ctx, qry)
	if cursorNextSucceeded = cr.Next(); cursorNextSucceeded {
		err = cr.Scan(&rowValue)
	}

	// Assert
	require.True(t, cursorNextSucceeded, "expected cursor.Next() to succeed")
	require.NoError(t, err, "expected no scan errors on execution")
	require.Equal(t, wantRowValue, rowValue)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestQueryWithFilter(t *testing.T) {
	// Arrange
	db, mock := createMatchEqualSQLMock(t)
	defer db.Close()

	adapter := Adapter{db: db}
	ctx := context.Background()

	// Arrange: Query
	wantArg := "bar"

	qry := query.New(
		"baz",
		query.Fields("foo"),
		query.WithFilters(
			NewFilter("foo", wantArg),
		),
	)

	// Arrange: Database mock and expectations
	fields := []string{"baz"}
	rows := sqlmock.NewRows(fields)

	mock.
		ExpectQuery(`
			SELECT DISTINCT foo
			FROM baz
			WHERE foo = $1
			LIMIT 30 OFFSET 0
		`).
		WithArgs(wantArg).
		WillReturnRows(rows)

	// Act
	_, err := adapter.Query(ctx, qry)

	// Assert
	require.NoError(t, err, "expected no query errors on execution")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestQueryError(t *testing.T) {
	// Arrange
	db, mock := createMatchEqualSQLMock(t)
	defer db.Close()

	adapter := Adapter{db: db}
	ctx := context.Background()

	// Arrange: Query
	qry := query.New("baz", query.Fields("foo"), query.WithoutPaging())

	// Arrange: Database mock and expectations
	wantErr := errors.New("expected error")

	mock.
		ExpectQuery("SELECT DISTINCT foo FROM baz WHERE true").
		WillReturnError(wantErr)

	// Act
	_, err := adapter.Query(ctx, qry)

	// Assert
	require.Equal(t, wantErr, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestQueryRowError(t *testing.T) {
	// Arrange
	db, mock := createMatchEqualSQLMock(t)
	defer db.Close()

	adapter := Adapter{db: db}
	ctx := context.Background()
	cols := []string{"name"}

	// Arrange: Query
	qry := query.New("baz", query.Fields(cols[0]), query.WithoutPaging())

	// Arrange: Database mock and expectations
	wantErr := errors.New("expected error")

	row := sqlmock.
		NewRows(cols).
		AddRow("foo").
		RowError(0, wantErr)

	mock.
		ExpectQuery("SELECT DISTINCT name FROM baz WHERE true").
		WillReturnRows(row)

	// Act
	cr, err := adapter.Query(ctx, qry)

	// Assert
	require.NoError(t, err, "expected no query errors on execution")
	require.False(t, cr.Next(), "expected cursor.Next() to fail")
	require.Equal(t, wantErr, cr.Err())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestEventQuery(t *testing.T) {
	// Arrange
	db, mock := createMatchEqualSQLMock(t)
	defer db.Close()

	adapter := Adapter{db: db}
	ctx := context.Background()

	// Arrange: Database mocks
	attrName := "foo"
	attrValue := []byte("42")
	event := query.Event{
		ID:     1,
		TXHash: "ABC123",
		Index:  0,
		Type:   "test",
		Attributes: []query.Attribute{
			query.NewAttribute(attrName, attrValue),
		},
		CreatedAt: time.Now(),
	}

	eventRows := sqlmock.
		NewRows(eventFields).
		AddRow(event.ID, event.Index, event.TXHash, event.Type, event.CreatedAt)
	eventAttrRows := sqlmock.
		NewRows(eventAttrFields).
		AddRow(event.ID, attrName, attrValue)

	mock.
		ExpectQuery(`
			SELECT event.id, event.index, event.tx_hash, event.type, event.created_at
			FROM event INNER JOIN tx ON event.tx_hash = tx.hash
			WHERE true
			ORDER BY tx.height, tx.index, event.index
			LIMIT 30 OFFSET 0
		`).
		WillReturnRows(eventRows)
	mock.
		ExpectQuery(`
			SELECT event_id, name, value FROM attribute
            WHERE event_id = ANY($1)
            ORDER BY event_id
		`).
		WillReturnRows(eventAttrRows).
		WithArgs(pq.Array([]int64{event.ID}))

	// Arrange: Expectations
	wantEvents := []query.Event{event}

	// Arrange: Query
	qry := query.NewEventQuery()

	// Act
	events, err := adapter.QueryEvents(ctx, qry)

	// Assert
	require.NoError(t, err, "expected no query errors on execution")
	require.Len(t, events, 1)
	require.Equal(t, wantEvents, events)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestEventQueryWithFilters(t *testing.T) {
	// Arrange
	db, mock := createMatchEqualSQLMock(t)
	defer db.Close()

	adapter := Adapter{db: db}
	ctx := context.Background()

	// Arrange: Database mocks
	indexValue := 2
	typeValue := "chain.test.Test"
	hashValues := []string{"HASH1", "HASH2"}
	eventRows := sqlmock.NewRows(eventFields)

	mock.
		ExpectQuery(`
			SELECT event.id, event.index, event.tx_hash, event.type, event.created_at
			FROM event INNER JOIN tx ON event.tx_hash = tx.hash
			WHERE event.index = $1 AND event.type = $2 AND event.tx_hash = ANY($3)
			ORDER BY tx.height, tx.index, event.index
			LIMIT 30 OFFSET 0
		`).
		WillReturnRows(eventRows).
		WithArgs(indexValue, typeValue, pq.Array(hashValues))

	// Arrange: Query
	qry := query.NewEventQuery(
		query.WithFilters(
			NewFilter("event.index", indexValue),
			FilterByEventType(typeValue),
			FilterByEventTXs(hashValues...),
		),
	)

	// Act
	events, err := adapter.QueryEvents(ctx, qry)

	// Assert
	require.NoError(t, err, "expected no query errors on execution")
	require.Len(t, events, 0)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestEventQueryWithEventAttrFilters(t *testing.T) {
	// Arrange
	db, mock := createMatchEqualSQLMock(t)
	defer db.Close()

	adapter := Adapter{db: db}
	ctx := context.Background()

	// Arrange: Database mocks
	attrNameValue := "foo"
	attrValue := int64(42)
	eventRows := sqlmock.NewRows(eventFields)

	mock.
		ExpectQuery(`
			SELECT DISTINCT events.*
			FROM (
				SELECT event.id, event.index, event.tx_hash, event.type, event.created_at
				FROM event
					INNER JOIN tx ON event.tx_hash = tx.hash
					INNER JOIN attribute ON event.id = attribute.event_id
				WHERE attribute.name = $1 AND attribute.name = $2 AND attribute.value::numeric = $3
				ORDER BY tx.height, tx.index, event.index
			) AS events
			LIMIT 30 OFFSET 0
		`).
		WillReturnRows(eventRows).
		WithArgs(attrNameValue, attrNameValue, attrValue)

	// Arrange: Query
	qry := query.NewEventQuery(
		query.WithFilters(
			NewFilter("attribute.name", attrNameValue),
			FilterByEventAttrName(attrNameValue),
			FilterByEventAttrValueInt(attrValue),
		),
	)

	// Act
	events, err := adapter.QueryEvents(ctx, qry)

	// Assert
	require.NoError(t, err, "expected no query errors on execution")
	require.Len(t, events, 0)
	require.NoError(t, mock.ExpectationsWereMet())
}

func createMatchEqualSQLMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	require.NoError(t, err)

	return db, mock
}
