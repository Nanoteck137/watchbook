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
var table = goqu.T("provider_cache")

type Provider struct {
	db *ember.Database
}

func OpenDatabase(dbFile string) (*ember.Database, error) {
	dbUrl := fmt.Sprintf("file:%s?_busy_timeout=5000&_journal_mode=WAL&_synchronous=NORMAL&_temp_store=2&_cache_size=-2000", dbFile)
	return ember.OpenDatabase("sqlite3", dbUrl)
}

func NewProvider(db *ember.Database) (*Provider, error) {
	_, err := db.Exec(context.Background(), ember.RawQuery{
		Sql: `
		CREATE TABLE IF NOT EXISTS provider_cache (
			key TEXT PRIMARY KEY,
			provider_name TEXT NOT NULL,
			value BLOB NOT NULL,
			expires_at TIMESTAMP NOT NULL
		)
		`,
	})
	if err != nil {
		return nil, err
	}
	// db.ErrorHandler = handleErr

	return &Provider{
		db: db,
	}, nil
}

func (c *Provider) Get(key string) ([]byte, bool) {
	var value []byte
	var expiresAt time.Time

	query := dialect.From(table).
		Select("value", "expires_at").
		Where(table.Col("key").Eq(key))
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
			Where(table.Col("cache.key").Eq(key))
		c.db.Exec(context.Background(), query)

		return nil, false
	}

	return value, true
}

func (c *Provider) Set(key, providerName string, value []byte, ttl time.Duration) error {
	expiresAt := time.Now().Add(ttl)

	query := dialect.Insert(table).Rows(goqu.Record{
		"key":           key,
		"provider_name": providerName,
		"value":         value,
		"expires_at":    expiresAt,
	}).
		OnConflict(
			goqu.DoUpdate("key", goqu.Record{
				"value":         value,
				"provider_name": providerName,
				"expires_at":    expiresAt,
			}),
		)

	_, err := c.db.Exec(context.Background(), query)
	return err
}

func (c *Provider) ClearByProviderName(providerName string) error {
	query := dialect.Delete(table).
		Where(table.Col("provider_name").Eq(providerName))
	_, err := c.db.Exec(context.Background(), query)
	return err
}

func (c *Provider) Clear() error {
	query := dialect.Delete(table)
	_, err := c.db.Exec(context.Background(), query)
	return err
}

var ErrNoData = errors.New("no data in cache")

func GetProviderJson[T any](cache *Provider, key string) (T, error) {
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

func SetProviderJson[T any](cache *Provider, key, providerName string, data T, ttl time.Duration) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = cache.Set(key, providerName, d, ttl)
	if err != nil {
		return err
	}

	return nil
}
