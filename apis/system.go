package apis

import (
	"net/http"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/watchbook"
	"github.com/nanoteck137/watchbook/core"
)

type GetSystemInfo struct {
	Version string `json:"version"`
}

func InstallSystemHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetSystemInfo",
			Path:         "/system/info",
			Method:       http.MethodGet,
			ResponseType: GetSystemInfo{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				return GetSystemInfo{
					Version: watchbook.Version,
				}, nil
			},
		},
	)
}
