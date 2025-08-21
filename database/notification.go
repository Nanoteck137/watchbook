package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/database/adapter"
	"github.com/nanoteck137/watchbook/filter"
	"github.com/nanoteck137/watchbook/kvstore"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type Notification struct {
	RowId int `db:"rowid"`

	Id     string `db:"id"`
	UserId string `db:"user_id"`

	Type     types.NotificationType `db:"type"`
	Title    string                 `db:"title"`
	Message  string                 `db:"message"`
	Metadata kvstore.Store          `db:"metadata"`
	IsRead   int                    `db:"is_read"`

	DedupKey string `db:"dedup_key"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Use goqu.T more
func NotificationQuery() *goqu.SelectDataset {
	query := dialect.From("notifications").
		Select(
			"notifications.rowid",

			"notifications.id",
			"notifications.user_id",

			"notifications.type",
			"notifications.title",
			"notifications.message",
			"notifications.metadata",
			"notifications.is_read",

			"notifications.dedup_key",

			"notifications.created",
			"notifications.updated",
		)

	return query
}

func (db *Database) GetPagedNotifications(ctx context.Context, userId string, filterStr, sortStr string, opts FetchOptions) ([]Notification, types.Page, error) {
	query := NotificationQuery().
		Where(goqu.I("notifications.user_id").Eq(userId))

	var err error

	a := adapter.NotificationResolverAdapter{}
	resolver := filter.New(&a)

	query, err = applyFilter(query, resolver, filterStr)
	if err != nil {
		return nil, types.Page{}, err
	}

	query, err = applySort(query, resolver, sortStr)
	if err != nil {
		return nil, types.Page{}, err
	}

	countQuery := query.
		Select(goqu.COUNT("notifications.id"))

	if opts.PerPage > 0 {
		query = query.
			Limit(uint(opts.PerPage)).
			Offset(uint(opts.Page * opts.PerPage))
	}

	totalItems, err := ember.Single[int](db.db, ctx, countQuery)
	if err != nil {
		return nil, types.Page{}, err
	}

	totalPages := utils.TotalPages(opts.PerPage, totalItems)
	page := types.Page{
		Page:       opts.Page,
		PerPage:    opts.PerPage,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	items, err := ember.Multiple[Notification](db.db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db *Database) GetAllNotifications(ctx context.Context) ([]Notification, error) {
	query := NotificationQuery()
	return ember.Multiple[Notification](db.db, ctx, query)
}

func (db *Database) GetNotificationById(ctx context.Context, id string) (Notification, error) {
	query := NotificationQuery().
		Where(goqu.I("notifications.id").Eq(id))

	return ember.Single[Notification](db.db, ctx, query)
}

type CreateNotificationParams struct {
	Id     string
	UserId string

	Type     types.NotificationType
	Title    string
	Message  string
	Metadata kvstore.Store
	IsRead   int

	DedupKey string

	Created int64
	Updated int64
}

func (db *Database) CreateNotification(ctx context.Context, params CreateNotificationParams) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = utils.CreateNotificationId()
	}

	if !types.IsValidNotificationType(params.Type) {
		params.Type = types.NotificationTypeUnknown
	}

	query := dialect.Insert("notifications").Rows(goqu.Record{
		"id":      params.Id,
		"user_id": params.UserId,

		"type":     params.Type,
		"title":    params.Title,
		"message":  params.Message,
		"metadata": params.Metadata,
		"is_read":  params.IsRead,

		"dedup_key": params.DedupKey,

		"created": params.Created,
		"updated": params.Updated,
	}).
		Returning("id")

	return ember.Single[string](db.db, ctx, query)
}

type NotificationChanges struct {
	Type Change[types.NotificationType]

	Title    Change[string]
	Message  Change[string]
	Metadata Change[kvstore.Store]
	IsRead   Change[int]

	DedupKey Change[string]

	Created Change[int64]
}

func (db *Database) UpdateNotification(ctx context.Context, id string, changes NotificationChanges) error {
	record := goqu.Record{}

	addToRecord(record, "type", changes.Type)

	addToRecord(record, "title", changes.Title)
	addToRecord(record, "message", changes.Message)
	addToRecord(record, "metadata", changes.Metadata)
	addToRecord(record, "is_read", changes.IsRead)

	addToRecord(record, "dedup_key", changes.DedupKey)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("notifications").
		Set(record).
		Where(goqu.I("notifications.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveNotification(ctx context.Context, id string) error {
	query := dialect.Delete("notifications").
		Where(goqu.I("notifications.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
