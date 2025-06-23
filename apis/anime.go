package apis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/transform"
	"github.com/nanoteck137/validate"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type AnimeStudio struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type AnimeTag struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type AnimeUser struct {
	List         *types.AnimeUserList `json:"list"`
	Score        *int64               `json:"score"`
	Episode      *int64               `json:"episode"`
	RewatchCount *int64               `json:"rewatchCount"`
	IsRewatching bool                 `json:"isRewatching"`
}

type AnimeImage struct {
	Hash    string `json:"hash"`
	Url     string `json:"url"`
	IsCover bool   `json:"isCover"`
}

type Anime struct {
	Id string `json:"id"`

	Title        string  `json:"title"`
	TitleEnglish *string `json:"titleEnglish"`

	Description *string `json:"description"`

	Type         types.AnimeType   `json:"type"`
	Score        *float64          `json:"score"`
	Status       types.AnimeStatus `json:"status"`
	Rating       types.AnimeRating `json:"rating"`
	EpisodeCount *int64            `json:"episodeCount"`
	AiringSeason *AnimeTag         `json:"airingSeason"`

	StartDate *string `json:"startDate"`
	EndDate   *string `json:"endDate"`

	Studios []AnimeStudio `json:"studios"`
	Tags    []AnimeTag    `json:"tags"`

	CoverUrl string       `json:"coverUrl"`
	Images   []AnimeImage `json:"images"`

	User *AnimeUser `json:"user,omitempty"`
}

type GetAnimes struct {
	Page   types.Page `json:"page"`
	Animes []Anime    `json:"animes"`
}

type GetAnimeById struct {
	Anime
}

// TODO(patrik): Move
func getPageOptions(q url.Values) database.FetchOptions {
	perPage := 100
	page := 0

	if s := q.Get("perPage"); s != "" {
		i, _ := strconv.Atoi(s)
		if i > 0 {
			perPage = i
		}
	}

	if s := q.Get("page"); s != "" {
		i, _ := strconv.Atoi(s)
		page = i
	}

	return database.FetchOptions{
		PerPage: perPage,
		Page:    page,
	}
}

func ConvertDBAnime(c pyrin.Context, hasUser bool, anime database.Anime) Anime {
	// TODO(patrik): Add default cover for animes without covers
	coverUrl := ""
	// for _, image := range anime.Images.Data {
	// 	if image.IsCover > 0 {
	// 		coverUrl = ConvertURL(c, fmt.Sprintf("/files/animes/%s/%s", anime.Id, image.Filename))
	// 		break
	// 	}
	// }

	// images := []AnimeImage{}
	// for _, image := range anime.Images.Data {
	// 	url := ConvertURL(c, fmt.Sprintf("/files/animes/%s/%s", anime.Id, image.Filename))
	//
	// 	images = append(images, AnimeImage{
	// 		Hash:    image.Hash,
	// 		Url:     url,
	// 		IsCover: image.IsCover > 0,
	// 	})
	// }

	// studios := make([]AnimeStudio, len(anime.Studios.Data))
	// for i, studio := range anime.Studios.Data {
	// 	studios[i] = AnimeStudio{
	// 		Slug: studio.Slug,
	// 		Name: studio.Name,
	// 	}
	// }

	// tags := make([]AnimeTag, len(anime.Tags.Data))
	// for i, tag := range anime.Tags.Data {
	// 	tags[i] = AnimeTag{
	// 		Slug: tag.Slug,
	// 		Name: tag.Name,
	// 	}
	// }

	// var user *AnimeUser
	// if hasUser {
	// 	user = &AnimeUser{}
	//
	// 	if anime.UserData.Valid {
	// 		val := anime.UserData.Data
	// 		user.List = val.List
	// 		user.Episode = val.Episode
	// 		user.RewatchCount = val.RewatchCount
	// 		user.Score = val.Score
	// 		user.IsRewatching = val.IsRewatching > 0
	// 	}
	// }

	// var airingSeason *AnimeTag
	// if anime.AiringSeason.Valid {
	// 	airingSeason = &AnimeTag{
	// 		Slug: anime.AiringSeason.Data.Slug,
	// 		Name: anime.AiringSeason.Data.Name,
	// 	}
	// }

	return Anime{
		Id:           anime.Id,
		Title:        anime.Title,
		TitleEnglish: utils.SqlNullToStringPtr(anime.TitleEnglish),
		Description:  utils.SqlNullToStringPtr(anime.Description),
		// Type:         anime.Type,
		Status: anime.Status,
		Rating: anime.Rating,
		// AiringSeason: airingSeason,
		// EpisodeCount: utils.SqlNullToInt64Ptr(anime.EpisodeCount),
		Score:     utils.SqlNullToFloat64Ptr(anime.Score),
		StartDate: utils.SqlNullToStringPtr(anime.StartDate),
		EndDate:   utils.SqlNullToStringPtr(anime.EndDate),
		// Studios:      studios,
		// Tags:         tags,
		CoverUrl: coverUrl,
		// Images:       images,
		// User:         user,
	}
}

