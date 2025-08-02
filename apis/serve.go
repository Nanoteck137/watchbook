package apis

import (
	"errors"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/watchbook"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
)

func RegisterHandlers(app core.App, router pyrin.Router) {
	g := router.Group("/api/v1")
	InstallAuthHandlers(app, g)
	InstallSystemHandlers(app, g)
	InstallUserHandlers(app, g)

	InstallMediaHandlers(app, g)
	InstallCollectionHandlers(app, g)
	InstallProviderHandlers(app, g)

	g = router.Group("/files")
	g.Register(
		pyrin.NormalHandler{
			Name:        "GetMediaImage",
			Method:      http.MethodGet,
			Path:        "/media/:id/images/:file",
			HandlerFunc: func(c pyrin.Context) error {
				id := c.Param("id")
				file := c.Param("file")

				mediaDir := app.WorkDir().MediaDirById(id)
				p := mediaDir.Images()
				f := os.DirFS(p)

				return pyrin.ServeFile(c, f, file)
			},
		},

		pyrin.NormalHandler{
			Name:        "GetCollectionImage",
			Method:      http.MethodGet,
			Path:        "/collections/:id/:image",
			HandlerFunc: func(c pyrin.Context) error {
				id := c.Param("id")
				image := c.Param("image")

				collection, err := app.DB().GetCollectionById(c.Request().Context(), nil, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return pyrin.NoContentNotFound()
					}
				}

				// TODO(patrik): Better handling
				name := strings.TrimSuffix(image, path.Ext(image))

				imageFile := ""
				switch name {
				case "cover":
					imageFile = collection.CoverFile.String
				case "logo":
					imageFile = collection.LogoFile.String
				case "banner":
					imageFile = collection.BannerFile.String
				}

				p := path.Dir(imageFile)

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
