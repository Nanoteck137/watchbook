package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
)

type Tag struct {
	Slug string `db:"slug"`
	Name string `db:"name"`
}

func TagQuery() *goqu.SelectDataset {
	query := dialect.From("tags").
		Select(
			"tags.slug",
			"tags.name",
		)

	return query
}

func (db *Database) GetAllTags(ctx context.Context) ([]Tag, error) {
	query := TagQuery()

	return ember.Multiple[Tag](db.db, ctx, query)
}

func (db *Database) GetTagBySlug(ctx context.Context, slug string) (Tag, error) {
	query := TagQuery().
		Where(goqu.I("tags.slug").Eq(slug))

	return ember.Single[Tag](db.db, ctx, query)
}

func (db *Database) CreateTag(ctx context.Context, slug, name string) error {
	query := dialect.Insert("tags").
		Rows(goqu.Record{
			"slug": slug,
			"name": name,
		})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
