package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/nanoteck137/watchbook/provider"
	"github.com/nanoteck137/watchbook/tools/cache"
)

type MovieSearchResult struct {
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	Id               int     `json:"id"`
	Title            string  `json:"title"`
	OriginalTitle    string  `json:"original_title"`
	Overview         string  `json:"overview"`
	PosterPath       string  `json:"poster_path"`
	MediaType        string  `json:"media_type"`
	OriginalLanguage string  `json:"original_language"`
	GenreIds         []int   `json:"genre_ids"`
	Popularity       float32 `json:"popularity"`
	ReleaseDate      string  `json:"release_date"`
	Video            bool    `json:"video"`
	VoteAverage      float32 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
}

type TvSearchResult struct {
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	Id               int     `json:"id"`
	Name             string  `json:"name"`
	OriginalName     string  `json:"original_name"`
	Overview         string  `json:"overview"`
	PosterPath       string  `json:"poster_path"`
	MediaType        string  `json:"media_type"`
	OriginalLanguage string  `json:"original_language"`
	GenreIds         []int   `json:"genre_ids"`
	Popularity       float32 `json:"popularity"`
	ReleaseDate      string  `json:"release_date"`
	Video            bool    `json:"video"`
	VoteAverage      float32 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
}

type SearchRequest[T any] struct {
	Page         int `json:"page"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
	Results      []T `json:"results"`
}

type ProductionCompany struct {
	Id            int    `json:"id"`
	LogoPath      string `json:"logo_path"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}

type Genre struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type MovieDetails struct {
	Adult               bool                `json:"adult"`                 //: false,
	BackdropPath        string              `json:"backdrop_path"`         //: "/i7UCf0ysjbYYaqcSKUox9BJz4Kp.jpg",
	BelongsToCollection any                 `json:"belongs_to_collection"` //: null,
	Budget              int                 `json:"budget"`                //: 28000000,
	Genres              []Genre             `json:"genres"`                //: [
	Homepage            string              `json:"homepage"`              //: "http://www.thebigshortmovie.com",
	Id                  int                 `json:"id"`                    //: 318846,
	ImdbId              string              `json:"imdb_id"`               //: "tt1596363",
	OriginCountry       any                 `json:"origin_country"`        //: [
	OriginalLanguage    string              `json:"original_language"`     //: "en",
	OriginalTitle       string              `json:"original_title"`        //: "The Big Short",
	Overview            string              `json:"overview"`              //: "The men who made millions from a global economic meltdown.",
	Popularity          float32             `json:"popularity"`            //: 7.7052,
	PosterPath          string              `json:"poster_path"`           //: "/scVEaJEwP8zUix8vgmMoJJ9Nq0w.jpg",
	ProductionCompanies []ProductionCompany `json:"production_companies"`  //: [
	ProductionCountries any                 `json:"production_countries"`  //: [
	ReleaseDate         string              `json:"release_date"`          //: "2015-12-11",
	Revenue             int                 `json:"revenue"`               //: 133346506,
	Runtime             int                 `json:"runtime"`               //: 131,
	SpokenLanguages     any                 `json:"spoken_languages"`      //: [
	Status              string              `json:"status"`                //: "Released",
	Tagline             string              `json:"tagline"`               //: "This is a true story.",
	Title               string              `json:"title"`                 //: "The Big Short",
	Video               bool                `json:"video"`                 //: false,
	VoteAverage         float64             `json:"vote_average"`          //: 7.357,
	VoteCount           int                 `json:"vote_count"`            //: 9313
}

