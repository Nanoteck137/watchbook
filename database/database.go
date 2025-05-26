package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/database/migrations"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

var ErrItemNotFound = errors.New("database: item not found")
var ErrItemAlreadyExists = errors.New("database: item already exists")

var dialect = ember.SqliteDialect()

type Database struct {
	db *ember.Database
}

func (db *Database) RunMigrateUp() error {
	return migrations.RunMigrateUp(db.db.DB.DB)
}

func (db *Database) RunMigrateDown() error {
	return migrations.RunMigrateDown(db.db.DB.DB)
}

func Open(dbFile string) (*Database, error) {
	dbUrl := fmt.Sprintf("file:%s?_foreign_keys=true", dbFile)
	db, err := ember.OpenDatabase("sqlite3", dbUrl)
	if err != nil {
		return nil, err
	}

	db.ErrorHandler = handleErr

	return &Database{
		db: db,
	}, nil
}

func handleErr(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrItemNotFound
	}

	var e sqlite3.Error
	if errors.As(err, &e) {
		switch e.ExtendedCode {
		case sqlite3.ErrConstraintPrimaryKey:
			return ErrItemAlreadyExists
		}
	}

	return err
}
