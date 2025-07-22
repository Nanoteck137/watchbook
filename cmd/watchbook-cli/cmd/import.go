package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/library"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use: "init",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat("media.toml")
		if err == nil {
			logger.Fatal("media.toml already exists")
			return
		}

		media := library.Media{
			Id:        utils.CreateMediaId(),
			MediaType: types.MediaTypeSeason,
			General: library.MediaGeneral{
				Title:        "Some title",
				Score:        0,
				Status:       types.MediaStatusUnknown,
				Rating:       types.MediaRatingUnknown,
				AiringSeason: "winter-2020",
				StartDate:    "2020-10-04",
				EndDate:      "2020-12-29",
				Studios:      []string{"some studio"},
				Tags:         []string{"some tag"},
			},
		}

		d, err := toml.Marshal(media)
		if err != nil {
			logger.Fatal("failed to marshal media", "err", err)
		}

		err = os.WriteFile("media.toml", d, 0644)
		if err != nil {
			logger.Fatal("failed to write media to file", "err", err)
		}
	},
}

var dialect = ember.SqliteDialect()

type AnimeDownloadType string

const (
	AnimeDownloadTypeNone    AnimeDownloadType = "none"
	AnimeDownloadTypeMal     AnimeDownloadType = "myanimelist"
	AnimeDownloadTypeAnilist AnimeDownloadType = "anilist"
)

type AnimeThemeSongType string

const (
	AnimeThemeSongTypeUnknown AnimeThemeSongType = "unknown"
	AnimeThemeSongTypeOpening AnimeThemeSongType = "opening"
	AnimeThemeSongTypeEnding  AnimeThemeSongType = "ending"
)

type AnimeType string

const (
	AnimeTypeUnknown   AnimeType = "unknown"
	AnimeTypeTV        AnimeType = "tv"
	AnimeTypeOVA       AnimeType = "original-video-anime"
	AnimeTypeMovie     AnimeType = "movie"
	AnimeTypeSpecial   AnimeType = "special"
	AnimeTypeONA       AnimeType = "original-network-anime"
	AnimeTypeMusic     AnimeType = "music"
	AnimeTypeCM        AnimeType = "commercial"
	AnimeTypePV        AnimeType = "promotional-video"
	AnimeTypeTVSpecial AnimeType = "tv-special"
)

type AnimeStatus string

const (
	AnimeStatusUnknown  AnimeStatus = "unknown"
	AnimeStatusAiring   AnimeStatus = "airing"
	AnimeStatusFinished AnimeStatus = "finished"
	AnimeStatusNotAired AnimeStatus = "not-aired"
)

type AnimeRating string

const (
	AnimeRatingUnknown     AnimeRating = "unknown"
	AnimeRatingAllAges     AnimeRating = "all-ages"
	AnimeRatingPG          AnimeRating = "pg"
	AnimeRatingPG13        AnimeRating = "pg-13"
	AnimeRatingR17         AnimeRating = "r-17"
	AnimeRatingRMildNudity AnimeRating = "r-mild-nudity"
	AnimeRatingRHentai     AnimeRating = "r-hentai"
)

type AnimeUserList string

const (
	AnimeUserListWatching    AnimeUserList = "watching"
	AnimeUserListCompleted   AnimeUserList = "completed"
	AnimeUserListOnHold      AnimeUserList = "on-hold"
	AnimeUserListDropped     AnimeUserList = "dropped"
	AnimeUserListPlanToWatch AnimeUserList = "plan-to-watch"
)

type AnimeStudio struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type AnimeTag struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type AnimeImageJson struct {
	Hash     string `json:"hash"`
	Filename string `json:"filename"`
	IsCover  int    `json:"is_cover"`
}

