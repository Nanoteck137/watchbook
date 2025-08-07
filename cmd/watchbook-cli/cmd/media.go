package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/kr/pretty"
	"github.com/nanoteck137/watchbook/cmd/watchbook-cli/api"
	"github.com/nanoteck137/watchbook/utils"
	"github.com/spf13/cobra"
)

var mediaCmd = &cobra.Command{
	Use: "media",
}

var mediaEditCmd = &cobra.Command{
	Use: "edit",
	Run: func(cmd *cobra.Command, args []string) {
		apiAddress, _ := cmd.Flags().GetString("api-address")
		client := api.New(apiAddress)

		var searchQuery string
		var id string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Search Query").
					Value(&searchQuery),
				huh.NewSelect[string]().
					Title("Results").
					Value(&id).
					OptionsFunc(func() []huh.Option[string] {
						media, err := client.GetMedia(api.Options{
							Query: url.Values{
								"filter":  {fmt.Sprintf("title %% \"%%%s%%\"", searchQuery)},
								"perPage": {"10"},
							},
						})
						if err != nil {
							logger.Error("failed to get media", "err", err)
							return []huh.Option[string]{}
						}

						options := make([]huh.Option[string], 0, len(media.Media))

						for _, media := range media.Media {
							options = append(options, huh.NewOption(media.Title, media.Id))
						}

						return options
					}, &searchQuery),
			),
		)

		err := form.Run()
		if err != nil {
			logger.Fatal("failed to run form", "err", err)
		}

		fmt.Printf("id: %v\n", id)

		media, err := client.GetMediaById(id, api.Options{})
		if err != nil {
			logger.Fatal("failed to get media by id", "err", err, "id", id)
		}

		pretty.Println(media)

		{
			var selected string
			form := huh.NewSelect[string]().
				Title("Main Menu").
				Options(
					huh.NewOption("Edit Info", "edit-info"),
					huh.NewOption("Change images", "change-images"),
				).
				Value(&selected)
			err := form.Run()
			if err != nil {
				logger.Fatal("failed to run form", "err", err)
			}

			switch selected {
			case "edit-info":
				// TODO(patrik): More properties
				valOrEmptyString := func(val *string) string {
					if val == nil {
						return ""
					}

					return *val 
				}

				title := media.Title
				description := valOrEmptyString(media.Description)

				form := huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Title").
							Value(&title),
						huh.NewText().
							Title("Description").
							Value(&description),
					),
				)

				err := form.Run()
				if err != nil {
					logger.Fatal("failed to run form", "err", err)
				}
			case "change-images":
				tempDir, err := os.MkdirTemp("", "watchbook-cli-*")
				if err != nil {
					logger.Fatal("failed to create temp dir", "err", err)
				}
				defer os.RemoveAll(tempDir)

				var coverUrl string
				var logoUrl string
				var bannerUrl string

				form := huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Cover URL").
							Value(&coverUrl),
						huh.NewInput().
							Title("Logo URL").
							Value(&logoUrl),
						huh.NewInput().
							Title("Banner URL").
							Value(&bannerUrl),
					),
				)

				err = form.Run()
				if err != nil {
					logger.Fatal("failed to run form", "err", err)
				}

				var coverPath string
				var logoPath string
				var bannerPath string

				if coverUrl != "" {
					coverPath, err = utils.DownloadImage(coverUrl, tempDir, "cover")
					if err != nil {
						logger.Fatal("failed to download cover image", "err", err)
					}
				}

				if logoUrl != "" {
					logoPath, err = utils.DownloadImage(logoUrl, tempDir, "logo")
					if err != nil {
						logger.Fatal("failed to download logo image", "err", err)
					}
				}

				if bannerUrl != "" {
					bannerPath, err = utils.DownloadImage(bannerUrl, tempDir, "banner")
					if err != nil {
						logger.Fatal("failed to download banner image", "err", err)
					}
				}

				images, err := createImageForm(coverPath, logoPath, bannerPath)
				if err != nil {
					logger.Fatal("failed to create image form", "err", err, "id", media.Id, "title", media.Title)
				}

				_, err = client.ChangeMediaImages(media.Id, images.Boundary, &images.Buf, api.Options{})
				if err != nil {
					logger.Fatal("failed to set images", "err", err, "id", media.Id, "title", media.Title)
				}
			}
		}
	},
}

func init() {
	mediaCmd.AddCommand(mediaEditCmd)

	rootCmd.AddCommand(mediaCmd)
}
