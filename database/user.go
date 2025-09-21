package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/utils"
)

type UserSettings struct {
	Id          string         `db:"id"`
	DisplayName sql.NullString `db:"display_name"`
}

type User struct {
	Id       string `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Role     string `db:"role"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	// NOTE(patrik): This needs to match UserSettings
	DisplayName sql.NullString `db:"display_name"`
}

func (u User) ToUserSettings() UserSettings {
	return UserSettings{
		Id:          u.Id,
		DisplayName: u.DisplayName,
	}
}

func UserQuery() *goqu.SelectDataset {
	query := dialect.From("users").
		Select(
			"users.id",
			"users.username",
			// TODO(patrik): Maybe don't include the password?
			"users.password",
			"users.role",

			"users.created",
			"users.updated",

			"users_settings.display_name",
		).
		LeftJoin(
			goqu.I("users_settings"),
			goqu.On(goqu.I("users.id").Eq(goqu.I("users_settings.id"))),
		)

	return query
}

func UserSettingsQuery() *goqu.SelectDataset {
	query := dialect.From("users_settings").
		Select(
			"users_settings.id",
			"users_settings.display_name",
		)

	return query
}

func (db DB) GetAllUsers(ctx context.Context) ([]User, error) {
	query := UserQuery()

	return ember.Multiple[User](db.db, ctx, query)
}

type CreateUserParams struct {
	Id       string
	Username string
	Password string
	Role     string

	Created int64
	Updated int64
}

func (db DB) CreateUser(ctx context.Context, params CreateUserParams) (string, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateUserId()
	}

	query := dialect.
		Insert("users").
		Rows(goqu.Record{
			"id":       id,
			"username": params.Username,
			"password": params.Password,
			"role":     params.Role,

			"created": created,
			"updated": updated,
		}).
		Returning("users.id")

	return ember.Single[string](db.db, ctx, query)
}

func (db DB) GetUserById(ctx context.Context, id string) (User, error) {
	query := UserQuery().
		Where(goqu.I("users.id").Eq(id))

	return ember.Single[User](db.db, ctx, query)
}

func (db DB) GetUserSettingsById(ctx context.Context, id string) (UserSettings, error) {
	query := UserSettingsQuery().
		Where(goqu.I("users_settings.id").Eq(id))

	return ember.Single[UserSettings](db.db, ctx, query)
}

func (db DB) GetUserByUsername(ctx context.Context, username string) (User, error) {
	query := UserQuery().
		Where(goqu.I("users.username").Eq(username))

	return ember.Single[User](db.db, ctx, query)
}

type UserChanges struct {
	Username Change[string]
	Password Change[string]

	Created Change[int64]
}

func (db DB) UpdateUser(ctx context.Context, id string, changes UserChanges) error {
	record := goqu.Record{}

	addToRecord(record, "username", changes.Username)
	addToRecord(record, "password", changes.Password)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("users").
		Set(record).
		Where(goqu.I("users.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) UpdateUserSettings(ctx context.Context, settings UserSettings) error {
	query := dialect.Insert("users_settings").
		Rows(goqu.Record{
			"id":           settings.Id,
			"display_name": settings.DisplayName,
		}).
		OnConflict(goqu.DoUpdate("id", goqu.Record{
			"display_name": settings.DisplayName,
		}))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
