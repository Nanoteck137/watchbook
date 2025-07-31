package colui

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/nanoteck137/watchbook/utils"
)

func runCollectionImagesForm(col *collection) error {
	for {
		clearScreen()
		prettyPrintCollection(col)

		var selected string
		form := huh.NewSelect[string]().
			Title("Edit Collection Images").
			Options(
				huh.NewOption("Find images on disk", "fs-search"),
				huh.NewOption("Get images from URLs", "get-from-urls"),
				huh.NewOption("Back", "back"),
			).
			Value(&selected)

		err := form.Run()
		if err != nil {
			return err
		}

		quit := false

		switch selected {
		case "fs-search":
			var cover string
			var logo string
			var banner string

			filepath.WalkDir(col.Dir, func(p string, d fs.DirEntry, err error) error {
				name := d.Name()
				if strings.HasPrefix(name, ".") {
					return nil
				}

				if d.IsDir() {
					return nil
				}

				ext := path.Ext(name)
				nameWithoutExt := strings.TrimSuffix(name, ext)

				switch ext {
				case ".png", ".jpeg", ".jpg":
					switch nameWithoutExt {
					case "cover":
						cover = name
					case "logo":
						logo = name
					case "banner":
						banner = name
					}
				}

				return nil
			})

			if col.Data.Images.Cover == "" && cover != "" {
				col.Data.Images.Cover = cover
			}

			if col.Data.Images.Logo == "" && logo != "" {
				col.Data.Images.Logo = logo
			}

			if col.Data.Images.Banner == "" && banner != "" {
				col.Data.Images.Banner = banner
			}

			quit = true

		case "get-from-urls":
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
				).
					Title("Image URLs"),
			)
			err = form.Run()
			if err != nil {
				return err
			}

			if coverUrl != "" {
				if col.Data.Images.Cover != "" {
					err = os.Remove(path.Join(col.Dir, col.Data.Images.Cover))
					if err != nil {
						return fmt.Errorf("failed to remove old cover image: %w", err)
					}
				}

				p, err := utils.DownloadImage(coverUrl, col.Dir, "cover")
				if err != nil {
					return fmt.Errorf("failed to donwload cover image: %w", err)
				}

				col.Data.Images.Cover = path.Base(p)
			}

			if logoUrl != "" {
				if col.Data.Images.Logo != "" {
					err = os.Remove(path.Join(col.Dir, col.Data.Images.Logo))
					if err != nil {
						return fmt.Errorf("failed to remove old logo image: %w", err)
					}
				}

				p, err := utils.DownloadImage(logoUrl, col.Dir, "logo")
				if err != nil {
					return fmt.Errorf("failed to donwload logo image: %w", err)
				}

				col.Data.Images.Logo = path.Base(p)
			}

			if bannerUrl != "" {
				if col.Data.Images.Banner != "" {
					err = os.Remove(path.Join(col.Dir, col.Data.Images.Banner))
					if err != nil {
						return fmt.Errorf("failed to remove old banner image: %w", err)
					}
				}

				p, err := utils.DownloadImage(bannerUrl, col.Dir, "banner")
				if err != nil {
					return fmt.Errorf("failed to donwload banner image: %w", err)
				}

				col.Data.Images.Banner = path.Base(p)
			}

		case "back":
			quit = true
		}

		if quit {
			break
		}
	}

	return nil
}
