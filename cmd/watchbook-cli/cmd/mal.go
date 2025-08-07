package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/kr/pretty"
	"github.com/nanoteck137/watchbook/cmd/watchbook-cli/api"
	"github.com/nanoteck137/watchbook/provider/myanimelist"
	"github.com/nanoteck137/watchbook/utils"
	"github.com/spf13/cobra"
)

// TODO(patrik): Return errors
func createMediaFromMalId(client *api.Client, tempDir, malId string) (string, error) {
	data, err := myanimelist.RawGetAnime(malId)
	if err != nil {
		logger.Fatal("failed to get anime data", "err", err, "id", malId)
	}

	title := data.Title
	if data.TitleEnglish != "" {
		title = data.TitleEnglish
	}

	score := 0.0
	if data.Score != nil {
		score = *data.Score
	}

	startDate := ""
	if data.StartDate != nil {
		startDate = *data.StartDate
	}

	endDate := ""
	if data.EndDate != nil {
		endDate = *data.EndDate
	}

	res, err := client.CreateMedia(api.CreateMediaBody{
		MediaType: string(data.Type),
		// TmdbId:         "",
		// ImdbId:         "",
		MalId: fmt.Sprintf("anime@%s", malId),
		// AnilistId:      "",
		Title:        title,
		Description:  data.Description,
		Score:        float32(score),
		Status:       string(data.Status),
		Rating:       string(data.Rating),
		AiringSeason: data.AiringSeason,
		StartDate:    startDate,
		EndDate:      endDate,
		// PartCount:      0,
		// CoverUrl:       "",
		// BannerUrl:      "",
		// LogoUrl:        "",
		Creators: data.Studios,
		Tags:    data.Tags,
		// CollectionId:   "",
		// CollectionName: "",
	}, api.Options{})
	if err != nil {
		logger.Fatal("failed to create media", "err", err, "id", malId)
	}

	if data.Type.IsMovie() {
		_, err := client.SetParts(res.Id, api.SetPartsBody{
			Parts: []api.PartBody{
				{
					Name: title,
				},
			},
		}, api.Options{})
		if err != nil {
			logger.Fatal("failed to set parts", "err", err, "id", malId)
		}
	} else if data.EpisodeCount != nil {
		parts := make([]api.PartBody, 0, *data.EpisodeCount)
		for i := range *data.EpisodeCount {
			parts = append(parts, api.PartBody{
				Name: fmt.Sprintf("Episode %d", i+1),
			})
		}

		_, err := client.SetParts(res.Id, api.SetPartsBody{
			Parts: parts,
		}, api.Options{})
		if err != nil {
			logger.Fatal("failed to set parts", "err", err, "id", malId)
		}
	}

	p, err := utils.DownloadImage(data.CoverImageUrl, tempDir, "cover")
	if err != nil {
		logger.Fatal("failed to download image", "err", err, "id", malId)
	}

	images, err := createImageForm(p, "", "")
	if err != nil {
		logger.Fatal("failed to create image form", "err", err, "id", malId)
	}

	_, err = client.ChangeMediaImages(res.Id, images.Boundary, &images.Buf, api.Options{})
	if err != nil {
		logger.Fatal("failed to set images", "err", err, "id", malId)
	}

	return res.Id, nil
}

var malCmd = &cobra.Command{
	Use: "mal",
}

var malGetCmd = &cobra.Command{
	Use:  "get",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.New("http://localhost:3000")

		dir, err := os.MkdirTemp("", "watchbook-cli-*")
		if err != nil {
			logger.Fatal("failed to create temp dir", "err", err)
		}
		defer os.RemoveAll(dir)

		for _, malId := range args {
			search, err := client.GetMedia(api.Options{
				Query: url.Values{
					"filter": {fmt.Sprintf(`malId=="anime@%s"`, malId)},
				},
			})
			if err != nil {
				logger.Fatal("failed to get media", "err", err)
			}

			if len(search.Media) > 0 {
				logger.Warn("entry with id already exists", "id", malId)
				continue
			}

			_, err = createMediaFromMalId(client, dir, malId)
			if err != nil {
				logger.Warn("failed to create media", "err", err, "id", malId)
			}
		}
	},
}

var malCreateCollectionCmd = &cobra.Command{
	Use:  "create-collection",
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		ids := args[1:]

		client := api.New("http://localhost:3000")

		dir, err := os.MkdirTemp("", "watchbook-cli-*")
		if err != nil {
			logger.Fatal("failed to create temp dir", "err", err)
		}
		defer os.RemoveAll(dir)

		res, err := client.CreateCollection(api.CreateCollectionBody{
			CollectionType: "anime",
			Name:           name,
		}, api.Options{})
		if err != nil {
			logger.Fatal("failed to create collection", "err", err)
		}

		_ = res

		var media []api.GetMediaById

		for _, id := range ids {
			search, err := client.GetMedia(api.Options{
				Query: url.Values{
					"filter": {fmt.Sprintf(`malId=="anime@%s"`, id)},
				},
			})
			if err != nil {
				logger.Fatal("failed to get media", "err", err)
			}

			realId := ""

			if len(search.Media) <= 0 {
				realId, err = createMediaFromMalId(client, dir, id)
				if err != nil {
					logger.Warn("failed to create media", "err", err, "id", id)
				}
			} else {
				realId = search.Media[0].Id
			}

			m, err := client.GetMediaById(realId, api.Options{})
			if err != nil {
				logger.Fatal("failed to get media", "err", err)
			}

			media = append(media, *m)
		}

		type Entry struct {
			Media api.GetMediaById

			Name       string
			SearchSlug string
			Order      string
		}

		var groups []*huh.Group
		var entries []*Entry
		for _, m := range media {
			entry := &Entry{
				Media: m,
			}

			g := huh.NewGroup(
				huh.NewInput().
					Title("Entry Name").
					Value(&entry.Name),
				huh.NewInput().
					Title("Entry Search Slug (defaults to name)").
					Value(&entry.SearchSlug),
				huh.NewInput().
					Title("Entry Order").
					Value(&entry.Order),
			).
				Title(m.Title)

			groups = append(groups, g)
			entries = append(entries, entry)
		}

		form := huh.NewForm(
			groups...,
		)
		err = form.Run()
		if err != nil {
			logger.Fatal("failed to run form", "err", err)
		}

		pretty.Println(entries)
		for _, entry := range entries {
			searchSlug := utils.Slug(entry.SearchSlug)
			if searchSlug == "" {
				searchSlug = utils.Slug(entry.Name)
			} 

			order, _ := strconv.Atoi(entry.Order)

			_, err = client.AddCollectionItem(res.Id, api.AddCollectionItemBody{
				MediaId:    entry.Media.Id,
				Name:       entry.Name,
				SearchSlug: searchSlug,
				Order:      order,
			}, api.Options{})
			if err != nil {
				logger.Fatal("failed to add media to collection", "err", err, "entry", entry.Media.Title)
			}
		}
	},
}

func init() {
	malCmd.AddCommand(malGetCmd)
	malCmd.AddCommand(malCreateCollectionCmd)

	rootCmd.AddCommand(malCmd)
}
