package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattn/go-sqlite3"
)

type Studio struct {
	Slug string `db:"slug"`
	Name string `db:"name"`
}

func StudioQuery() *goqu.SelectDataset {
	query := dialect.From("studios").
		Select(
			"studios.slug",
			"studios.name",
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllStudios(ctx context.Context) ([]Studio, error) {
	query := StudioQuery()

	var items []Studio
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetStudioBySlug(ctx context.Context, slug string) (Studio, error) {
	query := StudioQuery().
		Where(goqu.I("studios.slug").Eq(slug))

	var item Studio
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Studio{}, ErrItemNotFound
		}

		return Studio{}, err
	}

	return item, nil
}

func (db *Database) CreateStudio(ctx context.Context, slug, name string) error {
	query := dialect.Insert("studios").
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
