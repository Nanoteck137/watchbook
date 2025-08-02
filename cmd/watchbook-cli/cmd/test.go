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
	"strings"

	"github.com/kr/pretty"
	"github.com/nanoteck137/watchbook/cmd/watchbook-cli/api"
	"github.com/nanoteck137/watchbook/library"
	"github.com/nanoteck137/watchbook/utils"
	"github.com/spf13/cobra"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

var testCmd = &cobra.Command{
	Use: "test",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := library.ReadMedia("/Volumes/media/watch/mal/entries/11061-hunter-x-hunter")
		if err != nil {
			logger.Fatal("failed to read media", "err", err)
		}

		pretty.Println(m)

		client := api.New("http://localhost:3000")

		lel, err := client.GetMedia(api.Options{
			Query: url.Values{
				"filter": {fmt.Sprintf(`malId=="anime@%s"`, m.Ids.MyAnimeList)},
			},
		})
		if err != nil {
			logger.Fatal("failed to get media", "err", err)
		}

		pretty.Println(lel)

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
				Status:       string(m.General.Status),
				Rating:       string(m.General.Rating),
				AiringSeason: m.General.AiringSeason,
				StartDate:    m.General.StartDate,
				EndDate:      m.General.EndDate,
				PartCount:    len(m.Parts),
				// CoverUrl:       "",
				// BannerUrl:      "",
				// LogoUrl:        "",
				Tags:    m.General.Tags,
				Studios: m.General.Studios,
				// CollectionId:   "",
				// CollectionName: "",
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

		createFormPart := func(p, name string) {
			f, err := os.Open(p)
			if err != nil {
				logger.Fatal("failed", "err", err)
			}
			defer f.Close()

			contentType, err := utils.ImageExtToContentType(path.Ext(p))
			if err != nil {
				logger.Fatal("failed", "err", err)
			}

			formFile, err := createFileField(name, path.Base(p), contentType)
			if err != nil {
				logger.Fatal("failed", "err", err)
			}

			_, err = io.Copy(formFile, f)
			if err != nil {
				logger.Fatal("failed", "err", err)
			}

		}

		_ = createFormPart

		coverPath := m.GetCoverPath()
		if coverPath != "" {
			createFormPart(coverPath, "cover")
		}

		logoPath := m.GetLogoPath()
		if logoPath != "" {
			createFormPart(logoPath, "logo")
		}

		bannerPath := m.GetLogoPath()
		if bannerPath != "" {
			createFormPart(bannerPath, "banner")
		}

		w.Close()

		_, err = client.ChangeImages(id, w.Boundary(), &b, api.Options{})
		if err != nil {
			logger.Fatal("failed", "err", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