type SeasonDetailsEpisode struct {
	AirDate        string  `json:"air_date"`
	EpisodeNumber  int     `json:"episode_number"`
	EpisodeType    string  `json:"episode_type"`
	Id             int     `json:"id"`
	Name           string  `json:"name"`
	Overview       string  `json:"overview"`
	ProductionCode string  `json:"production_code"`
	Runtime        int     `json:"runtime"`
	SeasonNumber   int     `json:"season_number"`
	ShowId         int     `json:"show_id"`
	StillPath      string  `json:"still_path"`
	VoteAverage    float64 `json:"vote_average"`
	VoteCount      int     `json:"vote_count"`

	//	  "crew": [
	//	    {
	//	      "department": "Writing",
	//	      "job": "Writer",
	//	      "credit_id": "52542275760ee313280006ce",
	//	      "adult": false,
	//	      "gender": 2,
	//	      "id": 66633,
	//	      "known_for_department": "Writing",
	//	      "name": "Vince Gilligan",
	//	      "original_name": "Vince Gilligan",
	//	      "popularity": 0.8583,
	//	      "profile_path": "/z3E0DhBg1V1PZVEtS9vfFPzOWYB.jpg"
	//	    },
	//	  ],
	//	  "guest_stars": [
	//	    {
	//	      "character": "Steven Gomez",
	//	      "credit_id": "5271b489760ee35b3e0881a7",
	//	      "order": 8,
	//	      "adult": false,
	//	      "gender": 2,
	//	      "id": 61535,
	//	      "known_for_department": "Acting",
	//	      "name": "Steven Michael Quezada",
	//	      "original_name": "Steven Michael Quezada",
	//	      "popularity": 0.4711,
	//	      "profile_path": "/pVYrDkwI6GWvCNL2kJhpDJfBFyd.jpg"
	//	    },
	//	  ]
	//	},
}

type SeasonDetails struct {
	Id           string                 `json:"_id"`
	AirDate      string                 `json:"air_date"`
	Name         string                 `json:"name"`
	Episodes     []SeasonDetailsEpisode `json:"episodes"`
	Overview     string                 `json:"overview"`
	SerieId      int                    `json:"id"`
	PosterPath   string                 `json:"poster_path"`
	SeasonNumber int                    `json:"season_number"`
	VoteAverage  float64                `json:"vote_average"`
}

type ImageItem struct {
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	AspectRatio float64 `json:"aspect_ratio"`
	Iso6391     *string `json:"iso_639_1"`
	FilePath    string  `json:"file_path"`
	VoteAverge  float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}

type Images struct {
	Id        int         `json:"id"`
	Backdrops []ImageItem `json:"backdrops"`
	Logos     []ImageItem `json:"logos"`
	Posters   []ImageItem `json:"posters"`
}

type TvDetailsSeason struct {
	AirDate      string  `json:"air_date"`
	EpisodeCount int     `json:"episode_count"`
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	Overview     string  `json:"overview"`
	PosterPath   string  `json:"poster_path"`
	SeasonNumber int     `json:"season_number"`
	VoteAverage  float64 `json:"vote_average"`
}

type TvDetails struct {
	Adult               bool                `json:"adult"`
	BackdropPath        string              `json:"backdrop_path"`
	CreatedBy           any                 `json:"created_by"`
	EpisodeRunTime      any                 `json:"episode_run_time"`
	FirstAirDate        string              `json:"first_air_date"`
	Genres              []Genre             `json:"genres"`
	Homepage            string              `json:"homepage"`
	Id                  int                 `json:"id"`
	InProduction        bool                `json:"in_production"`
	Languages           any                 `json:"languages"`
	LastAirDate         string              `json:"last_air_date"`
	LastEpisodeToAir    any                 `json:"last_episode_to_air"`
	Name                string              `json:"name"`
	NextEpisodeToAir    any                 `json:"next_episode_to_air"`
	Networks            []ProductionCompany `json:"networks"`
	NumberOfEpisodes    int                 `json:"number_of_episodes"`
	NumberOfSeasons     int                 `json:"number_of_seasons"`
	OriginCountry       any                 `json:"origin_country"`
	OriginalLanguage    string              `json:"original_language"`
	OriginalName        string              `json:"original_name"`
	Overview            string              `json:"overview"`
	Popularity          float64             `json:"popularity"`
	PosterPath          string              `json:"poster_path"`
	ProductionCompanies []ProductionCompany `json:"production_companies"`
	ProductionCountries any                 `json:"production_countries"`
	Seasons             []TvDetailsSeason   `json:"seasons"`
	SpokenLanguages     any                 `json:"spoken_languages"`
	Status              string              `json:"status"`
	Tagline             string              `json:"tagline"`
	Type                string              `json:"type"`
	VoteAverage         float64             `json:"vote_average"`
	VoteCount           int                 `json:"vote_count"`
}

const requestTTL = 1 * time.Hour

type ApiClient struct {
	client *provider.HTTPClient
	cache  cache.Cache
}

