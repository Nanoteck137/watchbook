package sonarr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/nanoteck137/watchbook/provider"
	"github.com/nanoteck137/watchbook/tools/cache"
)

type Image struct {
	CoverType string `json:"coverType"`
	URL       string `json:"url"`
	RemoteURL string `json:"remoteUrl"`
}

type SerieAlternateTitle struct {
	Title             string `json:"title"`
	SceneSeasonNumber int    `json:"sceneSeasonNumber,omitempty"`
	SeasonNumber      int    `json:"seasonNumber,omitempty"`
}

type SerieOriginalLanguage struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type SerieSeason struct {
	SeasonNumber int  `json:"seasonNumber"`
	Monitored    bool `json:"monitored"`
	Statistics   struct {
		EpisodeFileCount  int      `json:"episodeFileCount"`
		EpisodeCount      int      `json:"episodeCount"`
		TotalEpisodeCount int      `json:"totalEpisodeCount"`
		SizeOnDisk        int64    `json:"sizeOnDisk"`
		ReleaseGroups     []string `json:"releaseGroups"`
		PercentOfEpisodes float32  `json:"percentOfEpisodes"`
	} `json:"statistics"`
	Images []Image `json:"images"`
}

type SerieRatings struct {
	Votes int     `json:"votes"`
	Value float64 `json:"value"`
}

type SerieStatistics struct {
	SeasonCount       int      `json:"seasonCount"`
	EpisodeFileCount  int      `json:"episodeFileCount"`
	EpisodeCount      int      `json:"episodeCount"`
	TotalEpisodeCount int      `json:"totalEpisodeCount"`
	SizeOnDisk        int64    `json:"sizeOnDisk"`
	ReleaseGroups     []string `json:"releaseGroups"`
	PercentOfEpisodes float32  `json:"percentOfEpisodes"`
}

type Serie struct {
	Id                int                   `json:"id"`
	Title             string                `json:"title"`
	AlternateTitles   []SerieAlternateTitle `json:"alternateTitles"`
	SortTitle         string                `json:"sortTitle"`
	Status            string                `json:"status"`
	Ended             bool                  `json:"ended"`
	Overview          string                `json:"overview"`
	Network           string                `json:"network"`
	AirTime           string                `json:"airTime"`
	Images            []Image               `json:"images"`
	OriginalLanguage  SerieOriginalLanguage `json:"originalLanguage"`
	Seasons           []SerieSeason         `json:"seasons"`
	Year              int                   `json:"year"`
	Path              string                `json:"path"`
	QualityProfileID  int                   `json:"qualityProfileId"`
	SeasonFolder      bool                  `json:"seasonFolder"`
	Monitored         bool                  `json:"monitored"`
	MonitorNewItems   string                `json:"monitorNewItems"`
	UseSceneNumbering bool                  `json:"useSceneNumbering"`
	Runtime           int                   `json:"runtime"`
	TvdbId            int                   `json:"tvdbId"`
	TvRageId          int                   `json:"tvRageId"`
	TvMazeId          int                   `json:"tvMazeId"`
	TmdbId            int                   `json:"tmdbId"`
	FirstAired        time.Time             `json:"firstAired"`
	LastAired         time.Time             `json:"lastAired"`
	SeriesType        string                `json:"seriesType"`
	CleanTitle        string                `json:"cleanTitle"`
	ImdbID            string                `json:"imdbId"`
	TitleSlug         string                `json:"titleSlug"`
	RootFolderPath    string                `json:"rootFolderPath"`
	Certification     string                `json:"certification"`
	Genres            []string              `json:"genres"`
	Tags              []int                 `json:"tags"`
	Added             time.Time             `json:"added"`
	Ratings           SerieRatings          `json:"ratings"`
	Statistics        SerieStatistics       `json:"statistics"`
	LanguageProfileId int                   `json:"languageProfileId"`
}

type Episode struct {
	Id                         int       `json:"id"`
	SeriesId                   int       `json:"seriesId"`
	TvdbId                     int       `json:"tvdbId"`
	EpisodeFileId              int       `json:"episodeFileId"`
	SeasonNumber               int       `json:"seasonNumber"`
	EpisodeNumber              int       `json:"episodeNumber"`
	Title                      string    `json:"title"`
	AirDate                    string    `json:"airDate"`
	AirDateUtc                 time.Time `json:"airDateUtc"`
	LastSearchTime             time.Time `json:"lastSearchTime"`
	Runtime                    int       `json:"runtime"`
	Overview                   string    `json:"overview"`
	HasFile                    bool      `json:"hasFile"`
	Monitored                  bool      `json:"monitored"`
	AbsoluteEpisodeNumber      int       `json:"absoluteEpisodeNumber"`
	SceneAbsoluteEpisodeNumber int       `json:"sceneAbsoluteEpisodeNumber"`
	SceneEpisodeNumber         int       `json:"sceneEpisodeNumber"`
	SceneSeasonNumber          int       `json:"sceneSeasonNumber"`
	UnverifiedSceneNumbering   bool      `json:"unverifiedSceneNumbering"`
	Images                     []Image   `json:"images"`
}

const requestTTL = 1 * time.Hour

type ApiClient struct {
	client *provider.HTTPClient
	cache  cache.Cache
	apiUrl string
	apiKey string
}

func NewApiClient(cache cache.Cache, apiUrl, apiKey string) *ApiClient {
	return &ApiClient{
		client: provider.NewHttpClient(apiUrl),
		cache:  cache,
		apiUrl: apiUrl,
		apiKey: apiKey,
	}
}

type requestData struct {
	client *ApiClient

	cacheKey string

	path  string
	query url.Values
}

func apiRequest[T any](ctx context.Context, req requestData) (T, error) {
	var res T

	if data, ok := cache.GetJson[T](req.client.cache, req.cacheKey); ok {
		fmt.Println("Using the cache", req.cacheKey)
		return data, nil
	}

	d, err := req.client.client.Get(ctx, req.path, provider.RequestOptions{
		Headers: http.Header{
			"accept":    {"application/json"},
			"X-Api-Key": {req.client.apiKey},
		},
		Query: req.query,
	})
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(d, &res)
	if err != nil {
		return res, err
	}

	// TODO(patrik): Log error?
	fmt.Println("Setting the cache", req.cacheKey)
	cache.SetJson(req.client.cache, req.cacheKey, res, requestTTL)

	return res, nil
}

func (c *ApiClient) GetSeries(ctx context.Context) ([]Serie, error) {
	return apiRequest[[]Serie](ctx, requestData{
		client:   c,
		cacheKey: "api:get-series",
		path:     "/api/v3/series",
		query:    url.Values{},
	})
}

func (c *ApiClient) GetSerieById(ctx context.Context, id string) (Serie, error) {
	return apiRequest[Serie](ctx, requestData{
		client:   c,
		cacheKey: "api:get-serie-by-id:" + id,
		path:     fmt.Sprintf("/api/v3/series/%s", id),
		query: url.Values{
			"includeSeasonImages": {"true"},
		},
	})
}

func (c *ApiClient) GetSeasonEpisodes(ctx context.Context, seriesId string, seasonNumber int) ([]Episode, error) {
	return apiRequest[[]Episode](ctx, requestData{
		client:   c,
		cacheKey: "api:get-season-episodes:" + seriesId + ":" + strconv.Itoa(seasonNumber),
		path:     "/api/v3/episode",
		query: url.Values{
			"seriesId":      {seriesId},
			"seasonNumber":  {strconv.Itoa(seasonNumber)},
			"includeImages": {"true"},
		},
	})
}
