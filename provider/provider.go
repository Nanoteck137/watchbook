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

	ExtraProviderIds map[string]string `json:"extraProviderIds"`
}

type Collection struct{}

type Provider interface {
	Name() string

	GetMedia(ctx context.Context, id string) (Media, error)
	SearchMedia(ctx context.Context, query string) ([]Media, error)

	GetCollection(ctx context.Context, id string) (Collection, error)
	SearchCollection(ctx context.Context, query string) ([]Collection, error)
}

var ErrNoProvider = errors.New("no provider")

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
		err = cache.SetJson(p.cache, cacheKey, m, 24 * time.Hour)
		if err != nil {
			return Media{}, err
		}

		return m, nil
	}

	return Media{}, err
}

func (p *ProviderManager) SearchMedia(ctx context.Context, providerName, query string) ([]Media, error) {
	return nil, nil
}

func (p *ProviderManager) GetCollection(ctx context.Context, providerName, id string) (Collection, error) {
	return Collection{}, nil
}

func (p *ProviderManager) SearchCollection(ctx context.Context, providerName, query string) ([]Collection, error) {
	return nil, nil
}
