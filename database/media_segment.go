package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/types"
)

type MediaSegment struct {
	RowId int `db:"rowid"`

	Idx     int    `db:"idx"`
	MediaId string `db:"media_id"`

	Type types.MediaSegmentType `db:"type"`

	Title       string         `db:"title"`
	Description sql.NullString `db:"description"`

	Score        sql.NullFloat64   `db:"score"`
	Status       types.MediaStatus `db:"status"`
	Rating       types.MediaRating `db:"rating"`
	AiringSeason sql.NullString    `db:"airing_season"`

	StartDate sql.NullString `db:"start_date"`
	EndDate   sql.NullString `db:"end_date"`

	CoverFile sql.NullString `db:"cover_file"`
	// TODO(patrik): Remove?
	LogoFile   sql.NullString `db:"logo_file"`
	BannerFile sql.NullString `db:"banner_file"`

	DefaultProvider sql.NullString `db:"default_provider"`
	Providers       ember.KVStore  `db:"providers"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	Creators ember.JsonColumn[[]string] `db:"creators"`
	Tags     ember.JsonColumn[[]string] `db:"tags"`
}

func MediaSegmentTagQuery() *goqu.SelectDataset {
	tbl := goqu.T("media_tags")

	return dialect.From(tbl).
		Select(
			tbl.Col("media_id").As("id"),
			goqu.Func(
				"json_group_array",
				tbl.Col("tag_slug"),
			).As("data"),
		).
		GroupBy(tbl.Col("media_id"))
}

func MediaSegmentCreatorQuery() *goqu.SelectDataset {
	tbl := goqu.T("media_creators")

	return dialect.From(tbl).
		Select(
			tbl.Col("media_id").As("id"),
			goqu.Func(
				"json_group_array",
				tbl.Col("tag_slug"),
			).As("data"),
		).
		GroupBy(tbl.Col("media_id"))
}

// TODO(patrik): Use goqu.T more
func MediaSegmentQuery() *goqu.SelectDataset {
	// creatorsQuery := MediaSegmentCreatorQuery()
	// tagsQuery := MediaSegmentTagQuery()

	query := dialect.From("media_segments").
		Select(
			"media_segments.rowid",

			"media_segments.idx",
			"media_segments.media_id",

			"media_segments.type",

			"media_segments.title",
			"media_segments.description",

			"media_segments.score",
			"media_segments.status",
			"media_segments.rating",
			"media_segments.airing_season",

			"media_segments.start_date",
			"media_segments.end_date",

			"media_segments.cover_file",
			"media_segments.logo_file",
			"media_segments.banner_file",

			"media_segments.default_provider",
			"media_segments.providers",

			"media_segments.created",
			"media_segments.updated",

			// goqu.I("creators.data").As("creators"),
			// goqu.I("tags.data").As("tags"),
		)
		// LeftJoin(
		// 	creatorsQuery.As("creators"),
		// 	goqu.On(goqu.I("media_segments.id").Eq(goqu.I("creators.id"))),
		// ).
		// LeftJoin(
		// 	tagsQuery.As("tags"),
		// 	goqu.On(goqu.I("media_segments.id").Eq(goqu.I("tags.id"))),
		// ).

	return query
}

func (db DB) GetAllMediaSegment(ctx context.Context) ([]MediaSegment, error) {
	query := MediaSegmentQuery()
	return ember.Multiple[MediaSegment](db.db, ctx, query)
}

func (db DB) GetAllMediaSegmentForUnknownUpdate(ctx context.Context) ([]MediaSegment, error) {
	query := MediaSegmentQuery().
		Where(
			goqu.I("media_segments.type").In(types.MediaSegmentTypeUnknown),
		)
	return ember.Multiple[MediaSegment](db.db, ctx, query)
}

func (db DB) GetMediaSegmentById(ctx context.Context, id string) (MediaSegment, error) {
	query := MediaSegmentQuery().
		Where(goqu.I("media_segments.id").Eq(id))

	return ember.Single[MediaSegment](db.db, ctx, query)
}

func (db DB) GetMediaSegmentByProviderId(ctx context.Context, providerName, value string) (MediaSegment, error) {
	query := MediaSegmentQuery().
		Where(
			goqu.Func("json_extract", goqu.I("media_segments.providers"), "$."+providerName).Eq(value),
		)

	return ember.Single[MediaSegment](db.db, ctx, query)
}

type CreateMediaSegmentParams struct {
	Idx     int
	MediaId string

	Type types.MediaSegmentType

	Title string

	Description sql.NullString

	Score        sql.NullFloat64
	Status       types.MediaStatus
	Rating       types.MediaRating
	AiringSeason sql.NullString

	StartDate sql.NullString
	EndDate   sql.NullString

	CoverFile  sql.NullString
	LogoFile   sql.NullString
	BannerFile sql.NullString

	DefaultProvider sql.NullString
	Providers       ember.KVStore

	Created int64
	Updated int64
}

func (db DB) CreateMediaSegment(ctx context.Context, params CreateMediaSegmentParams) error {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	if params.Type == "" {
		params.Type = types.MediaSegmentTypeUnknown
	}

	if params.Status == "" {
		params.Status = types.MediaStatusUnknown
	}

	if params.Rating == "" {
		params.Rating = types.MediaRatingUnknown
	}

	query := dialect.Insert("media_segments").Rows(goqu.Record{
		"idx":   params.Idx,
		"media_id":   params.MediaId,

		"type": params.Type,

		"title":       params.Title,
		"description": params.Description,

		"score":         params.Score,
		"status":        params.Status,
		"rating":        params.Rating,
		"airing_season": params.AiringSeason,

		"start_date": params.StartDate,
		"end_date":   params.EndDate,

		"cover_file":  params.CoverFile,
		"logo_file":   params.LogoFile,
		"banner_file": params.BannerFile,

		"default_provider": params.DefaultProvider,
		"providers":        params.Providers,

		"created": created,
		"updated": updated,
	})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type MediaSegmentChanges struct {
	Type Change[types.MediaSegmentType]

	Title       Change[string]
	Description Change[sql.NullString]

	Score        Change[sql.NullFloat64]
	Status       Change[types.MediaStatus]
	Rating       Change[types.MediaRating]
	AiringSeason Change[sql.NullString]

	StartDate Change[sql.NullString]
	EndDate   Change[sql.NullString]

	CoverFile  Change[sql.NullString]
	LogoFile   Change[sql.NullString]
	BannerFile Change[sql.NullString]

	DefaultProvider Change[sql.NullString]
	Providers       Change[ember.KVStore]

	Created Change[int64]
}

func (db DB) UpdateMediaSegment(ctx context.Context, id string, changes MediaChanges) error {
	record := goqu.Record{}

	addToRecord(record, "type", changes.Type)

	addToRecord(record, "title", changes.Title)
	addToRecord(record, "description", changes.Description)

	addToRecord(record, "score", changes.Score)
	addToRecord(record, "status", changes.Status)
	addToRecord(record, "rating", changes.Rating)
	addToRecord(record, "airing_season", changes.AiringSeason)

	addToRecord(record, "start_date", changes.StartDate)
	addToRecord(record, "end_date", changes.EndDate)

	addToRecord(record, "cover_file", changes.CoverFile)
	addToRecord(record, "logo_file", changes.LogoFile)
	addToRecord(record, "banner_file", changes.BannerFile)

	addToRecord(record, "default_provider", changes.DefaultProvider)
	addToRecord(record, "providers", changes.Providers)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("media_segments").
		Set(record).
		Where(goqu.I("media_segments.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) RemoveMediaSegment(ctx context.Context, id string) error {
	query := dialect.Delete("media_segments").
		Where(goqu.I("media_segments.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// func (db DB) AddTagToMediaSegment(ctx context.Context, mediaId, tagSlug string) error {
// 	query := dialect.Insert("media_tags").
// 		Rows(goqu.Record{
// 			"media_id": mediaId,
// 			"tag_slug": tagSlug,
// 		})
//
// 	_, err := db.db.Exec(ctx, query)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func (db DB) RemoveAllTagsFromMediaSegment(ctx context.Context, mediaId string) error {
// 	query := dialect.Delete("media_tags").
// 		Where(goqu.I("media_tags.media_id").Eq(mediaId))
//
// 	_, err := db.db.Exec(ctx, query)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func (db DB) AddCreatorToMediaSegment(ctx context.Context, mediaId, tagSlug string) error {
// 	query := dialect.Insert("media_creators").
// 		Rows(goqu.Record{
// 			"media_id": mediaId,
// 			"tag_slug": tagSlug,
// 		})
//
// 	_, err := db.db.Exec(ctx, query)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func (db DB) RemoveAllCreatorsFromMediaSegment(ctx context.Context, mediaId string) error {
// 	query := dialect.Delete("media_creators").
// 		Where(goqu.I("media_creators.media_id").Eq(mediaId))
//
// 	_, err := db.db.Exec(ctx, query)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
