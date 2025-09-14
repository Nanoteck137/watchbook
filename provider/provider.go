package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nanoteck137/watchbook/tools/cache"
	"github.com/nanoteck137/watchbook/types"
)

const (
	ProviderNameMyAnimeListAnime string = "myanimelist-anime"
	ProviderNameAnilistAnime     string = "anilist-anime"
	ProviderNameTheMovieDbMovie  string = "tmdb-movie"
	ProviderNameTheMovieDbTv     string = "tmdb-tv"
)

type SearchResultType string

const (
	SearchResultTypeMedia      SearchResultType = "media"
	SearchResultTypeCollection SearchResultType = "collection"
)

type SearchResult struct {
	SearchType SearchResultType `json:"searchType"`
	ProviderId string           `json:"providerId"`
	Title      string           `json:"title"`
	MediaType  types.MediaType  `json:"mediaType"`
	ImageUrl   string           `json:"imageUrl"`
}

type MediaPart struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
}

type Media struct {
	ProviderId string          `json:"id"`
	Type       types.MediaType `json:"type"`

	Title       string  `json:"title"`
	Description *string `json:"description"`

	Score        *float64          `json:"score"`
	Status       types.MediaStatus `json:"status"`
	Rating       types.MediaRating `json:"rating"`
	AiringSeason *string           `json:"airingSeason"`

	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`

	CoverUrl  *string `json:"coverUrl"`
	LogoUrl   *string `json:"logoUrl"`
	BannerUrl *string `json:"bannerUrl"`

	Creators []string `json:"creators"`
	Tags     []string `json:"tags"`

	Parts []MediaPart `json:"parts"`

	ExtraProviderIds map[string]string `json:"extraProviderIds"`
}

type Collection struct{}

type Provider interface {
	Name() string

	GetMedia(ctx context.Context, id string) (Media, error)
	SearchMedia(ctx context.Context, query string) ([]SearchResult, error)

	GetCollection(ctx context.Context, id string) (Collection, error)
	SearchCollection(ctx context.Context, query string) ([]SearchResult, error)
}

var ErrNoProvider = errors.New("no provider")

type ProviderInfo struct {
	Name string
}

type ProviderManager struct {
	providers map[string]Provider
	cache     *cache.Cache
}

func NewProviderManager(cache *cache.Cache) *ProviderManager {
	return &ProviderManager{
		providers: map[string]Provider{},
		cache:     cache,
	}
}

func (p *ProviderManager) RegisterProvider(provider Provider) {
	p.providers[provider.Name()] = provider
}

func (p *ProviderManager) IsValidProvider(name string) bool {
	_, ok := p.providers[name]
	return ok
}

func (p *ProviderManager) GetProviders() []ProviderInfo {
	res := make([]ProviderInfo, 0, len(p.providers))

	for _, p := range p.providers {
		res = append(res, ProviderInfo{
			Name: p.Name(),
		})
	}

	return res
}

func (p *ProviderManager) GetMedia(ctx context.Context, providerName, id string) (Media, error) {
	if !p.IsValidProvider(providerName) {
		return Media{}, ErrNoProvider
	}

	provider := p.providers[providerName]
	cacheKey := fmt.Sprintf("%s:media:%s", provider.Name(), id)

	media, err := cache.GetJson[Media](p.cache, cacheKey)
	if err == nil {
		fmt.Println("Using the cached version")
		return media, nil
	}

	if errors.Is(err, cache.ErrNoData) {
		fmt.Println("Data not found in cache, fetching new data")

		m, err := provider.GetMedia(ctx, id)
		if err != nil {
			return Media{}, err
		}

		// TODO(patrik): Make ttl to const
		err = cache.SetJson(p.cache, cacheKey, m, 24*time.Hour)
		if err != nil {
			return Media{}, err
		}

		return m, nil
	}

	return Media{}, err
}

func (p *ProviderManager) SearchMedia(ctx context.Context, providerName, query string) ([]SearchResult, error) {
	if !p.IsValidProvider(providerName) {
		return nil, ErrNoProvider
	}

	provider := p.providers[providerName]
	cacheKey := fmt.Sprintf("%s:media-search:%s", provider.Name(), query)

	media, err := cache.GetJson[[]SearchResult](p.cache, cacheKey)
	if err == nil {
		fmt.Println("Using the cached version")
		return media, nil
	}

	if errors.Is(err, cache.ErrNoData) {
		fmt.Println("Data not found in cache, fetching new data")

		items, err := provider.SearchMedia(ctx, query)
		if err != nil {
			return nil, err
		}

		// TODO(patrik): Make ttl to const
		err = cache.SetJson(p.cache, cacheKey, items, 1*time.Second)
		if err != nil {
			return nil, err
		}

		return items, nil
	}

	return nil, err
}

func (p *ProviderManager) GetCollection(ctx context.Context, providerName, id string) (Collection, error) {
	return Collection{}, nil
}

func (p *ProviderManager) SearchCollection(ctx context.Context, providerName, query string) ([]Collection, error) {
	return nil, nil
}
