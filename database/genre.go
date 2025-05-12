package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattn/go-sqlite3"
)

type Genre struct {
	Slug string `db:"slug"`
	Name string `db:"name"`
}

func GenreQuery() *goqu.SelectDataset {
	query := dialect.From("genres").
		Select(
			"genres.slug",
			"genres.name",
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllGenres(ctx context.Context) ([]Genre, error) {
	query := GenreQuery()

	var items []Genre
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetGenreBySlug(ctx context.Context, slug string) (Genre, error) {
	query := GenreQuery().
		Where(goqu.I("genres.slug").Eq(slug))

	var item Genre
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Genre{}, ErrItemNotFound
		}

		return Genre{}, err
	}

	return item, nil
}

func (db *Database) CreateGenre(ctx context.Context, slug, name string) error {
	query := dialect.Insert("genres").
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
