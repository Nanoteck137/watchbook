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

type CollectionItem struct {
	Id   string
	Name string
}

type Collection struct {
	Name string

	CoverUrl  *string
	LogoUrl   *string
	BannerUrl *string

	Items []CollectionItem
}

type ProviderInfo struct {
	Name        string
	DisplayName string

	SupportGetMedia    bool
	SupportSearchMedia bool

	SupportGetCollection    bool
	SupportSearchCollection bool
}

type Provider interface {
	Info() ProviderInfo

	GetMedia(ctx context.Context, id string) (Media, error)
	SearchMedia(ctx context.Context, query string) ([]SearchResult, error)

	GetCollection(ctx context.Context, id string) (Collection, error)
	SearchCollection(ctx context.Context, query string) ([]SearchResult, error)
}

const (
	mediaTTL  = 24 * time.Hour
	searchTTL = 1 * time.Hour
)

var ErrNoProvider = errors.New("no provider")

type ProviderManager struct {
	providers map[string]Provider
	cache     *cache.Provider
}

func NewProviderManager(cache *cache.Provider) *ProviderManager {
	return &ProviderManager{
		providers: map[string]Provider{},
		cache:     cache,
	}
}

func (p *ProviderManager) RegisterProvider(provider Provider) {
	name := provider.Info().Name
	if name != "" {
		p.providers[provider.Info().Name] = provider
	}
}

func (p *ProviderManager) IsValidProvider(name string) bool {
	_, ok := p.providers[name]
	return ok
}

func (p *ProviderManager) GetProviders() []ProviderInfo {
	res := make([]ProviderInfo, 0, len(p.providers))

	for _, p := range p.providers {
		res = append(res, p.Info())
	}

	return res
}

func (p *ProviderManager) GetMedia(ctx context.Context, providerName, id string) (Media, error) {
	if !p.IsValidProvider(providerName) {
		return Media{}, ErrNoProvider
	}

	provider := p.providers[providerName]
	cacheKey := fmt.Sprintf("%s:media:%s", providerName, id)

	media, err := cache.GetProviderJson[Media](p.cache, cacheKey)
	if err == nil {
		return media, nil
	}

	if errors.Is(err, cache.ErrNoData) {
		m, err := provider.GetMedia(ctx, id)
		if err != nil {
			return Media{}, err
		}

		err = cache.SetProviderJson(p.cache, cacheKey, providerName, m, mediaTTL)
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
	cacheKey := fmt.Sprintf("%s:media-search:%s", providerName, query)

	media, err := cache.GetProviderJson[[]SearchResult](p.cache, cacheKey)
	if err == nil {
		return media, nil
	}

	if errors.Is(err, cache.ErrNoData) {
		items, err := provider.SearchMedia(ctx, query)
		if err != nil {
			return nil, err
		}

		err = cache.SetProviderJson(p.cache, cacheKey, providerName, items, searchTTL)
		if err != nil {
			return nil, err
		}

		return items, nil
	}

	return nil, err
}

func (p *ProviderManager) GetCollection(ctx context.Context, providerName, id string) (Collection, error) {
	if !p.IsValidProvider(providerName) {
		return Collection{}, ErrNoProvider
	}

	provider := p.providers[providerName]
	cacheKey := fmt.Sprintf("%s:collections:%s", providerName, id)

	col, err := cache.GetProviderJson[Collection](p.cache, cacheKey)
	if err == nil {
		return col, nil
	}

	if !errors.Is(err, cache.ErrNoData) {
		return Collection{}, err
	}

	col, err = provider.GetCollection(ctx, id)
	if err != nil {
		return Collection{}, err
	}

	err = cache.SetProviderJson(p.cache, cacheKey, providerName, col, mediaTTL)
	if err != nil {
		return Collection{}, err
	}

	return col, nil
}

func (p *ProviderManager) SearchCollection(ctx context.Context, providerName, query string) ([]SearchResult, error) {
	if !p.IsValidProvider(providerName) {
		return nil, ErrNoProvider
	}

	provider := p.providers[providerName]
	cacheKey := fmt.Sprintf("%s:collections-search:%s", providerName, query)

	media, err := cache.GetProviderJson[[]SearchResult](p.cache, cacheKey)
	if err == nil {
		return media, nil
	}

	if errors.Is(err, cache.ErrNoData) {
		items, err := provider.SearchCollection(ctx, query)
		if err != nil {
			return nil, err
		}

		err = cache.SetProviderJson(p.cache, cacheKey, providerName, items, searchTTL)
		if err != nil {
			return nil, err
		}

		return items, nil
	}

	return nil, err
}
