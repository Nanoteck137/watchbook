package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/types"
)

type FolderItem struct {
	RowId int `db:"rowid"`

	FolderId string `db:"folder_id"`
	MediaId  string `db:"media_id"`

	Position int `db:"position"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Update this to match Media
type FullFolderItem struct {
	RowId int `db:"rowid"`

	FolderId string `db:"folder_id"`
	MediaId  string `db:"media_id"`

	Position int `db:"position"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	MediaType types.MediaType `db:"media_type"`

	MediaTitle       string         `db:"media_title"`
	MediaDescription sql.NullString `db:"media_description"`

	MediaScore        sql.NullFloat64   `db:"media_score"`
	MediaStatus       types.MediaStatus `db:"media_status"`
	MediaRating       types.MediaRating `db:"media_rating"`
	MediaAiringSeason sql.NullString    `db:"media_airing_season"`

	MediaStartDate sql.NullString `db:"media_start_date"`
	MediaEndDate   sql.NullString `db:"media_end_date"`

	MediaCoverFile  sql.NullString `db:"media_cover_file"`
	MediaLogoFile   sql.NullString `db:"media_logo_file"`
	MediaBannerFile sql.NullString `db:"media_banner_file"`

	MediaProviders ember.KVStore `db:"media_providers"`

	MediaCreated int64 `db:"media_created"`
	MediaUpdated int64 `db:"media_updated"`

	MediaPartCount sql.NullInt64 `db:"media_part_count"`

	MediaCreators ember.JsonColumn[[]string] `db:"media_creators"`
	MediaTags     ember.JsonColumn[[]string] `db:"media_tags"`

	MediaUserData ember.JsonColumn[MediaUserData] `db:"media_user_data"`
}

// TODO(patrik): Use goqu.T more
func FolderItemQuery(userId *string) *goqu.SelectDataset {
	query := dialect.From("folder_items").
		Select(
			"folder_items.rowid",

			"folder_items.folder_id",
			"folder_items.media_id",

			"folder_items.position",

			"folder_items.created",
			"folder_items.updated",
		)

	return query
}

// TODO(patrik): Use goqu.T more
func FullFolderItemQuery(userId *string) *goqu.SelectDataset {
	mediaQuery := MediaQuery(userId)

	query := dialect.From("folder_items").
		Select(
			"folder_items.rowid",

			"folder_items.folder_id",
			"folder_items.media_id",

			"folder_items.position",

			goqu.I("media.type").As("media_type"),

			goqu.I("media.title").As("media_title"),
			goqu.I("media.description").As("media_description"),

			goqu.I("media.score").As("media_score"),
			goqu.I("media.status").As("media_status"),
			goqu.I("media.rating").As("media_rating"),
			goqu.I("media.airing_season").As("media_airing_season"),

			goqu.I("media.start_date").As("media_start_date"),
			goqu.I("media.end_date").As("media_end_date"),

			goqu.I("media.cover_file").As("media_cover_file"),
			goqu.I("media.logo_file").As("media_logo_file"),
			goqu.I("media.banner_file").As("media_banner_file"),

			goqu.I("media.providers").As("media_providers"),

			goqu.I("media.created").As("media_created"),
			goqu.I("media.updated").As("media_updated"),

			goqu.I("media.part_count").As("media_part_count"),

			goqu.I("media.creators").As("media_creators"),
			goqu.I("media.tags").As("media_tags"),

			goqu.I("media.user_data").As("media_user_data"),
		).
		Join(
			mediaQuery.As("media"),
			goqu.On(goqu.I("folder_items.media_id").Eq(goqu.I("media.id"))),
		)

	return query
}

func (db DB) GetAllFolderItems(ctx context.Context) ([]FolderItem, error) {
	query := FolderItemQuery(nil)
	return ember.Multiple[FolderItem](db.db, ctx, query)
}

func (db DB) GetFullAllFolderItemsByFolder(ctx context.Context, userId *string, folderId string) ([]FullFolderItem, error) {
	query := FullFolderItemQuery(userId).
		Where(goqu.I("folder_items.folder_id").Eq(folderId))
	return ember.Multiple[FullFolderItem](db.db, ctx, query)
}

