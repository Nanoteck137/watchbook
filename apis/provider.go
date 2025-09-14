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
	)
}
