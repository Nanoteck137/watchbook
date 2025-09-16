package tmdb

import (
	"strconv"
	"time"

	"github.com/nanoteck137/watchbook/provider"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

var _ provider.Provider = (*TmdbMovieProvider)(nil)

const MovieProviderName = "tmdb-movie"

type TmdbMovieProvider struct {
}

func (t *TmdbMovieProvider) Info() provider.Info {
	return provider.Info{
		Name:                    MovieProviderName,
		DisplayName:             "TheMovieDB Movie",
		SupportGetMedia:         true,
		SupportSearchMedia:      true,
		SupportGetCollection:    false,
		SupportSearchCollection: false,
	}
}

func (t *TmdbMovieProvider) GetCollection(c provider.Context, id string) (provider.Collection, error) {
	panic("unsupported")
}

func (t *TmdbMovieProvider) GetMedia(c provider.Context, id string) (provider.Media, error) {
	apiClient := NewApiClient(c.Cache())

	details, err := apiClient.GetMovieDetails(c.Context(), id)
	if err != nil {
		return provider.Media{}, err
	}

	images, err := apiClient.GetMovieImages(c.Context(), id)
	if err != nil {
		return provider.Media{}, err
	}

	status := types.MediaStatusUpcoming
	switch details.Status {
	case "Released":
		status = types.MediaStatusCompleted
	}

	creators := make([]string, len(details.ProductionCompanies))

	for i, company := range details.ProductionCompanies {
		creators[i] = company.Name
	}

	tags := make([]string, 0, len(details.Genres))

	for _, genre := range details.Genres {
		tags = append(tags, genre.Name)
	}

	airingSeason := types.GetAiringSeason(details.ReleaseDate)

	coverUrl := "http://image.tmdb.org/t/p/original" + details.PosterPath
	bannerUrl := "http://image.tmdb.org/t/p/original" + details.BackdropPath
	var logoUrl *string

	if len(images.Logos) > 0 {
		logo := images.Logos[0]
		u := "http://image.tmdb.org/t/p/original" + logo.FilePath
		logoUrl = &u
	}

	var description *string
	if details.Overview != "" {
		description = &details.Overview
	}

	score := utils.RoundFloat(details.VoteAverage, 2)

	var releaseDate *time.Time
	d, err := time.Parse(types.MediaDateLayout, details.ReleaseDate)
	if err == nil {
		releaseDate = &d
	}

	return provider.Media{
		ProviderId:       id,
		Type:             types.MediaTypeMovie,
		Title:            details.Title,
		Description:      description,
		Score:            &score,
		Status:           status,
		Rating:           "",
		AiringSeason:     &airingSeason,
		StartDate:        releaseDate,
		EndDate:          releaseDate,
		CoverUrl:         &coverUrl,
		LogoUrl:          logoUrl,
		BannerUrl:        &bannerUrl,
		Creators:         creators,
		Tags:             tags,
		Parts:            []provider.MediaPart{},
		ExtraProviderIds: map[string]string{},
	}, nil
}

func (t *TmdbMovieProvider) SearchCollection(c provider.Context, query string) ([]provider.SearchResult, error) {
	panic("unsupported")
}

func (t *TmdbMovieProvider) SearchMedia(c provider.Context, query string) ([]provider.SearchResult, error) {
	apiClient := NewApiClient(c.Cache())

	search, err := apiClient.MovieSearch(c.Context(), query)
	if err != nil {
		return nil, err
	}

	res := make([]provider.SearchResult, len(search.Results))

	for i, result := range search.Results {
		res[i] = provider.SearchResult{
			SearchType: provider.SearchResultTypeMedia,
			ProviderId: strconv.Itoa(result.Id),
			Title:      result.Title,
			MediaType:  types.MediaTypeMovie,
			ImageUrl:   "http://image.tmdb.org/t/p/original" + result.PosterPath,
		}
	}

	return res, nil
}