type SetAnimeUserData struct {
	List         *types.AnimeUserList `json:"list,omitempty"`
	Score        *int64               `json:"score,omitempty"`
	Episode      *int64               `json:"episode,omitempty"`
	RewatchCount *int64               `json:"rewatchCount,omitempty"`
	IsRewatching *bool                `json:"isRewatching,omitempty"`
}

func (b *SetAnimeUserData) Transform() {
	if b.Score != nil {
		*b.Score = utils.Clamp(*b.Score, 0, 10)
	}
}

type CreateAnime struct {
	Id string `json:"id"`
}

type CreateAnimeBody struct {
	Type string `json:"type"`

	TmdbId    string `json:"tmdbId"`
	MalId     string `json:"malId"`
	AnilistId string `json:"anilistId"`

	Title       string `json:"title"`
	Description string `json:"description"`

	Score  float64 `json:"score"`
	Status string  `json:"status"`
	Rating string  `json:"rating"`

	EpisodeCount int `json:"episodeCount"`
}

func (b *CreateAnimeBody) Transform() {
	b.TmdbId = transform.String(b.TmdbId)
	b.MalId = transform.String(b.MalId)
	b.AnilistId = transform.String(b.AnilistId)

	b.Title = transform.String(b.Title)
	b.Description = transform.String(b.Description)

	b.Score = utils.Clamp(b.Score, 0.0, 10.0)

	b.EpisodeCount = utils.Min(b.EpisodeCount, 0)
}

func (b CreateAnimeBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Type, validate.Required, validate.By(types.ValidateAnimeType)),

		validate.Field(&b.Status, validate.By(types.ValidateAnimeStatus)),
		validate.Field(&b.Rating, validate.By(types.ValidateAnimeRating)),
	)
}

