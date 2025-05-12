package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattn/go-sqlite3"
)

type Theme struct {
	Slug string `db:"slug"`
	Name string `db:"name"`
}

func ThemeQuery() *goqu.SelectDataset {
	query := dialect.From("themes").
		Select(
			"themes.slug",
			"themes.name",
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllThemes(ctx context.Context) ([]Theme, error) {
	query := ThemeQuery()

	var items []Theme
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetThemeBySlug(ctx context.Context, slug string) (Theme, error) {
	query := ThemeQuery().
		Where(goqu.I("themes.slug").Eq(slug))

	var item Theme
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Theme{}, ErrItemNotFound
		}

		return Theme{}, err
	}

	return item, nil
}

func (db *Database) CreateTheme(ctx context.Context, slug, name string) error {
	query := dialect.Insert("themes").
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