func (db DB) GetFolderItemById(ctx context.Context, folderId, mediaId string) (FolderItem, error) {
	query := FolderItemQuery(nil).
		Where(
			goqu.I("folder_items.folder_id").Eq(folderId),
			goqu.I("folder_items.media_id").Eq(mediaId),
		)

	return ember.Single[FolderItem](db.db, ctx, query)
}

func (db DB) GetLastFolderItemPosition(ctx context.Context, folderId string) (int, error) {
	query := FolderItemQuery(nil).
		Select(goqu.I("folder_items.position")).
		Where(
			goqu.I("folder_items.folder_id").Eq(folderId),
		).
		Order(goqu.I("folder_items.position").Desc())

	return ember.Single[int](db.db, ctx, query)
}

type CreateFolderItemParams struct {
	FolderId string
	MediaId  string

	Position int

	Created int64
	Updated int64
}

func (db DB) CreateFolderItem(ctx context.Context, params CreateFolderItemParams) error {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	query := dialect.Insert("folder_items").Rows(goqu.Record{
		"folder_id": params.FolderId,
		"media_id":  params.MediaId,

		"position": params.Position,

		"created": params.Created,
		"updated": params.Updated,
	})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type FolderItemChanges struct {
	Position Change[int]

	Created Change[int64]
}

func (db DB) UpdateFolderItem(ctx context.Context, folderId, mediaId string, changes FolderItemChanges) error {
	record := goqu.Record{}

	addToRecord(record, "position", changes.Position)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("folder_items").
		Set(record).
		Where(
			goqu.I("folder_items.folder_id").Eq(folderId),
			goqu.I("folder_items.media_id").Eq(mediaId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) RemoveFolderItem(ctx context.Context, folderId, mediaId string) error {
	query := dialect.Delete("folder_items").
		Where(
			goqu.I("folder_items.folder_id").Eq(folderId),
			goqu.I("folder_items.media_id").Eq(mediaId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) RemoveAllFolderItems(ctx context.Context, folderId string) error {
	query := dialect.Delete("folder_items").
		Where(
			goqu.I("folder_items.folder_id").Eq(folderId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) ShiftFolderItemsDown(ctx context.Context, folderId string, newPos, oldPos int) error {
	query := dialect.Update("folder_items").
		Set(goqu.Record{
			"position": goqu.L("? + 1", goqu.I("position")),
		}).
		Where(
			goqu.I("folder_items.folder_id").Eq(folderId),
			goqu.I("folder_items.position").Gte(newPos),
			goqu.I("folder_items.position").Lt(oldPos),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) ShiftFolderItemsUp(ctx context.Context, folderId string, newPos, oldPos int) error {
	query := dialect.Update("folder_items").
		Set(goqu.Record{
			"position": goqu.L("? - 1", goqu.I("position")),
		}).
		Where(
			goqu.I("folder_items.folder_id").Eq(folderId),
			goqu.I("folder_items.position").Gt(oldPos),
			goqu.I("folder_items.position").Lte(newPos),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) MoveFolderItem(ctx context.Context, folderId string, mediaId string, newPos int) error {
	item, err := db.GetFolderItemById(ctx, folderId, mediaId)
	if err != nil {
		return err
	}

	if item.Position == newPos {
		return nil
	}

	if newPos < item.Position {
		err := db.ShiftFolderItemsDown(ctx, folderId, newPos, item.Position)
		if err != nil {
			return err
		}
	} else {
		err := db.ShiftFolderItemsUp(ctx, folderId, newPos, item.Position)
		if err != nil {
			return err
		}
	}

	err = db.UpdateFolderItem(ctx, folderId, mediaId, FolderItemChanges{
		Position: Change[int]{
			Value:   newPos,
			Changed: newPos != item.Position,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (db DB) RepackFolderItems(ctx context.Context, folderId string) error {
	ordered := goqu.From("folder_items").
		Select(
			goqu.I("media_id"), 
			goqu.ROW_NUMBER().Over(goqu.W().OrderBy(goqu.I("position"))).As("new_pos"),
		).
		Where(goqu.I("folder_items.folder_id").Eq(folderId)).
		As("ordered")

	newPosQuery := goqu.Select("new_pos").
		From(ordered).
		Where(
			goqu.I("ordered.media_id").Eq(goqu.I("folder_items.media_id")),
		)

	query := goqu.Update("folder_items").
		Set(goqu.Record{
			"position": newPosQuery,
		}).
		Where(goqu.I("folder_items.folder_id").Eq(folderId))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
