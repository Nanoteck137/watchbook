package apis

import (
	"net/http"
	"os"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/watchbook"
	"github.com/nanoteck137/watchbook/core"
)

func RegisterHandlers(app core.App, router pyrin.Router) {
	g := router.Group("/api/v1")
	InstallAuthHandlers(app, g)
	InstallSystemHandlers(app, g)
	InstallUserHandlers(app, g)

	InstallMediaHandlers(app, g)
	InstallProviderHandlers(app, g)

	g = router.Group("/files")
	g.Register(
		pyrin.NormalHandler{
			Name:        "GetMediaImage",
			Method:      http.MethodGet,
			Path:        "/media/:id/:image",
			HandlerFunc: func(c pyrin.Context) error {
				id := c.Param("id")
				image := c.Param("image")

				mediaDir := app.WorkDir().MediaDir()
				p := mediaDir.MediaImageDir(id)

				f := os.DirFS(p)
				return pyrin.ServeFile(c, f, image)
			},
		},
	)
}

func Server(app core.App) (*pyrin.Server, error) {
	s := pyrin.NewServer(&pyrin.ServerConfig{
		LogName: watchbook.AppName,
		RegisterHandlers: func(router pyrin.Router) {
			RegisterHandlers(app, router)
		},
	})

	return s, nil
}
