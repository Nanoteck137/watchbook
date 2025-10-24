package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/kr/pretty"
	"github.com/nanoteck137/watchbook/cmd/watchbook-cli/api"
	"github.com/nanoteck137/watchbook/utils"
	"github.com/spf13/cobra"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
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

type ProviderValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CollectionItem struct {
	MediaId string `json:"mediaId"`
	Name    string `json:"name"`
}

type Collection struct {
	Id              string           `json:"id"`
	Type            string           `json:"type"`
	Name            string           `json:"name"`
	CoverFile       string           `json:"coverFile"`
	LogoFile        string           `json:"logoFile"`
	BannerFile      string           `json:"bannerFile"`
	DefaultProvider string           `json:"defaultProvider"`
	Providers       []ProviderValue  `json:"providers"`
	Items           []CollectionItem `json:"items"`
}

type ExportData struct {
	Collections []Collection `json:"collections"`
}

var exportCollectionsCmd = &cobra.Command{
	Use: "export-collections",
	Run: func(cmd *cobra.Command, args []string) {
		apiAddress, _ := cmd.Flags().GetString("api-address")
		output, _ := cmd.Flags().GetString("output")

		client := api.New(apiAddress)

		res, err := client.GetCollections(api.Options{
			Query: url.Values{
				"perPage": {"1000000"},
			},
		})
		if err != nil {
			logger.Fatal("failed", "err", err)
		}

		exportData := ExportData{}

		downloadImage := func(url, name, outDir string) (string, error) {
			resp, err := http.Get(url)
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()

			ext := path.Ext(url)
			n := name + ext
			out := path.Join(outDir, n)

			f, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return "", err
			}
			defer f.Close()

			_, err = io.Copy(f, resp.Body)
			if err != nil {
				return "", err
			}

			return n, nil
		}

		for _, collection := range res.Collections {
			dir := path.Join(output, collection.Id)

			err := os.Mkdir(dir, 0755)
			if err != nil && !os.IsExist(err) {
				logger.Fatal("failed", "err", err)
			}

			var coverFile string
			var bannerFile string
			var logoFile string

			if collection.CoverUrl != nil {
				out, err := downloadImage(*collection.CoverUrl, "cover", dir)
				if err != nil {
					logger.Fatal("failed", "err", err)
				}

				coverFile = out
			}

			if collection.BannerUrl != nil {
				out, err := downloadImage(*collection.BannerUrl, "banner", dir)
				if err != nil {
					logger.Fatal("failed", "err", err)
				}

				bannerFile = out
			}

			if collection.LogoUrl != nil {
				out, err := downloadImage(*collection.LogoUrl, "logo", dir)
				if err != nil {
					logger.Fatal("failed", "err", err)
				}

				logoFile = out
			}

			items, err := client.GetCollectionItems(collection.Id, api.Options{})
			if err != nil {
				logger.Fatal("failed", "err", err)
			}

			defaultProvider := ""
			if collection.DefaultProvider != nil {
				defaultProvider = *collection.DefaultProvider
			}

			col := Collection{
				Id:              collection.Id,
				Type:            collection.Type,
				Name:            collection.Name,
				CoverFile:       coverFile,
				LogoFile:        logoFile,
				BannerFile:      bannerFile,
				DefaultProvider: defaultProvider,
				Providers:       []ProviderValue{},
				Items:           []CollectionItem{},
			}

			for _, val := range collection.Providers {
				col.Providers = append(col.Providers, ProviderValue{
					Name:  val.Name,
					Value: val.Value,
				})
			}

			for _, item := range items.Items {
				col.Items = append(col.Items, CollectionItem{
					MediaId: item.MediaId,
					Name:    item.CollectionName,
				})
			}

			exportData.Collections = append(exportData.Collections, col)
		}

		d, err := json.MarshalIndent(exportData, "", "  ")
		if err != nil {
			logger.Fatal("failed", "err", err)
		}

		err = os.WriteFile(path.Join(output, "export.json"), d, 0644)
		if err != nil {
			logger.Fatal("failed", "err", err)
		}
	},
}

var importShowCmd = &cobra.Command{
	Use:  "import-show <DIR>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]

		apiAddress, _ := cmd.Flags().GetString("api-address")
		authToken, _ := cmd.Flags().GetString("auth-token")

		client := api.New(apiAddress)
		client.Headers.Add("X-Api-Token", authToken)

		d, err := os.ReadFile(path.Join(dir, "export.json"))
		if err != nil {
			logger.Fatal("failed", "err", err)
		}

		var exportData ExportData
		err = json.Unmarshal(d, &exportData)
		if err != nil {
			logger.Fatal("failed", "err", err)
		}

		pretty.Println(exportData)

		for _, col := range exportData.Collections {
			ty := "unknown"
			switch col.Type {
			case "anime":
				ty = "anime"
			case "series":
				ty = "tv-series"
			}

			show, err := client.CreateShow(api.CreateShowBody{
				Type: ty,
				Name: col.Name,
			}, api.Options{})
			if err != nil {
				logger.Fatal("failed", "err", err)
			}

			pretty.Println(show)

			var coverFile string
			var bannerFile string
			var logoFile string

			if col.CoverFile != "" {
				coverFile = path.Join(dir, col.Id, col.CoverFile)
			}

			if col.BannerFile != "" {
				bannerFile = path.Join(dir, col.Id, col.BannerFile)
			}

			if col.LogoFile != "" {
				logoFile = path.Join(dir, col.Id, col.LogoFile)
			}

			res, err := createImageForm(coverFile, logoFile, bannerFile)
			if err != nil {
				logger.Fatal("failed", "err", err)
			}

			_, err = client.ChangeShowImages(show.Id, res.Boundary, &res.Buf, api.Options{})
			if err != nil {
				logger.Fatal("failed", "err", err)
			}

		}
	},
}

func init() {
	exportCollectionsCmd.Flags().StringP("output", "o", "", "output directory")
	exportCollectionsCmd.MarkFlagRequired("output")

	importShowCmd.Flags().StringP("auth-token", "t", "", "auth token")
	importShowCmd.MarkFlagRequired("auth-token")

	rootCmd.AddCommand(exportCollectionsCmd, importShowCmd)
}
