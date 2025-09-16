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

// TODO(patrik): Add date aired
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
	Id       string
	Name     string
	Position int
}

type Collection struct {
	ProviderId string
	Type       types.CollectionType

	Name string

	CoverUrl  *string
	LogoUrl   *string
	BannerUrl *string

	Items []CollectionItem
}

type Info struct {
	Name        string
	DisplayName string

	SupportGetMedia    bool
	SupportSearchMedia bool

	SupportGetCollection    bool
	SupportSearchCollection bool
}

func (i Info) GetDisplayName() string {
	if i.DisplayName != "" {
		return i.DisplayName
	}

	return i.Name
}

type Context struct {
	ctx   context.Context
	cache cache.Cache
}

func (c *Context) Context() context.Context {
	return c.ctx
}

func (c *Context) Cache() cache.Cache {
	return c.cache
}

type Provider interface {
	Info() Info

	GetMedia(c Context, id string) (Media, error)
	SearchMedia(c Context, query string) ([]SearchResult, error)

	GetCollection(c Context, id string) (Collection, error)
	SearchCollection(c Context, query string) ([]SearchResult, error)
}

const (
	mediaTTL  = 24 * time.Hour
	searchTTL = 1 * time.Hour
)

var ErrNoProvider = errors.New("no provider")

type ProviderManager struct {
	providers     map[string]Provider
	providerInfos map[string]Info
	cache         *cache.ProviderCache
}

func NewProviderManager(cache *cache.ProviderCache) *ProviderManager {
	return &ProviderManager{
		providers:     map[string]Provider{},
		providerInfos: map[string]Info{},
		cache:         cache,
	}
}

func (p *ProviderManager) RegisterProvider(provider Provider) {
	info := provider.Info()
	name := info.Name
	if name != "" {
		p.providers[info.Name] = provider
		p.providerInfos[info.Name] = info
	}
}

func (p *ProviderManager) GetProviderInfo(name string) (Info, bool) {
	info, ok := p.providerInfos[name]
	return info, ok
}

func (p *ProviderManager) IsValidProvider(name string) bool {
	_, ok := p.providers[name]
	return ok
}

func (p *ProviderManager) GetProviders() []Info {
	res := make([]Info, 0, len(p.providers))

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
	cacheKey := fmt.Sprintf("media:%s", id)

	providerCache := p.cache.WithName(providerName)

	if data, ok := cache.GetJson[Media](providerCache, cacheKey); ok {
		return data, nil
	}

	c := Context{
		ctx:   ctx,
		cache: providerCache,
	}

	m, err := provider.GetMedia(c, id)
	if err != nil {
		return Media{}, err
	}

	err = cache.SetJson(providerCache, cacheKey, m, mediaTTL)
	if err != nil {
		return Media{}, err
	}

	return m, nil
}

func (p *ProviderManager) SearchMedia(ctx context.Context, providerName, query string) ([]SearchResult, error) {
	if !p.IsValidProvider(providerName) {
		return nil, ErrNoProvider
	}

	provider := p.providers[providerName]
	cacheKey := fmt.Sprintf("media-search:%s", query)

	providerCache := p.cache.WithName(providerName)

	if data, ok := cache.GetJson[[]SearchResult](providerCache, cacheKey); ok {
		return data, nil
	}

	c := Context{
		ctx:   ctx,
		cache: providerCache,
	}

	items, err := provider.SearchMedia(c, query)
	if err != nil {
		return nil, err
	}

	err = cache.SetJson(providerCache, cacheKey, items, searchTTL)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (p *ProviderManager) GetCollection(ctx context.Context, providerName, id string) (Collection, error) {
	if !p.IsValidProvider(providerName) {
		return Collection{}, ErrNoProvider
	}

	provider := p.providers[providerName]
	cacheKey := fmt.Sprintf("collections:%s", id)

	providerCache := p.cache.WithName(providerName)

	if data, ok := cache.GetJson[Collection](providerCache, cacheKey); ok {
		return data, nil
	}

	c := Context{
		ctx:   ctx,
		cache: providerCache,
	}

	col, err := provider.GetCollection(c, id)
	if err != nil {
		return Collection{}, err
	}

	err = cache.SetJson(providerCache, cacheKey, col, mediaTTL)
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
	cacheKey := fmt.Sprintf("collections-search:%s", query)

	providerCache := p.cache.WithName(providerName)

	if data, ok := cache.GetJson[[]SearchResult](providerCache, cacheKey); ok {
		return data, nil
	}

	c := Context{
		ctx:   ctx,
		cache: providerCache,
	}

	items, err := provider.SearchCollection(c, query)
	if err != nil {
		return nil, err
	}

	err = cache.SetJson(providerCache, cacheKey, items, searchTTL)
	if err != nil {
		return nil, err
	}

	return items, nil
}
