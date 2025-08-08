package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"os"
	"path"

	"github.com/nanoteck137/watchbook/cmd/watchbook-cli/api"
	"github.com/nanoteck137/watchbook/cmd/watchbook-cli/library"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
	"github.com/spf13/cobra"
)

var oldCmd = &cobra.Command{
	Use: "old",
}

type ImageRes struct {
	Buf      bytes.Buffer
	Boundary string
}

func createImageForm(coverPath, logoPath, bannerPath string) (ImageRes, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	createFileField := func(fieldname, filename, contentType string) (io.Writer, error) {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
				escapeQuotes(fieldname), escapeQuotes(filename)))
		h.Set("Content-Type", contentType)
		return w.CreatePart(h)
	}

	createFormPart := func(p, name string) error {
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()

		contentType, err := utils.ImageExtToContentType(path.Ext(p))
		if err != nil {
			return err
		}

		formFile, err := createFileField(name, path.Base(p), contentType)
		if err != nil {
			return err
		}

		_, err = io.Copy(formFile, f)
		if err != nil {
			return err
		}

		return nil
	}

	if coverPath != "" {
		err := createFormPart(coverPath, "cover")
		if err != nil {
			return ImageRes{}, err
		}
	}

	if logoPath != "" {
		err := createFormPart(logoPath, "logo")
		if err != nil {
			return ImageRes{}, err
		}
	}

	if bannerPath != "" {
		err := createFormPart(bannerPath, "banner")
		if err != nil {
			return ImageRes{}, err
		}
	}

	w.Close()

	return ImageRes{
		Buf:      b,
		Boundary: w.Boundary(),
	}, nil
}

var oldImportCmd = &cobra.Command{
	Use: "import",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")

		apiAddress, _ := cmd.Flags().GetString("api-address")
		client := api.New(apiAddress)

		lib, err := library.SearchLibrary(dir)
		if err != nil {
			logger.Fatal("failed to read media", "err", err)
		}

		for _, m := range lib.Media {
			search, err := client.GetMedia(api.Options{
				Query: url.Values{
					"filter": {fmt.Sprintf(`malId=="anime@%s"`, m.Ids.MyAnimeList)},
				},
			})
			if err != nil {
				logger.Fatal("failed to get media", "err", err)
			}

			if len(search.Media) > 0 {
				logger.Warn("entry with id already exists", "id", m.Ids.MyAnimeList)
				continue
			}

			status := types.MediaStatusUnknown
			switch m.General.Status {
			case "finished":
				status = types.MediaStatusCompleted
			case "not-aired":
				status = types.MediaStatusUpcoming
			case "airing":
				status = types.MediaStatusOngoing
			default:
				logger.Warn("unknown status", "status", m.General.Status)
			}

			res, err := client.CreateMedia(api.CreateMediaBody{
				MediaType:    string(m.MediaType),
				TmdbId:       "",
				MalId:        "anime@" + m.Ids.MyAnimeList,
				AnilistId:    "",
				Title:        m.General.Title,
				Description:  m.General.Description,
				Score:        float32(m.General.Score),
				Status:       string(status),
				Rating:       string(m.General.Rating),
				AiringSeason: m.General.AiringSeason,
				StartDate:    m.General.StartDate,
				EndDate:      m.General.EndDate,
				Tags:         m.General.Tags,
				Creators:     m.General.Studios,
			}, api.Options{})

			if err != nil {
				logger.Fatal("failed to create media", "err", err)
			}

			images, err := createImageForm(m.GetCoverPath(), m.GetLogoPath(), m.GetBannerPath())
			if err != nil {
				logger.Fatal("failed create image form", "err", err)
			}

			_, err = client.ChangeMediaImages(res.Id, images.Boundary, &images.Buf, api.Options{})
			if err != nil {
				logger.Fatal("failed to change images", "err", err)
			}

			numEpisodes := len(m.Parts)

			var parts []api.PartBody
			for i := range numEpisodes {
				parts = append(parts, api.PartBody{
					Name: fmt.Sprintf("Episode %d", i+1),
				})
			}

			_, err = client.SetParts(res.Id, api.SetPartsBody{
				Parts: parts,
			}, api.Options{})
			if err != nil {
				logger.Fatal("failed to set parts", "err", err)
			}
		}

	},
}

var oldColCmd = &cobra.Command{
	Use: "col",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")

		apiAddress, _ := cmd.Flags().GetString("api-address")
		client := api.New(apiAddress)

		col, err := library.ReadCollection(dir)
		if err != nil {
			logger.Fatal("failed to read collection", "err", err)
		}

		c, err := client.CreateCollection(api.CreateCollectionBody{
			CollectionType: string(col.Type),
			Name:           col.General.Name,
		}, api.Options{})
		if err != nil {
			logger.Fatal("failed", "err", err)
		}

		images, err := createImageForm(col.GetCoverPath(), col.GetLogoPath(), col.GetBannerPath())
		if err != nil {
			logger.Fatal("failed", "err", err)
		}

		_, err = client.ChangeCollectionImages(c.Id, images.Boundary, &images.Buf, api.Options{})
		if err != nil {
			logger.Fatal("failed", "err", err)
		}

		for _, group := range col.Groups {
			for _, entry := range group.Entries {
				p := path.Join(col.Path, entry.Path)

				media, err := library.ReadMedia(p)
				if err != nil {
					logger.Fatal("failed to read media", "path", p)
				}

				search, err := client.GetMedia(api.Options{
					Query: url.Values{
						"filter": {fmt.Sprintf(`malId=="anime@%s"`, media.Ids.MyAnimeList)},
					},
				})
				if err != nil {
					logger.Fatal("failed to get media", "err", err)
				}

				if len(search.Media) > 0 {
					cm := search.Media[0]

					_, err = client.AddCollectionItem(c.Id, api.AddCollectionItemBody{
						MediaId:    cm.Id,
						Name:       entry.Name,
						SearchSlug: entry.SearchSlug,
						Order:      entry.Order,
					}, api.Options{})
					if err != nil {
						logger.Fatal("failed to add item to collection", "err", err)
					}
				}
			}
		}
	},
}

func init() {
	oldImportCmd.Flags().StringP("dir", "d", ".", "directory to import")
	oldColCmd.Flags().StringP("dir", "d", ".", "directory to import")

	oldCmd.AddCommand(oldImportCmd)
	oldCmd.AddCommand(oldColCmd)

	rootCmd.AddCommand(oldCmd)
}
