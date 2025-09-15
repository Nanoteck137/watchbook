package tmdb

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kr/pretty"
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

func (t *TmdbTvProvider) GetCollection(ctx context.Context, id string) (provider.Collection, error) {
	details, err := getTvDetails(id)
	if err != nil {
		return provider.Collection{}, err
	}

	images, err := getTvImages(id)
	if err != nil {
		return provider.Collection{}, err
	}

	pretty.Println(details)

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
			Id:   strconv.Itoa(details.Id) + "@" + strconv.Itoa(season.SeasonNumber),
			Name: season.Name,
		}
	}

	return res, nil
}

func (t *TmdbTvProvider) GetMedia(ctx context.Context, id string) (provider.Media, error) {
	splits := strings.Split(id, "@")
	if len(splits) != 2 {
		return provider.Media{}, errors.New("not found")
	}

	serieId := splits[0]
	seasonNumber := splits[1]

	details, err := getTvDetails(serieId)
	if err != nil {
		return provider.Media{}, err
	}

	_ = details

	seasonDetails, err := getTvSeasonDetails(serieId, seasonNumber)
	if err != nil {
		return provider.Media{}, err
	}

	pretty.Println(seasonDetails)

	fmt.Printf("serieId: %v\n", serieId)
	fmt.Printf("seasonNumber: %v\n", seasonNumber)

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
		res.Parts[i] = provider.MediaPart{
			Name:   episode.Name,
			Number: episode.EpisodeNumber,
		}
	}

	return res, nil
}

func (t *TmdbTvProvider) SearchCollection(ctx context.Context, query string) ([]provider.SearchResult, error) {
	search, err := tvSearch(query)
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

func (t *TmdbTvProvider) SearchMedia(ctx context.Context, query string) ([]provider.SearchResult, error) {
	panic("unsupported")
}
