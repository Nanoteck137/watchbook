package apis

import (
	"net/http"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/provider/myanimelist"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type ProviderMyAnimeListAnime struct {
	MalId string `json:"malId"`

	MediaType types.MediaType `json:"mediaType"`

	Title        string `json:"title"`
	TitleEnglish string `json:"titleEnglish"`

	Description string `json:"description"`

	Score        *float64          `json:"score"`
	Status       types.MediaStatus `json:"status"`
	Rating       types.MediaRating `json:"rating"`
	AiringSeason string            `json:"airingSeason"`

	StartDate *string `json:"startDate"`
	EndDate   *string `json:"endDate"`

	Studios []string `json:"studios"`
	Tags    []string `json:"tags"`

	CoverImageUrl string `json:"coverImageUrl"`

	EpisodeCount *int64 `json:"episodeCount"`

	UsingCache bool `json:"usingCache"`
}

func InstallProviderHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "ProviderMyAnimeListGetAnime",
			Method:       http.MethodGet,
			Path:         "/provider/myanimelist/anime/:id",
			ResponseType: ProviderMyAnimeListAnime{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				url := c.Request().URL
				q := url.Query()

				invalidateCache := q.Get("invalidateCache") == "true"

				entry, err := myanimelist.GetAnime(app.WorkDir(), id, !invalidateCache)
				if err != nil {
					return nil, err
				}

				return ProviderMyAnimeListAnime{
					MalId:         id,
					MediaType:     entry.Type,
					Title:         entry.Title,
					TitleEnglish:  entry.TitleEnglish,
					Description:   entry.Description,
					Score:         entry.Score,
					Status:        entry.Status,
					Rating:        entry.Rating,
					AiringSeason:  entry.AiringSeason,
					StartDate:     entry.StartDate,
					EndDate:       entry.EndDate,
					Studios:       utils.FixNilArrayToEmpty(entry.Studios),
					Tags:          utils.FixNilArrayToEmpty(entry.Tags),
					CoverImageUrl: entry.CoverImageUrl,
					EpisodeCount:  entry.EpisodeCount,
					UsingCache:    entry.UsingCache,
				}, nil
			},
		},
	)
}
