package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
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

	return ember.Multiple[Studio](db.db, ctx, query)
}

func (db *Database) GetStudioBySlug(ctx context.Context, slug string) (Studio, error) {
	query := StudioQuery().
		Where(goqu.I("studios.slug").Eq(slug))

	return ember.Single[Studio](db.db, ctx, query)
}

func (db *Database) CreateStudio(ctx context.Context, slug, name string) error {
	query := dialect.Insert("studios").
		Rows(goqu.Record{
			"slug": slug,
			"name": name,
		}).
		Prepared(true)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
