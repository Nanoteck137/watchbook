package apis

import (
	"net/http"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/watchbook/core"
)

type Provider struct {
	Name string `json:"name"`
}

type GetProviders struct {
	Providers []Provider `json:"providers"`
}

type ProviderSearchResult struct {
	Provider   string `json:"provider"`
	ProviderId string `json:"providerId"`
	Title      string `json:"title"`
	ImageUrl   string `json:"imageUrl"`
}

type GetProviderSearch struct {
	SearchResults []ProviderSearchResult `json:"searchResults"`
}

func InstallProviderHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetProviders",
			Method:       http.MethodGet,
			Path:         "/providers",
			ResponseType: GetProviders{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				providers := app.ProviderManager().GetProviders()

				res := GetProviders{
					Providers: make([]Provider, len(providers)),
				}

				for i, p := range providers {
					res.Providers[i] = Provider{
						Name: p.Name,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "ProviderSearchMedia",
			Method:       http.MethodGet,
			Path:         "/providers/:providerName",
			ResponseType: GetProviderSearch{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				providerName := c.Param("providerName")
				url := c.Request().URL
				query := url.Query().Get("query")
				_ = query

				pm := app.ProviderManager()
				_ = pm

				// items, err := pm.SearchMedia(context.Background(), providerName, query)
				// if err != nil {
				// 	// TODO(patrik): Handle some of the errors the provider gives, like the provider not found
				// 	return nil, err
				// }

				res := GetProviderSearch{
					// SearchResults: make([]ProviderSearchResult, len(items)),
					SearchResults: []ProviderSearchResult{
						{
							Provider:   providerName,
							ProviderId: "21",
							Title:      "One Piece",
							ImageUrl:   "https://cdn.myanimelist.net/images/anime/1244/138851.jpg",
						},
					},
				}

				// for i, m := range items {
				// 	res.Providers[i] = Provider{
				// 		Name: p.Name,
				// 	}
				// }

				return res, nil
			},
		},
	)
}
