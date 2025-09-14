package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/utils"
)

type ApiToken struct {
	Id     string `db:"id"`
	UserId string `db:"user_id"`

	Name string `db:"name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func ApiTokenQuery() *goqu.SelectDataset {
	query := dialect.From("api_tokens").
		Select(
			"api_tokens.id",
			"api_tokens.user_id",

			"api_tokens.name",

			"api_tokens.updated",
			"api_tokens.created",
		)

	return query
}

func (db *Database) GetApiTokenById(ctx context.Context, id string) (ApiToken, error) {
	query := ApiTokenQuery().
		Where(goqu.I("api_tokens.id").Eq(id))

	return ember.Single[ApiToken](db.db, ctx, query)
}

func (db *Database) GetAllApiTokensForUser(ctx context.Context, userId string) ([]ApiToken, error) {
	query := ApiTokenQuery().
		Where(goqu.I("api_tokens.user_id").Eq(userId))

	return ember.Multiple[ApiToken](db.db, ctx, query)
}

type CreateApiTokenParams struct {
	Id     string
	UserId string
	Name   string

	Created int64
	Updated int64
}

func (db *Database) CreateApiToken(ctx context.Context, params CreateApiTokenParams) (string, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateApiTokenId()
	}

	query := dialect.Insert("api_tokens").Rows(goqu.Record{
		"id":      id,
		"user_id": params.UserId,

		"name": params.Name,

		"created": created,
		"updated": updated,
	}).
		Returning("api_tokens.id")

	return ember.Single[string](db.db, ctx, query)
}

func (db *Database) DeleteApiToken(ctx context.Context, id string) error {
	query := dialect.Delete("api_tokens").
		Where(goqu.I("api_tokens.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
