package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/utils"
)

type ShowSeason struct {
	RowId int `db:"rowid"`

	Num    int    `db:"num"`
	ShowId string `db:"show_id"`

	MediaId string `db:"media_id"`

	Name       string `db:"name"`
	SearchSlug string `db:"search_slug"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Use goqu.T more
func ShowSeasonQuery() *goqu.SelectDataset {
	query := dialect.From("show_seasons").
		Select(
			"show_seasons.rowid",

			"show_seasons.num",
			"show_seasons.show_id",

			"show_seasons.media_id",

			"show_seasons.name",
			"show_seasons.search_slug",

			"show_seasons.created",
			"show_seasons.updated",
		)

	return query
}

func (db DB) GetAllShowSeasons(ctx context.Context) ([]ShowSeason, error) {
	query := ShowSeasonQuery()
	return ember.Multiple[ShowSeason](db.db, ctx, query)
}

func (db DB) GetAllShowSeasonsByShowId(ctx context.Context, showId string) ([]ShowSeason, error) {
	query := ShowSeasonQuery().
		Where(goqu.I("show_seasons.show_id").Eq(showId)).
		Order(goqu.I("show_seasons.num").Asc())
	return ember.Multiple[ShowSeason](db.db, ctx, query)
}

func (db DB) GetShowSeasonById(ctx context.Context, num int, showId string) (ShowSeason, error) {
	query := ShowSeasonQuery().
		Where(
			goqu.I("show_seasons.num").Eq(num),
			goqu.I("show_seasons.show_id").Eq(showId),
		)

	return ember.Single[ShowSeason](db.db, ctx, query)
}

type CreateShowSeasonParams struct {
	Num    int
	ShowId string

	MediaId string

	Name       string
	SearchSlug string

	Created int64
	Updated int64
}

func (db DB) CreateShowSeason(ctx context.Context, params CreateShowSeasonParams) error {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.SearchSlug == "" {
		params.SearchSlug = utils.Slug(params.Name)
	}

	query := dialect.Insert("show_seasons").Rows(goqu.Record{
		"num":     params.Num,
		"show_id": params.ShowId,

		"media_id": params.MediaId,

		"name":        params.Name,
		"search_slug": params.SearchSlug,

		"created": params.Created,
		"updated": params.Updated,
	})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type ShowSeasonChanges struct {
	Num Change[int]

	MediaId Change[string]

	Name       Change[string]
	SearchSlug Change[string]

	Created Change[int64]
}

func (db DB) UpdateShowSeason(ctx context.Context, num int, showId string, changes ShowSeasonChanges) error {
	record := goqu.Record{}

	addToRecord(record, "num", changes.Num)

	addToRecord(record, "media_id", changes.MediaId)

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "search_slug", changes.SearchSlug)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("show_seasons").
		Set(record).
		Where(
			goqu.I("show_seasons.num").Eq(num),
			goqu.I("show_seasons.show_id").Eq(showId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) RemoveShowSeason(ctx context.Context, num int, showId string) error {
	query := dialect.Delete("show_seasons").
		Where(
			goqu.I("show_seasons.num").Eq(num),
			goqu.I("show_seasons.show_id").Eq(showId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
