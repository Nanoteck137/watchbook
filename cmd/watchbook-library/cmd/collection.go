package cmd

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/nanoteck137/watchbook/library"
	"github.com/nanoteck137/watchbook/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var collectionCmd = &cobra.Command{
	Use: "collection",
}

var collectionInitCmd = &cobra.Command{
	Use: "init",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")

		entries, err := os.ReadDir(dir)
		if err != nil {
			logger.Fatal("failed to read dir", "err", err)
		}

		var dirs []string

		for _, entry := range entries {
			name := entry.Name()
			if strings.HasPrefix(name, ".") {
				continue
			}

			if entry.IsDir() {
				dirs = append(dirs, name)
			}
		}

		name := ""
		if a, err := filepath.Abs(dir); err == nil {
			name = path.Base(a)
		}

		col := library.Collection{
			Id: utils.CreateCollectionId(),
			General: library.CollectionGeneral{
				Name: name,
			},
			Images:  library.Images{},
			Entries: []library.CollectionEntry{},
		}

		for _, dir := range dirs {
			col.Entries = append(col.Entries, library.CollectionEntry{
				Path:       dir,
				Name:       "",
				SearchSlug: "",
				Order:      0,
				SubOrder:   0,
			})
		}

		d, err := toml.Marshal(col)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(path.Join(dir, "collection.toml"), d, 0644)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var collectionImagesCmd = &cobra.Command{
	Use: "images",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")

		d, err := os.ReadFile(path.Join(dir, "collection.toml"))
		if err != nil {
			log.Fatal(err)
		}

		var col library.Collection
		err = toml.Unmarshal(d, &col)
		if err != nil {
			log.Fatal(err)
		}

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
			log.Fatal(err)
		}

		if coverUrl != "" {
			p, err := downloadImage(coverUrl, dir, "cover")
			if err != nil {
				logger.Fatal("failed to donwload cover image", "err", err)
			}

			col.Images.Cover = path.Base(p)
		}

		if logoUrl != "" {
			p, err := downloadImage(logoUrl, dir, "logo")
			if err != nil {
				logger.Fatal("failed to donwload logo image", "err", err)
			}

			col.Images.Logo = path.Base(p)
		}

		if bannerUrl != "" {
			p, err := downloadImage(bannerUrl, dir, "banner")
			if err != nil {
				logger.Fatal("failed to donwload banner image", "err", err)
			}

			col.Images.Banner = path.Base(p)
		}

		d, err = toml.Marshal(col)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(path.Join(dir, "collection.toml"), d, 0644)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	collectionInitCmd.Flags().StringP("dir", "d", ".", "")
	collectionImagesCmd.Flags().StringP("dir", "d", ".", "")

	collectionCmd.AddCommand(collectionInitCmd)
	collectionCmd.AddCommand(collectionImagesCmd)

	rootCmd.AddCommand(collectionCmd)
}