func NewApiClient(cache cache.Cache) *ApiClient {
	return &ApiClient{
		client: provider.NewHttpClient("https://api.themoviedb.org"),
		cache:  cache,
	}
}

type requestData struct {
	client *provider.HTTPClient
	cache  cache.Cache

	cacheKey string

	path  string
	query url.Values
}

func apiRequest[T any](ctx context.Context, req requestData) (T, error) {
	var res T

	if data, ok := cache.GetJson[T](req.cache, req.cacheKey); ok {
		fmt.Println("Using the cache", req.cacheKey)
		return data, nil
	}

	d, err := req.client.Get(ctx, req.path, provider.RequestOptions{
		Headers: http.Header{
			"accept":        {"application/json"},
			"Authorization": {"Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJhMWVlNGNjODJiM2MxNDI3MTI4MjQ0Zjg2MmRmYzdmNCIsIm5iZiI6MTc0MDUzMDU2MC4xNzEsInN1YiI6IjY3YmU2MzgwMGZmNjY0M2U3MTNkM2FiYiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.95qulxC-Skebj9avAqPtxgEuKZwlT_Bqp36iuqOF9wc"},
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
	cache.SetJson(req.cache, req.cacheKey, res, requestTTL)

	return res, nil
}

func (c *ApiClient) MovieSearch(ctx context.Context, query string) (SearchRequest[MovieSearchResult], error) {
	return apiRequest[SearchRequest[MovieSearchResult]](ctx, requestData{
		client:   c.client,
		cache:    c.cache,
		cacheKey: "api:movie-search:" + query,
		path:     "/3/search/movie",
		query: url.Values{
			"query":         {query},
			"include_adult": {"true"},
			"language":      {"en-US"},
			"page":          {"1"},
		},
	})
}

func (c *ApiClient) TvSearch(ctx context.Context, query string) (SearchRequest[TvSearchResult], error) {
	return apiRequest[SearchRequest[TvSearchResult]](ctx, requestData{
		client:   c.client,
		cache:    c.cache,
		cacheKey: "api:tv-search:" + query,
		path:     "/3/search/tv",
		query: url.Values{
			"query":         {query},
			"include_adult": {"true"},
			"language":      {"en-US"},
			"page":          {"1"},
		},
	})
}

func (c *ApiClient) GetMovieDetails(ctx context.Context, id string) (MovieDetails, error) {
	return apiRequest[MovieDetails](ctx, requestData{
		client:   c.client,
		cache:    c.cache,
		cacheKey: "api:movie-details:" + id,
		path:     fmt.Sprintf("/3/movie/%s", id),
		query: url.Values{
			"language": {"en-US"},
		},
	})
}

func (c *ApiClient) GetTvDetails(ctx context.Context, id string) (TvDetails, error) {
	return apiRequest[TvDetails](ctx, requestData{
		client:   c.client,
		cache:    c.cache,
		cacheKey: "api:tv-details:" + id,
		path:     fmt.Sprintf("/3/tv/%s", id),
		query: url.Values{
			"language": {"en-US"},
		},
	})
}

func (c *ApiClient) GetSeasonDetails(ctx context.Context, tvId, seasonNumber string) (SeasonDetails, error) {
	return apiRequest[SeasonDetails](ctx, requestData{
		client:   c.client,
		cache:    c.cache,
		cacheKey: "api:season-details:" + tvId + ":" + seasonNumber,
		path:     fmt.Sprintf("/3/tv/%s/season/%s", tvId, seasonNumber),
		query: url.Values{
			"language": {"en-US"},
		},
	})
}

func (c *ApiClient) GetMovieImages(ctx context.Context, id string) (Images, error) {
	return apiRequest[Images](ctx, requestData{
		client:   c.client,
		cache:    c.cache,
		cacheKey: "api:movie-images:" + id,
		path:     fmt.Sprintf("/3/movie/%s/images", id),
		query: url.Values{
			"include_image_language": {"en"},
		},
	})
}

func (c *ApiClient) GetTvImages(ctx context.Context, id string) (Images, error) {
	return apiRequest[Images](ctx, requestData{
		client:   c.client,
		cache:    c.cache,
		cacheKey: "api:tv-images:" + id,
		path:     fmt.Sprintf("/3/tv/%s/images", id),
		query: url.Values{
			"include_image_language": {"en"},
		},
	})
}
