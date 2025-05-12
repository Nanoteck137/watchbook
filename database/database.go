package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	goqusqlite3 "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	"github.com/jmoiron/sqlx"
	"github.com/nanoteck137/watchbook/database/migrations"
	"github.com/nanoteck137/watchbook/types"
	"github.com/pressly/goose/v3"

	_ "github.com/mattn/go-sqlite3"
)

var ErrItemNotFound = errors.New("database: item not found")
var ErrItemAlreadyExists = errors.New("database: item already exists")

var dialect = goqu.Dialect("sqlite_returning")

type ToSQL interface {
	ToSQL() (string, []interface{}, error)
}

type Connection interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

type Database struct {
	NewRawConn *sqlx.DB
	RawConn    *sql.DB
	Conn       Connection
}

func New(conn *sql.DB) *Database {
	return &Database{
		NewRawConn: sqlx.NewDb(conn, "sqlite3"),
		RawConn:    conn,
		Conn:       conn,
	}
}

func Open(workDir types.WorkDir) (*Database, error) {
	dbUrl := fmt.Sprintf("file:%s?_foreign_keys=true", workDir.DatabaseFile())

	conn, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		return nil, err
	}

	return New(conn), nil
}

func (db *Database) RunMigrateUp() error {
	return migrations.RunMigrateUp(db.RawConn)
}

func (db *Database) RunMigrateDown() error {
	return migrations.RunMigrateDown(db.RawConn)
}

func (db *Database) Begin() (*Database, *sqlx.Tx, error) {
	tx, err := db.NewRawConn.Beginx()
	if err != nil {
		return nil, nil, err
	}

	return &Database{
		NewRawConn: db.NewRawConn,
		RawConn:    db.RawConn,
		Conn:       tx,
	}, tx, nil
}

func (db *Database) Query(ctx context.Context, s ToSQL) (*sql.Rows, error) {
	sql, params, err := s.ToSQL()
	if err != nil {
		return nil, err
	}

	return db.Conn.Query(sql, params...)
}

func (db *Database) QueryRow(ctx context.Context, s ToSQL) (*sql.Row, error) {
	sql, params, err := s.ToSQL()
	if err != nil {
		return nil, err
	}

	row := db.Conn.QueryRow(sql, params...)
	return row, nil
}

func (db *Database) Exec(ctx context.Context, s ToSQL) (sql.Result, error) {
	sql, params, err := s.ToSQL()
	if err != nil {
		return nil, err
	}

	return db.Conn.Exec(sql, params...)
}

func (db *Database) Select(dest any, s ToSQL) error {
	sql, params, err := s.ToSQL()
	if err != nil {
		return err
	}

	return db.NewRawConn.Select(dest, sql, params...)
}

func (db *Database) Get(dest any, s ToSQL) error {
	sql, params, err := s.ToSQL()
	if err != nil {
		return err
	}

	return db.NewRawConn.Get(dest, sql, params...)
}

func init() {
	opts := goqusqlite3.DialectOptions()
	opts.SupportsReturn = true
	goqu.RegisterDialect("sqlite_returning", opts)
}

func RunMigrateUp(db *Database) error {
	return goose.Up(db.RawConn, ".")
}
