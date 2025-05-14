package apis

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/kr/pretty"
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

type AnimeProducer struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type AnimeTheme struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type AnimeGenre struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type Anime struct {
	Id string `json:"id"`

	Title        string  `json:"title"`
	TitleEnglish *string `json:"titleEnglish"`

	Description *string `json:"description"`

	Type         types.AnimeType   `json:"type"`
	Status       types.AnimeStatus `json:"status"`
	Rating       types.AnimeRating `json:"rating"`
	AiringSeason string            `json:"airingSeason"`
	EpisodeCount int64             `json:"episodeCount"`

	Score *float64 `json:"score"`

	StartDate *string `json:"startDate"`
	EndDate   *string `json:"endDate"`

	Studios   []AnimeStudio   `json:"studios"`
	Producers []AnimeProducer `json:"producers"`
	Themes    []AnimeTheme    `json:"themes"`
	Genres    []AnimeGenre    `json:"genres"`

	CoverUrl string `json:"coverUrl"`
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

func ConvertDBAnime(c pyrin.Context, anime database.Anime) Anime {
	coverUrl := ""
	if anime.CoverFilename.Valid {
		coverUrl = ConvertURL(c, fmt.Sprintf("/files/animes/%s/%s", anime.Id, anime.CoverFilename.String))
	}

	studios := make([]AnimeStudio, len(anime.Studios.Val))
	for i, studio := range anime.Studios.Val {
		studios[i] = AnimeStudio{
			Slug: studio.Slug,
			Name: studio.Name,
		}
	}

	producers := make([]AnimeProducer, len(anime.Producers.Val))
	for i, producer := range anime.Producers.Val {
		producers[i] = AnimeProducer{
			Slug: producer.Slug,
			Name: producer.Name,
		}
	}

	themes := make([]AnimeTheme, len(anime.Themes.Val))
	for i, theme := range anime.Themes.Val {
		themes[i] = AnimeTheme{
			Slug: theme.Slug,
			Name: theme.Name,
		}
	}

	genres := make([]AnimeGenre, len(anime.Genres.Val))
	for i, genre := range anime.Genres.Val {
		genres[i] = AnimeGenre{
			Slug: genre.Slug,
			Name: genre.Name,
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
		AiringSeason: anime.AiringSeason,
		EpisodeCount: anime.EpisodeCount.Int64,
		Score:        utils.SqlNullToFloat64Ptr(anime.Score),
		StartDate:    utils.SqlNullToStringPtr(anime.StartDate),
		EndDate:      utils.SqlNullToStringPtr(anime.EndDate),
		Studios:      studios,
		Producers:    producers,
		Themes:       themes,
		Genres:       genres,
		CoverUrl:     coverUrl,
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

				animes, p, err := app.DB().GetPagedAnimes(ctx, opts)
				if err != nil {
					return nil, err
				}

				pretty.Println(animes)

				res := GetAnimes{
					Page:   p,
					Animes: make([]Anime, len(animes)),
				}

				for i, anime := range animes {
					res.Animes[i] = ConvertDBAnime(c, anime)
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

				anime, err := app.DB().GetAnimeById(c.Request().Context(), id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AnimeNotFound()
					}

					return nil, err
				}

				return GetAnimeById{
					Anime: ConvertDBAnime(c, anime),
				}, nil
			},
		},
	)
}