type EditAnimeBody struct {
	Type *string `json:"type,omitempty"`

	TmdbId    *string `json:"tmdbId,omitempty"`
	MalId     *string `json:"malId,omitempty"`
	AnilistId *string `json:"anilistId,omitempty"`

	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`

	Score  *float64 `json:"score,omitempty"`
	Status *string  `json:"status,omitempty"`
	Rating *string  `json:"rating,omitempty"`

	AdminStatus *string `json:"adminStatus,omitempty"`
}

func (b *EditAnimeBody) Transform() {
	b.TmdbId = transform.StringPtr(b.TmdbId)
	b.MalId = transform.StringPtr(b.MalId)
	b.AnilistId = transform.StringPtr(b.AnilistId)

	b.Title = transform.StringPtr(b.Title)
	b.Description = transform.StringPtr(b.Description)

	if b.Score != nil {
		*b.Score = utils.Clamp(*b.Score, 0.0, 10.0)
	}
}

func (b EditAnimeBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Type, validate.Required.When(b.Type != nil), validate.By(types.ValidateAnimeType)),

		validate.Field(&b.Status, validate.Required.When(b.Status != nil), validate.By(types.ValidateAnimeStatus)),
		validate.Field(&b.Rating, validate.Required.When(b.Rating != nil), validate.By(types.ValidateAnimeRating)),

		validate.Field(&b.AdminStatus, validate.Required.When(b.AdminStatus != nil), validate.By(types.ValidateEntryAdminStatus)),
	)
}

type AnimeEpisode struct {
	Index   int64  `json:"index"`
	AnimeId string `json:"animeId"`

	Name string `json:"name"`
}

type GetAnimeEpisodes struct {
	Episodes []AnimeEpisode `json:"episodes"`
}

type AddEpisodesBody struct {
	Count int `json:"count"`
}

func (b *AddEpisodesBody) Transform() {
	b.Count = utils.Min(b.Count, 0)
}

type EditEpisodeBody struct {
	Name *string `json:"name"`
}

func (b *EditEpisodeBody) Transform() {
	b.Name = transform.StringPtr(b.Name)
}

func (b EditEpisodeBody) Validate() error {
	return validate.ValidateStruct(&b, 
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),
	)
}

func InstallAnimeHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetAnimes",
			Method:       http.MethodGet,
			Path:         "/animes",
			ResponseType: GetAnimes{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				opts := getPageOptions(q)

				ctx := context.TODO()

				var userId *string
				if user, err := User(app, c); err == nil {
					userId = &user.Id
				}

				filterStr := q.Get("filter")
				sortStr := q.Get("sort")
				animes, p, err := app.DB().GetPagedAnimes(ctx, userId, filterStr, sortStr, opts)
				if err != nil {
					return nil, err
				}

				res := GetAnimes{
					Page:   p,
					Animes: make([]Anime, len(animes)),
				}

				for i, anime := range animes {
					res.Animes[i] = ConvertDBAnime(c, userId != nil, anime)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetAnimeById",
			Method:       http.MethodGet,
			Path:         "/animes/:id",
			ResponseType: GetAnimeById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				var userId *string
				if user, err := User(app, c); err == nil {
					userId = &user.Id
				}

				anime, err := app.DB().GetAnimeById(c.Request().Context(), userId, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AnimeNotFound()
					}

					return nil, err
				}

				return GetAnimeById{
					Anime: ConvertDBAnime(c, userId != nil, anime),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreateAnime",
			Method:       http.MethodPost,
			Path:         "/animes",
			ResponseType: CreateAnime{},
			BodyType:     CreateAnimeBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// TODO(patrik): Add admin check

				body, err := pyrin.Body[CreateAnimeBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				ty := types.AnimeType(body.Type)

				id, err := app.DB().CreateAnime(ctx, database.CreateAnimeParams{
					Type: ty,
					TmdbId: sql.NullString{
						String: body.TmdbId,
						Valid:  body.TmdbId != "",
					},
					MalId: sql.NullString{
						String: body.MalId,
						Valid:  body.MalId != "",
					},
					AnilistId: sql.NullString{
						String: body.AnilistId,
						Valid:  body.AnilistId != "",
					},
					Title: body.Title,
					Description: sql.NullString{
						String: body.Description,
						Valid:  body.Description != "",
					},
					Score: sql.NullFloat64{
						Float64: body.Score,
						Valid:   body.Score != 0.0,
					},
					Status: types.AnimeStatus(body.Status),
					Rating: types.AnimeRating(body.Rating),
					// StartDate: sql.NullString{},
					// EndDate:   sql.NullString{},
				})
				if err != nil {
					return nil, err
				}

				if ty.IsMovie() {
					err := app.DB().CreateAnimeEpisode(ctx, database.CreateAnimeEpisodeParams{
						AnimeId: id,
						Name:    body.Title,
						Index:   1,
					})
					if err != nil {
						return nil, err
					}
				} else {
					for i := range body.EpisodeCount {
						err := app.DB().CreateAnimeEpisode(ctx, database.CreateAnimeEpisodeParams{
							AnimeId: id,
							Name:    fmt.Sprintf("Episode %d", i+1),
							Index:   int64(i + 1),
						})
						if err != nil {
							return nil, err
						}
					}
				}

				return CreateAnime{
					Id: id,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "EditAnime",
			Method:       http.MethodPatch,
			Path:         "/animes/:id",
			ResponseType: nil,
			BodyType:     EditAnimeBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				// TODO(patrik): Add admin check

				body, err := pyrin.Body[EditAnimeBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbAnime, err := app.DB().GetAnimeById(ctx, nil, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AnimeNotFound()
					}

					return nil, err
				}

				changes := database.AnimeChanges{}

				if body.Type != nil {
					t := types.AnimeType(*body.Type)

					changes.Type = database.Change[types.AnimeType]{
						Value:   t,
						Changed: t != dbAnime.Type,
					}
				}

				if body.TmdbId != nil {
					changes.TmdbId = database.Change[sql.NullString]{
						Value: sql.NullString{
							String: *body.TmdbId,
							Valid:  *body.TmdbId != "",
						},
						Changed: *body.TmdbId != dbAnime.TmdbId.String,
					}
				}

				if body.MalId != nil {
					changes.MalId = database.Change[sql.NullString]{
						Value: sql.NullString{
							String: *body.MalId,
							Valid:  *body.MalId != "",
						},
						Changed: *body.MalId != dbAnime.MalId.String,
					}
				}

				if body.AnilistId != nil {
					changes.AnilistId = database.Change[sql.NullString]{
						Value: sql.NullString{
							String: *body.AnilistId,
							Valid:  *body.AnilistId != "",
						},
						Changed: *body.AnilistId != dbAnime.AnilistId.String,
					}
				}

				if body.Title != nil {
					changes.Title = database.Change[string]{
						Value:   *body.Title,
						Changed: *body.Title != dbAnime.Title,
					}
				}

				if body.Description != nil {
					changes.Description = database.Change[sql.NullString]{
						Value: sql.NullString{
							String: *body.Description,
							Valid:  *body.Description != "",
						},
						Changed: *body.Description != dbAnime.Description.String,
					}
				}

				if body.Score != nil {
					changes.Score = database.Change[sql.NullFloat64]{
						Value: sql.NullFloat64{
							Float64: *body.Score,
							Valid:   *body.Score != 0.0,
						},
						Changed: *body.Score != dbAnime.Score.Float64,
					}
				}

				if body.Status != nil {
					s := types.AnimeStatus(*body.Status)
					changes.Status = database.Change[types.AnimeStatus]{
						Value:   s,
						Changed: s != dbAnime.Status,
					}
				}

				if body.Rating != nil {
					r := types.AnimeRating(*body.Rating)
					changes.Rating = database.Change[types.AnimeRating]{
						Value:   r,
						Changed: r != dbAnime.Rating,
					}
				}

				if body.AdminStatus != nil {
					s := types.EntryAdminStatus(*body.AdminStatus)
					changes.AdminStatus = database.Change[types.EntryAdminStatus]{
						Value:   s,
						Changed: s != dbAnime.AdminStatus,
					}
				}

				err = app.DB().UpdateAnime(ctx, dbAnime.Id, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetAnimeEpisodes",
			Method:       http.MethodGet,
			Path:         "/animes/:id/episodes",
			ResponseType: GetAnimeEpisodes{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.Background()

				dbAnime, err := app.DB().GetAnimeById(ctx, nil, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AnimeNotFound()
					}

					return nil, err
				}

				episodes, err := app.DB().GetAnimeEpisodesByAnimeId(ctx, dbAnime.Id)
				if err != nil {
					return nil, err
				}

				res := make([]AnimeEpisode, len(episodes))

				for i, episode := range episodes {
					res[i] = AnimeEpisode{
						Index:   episode.Index,
						AnimeId: episode.AnimeId,
						Name:    episode.Name,
					}
				}

				return GetAnimeEpisodes{
					Episodes: res,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "AddEpisodes",
			Method:       http.MethodPost,
			Path:         "/animes/:id/episodes",
			ResponseType: nil,
			BodyType:     AddEpisodesBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// TODO(patrik): Add admin check

				id := c.Param("id")

				body, err := pyrin.Body[AddEpisodesBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbAnime, err := app.DB().GetAnimeById(ctx, nil, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AnimeNotFound()
					}

					return nil, err
				}

				// TODO(patrik): A better implementation would be getting
				// the last episode from the database
				episodes, err := app.DB().GetAnimeEpisodesByAnimeId(ctx, dbAnime.Id)
				if err != nil {
					return nil, err
				}

				lastIndex := int64(0)
				if len(episodes) > 0 {
					episode := episodes[len(episodes)-1]
					lastIndex = episode.Index
				}

				for i := range body.Count {
					idx := lastIndex + int64(i) + 1

					err := app.DB().CreateAnimeEpisode(ctx, database.CreateAnimeEpisodeParams{
						Index:   idx,
						AnimeId: dbAnime.Id,
						Name:    fmt.Sprintf("Episode %d", idx),
					})
					if err != nil {
						return nil, err
					}
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "EditEpisode",
			Method:       http.MethodPatch,
			Path:         "/animes/:id/episodes/:index",
			ResponseType: nil,
			BodyType:     EditEpisodeBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// TODO(patrik): Add admin check

				id := c.Param("id")
				index, err := strconv.ParseInt(c.Param("index"), 10, 64)
				if err != nil {
					// TODO(patrik): Handle error better
					return nil, errors.New("failed to parse 'index' path param as integer")
				}

				body, err := pyrin.Body[EditEpisodeBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbAnime, err := app.DB().GetAnimeById(ctx, nil, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AnimeNotFound()
					}

					return nil, err
				}

				dbEpisode, err := app.DB().GetAnimeEpisodeByIndexAnimeId(ctx, index, dbAnime.Id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, EpisodeNotFound()
					}

					return nil, err
				}

				changes := database.AnimeEpisodeChanges{}

				if body.Name != nil {
					changes.Name = database.Change[string]{
						Value:   *body.Name,
						Changed: *body.Name != dbEpisode.Name,
					}
				}

				err = app.DB().UpdateAnimeEpisode(ctx, dbEpisode.Index, dbEpisode.AnimeId, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "SetAnimeUserData",
			Method:       http.MethodPost,
			Path:         "/animes/:id/user",
			ResponseType: nil,
			BodyType:     SetAnimeUserData{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.TODO()

				body, err := pyrin.Body[SetAnimeUserData](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				anime, err := app.DB().GetAnimeById(ctx, &user.Id, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AnimeNotFound()
					}

					return nil, err
				}

				val := anime.UserData.Data

				data := database.SetAnimeUserData{
					List:         utils.AnimeUserListPtrToSqlNull(val.List),
					Episode:      utils.Int64PtrToSqlNull(val.Episode),
					IsRewatching: val.IsRewatching > 0,
					Score:        utils.Int64PtrToSqlNull(val.Score),
				}

				if body.List != nil {
					data.List = sql.NullString{
						String: string(*body.List),
						Valid:  *body.List != "",
					}
				}

				if body.Episode != nil {
					data.Episode = sql.NullInt64{
						Int64: *body.Episode,
						Valid: *body.Episode != 0,
					}
				}

				if body.RewatchCount != nil {
					data.RewatchCount = sql.NullInt64{
						Int64: *body.RewatchCount,
						Valid: *body.RewatchCount != 0,
					}
				}

				if body.IsRewatching != nil {
					data.IsRewatching = *body.IsRewatching
				}

				if body.Score != nil {
					data.Score = sql.NullInt64{
						Int64: *body.Score,
						Valid: *body.Score != 0,
					}
				}

				err = app.DB().SetAnimeUserData(ctx, anime.Id, user.Id, data)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserAnimeList",
			Method:       http.MethodGet,
			Path:         "/animes/user/list/:id",
			ResponseType: GetAnimes{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				q := c.Request().URL.Query()
				opts := getPageOptions(q)

				ctx := context.TODO()

				user, err := app.DB().GetUserById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, UserNotFound()
					}

					return nil, err
				}

				filterStr := q.Get("filter")
				sortStr := q.Get("sort")
				animes, p, err := app.DB().GetPagedAnimes(ctx, &user.Id, filterStr, sortStr, opts)
				if err != nil {
					return nil, err
				}

				res := GetAnimes{
					Page:   p,
					Animes: make([]Anime, len(animes)),
				}

				for i, anime := range animes {
					res.Animes[i] = ConvertDBAnime(c, true, anime)
				}

				return res, nil
			},
		},
	)
}
