package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
)

type AnimeEpisode struct {
	RowId int `db:"rowid"`

	Index   int64  `db:"idx"`
	AnimeId string `db:"anime_id"`

	Name string `db:"name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Use goqu.T more
func AnimeEpisodeQuery() *goqu.SelectDataset {
	query := dialect.From("anime_episodes").
		Select(
			"anime_episodes.rowid",

			"anime_episodes.idx",
			"anime_episodes.anime_id",

			"anime_episodes.name",

			"anime_episodes.created",
			"anime_episodes.updated",
		)

	return query
}

func (db *Database) GetAllAnimeEpisodes(ctx context.Context) ([]AnimeEpisode, error) {
	query := AnimeEpisodeQuery()
	return ember.Multiple[AnimeEpisode](db.db, ctx, query)
}

func (db *Database) GetAnimeEpisodesByAnimeId(ctx context.Context, animeId string) ([]AnimeEpisode, error) {
	query := AnimeEpisodeQuery().
		Where(goqu.I("anime_episodes.anime_id").Eq(animeId)).
		Order(goqu.I("anime_episodes.idx").Asc())

	return ember.Multiple[AnimeEpisode](db.db, ctx, query)
}

type CreateAnimeEpisodeParams struct {
	Index int64
	AnimeId string

	Name  string

	Created int64
	Updated int64
}

func (db *Database) CreateAnimeEpisode(ctx context.Context, params CreateAnimeEpisodeParams) error {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	query := dialect.Insert("anime_episodes").Rows(goqu.Record{
		"idx":  params.Index,
		"anime_id": params.AnimeId,

		"name": params.Name,

		"created": created,
		"updated": updated,
	})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type AnimeEpisodeChanges struct {
	Name  Change[string]

	Created Change[int64]
}

func (db *Database) UpdateAnimeEpisode(ctx context.Context, id string, changes AnimeEpisodeChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("anime_episodes").
		Set(record).
		Where(goqu.I("animes.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveAnimeEpisode(ctx context.Context, idx int64, animeId string) error {
	query := dialect.Delete("anime_episodes").
		Where(
			goqu.I("anime_episodes.idx").Eq(idx),
			goqu.I("anime_episodes.anime_id").Eq(animeId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
