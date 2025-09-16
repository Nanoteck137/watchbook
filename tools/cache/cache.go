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

type Cache interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte, ttl time.Duration) error
}

var dialect = ember.SqliteDialect()
var table = goqu.T("provider_cache")

type ProviderCache struct {
	db *ember.Database
}

func OpenDatabase(dbFile string) (*ember.Database, error) {
	dbUrl := fmt.Sprintf("file:%s?_busy_timeout=5000&_journal_mode=WAL&_synchronous=NORMAL&_temp_store=2&_cache_size=-2000", dbFile)
	return ember.OpenDatabase("sqlite3", dbUrl)
}

func NewProvider(db *ember.Database) (*ProviderCache, error) {
	_, err := db.Exec(context.Background(), ember.RawQuery{
		Sql: `
		CREATE TABLE IF NOT EXISTS provider_cache (
			name TEXT NOT NULL,
			key TEXT NOT NULL,
			value BLOB NOT NULL,
			expires_at TIMESTAMP NOT NULL,

			PRIMARY KEY(name, key)
		)
		`,
	})
	if err != nil {
		return nil, err
	}
	// db.ErrorHandler = handleErr

	return &ProviderCache{
		db: db,
	}, nil
}

func (c *ProviderCache) WithName(name string) NamedProviderCache {
	return NamedProviderCache{
		cache: c,
		name:  name,
	}
}

func (c *ProviderCache) Get(name, key string) ([]byte, bool) {
	var value []byte
	var expiresAt time.Time

	query := dialect.From(table).
		Select("value", "expires_at").
		Where(
			table.Col("name").Eq(name),
			table.Col("key").Eq(key),
		)
	row, err := c.db.QueryRow(context.Background(), query)
	if err != nil {
		return nil, false
	}

	err = row.Scan(&value, &expiresAt)
	if err != nil {
		return nil, false
	}

	if time.Now().After(expiresAt) {
		query := dialect.Delete(table).
			Where(
				table.Col("name").Eq(name),
				table.Col("key").Eq(key),
			)
		c.db.Exec(context.Background(), query)

		return nil, false
	}

	return value, true
}

func (c *ProviderCache) Set(name, key string, value []byte, ttl time.Duration) error {
	expiresAt := time.Now().Add(ttl)

	query := dialect.Insert(table).Rows(goqu.Record{
		"name":       name,
		"key":        key,
		"value":      value,
		"expires_at": expiresAt,
	}).
		OnConflict(
			goqu.DoUpdate("name, key", goqu.Record{
				"value":      value,
				"expires_at": expiresAt,
			}),
		)

	_, err := c.db.Exec(context.Background(), query)
	return err
}

func (c *ProviderCache) ClearByProviderName(name string) error {
	query := dialect.Delete(table).
		Where(
			table.Col("name").Eq(name),
		)
	_, err := c.db.Exec(context.Background(), query)
	return err
}

func (c *ProviderCache) Clear() error {
	query := dialect.Delete(table)
	_, err := c.db.Exec(context.Background(), query)
	return err
}

var _ Cache = (*NamedProviderCache)(nil)

type NamedProviderCache struct {
	cache *ProviderCache
	name  string
}

func (p NamedProviderCache) Get(key string) ([]byte, bool) {
	return p.cache.Get(p.name, key)
}

func (p NamedProviderCache) Set(key string, value []byte, ttl time.Duration) error {
	return p.cache.Set(p.name, key, value, ttl)
}

var ErrNoData = errors.New("no data in cache")

func GetJson[T any](cache Cache, key string) (T, bool) {
	var res T

	d, hasData := cache.Get(key)
	if !hasData {
		return res, false
	}

	err := json.Unmarshal(d, &res)
	if err != nil {
		return res, false
	}

	return res, true
}

func SetJson[T any](cache Cache, key string, data T, ttl time.Duration) error {
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
