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
	for _, image := range anime.Images.Data {
		if image.IsCover > 0 {
			coverUrl = ConvertURL(c, fmt.Sprintf("/files/animes/%s/%s", anime.Id, image.Filename))
			break
		}
	}

	images := []AnimeImage{}
	for _, image := range anime.Images.Data {
		url := ConvertURL(c, fmt.Sprintf("/files/animes/%s/%s", anime.Id, image.Filename))

		images = append(images, AnimeImage{
			Hash:    image.Hash,
			Url:     url,
			IsCover: image.IsCover > 0,
		})
	}

	studios := make([]AnimeStudio, len(anime.Studios.Data))
	for i, studio := range anime.Studios.Data {
		studios[i] = AnimeStudio{
			Slug: studio.Slug,
			Name: studio.Name,
		}
	}

	tags := make([]AnimeTag, len(anime.Tags.Data))
	for i, tag := range anime.Tags.Data {
		tags[i] = AnimeTag{
			Slug: tag.Slug,
			Name: tag.Name,
		}
	}

	var user *AnimeUser
	if hasUser {
		user = &AnimeUser{}

		if anime.UserData.Valid {
			val := anime.UserData.Data
			user.List = val.List
			user.Episode = val.Episode
			user.RewatchCount = val.RewatchCount
			user.Score = val.Score
			user.IsRewatching = val.IsRewatching > 0
		}
	}

	var airingSeason *AnimeTag
	if anime.AiringSeason.Valid {
		airingSeason = &AnimeTag{
			Slug: anime.AiringSeason.Data.Slug,
			Name: anime.AiringSeason.Data.Name,
		}
	}

	return Anime{
		Id:           anime.Id,
		Title:        anime.Title,
		TitleEnglish: utils.SqlNullToStringPtr(anime.TitleEnglish),
		Description:  utils.SqlNullToStringPtr(anime.Description),
		Type:         anime.Type,
		Status:       anime.Status,
		Rating:       anime.Rating,
		AiringSeason: airingSeason,
		EpisodeCount: utils.SqlNullToInt64Ptr(anime.EpisodeCount),
		Score:        utils.SqlNullToFloat64Ptr(anime.Score),
		StartDate:    utils.SqlNullToStringPtr(anime.StartDate),
		EndDate:      utils.SqlNullToStringPtr(anime.EndDate),
		Studios:      studios,
		Tags:         tags,
		CoverUrl:     coverUrl,
		Images:       images,
		User:         user,
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
