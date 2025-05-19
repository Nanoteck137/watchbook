package apis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/transform"
	"github.com/nanoteck137/validate"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/mal"
	"github.com/nanoteck137/watchbook/types"
)

type UpdateUserSettingsBody struct {
	DisplayName *string `json:"displayName,omitempty"`
}

func (b *UpdateUserSettingsBody) Transform() {
	b.DisplayName = transform.StringPtr(b.DisplayName)
}

func (b UpdateUserSettingsBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.DisplayName,
			validate.Required.When(b.DisplayName != nil),
		),
	)
}

type CreateApiToken struct {
	Token string `json:"token"`
}

type CreateApiTokenBody struct {
	Name string `json:"name"`
}

func (b *CreateApiTokenBody) Transform() {
	b.Name = transform.String(b.Name)
}

func (b CreateApiTokenBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
	)
}

type ApiToken struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetAllApiTokens struct {
	Tokens []ApiToken `json:"tokens"`
}

type ImportMalListBody struct {
	Username                string `json:"username"`
	OverrideExistingEntries bool   `json:"overrideExistingEntries,omitempty"`
}

func (b *ImportMalListBody) Transform() {
	b.Username = transform.String(b.Username)
}

func (b ImportMalListBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Username, validate.Required),
	)
}

func InstallUserHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:     "UpdateUserSettings",
			Method:   http.MethodPatch,
			Path:     "/user/settings",
			BodyType: UpdateUserSettingsBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[UpdateUserSettingsBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				settings := user.ToUserSettings()

				if body.DisplayName != nil {
					settings.DisplayName = sql.NullString{
						String: *body.DisplayName,
						Valid:  true,
					}
				}

				err = app.DB().UpdateUserSettings(context.TODO(), settings)
				if err != nil {
					// TODO(patrik): Handle error
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "ImportMalList",
			Method:       http.MethodPost,
			Path:         "/user/import/mal",
			ResponseType: nil,
			BodyType:     ImportMalListBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[ImportMalListBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				entries, err := mal.GetUserWatchlist(dl, body.Username)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				for _, entry := range entries {
					malId := strconv.Itoa(entry.AnimeId)

					_, err := app.DB().GetAnimeByMalId(ctx, &user.Id, malId)
					if err != nil && errors.Is(err, database.ErrItemNotFound) {
						_, err := app.DB().CreateAnime(ctx, database.CreateAnimeParams{
							MalId: sql.NullString{
								String: malId,
								Valid:  true,
							},
							Title:           string(entry.AnimeTitle),
							ShouldFetchData: true,
						})
						if err != nil {
							return nil, err
						}
					}
				}

				for _, entry := range entries {
					malId := strconv.Itoa(entry.AnimeId)

					anime, err := app.DB().GetAnimeByMalId(ctx, &user.Id, malId)
					if err != nil {
						return nil, err
					}

					if anime.Id == "c57icqr7" {
						fmt.Printf("body.OverrideExistingEntries: %v\n", body.OverrideExistingEntries)
						fmt.Printf("anime.UserData.Valid: %v\n", anime.UserData.Valid)
						fmt.Printf("(anime.UserData.Valid && !body.OverrideExistingEntries): %v\n", (anime.UserData.Valid && !body.OverrideExistingEntries))
					}

					if anime.UserData.Valid && !body.OverrideExistingEntries {
						continue
					}

					noList := false
					list := types.AnimeUserListWatching
					switch entry.Status {
					case mal.WatchlistStatusCurrentlyWatching:
						list = types.AnimeUserListWatching
					case mal.WatchlistStatusCompleted:
						list = types.AnimeUserListCompleted
					case mal.WatchlistStatusOnHold:
						list = types.AnimeUserListOnHold
					case mal.WatchlistStatusDropped:
						list = types.AnimeUserListDropped
					case mal.WatchlistStatusPlanToWatch:
						list = types.AnimeUserListPlanToWatch
					default:
						noList = true
					}

					err = app.DB().SetAnimeUserData(ctx, anime.Id, user.Id, database.SetAnimeUserData{
						List: sql.NullString{
							String: string(list),
							Valid:  !noList,
						},
						Episode: sql.NullInt64{
							Int64: int64(entry.NumWatchedEpisodes),
							Valid: entry.NumWatchedEpisodes != 0,
						},
						IsRewatching: entry.IsRewatching > 0,
						Score: sql.NullInt64{
							Int64: int64(entry.Score),
							Valid: entry.Score != 0,
						},
					})
					if err != nil {
						return nil, err
					}
				}

				return nil, nil
			},
		},

		// TODO(patrik): Add back api tokens
		// pyrin.ApiHandler{
		// 	Name:         "CreateApiToken",
		// 	Method:       http.MethodPost,
		// 	Path:         "/user/apitoken",
		// 	ResponseType: CreateApiToken{},
		// 	BodyType:     CreateApiTokenBody{},
		// 	HandlerFunc: func(c pyrin.Context) (any, error) {
		// 		user, err := User(app, c)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		body, err := pyrin.Body[CreateApiTokenBody](c)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		ctx := context.TODO()
		//
		// 		token, err := app.DB().CreateApiToken(ctx, database.CreateApiTokenParams{
		// 			UserId: user.Id,
		// 			Name:   body.Name,
		// 		})
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		return CreateApiToken{
		// 			Token: token.Id,
		// 		}, nil
		// 	},
		// },
		//
		// pyrin.ApiHandler{
		// 	Name:         "GetAllApiTokens",
		// 	Method:       http.MethodGet,
		// 	Path:         "/user/apitoken",
		// 	ResponseType: GetAllApiTokens{},
		// 	HandlerFunc: func(c pyrin.Context) (any, error) {
		// 		user, err := User(app, c)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		ctx := context.TODO()
		//
		// 		tokens, err := app.DB().GetAllApiTokensForUser(ctx, user.Id)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		res := GetAllApiTokens{
		// 			Tokens: make([]ApiToken, len(tokens)),
		// 		}
		//
		// 		for i, token := range tokens {
		// 			res.Tokens[i] = ApiToken{
		// 				Id:   token.Id,
		// 				Name: token.Name,
		// 			}
		// 		}
		//
		// 		return res, nil
		// 	},
		// },
		//
		// pyrin.ApiHandler{
		// 	Name:   "DeleteApiToken",
		// 	Method: http.MethodDelete,
		// 	Path:   "/user/apitoken/:id",
		// 	Errors: []pyrin.ErrorType{ErrTypeApiTokenNotFound},
		// 	HandlerFunc: func(c pyrin.Context) (any, error) {
		// 		tokenId := c.Param("id")
		//
		// 		user, err := User(app, c)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		ctx := context.TODO()
		//
		// 		token, err := app.DB().GetApiTokenById(ctx, tokenId)
		// 		if err != nil {
		// 			if errors.Is(err, database.ErrItemNotFound) {
		// 				return nil, ApiTokenNotFound()
		// 			}
		//
		// 			return nil, err
		// 		}
		//
		// 		if token.UserId != user.Id {
		// 			return nil, ApiTokenNotFound()
		// 		}
		//
		// 		err = app.DB().DeleteApiToken(ctx, tokenId)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		return nil, nil
		// 	},
		// },
	)
}
