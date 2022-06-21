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

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
	"github.com/ignite-hq/cli/ignite/pkg/cosmostxcollector/query"
	"github.com/ignite-hq/cli/ignite/pkg/cosmostxcollector/query/call"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
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

func TestQueryWithArg(t *testing.T) {
	// Arrange
	var (
		cursorNextSucceded bool
		rowValue           string
	)

	db, mock := createMatchEqualSQLMock(t)
	defer db.Close()

	adapter := Adapter{db: db}
	ctx := context.Background()

	// Arrange: Call Query
	wantArg := "bar"

	c := call.New("baz", call.WithFields("foo"))
	qry := query.
		NewCall(c).
		AppendFilters(
			NewFilter("foo", wantArg),
		).
		WithoutPaging()

	// Arrange: Database mock and expectations
	wantRowValue := "expected"

	mock.
		ExpectQuery("SELECT DISTINCT foo FROM baz WHERE foo = $1").
		WithArgs(wantArg).
		WillReturnRows(
			sqlmock.NewRows([]string{"foo"}).AddRow(wantRowValue),
		)

	// Act
	cr, err := adapter.Query(ctx, qry)
	require.NoError(t, err, "expected no query errors on execution")

	if cursorNextSucceded = cr.Next(); cursorNextSucceded {
		err = cr.Scan(&rowValue)
	}

	// Assert
	require.NoError(t, err)
	require.True(t, cursorNextSucceded, "expected cursor.Next() to succeed")
	require.Equal(t, wantRowValue, rowValue)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestQueryWithoutArgs(t *testing.T) {
	// Arrange
	var (
		cursorNextSucceded bool
		rowValue           string
	)

	db, mock := createMatchEqualSQLMock(t)
	defer db.Close()

	adapter := Adapter{db: db}
	ctx := context.Background()

	// Arrange: Call Query
	c := call.New("baz", call.WithFields("foo"))
	qry := query.NewCall(c).WithoutPaging()

	// Arrange: Database mock and expectations
	wantRowValue := "expected"

	mock.
		ExpectQuery("SELECT DISTINCT foo FROM baz").
		WillReturnRows(
			sqlmock.NewRows([]string{"foo"}).AddRow(wantRowValue),
		)

	// Act
	cr, err := adapter.Query(ctx, qry)
	require.NoError(t, err, "expected no query errors on execution")

	if cursorNextSucceded = cr.Next(); cursorNextSucceded {
		err = cr.Scan(&rowValue)
	}

	// Assert
	require.NoError(t, err)
	require.True(t, cursorNextSucceded, "expected cursor.Next() to succeed")
	require.Equal(t, wantRowValue, rowValue)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestQueryError(t *testing.T) {
	// Arrange
	db, mock := createMatchEqualSQLMock(t)
	defer db.Close()

	adapter := Adapter{db: db}
	ctx := context.Background()

	// Arrange: Call Query
	c := call.New("baz")
	qry := query.NewCall(c).WithoutPaging()

	// Arrange: Database mock and expectations
	wantErr := errors.New("expected error")

	mock.
		ExpectQuery("SELECT * FROM baz").
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

	// Arrange: Call Query
	c := call.New("baz", call.WithFields(cols[0]))
	qry := query.NewCall(c).WithoutPaging()

	// Arrange: Database mock and expectations
	wantErr := errors.New("expected error")

	row := sqlmock.
		NewRows(cols).
		AddRow("foo").
		RowError(0, wantErr)

	mock.
		ExpectQuery("SELECT DISTINCT name FROM baz").
		WillReturnRows(row)

	// Act
	cr, err := adapter.Query(ctx, qry)
	require.NoError(t, err, "expected no query errors on execution")

	cursorNextSucceded := cr.Next()

	// Assert
	require.NoError(t, err)
	require.False(t, cursorNextSucceded, "expected cursor.Next() to fail")
	require.Equal(t, wantErr, cr.Err())
	require.NoError(t, mock.ExpectationsWereMet())
}

func createMatchEqualSQLMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	require.NoError(t, err)

	return db, mock
}
