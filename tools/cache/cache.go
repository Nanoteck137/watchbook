package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
)

var dialect = ember.SqliteDialect()

type Cache struct {
	db *ember.Database
}

func Open(dbFile string) (*Cache, error) {
	dbUrl := fmt.Sprintf("file:%s?_busy_timeout=5000&_journal_mode=WAL&_synchronous=NORMAL&_temp_store=2&_cache_size=-2000", dbFile)
	db, err := ember.OpenDatabase("sqlite3", dbUrl)
	if err != nil {
		return nil, err
	}

	db.Exec(context.Background(), ember.RawQuery{
		Sql: `
		CREATE TABLE IF NOT EXISTS cache (
			key TEXT PRIMARY KEY,
			value BLOB NOT NULL,
			expires_at TIMESTAMP NOT NULL
		)
		`,
	})
	// db.ErrorHandler = handleErr

	return &Cache{
		db: db,
	}, nil
}

func (c *Cache) Get(key string) ([]byte, bool) {
	var value []byte
	var expiresAt time.Time

	query := dialect.From("cache").
		Select("value", "expires_at").
		Where(goqu.I("cache.key").Eq(key))
	// err := c.db.QueryRow(context.Background(), `SELECT value, expires_at FROM cache WHERE key = ?`, key).Scan(&value, &expiresAt)
	row, err := c.db.QueryRow(context.Background(), query)
	if err != nil {
		return nil, false
	}

	err = row.Scan(&value, &expiresAt)
	if err != nil {
		return nil, false
	}

	if time.Now().After(expiresAt) {
		query := dialect.Delete("cache").Where(goqu.I("cache.key").Eq(key))
		c.db.Exec(context.Background(), query)

		return nil, false
	}

	return value, true
}

func (c *Cache) Set(key string, value []byte, ttl time.Duration) error {
	expiresAt := time.Now().Add(ttl)

	query := dialect.Insert("cache").Rows(goqu.Record{
		"key": key,
		"value": value,
		"expires_at": expiresAt,
	}).
	OnConflict(
		goqu.DoUpdate("key", goqu.Record{
			"value": value,
			"expires_at": expiresAt,
		}),
	)

	_, err := c.db.Exec(context.Background(), query)
	return err
}

var ErrNoData = errors.New("no data in cache")

func GetJson[T any](cache *Cache, key string) (T, error) {
	var res T

	d, hasData := cache.Get(key)
	if !hasData {
		return res, ErrNoData
	}

	err := json.Unmarshal(d, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func SetJson[T any](cache *Cache, key string, data T, ttl time.Duration) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = cache.Set(key, d, ttl)
	if err != nil {
		return err
	}

	return nil
}
