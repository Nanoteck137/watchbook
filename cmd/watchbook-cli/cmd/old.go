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

	"github.com/kr/pretty"
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
		apiAddress, _ := cmd.Flags().GetString("api-address")
		client := api.New(apiAddress)

		lib, err := library.SearchLibrary("/Volumes/media/watch/mal")
		if err != nil {
			logger.Fatal("failed to read media", "err", err)
		}

		for _, m := range lib.Media {
			lel, err := client.GetMedia(api.Options{
				Query: url.Values{
					"filter": {fmt.Sprintf(`malId=="anime@%s"`, m.Ids.MyAnimeList)},
				},
			})
			if err != nil {
				logger.Fatal("failed to get media", "err", err)
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

			id := ""

			if len(lel.Media) <= 0 {
				logger.Info("media not found, creating new entry", "title", m.General.Title)

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

				pretty.Println(res)
				id = res.Id
			} else {
				id = lel.Media[0].Id
			}

			fmt.Printf("id: %v\n", id)

			images, err := createImageForm(m.GetCoverPath(), m.GetLogoPath(), m.GetBannerPath())
			if err != nil {
				logger.Fatal("failed", "err", err)
			}

			_, err = client.ChangeMediaImages(id, images.Boundary, &images.Buf, api.Options{})
			if err != nil {
				logger.Fatal("failed", "err", err)
			}

			numEpisodes := len(m.Parts)

			var parts []api.PartBody
			for i := range numEpisodes {
				parts = append(parts, api.PartBody{
					Name: fmt.Sprintf("Episode %d", i+1),
				})
			}

			_, err = client.SetParts(id, api.SetPartsBody{
				Parts: parts,
			}, api.Options{})
			if err != nil {
				pretty.Println(err)
				logger.Fatal("failed", "err", err)
			}
		}

	},
}

var oldColCmd = &cobra.Command{
	Use: "col",
	Run: func(cmd *cobra.Command, args []string) {
		apiAddress, _ := cmd.Flags().GetString("api-address")
		client := api.New(apiAddress)

		dir := "/Volumes/media/watch/mal/series"
		lib, err := library.SearchLibrary(dir)
		if err != nil {
			logger.Fatal("failed to read media", "err", err)
		}

		_ = lib
		_ = client

		mediaPathMapping := make(map[string]library.Media)

		for _, m := range lib.Media {
			mediaPathMapping[m.Path] = m
		}

		for _, col := range lib.Collections {
			pretty.Println(col)
			// col.General.Name

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
					fmt.Printf("p: %v\n", p)
					media, ok := mediaPathMapping[p]
					if !ok {
						logger.Fatal("failed to map path to media", "path", p)
					}

					fmt.Printf("media.Ids.MyAnimeList: %v\n", media.Ids.MyAnimeList)

					lel, err := client.GetMedia(api.Options{
						Query: url.Values{
							"filter": {fmt.Sprintf(`malId=="anime@%s"`, media.Ids.MyAnimeList)},
						},
					})
					if err != nil {
						logger.Fatal("failed to get media", "err", err)
					}

					pretty.Println(lel)

					if len(lel.Media) > 0 {
						cm := lel.Media[0]

						_, err = client.AddCollectionItem(c.Id, api.AddCollectionItemBody{
							MediaId:    cm.Id,
							Name:       entry.Name,
							SearchSlug: entry.SearchSlug,
							Order:      entry.Order,
						}, api.Options{})
						if err != nil {
							logger.Fatal("failed", "err", err)
						}
					}

					//
					// err := db.CreateCollectionMediaItem(ctx, database.CreateCollectionMediaItemParams{
					// 	CollectionId:   dbCollection.Id,
					// 	MediaId:        mediaId,
					// 	GroupName:      group.Name,
					// 	GroupOrder:     int64(group.Order),
					// 	Name:           entry.Name,
					// 	OrderNumber:    int64(entry.Order),
					// 	SubOrderNumber: int64(0),
					// 	SearchSlug:     entry.SearchSlug,
					// })
					// if err != nil {
					// 	return fmt.Errorf("failed to add media to collection: %w", err)
					// }
				}
			}

		}
	},
}

func init() {
	oldCmd.AddCommand(oldImportCmd)
	oldCmd.AddCommand(oldColCmd)

	rootCmd.AddCommand(oldCmd)
}
