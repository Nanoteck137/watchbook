package tmdb

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nanoteck137/watchbook/provider"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

var _ provider.Provider = (*TmdbTvProvider)(nil)

const TvProviderName = "tmdb-tv"

type TmdbTvProvider struct {
}

func (t *TmdbTvProvider) Info() provider.Info {
	return provider.Info{
		Name:                    TvProviderName,
		DisplayName:             "TheMovieDB TV",
		SupportGetMedia:         true,
		SupportSearchMedia:      false,
		SupportGetCollection:    true,
		SupportSearchCollection: true,
	}
}

func (t *TmdbTvProvider) GetCollection(c provider.Context, id string) (provider.Collection, error) {
	apiClient := NewApiClient(c.Cache())

	details, err := apiClient.GetTvDetails(c.Context(), id)
	if err != nil {
		return provider.Collection{}, err
	}

	images, err := apiClient.GetTvImages(c.Context(), id)
	if err != nil {
		return provider.Collection{}, err
	}

	coverUrl := "http://image.tmdb.org/t/p/original" + details.PosterPath
	bannerUrl := "http://image.tmdb.org/t/p/original" + details.BackdropPath
	var logoUrl *string

	if len(images.Logos) > 0 {
		logo := images.Logos[0]
		u := "http://image.tmdb.org/t/p/original" + logo.FilePath
		logoUrl = &u
	}

	res := provider.Collection{
		ProviderId: strconv.Itoa(details.Id),
		Type:       types.CollectionTypeSeries,
		Name:       details.Name,
		CoverUrl:   &coverUrl,
		LogoUrl:    logoUrl,
		BannerUrl:  &bannerUrl,
		Items:      make([]provider.CollectionItem, len(details.Seasons)),
	}

	for i, season := range details.Seasons {
		res.Items[i] = provider.CollectionItem{
			Id:       strconv.Itoa(details.Id) + "@" + strconv.Itoa(season.SeasonNumber),
			Name:     season.Name,
			Position: season.SeasonNumber,
		}
	}

	return res, nil
}

func (t *TmdbTvProvider) GetMedia(c provider.Context, id string) (provider.Media, error) {
	apiClient := NewApiClient(c.Cache())

	splits := strings.Split(id, "@")
	if len(splits) != 2 {
		return provider.Media{}, errors.New("not found")
	}

	serieId := splits[0]
	seasonNumber := splits[1]

	details, err := apiClient.GetTvDetails(c.Context(), serieId)
	if err != nil {
		return provider.Media{}, err
	}

	seasonDetails, err := apiClient.GetSeasonDetails(c.Context(), serieId, seasonNumber)
	if err != nil {
		return provider.Media{}, err
	}

	var description *string
	if seasonDetails.Overview != "" {
		description = &seasonDetails.Overview
	}

	score := utils.RoundFloat(seasonDetails.VoteAverage, 2)

	creators := make([]string, len(details.ProductionCompanies))
	for i, company := range details.ProductionCompanies {
		creators[i] = company.Name
	}

	tags := make([]string, 0, len(details.Genres))
	for _, genre := range details.Genres {
		tags = append(tags, genre.Name)
	}

	var startDate *time.Time
	var endDate *time.Time

	if len(seasonDetails.Episodes) > 0 {
		first := seasonDetails.Episodes[0]
		last := seasonDetails.Episodes[len(seasonDetails.Episodes)-1]

		{
			d, err := time.Parse(types.MediaDateLayout, first.AirDate)
			if err == nil {
				startDate = &d
			}
		}

		{
			d, err := time.Parse(types.MediaDateLayout, last.AirDate)
			if err == nil {
				endDate = &d
			}
		}
	} else {
		d, err := time.Parse(types.MediaDateLayout, seasonDetails.AirDate)
		if err == nil {
			startDate = &d
			endDate = &d
		}
	}

	var airingSeason *string
	if startDate != nil {
		s := types.GetAiringSeason(startDate.Format(types.MediaDateLayout))
		airingSeason = &s
	}

	status := types.MediaStatusUnknown
	if endDate != nil {
		if endDate.Before(time.Now()) {
			status = types.MediaStatusCompleted
		} else {
			status = types.MediaStatusOngoing
		}

	}

	coverUrl := "http://image.tmdb.org/t/p/original" + seasonDetails.PosterPath

	res := provider.Media{
		ProviderId:       id,
		Type:             types.MediaTypeTV,
		Title:            fmt.Sprintf("%s (%s)", details.Name, seasonDetails.Name),
		Description:      description,
		Score:            &score,
		Status:           status,
		Rating:           types.MediaRatingUnknown,
		AiringSeason:     airingSeason,
		StartDate:        startDate,
		EndDate:          endDate,
		CoverUrl:         &coverUrl,
		Creators:         creators,
		Tags:             tags,
		Parts:            make([]provider.MediaPart, len(seasonDetails.Episodes)),
		ExtraProviderIds: map[string]string{},
	}

	for i, episode := range seasonDetails.Episodes {
		var releaseDate *time.Time

		d, err := time.Parse(types.MediaDateLayout, episode.AirDate)
		if err == nil {
			releaseDate = &d
		}

		res.Parts[i] = provider.MediaPart{
			Name:        episode.Name,
			Number:      episode.EpisodeNumber,
			ReleaseDate: releaseDate,
		}
	}

	return res, nil
}

func (t *TmdbTvProvider) SearchCollection(c provider.Context, query string) ([]provider.SearchResult, error) {
	apiClient := NewApiClient(c.Cache())

	search, err := apiClient.TvSearch(c.Context(), query)
	if err != nil {
		return nil, err
	}

	res := make([]provider.SearchResult, len(search.Results))

	for i, result := range search.Results {
		res[i] = provider.SearchResult{
			SearchType: provider.SearchResultTypeCollection,
			ProviderId: strconv.Itoa(result.Id),
			Title:      result.Name,
			MediaType:  types.MediaTypeTV,
			ImageUrl:   "http://image.tmdb.org/t/p/original" + result.PosterPath,
		}
	}

	return res, nil
}

func (t *TmdbTvProvider) SearchMedia(c provider.Context, query string) ([]provider.SearchResult, error) {
	panic("unsupported")
}