type Anime struct {
	RowId int `db:"rowid"`

	Id string `db:"id"`

	DownloadType       AnimeDownloadType `db:"download_type"`
	MalId              sql.NullString    `db:"mal_id"`
	AniDbId            sql.NullString    `db:"ani_db_id"`
	AnilistId          sql.NullString    `db:"anilist_id"`
	AnimeNewsNetworkId sql.NullString    `db:"anime_news_network_id"`

	Title        string         `db:"title"`
	TitleEnglish sql.NullString `db:"title_english"`

	Description sql.NullString `db:"description"`

	Type         AnimeType       `db:"type"`
	Score        sql.NullFloat64 `db:"score"`
	Status       AnimeStatus     `db:"status"`
	Rating       AnimeRating     `db:"rating"`
	EpisodeCount sql.NullInt64   `db:"episode_count"`
	AiringSeason sql.NullString  `db:"airing_season"`

	StartDate sql.NullString `db:"start_date"`
	EndDate   sql.NullString `db:"end_date"`

	CoverFilename sql.NullString `db:"cover_filename"`

	LastDataFetch sql.NullInt64 `db:"last_data_fetch"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	Studios database.JsonColumn[[]AnimeStudio]    `db:"studios"`
	Tags    database.JsonColumn[[]AnimeTag]       `db:"tags"`
	Images  database.JsonColumn[[]AnimeImageJson] `db:"images"`
}

func AnimeStudioQuery() *goqu.SelectDataset {
	tbl := goqu.T("anime_studios")

	return dialect.From(tbl).
		Select(
			tbl.Col("anime_id").As("id"),
			goqu.Func(
				"json_group_array",
				goqu.Func(
					"json_object",

					"slug",
					goqu.I("studios.slug"),
					"name",
					goqu.I("studios.name"),
				),
			).As("studios"),
		).
		Join(
			goqu.I("studios"),
			goqu.On(tbl.Col("studio_slug").Eq(goqu.I("studios.slug"))),
		).
		GroupBy(tbl.Col("anime_id"))
}

func AnimeAiringSeasonQuery() *goqu.SelectDataset {
	tbl := goqu.T("tags")

	return dialect.From(tbl).
		Select(
			tbl.Col("slug").As("slug"),
			goqu.Func(
				"json_group_array",
				goqu.Func(
					"json_object",

					"slug",
					tbl.Col("slug"),
					"name",
					tbl.Col("name"),
				),
			).As("data"),
		)
}

func AnimeTagQuery() *goqu.SelectDataset {
	tbl := goqu.T("anime_tags")

	return dialect.From(tbl).
		Select(
			tbl.Col("anime_id").As("id"),
			goqu.Func(
				"json_group_array",
				goqu.Func(
					"json_object",

					"slug",
					goqu.I("tags.slug"),
					"name",
					goqu.I("tags.name"),
				),
			).As("data"),
		).
		Join(
			goqu.I("tags"),
			goqu.On(tbl.Col("tag_slug").Eq(goqu.I("tags.slug"))),
		).
		GroupBy(tbl.Col("anime_id"))
}

func AnimeImageJsonQuery() *goqu.SelectDataset {
	tbl := goqu.T("anime_images")

	return dialect.From(tbl).
		Select(
			tbl.Col("anime_id").As("id"),
			goqu.Func(
				"json_group_array",
				goqu.Func(
					"json_object",

					"hash",
					goqu.I("anime_images.hash"),
					"filename",
					goqu.I("anime_images.filename"),
					"is_cover",
					goqu.I("anime_images.is_cover"),
				),
			).As("data"),
		).
		GroupBy(tbl.Col("anime_id"))
}

func AnimeUserDataQuery(userId *string) *goqu.SelectDataset {
	tbl := goqu.T("anime_user_data")

	query := dialect.From(tbl).
		Select(
			tbl.Col("anime_id").As("id"),

			tbl.Col("list"),
			tbl.Col("episode"),
			tbl.Col("rewatch_count"),
			tbl.Col("is_rewatching"),
			tbl.Col("score"),

			goqu.Func(
				"json_object",

				"list",
				tbl.Col("list"),

				"episode",
				tbl.Col("episode"),

				"rewatch_count",
				tbl.Col("rewatch_count"),

				"is_rewatching",
				tbl.Col("is_rewatching"),

				"score",
				tbl.Col("score"),
			).As("data"),
		)

	if userId != nil {
		query = query.Where(tbl.Col("user_id").Eq(*userId))
	} else {
		query = query.Where(goqu.L("false"))
	}

	return query
}

// TODO(patrik): Use goqu.T more
func AnimeQuery(userId *string) *goqu.SelectDataset {
	studiosQuery := AnimeStudioQuery()
	tagsQuery := AnimeTagQuery()
	imagesQuery := AnimeImageJsonQuery()
	// airingSeasonQuery := AnimeAiringSeasonQuery()

	query := dialect.From("animes").
		Select(
			"animes.rowid",

			"animes.id",

			"animes.download_type",
			"animes.mal_id",
			"animes.ani_db_id",
			"animes.anilist_id",
			"animes.anime_news_network_id",

			"animes.title",
			"animes.title_english",

			"animes.description",

			"animes.type",
			"animes.score",
			"animes.status",
			"animes.rating",
			"animes.episode_count",
			"animes.airing_season",

			"animes.start_date",
			"animes.end_date",

			"animes.last_data_fetch",

			"animes.created",
			"animes.updated",

			goqu.I("studios.studios").As("studios"),
			goqu.I("tags.data").As("tags"),
			goqu.I("images.data").As("images"),
		).
		LeftJoin(
			studiosQuery.As("studios"),
			goqu.On(goqu.I("animes.id").Eq(goqu.I("studios.id"))),
		).
		LeftJoin(
			tagsQuery.As("tags"),
			goqu.On(goqu.I("animes.id").Eq(goqu.I("tags.id"))),
		).
		LeftJoin(
			imagesQuery.As("images"),
			goqu.On(goqu.I("animes.id").Eq(goqu.I("images.id"))),
		)

	return query
}

var importCmd = &cobra.Command{
	Use: "import",
	Run: func(cmd *cobra.Command, args []string) {
		dir := "./work/library"
		_ = dir

		workDir := "/Volumes/old/watchbook"

		dbFile := path.Join(workDir, "data.db")
		dbUrl := fmt.Sprintf("file:%s?_foreign_keys=true", dbFile)
		db, err := ember.OpenDatabase("sqlite3", dbUrl)
		if err != nil {
			logger.Fatal("failed to open old database", "err", err)
		}

		ctx := context.Background()

		animeQuery := AnimeQuery(nil)
		animes, err := ember.Multiple[Anime](db, ctx, animeQuery)
		if err != nil {
			logger.Fatal("failed to get all animes", "err", err)
		}

		for _, anime := range animes {
			title := anime.Title
			if anime.TitleEnglish.Valid {
				title = anime.TitleEnglish.String
			}

			tags := make([]string, 0, len(anime.Tags.Data))
			for _, tag := range anime.Tags.Data {
				tags = append(tags, tag.Name)
			}

			studios := make([]string, 0, len(anime.Studios.Data))
			for _, studio := range anime.Studios.Data {
				studios = append(studios, studio.Name)
			}

			status := types.MediaStatusUnknown
			switch anime.Status {
			case AnimeStatusAiring:
				status = types.MediaStatusAiring
			case AnimeStatusFinished:
				status = types.MediaStatusFinished
			case AnimeStatusNotAired:
				status = types.MediaStatusNotAired
			}

			rating := types.MediaRatingUnknown
			switch anime.Rating {
			case AnimeRatingAllAges:
				rating = types.MediaRatingAllAges
			case AnimeRatingPG:
				rating = types.MediaRatingPG
			case AnimeRatingPG13:
				rating = types.MediaRatingPG13
			case AnimeRatingR17:
				rating = types.MediaRatingR17
			case AnimeRatingRMildNudity:
				rating = types.MediaRatingRMildNudity
			case AnimeRatingRHentai:
				rating = types.MediaRatingRHentai
			}

			mediaType := types.MediaTypeUnknown
			switch anime.Type {
			case AnimeTypeTV:
				mediaType = types.MediaTypeAnimeSeason
			case AnimeTypeOVA:
			case AnimeTypeMovie:
				mediaType = types.MediaTypeAnimeMovie
			case AnimeTypeSpecial:
			case AnimeTypeONA:
			case AnimeTypeMusic:
			case AnimeTypeCM:
			case AnimeTypePV:
			case AnimeTypeTVSpecial:
				mediaType = types.MediaTypeAnimeSeason
			}

			media := library.Media{
				Id:        utils.CreateMediaId(),
				MediaType: mediaType,
				Ids: library.MediaIds{
					TheMovieDB:  "",
					MyAnimeList: anime.MalId.String,
					Anilist:     "",
				},
				Images: library.MediaImages{},
				General: library.MediaGeneral{
					Title:        title,
					Score:        anime.Score.Float64,
					Status:       status,
					Rating:       rating,
					AiringSeason: anime.AiringSeason.String,
					StartDate:    anime.StartDate.String,
					EndDate:      anime.EndDate.String,
					Studios:      studios,
					Tags:         tags,
				},
				Path: "",
			}

			for i := range anime.EpisodeCount.Int64 {
				media.Parts = append(media.Parts, library.MediaPart{
					Name: fmt.Sprintf("Episode %d", i + 1),
				})
			}

			cover := ""

			for _, images := range anime.Images.Data {
				if images.IsCover > 0 {
					cover = images.Filename
				}
			}

			p := path.Join(workDir, "animes", "images", anime.Id, cover)
			ext := path.Ext(p)

			out := path.Join(dir, anime.MalId.String+"-"+utils.Slug(media.General.Title))
			err = os.Mkdir(out, 0755)
			if err != nil {
				logger.Fatal("failed to create dir for anime", "err", err, "title", media.General.Title)
			}

			dst := path.Join(out, "cover"+ext)
			_, err = utils.CopyFile(p, dst)
			if err != nil {
				logger.Fatal("failed to copy cover for anime", "err", err, "title", media.General.Title, "coverPath", p, "destPath", dst)
			}

			media.Images.Cover = "cover" + ext

			d, err := toml.Marshal(media)
			if err != nil {
				logger.Fatal("failed to marshal media", "err", err, "title", media.General.Title)
			}

			dst = path.Join(out, "media.toml")
			err = os.WriteFile(dst, d, 0644)
			if err != nil {
				logger.Fatal("failed to write media for anime", "err", err, "title", media.General.Title, "dstPath", dst)
			}

			// break
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(importCmd)
}
