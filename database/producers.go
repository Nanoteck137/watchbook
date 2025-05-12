package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattn/go-sqlite3"
)

type Producer struct {
	Slug string `db:"slug"`
	Name string `db:"name"`
}

func ProducerQuery() *goqu.SelectDataset {
	query := dialect.From("producers").
		Select(
			"producers.slug",
			"producers.name",
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllProducers(ctx context.Context) ([]Producer, error) {
	query := ProducerQuery()

	var items []Producer
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetProducerBySlug(ctx context.Context, slug string) (Producer, error) {
	query := ProducerQuery().
		Where(goqu.I("producers.slug").Eq(slug))

	var item Producer
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Producer{}, ErrItemNotFound
		}

		return Producer{}, err
	}

	return item, nil
}

func (db *Database) CreateProducer(ctx context.Context, slug, name string) error {
	query := dialect.Insert("producers").
		Rows(goqu.Record{
			"slug": slug,
			"name": name,
		}).
		Prepared(true)

	_, err := db.Exec(ctx, query)
	if err != nil {
		var e sqlite3.Error
		if errors.As(err, &e) {
			if e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return ErrItemAlreadyExists
			}
		}

		return err
	}

	return nil
}
