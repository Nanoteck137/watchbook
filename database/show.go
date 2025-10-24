package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/database/adapter"
	"github.com/nanoteck137/watchbook/filter"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type Show struct {
	RowId int `db:"rowid"`

	Id   string         `db:"id"`
	Type types.ShowType `db:"type"`

	Name       string `db:"name"`
	SearchSlug string `db:"search_slug"`

	CoverFile  sql.NullString `db:"cover_file"`
	LogoFile   sql.NullString `db:"logo_file"`
	BannerFile sql.NullString `db:"banner_file"`

	DefaultProvider sql.NullString `db:"default_provider"`
	Providers       ember.KVStore  `db:"providers"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Use goqu.T more
// TODO(patrik): Add season count
func ShowQuery() *goqu.SelectDataset {
	query := dialect.From("shows").
		Select(
			"shows.rowid",

			"shows.id",
			"shows.type",

			"shows.name",
			"shows.search_slug",

			"shows.cover_file",
			"shows.logo_file",
			"shows.banner_file",

			"shows.default_provider",
			"shows.providers",

			"shows.created",
			"shows.updated",
		)

	return query
}

func (db DB) GetPagedShows(ctx context.Context, filterStr, sortStr string, opts FetchOptions) ([]Show, types.Page, error) {
	query := ShowQuery()

	var err error

	a := adapter.ShowResolverAdapter{}
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
		Select(goqu.COUNT("shows.id"))

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

	items, err := ember.Multiple[Show](db.db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db DB) GetAllShows(ctx context.Context) ([]Show, error) {
	query := ShowQuery()
	return ember.Multiple[Show](db.db, ctx, query)
}

func (db DB) GetShowById(ctx context.Context, id string) (Show, error) {
	query := ShowQuery().
		Where(goqu.I("shows.id").Eq(id))

	return ember.Single[Show](db.db, ctx, query)
}

func (db DB) GetShowByProviderId(ctx context.Context, providerName, value string) (Show, error) {
	query := ShowQuery().
		Where(
			goqu.Func("json_extract", goqu.I("shows.providers"), "$."+providerName).Eq(value),
		)

	return ember.Single[Show](db.db, ctx, query)
}

type CreateShowParams struct {
	Id   string
	Type types.ShowType

	Name       string
	SearchSlug string

	CoverFile  sql.NullString
	LogoFile   sql.NullString
	BannerFile sql.NullString

	DefaultProvider sql.NullString
	Providers       ember.KVStore

	Created int64
	Updated int64
}

func (db DB) CreateShow(ctx context.Context, params CreateShowParams) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = utils.CreateShowId()
	}

	if params.Type == "" {
		params.Type = types.ShowTypeUnknown
	}

	if params.SearchSlug == "" {
		params.SearchSlug = utils.Slug(params.Name)
	}

	query := dialect.Insert("shows").Rows(goqu.Record{
		"id":   params.Id,
		"type": params.Type,

		"name":        params.Name,
		"search_slug": params.SearchSlug,

		"cover_file":  params.CoverFile,
		"logo_file":   params.LogoFile,
		"banner_file": params.BannerFile,

		"default_provider": params.DefaultProvider,
		"providers":        params.Providers,

		"created": params.Created,
		"updated": params.Updated,
	}).
		Returning("id")

	return ember.Single[string](db.db, ctx, query)
}

type ShowChanges struct {
	Type Change[types.ShowType]

	Name       Change[string]
	SearchSlug Change[string]

	CoverFile  Change[sql.NullString]
	LogoFile   Change[sql.NullString]
	BannerFile Change[sql.NullString]

	DefaultProvider Change[sql.NullString]
	Providers       Change[ember.KVStore]

	Created Change[int64]
}

func (db DB) UpdateShow(ctx context.Context, id string, changes ShowChanges) error {
	record := goqu.Record{}

	addToRecord(record, "type", changes.Type)

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "search_slug", changes.SearchSlug)

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

	query := dialect.Update("shows").
		Set(record).
		Where(goqu.I("shows.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) RemoveShow(ctx context.Context, id string) error {
	query := dialect.Delete("shows").
		Where(goqu.I("shows.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
