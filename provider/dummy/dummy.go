package dummy

import (
	"context"
	"errors"
	"net/url"

	"github.com/nanoteck137/watchbook/provider"
	"github.com/nanoteck137/watchbook/types"
)

const ProviderName = "dummy"

var _ provider.Provider = (*DummyProvider)(nil)

type DummyProvider struct {
}

func (d *DummyProvider) Info() provider.Info {
	return provider.Info{
		Name:                    ProviderName,
		DisplayName:             "Dummy",
		SupportGetMedia:         true,
		SupportSearchMedia:      false,
		SupportGetCollection:    true,
		SupportSearchCollection: true,
	}
}

func (d *DummyProvider) GetCollection(ctx context.Context, id string) (provider.Collection, error) {
	if id == "1" {
		cover := "https://placehold.co/300x450/png?text=" + url.QueryEscape("Attack on Titan")

		return provider.Collection{
			Name:      "Attack on Titan",
			CoverUrl:  &cover,
			LogoUrl:   nil,
			BannerUrl: nil,
			Items:     []provider.CollectionItem{
				{
					Id:   "1@1",
					Name: "Season 1",
				},
				{
					Id:   "1@2",
					Name: "Season 2",
				},
				{
					Id:   "1@3",
					Name: "Season 3",
				},
				{
					Id:   "1@4",
					Name: "Season 4",
				},
			},
		}, nil
	}

	return provider.Collection{}, errors.New("not found")
}

func (d *DummyProvider) GetMedia(ctx context.Context, id string) (provider.Media, error) {
	switch id {
	case "1@1":
		title := "Attack on Titan Season 1" 
		cover := "https://placehold.co/300x450/png?text=" + url.QueryEscape(title)
		return provider.Media{
			ProviderId:       id,
			Type:             types.MediaTypeAnimeSeason,
			Title:            title,
			CoverUrl:         &cover,
			ExtraProviderIds: map[string]string{
				"myanimelist-anime": "21",
			},
		}, nil
	case "1@2":
		title := "Attack on Titan Season 2" 
		cover := "https://placehold.co/300x450/png?text=" + url.QueryEscape(title)
		return provider.Media{
			ProviderId:       id,
			Type:             types.MediaTypeAnimeSeason,
			Title:            title,
			CoverUrl:         &cover,
			ExtraProviderIds: map[string]string{
				"myanimelist-anime": "21",
			},
		}, nil
	case "1@3":
		title := "Attack on Titan Season 3" 
		cover := "https://placehold.co/300x450/png?text=" + url.QueryEscape(title)
		return provider.Media{
			ProviderId:       id,
			Type:             types.MediaTypeAnimeSeason,
			Title:            title,
			CoverUrl:         &cover,
			ExtraProviderIds: map[string]string{
				"myanimelist-anime": "21",
			},
		}, nil
	case "1@4":
		title := "Attack on Titan Final Season" 
		cover := "https://placehold.co/300x450/png?text=" + url.QueryEscape(title)
		return provider.Media{
			ProviderId:       id,
			Type:             types.MediaTypeAnimeSeason,
			Title:            title,
			CoverUrl:         &cover,
			ExtraProviderIds: map[string]string{
				"myanimelist-anime": "21",
			},
		}, nil
	}

	return provider.Media{}, errors.New("not found")
}

func (d *DummyProvider) SearchCollection(ctx context.Context, query string) ([]provider.SearchResult, error) {
	return []provider.SearchResult{
		{
			SearchType: provider.SearchResultTypeCollection,
			ProviderId: "1",
			Title:      "Attack on Titan",
			MediaType:  types.MediaTypeUnknown,
			ImageUrl:   "https://placehold.co/300x450/png?text=" + url.QueryEscape("Attack on Titan"),
		},
		{
			SearchType: provider.SearchResultTypeCollection,
			ProviderId: "2",
			Title:      "Lycoris Recoil",
			MediaType:  types.MediaTypeUnknown,
			ImageUrl:   "https://placehold.co/300x450/png?text=" + url.QueryEscape("Lycoris Recoil"),
		},
	}, nil
}

func (d *DummyProvider) SearchMedia(ctx context.Context, query string) ([]provider.SearchResult, error) {
	panic("unsupported")
}

