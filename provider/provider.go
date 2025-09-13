package provider

import (
	"context"
	"time"

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
	GetMedia(ctx context.Context, id string) (Media, error)
	SearchMedia(ctx context.Context, query string) ([]Media, error)

	GetCollection(ctx context.Context, id string) (Collection, error)
	SearchCollection(ctx context.Context, query string) ([]Collection, error)
}
