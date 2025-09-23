package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/nanoteck137/pyrin/ember"
	old "github.com/nanoteck137/watchbook/cmd/watchbook/database"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/provider/myanimelist"
	"github.com/nanoteck137/watchbook/provider/tmdb"
	"github.com/nanoteck137/watchbook/utils"
	"github.com/spf13/cobra"
)

var oldCmd = &cobra.Command{
	Use:  "old <OLD_DB_FILE> <NEW_DB_FILE>",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		oldDbFile := args[0]
		newDbFile := args[1]

		db, err := database.Open(newDbFile)
		if err != nil {
			logger.Fatal("failed to open new database", "err", err)
		}

		err = db.RunMigrateUp()
		if err != nil {
			logger.Fatal("failed to open run migrate up on new database", "err", err)
		}

		oldDb, err := old.Open(oldDbFile)
		if err != nil {
			logger.Fatal("failed to open old database", "err", err)
		}

		_ = oldDb

		ctx := context.Background()

		users, err := oldDb.GetAllUsers(ctx)
		if err != nil {
			logger.Fatal("failed", "err", err)
		}

		for _, user := range users {
			newUserId := utils.CreateUserId()

			_, err := db.CreateUser(ctx, database.CreateUserParams{
				Id:       newUserId,
				Username: user.Username,
				Password: user.Password,
				Role:     user.Role,
				Created:  user.Created,
				Updated:  user.Updated,
			})
			if err != nil {
				logger.Fatal("failed", "err", err)
			}

			err = db.UpdateUserSettings(ctx, database.UserSettings{
				Id:          newUserId,
				DisplayName: user.DisplayName,
			})
			if err != nil {
				logger.Fatal("failed", "err", err)
			}

			tokens, err := oldDb.GetAllApiTokensForUser(ctx, user.Id)
			if err != nil {
				logger.Fatal("failed", "err", err)
			}

			for _, token := range tokens {
				_, err := db.CreateApiToken(ctx, database.CreateApiTokenParams{
					Id:      token.Id,
					UserId:  newUserId,
					Name:    token.Name,
					Created: token.Created,
					Updated: token.Updated,
				})
				if err != nil {
					logger.Fatal("failed", "err", err)
				}
			}

			tags, err := oldDb.GetAllTags(ctx)
			if err != nil {
				logger.Fatal("failed", "err", err)
			}

			for _, tag := range tags {
				err := db.CreateTag(ctx, tag.Slug, tag.Name)
				if err != nil {
					logger.Fatal("failed", "err", err)
				}
			}

			media, err := oldDb.GetAllMedia(ctx)
			if err != nil {
				logger.Fatal("failed", "err", err)
			}

			for _, media := range media {
				providers := ember.KVStore{}
				providerUsed := ""

				if media.TmdbId.Valid {
					splits := strings.Split(media.TmdbId.String, "@")
					if len(splits) == 2 {
						ty := splits[0]
						id := splits[1]

						switch ty {
						case "movie":
							providers[tmdb.MovieProviderName] = id
							providerUsed = tmdb.MovieProviderName
						case "tv":
							providers[tmdb.TvProviderName] = id
							providerUsed = tmdb.TvProviderName
						default:
							fmt.Println("WARNING: Invalid TMDB-ID TYPE", media.TmdbId.String)
						}
					} else {
						fmt.Println("WARNING: Invalid TMDB-ID", media.TmdbId.String)
					}

					// providers[tmdb.
				}

				if media.MalId.Valid {
					splits := strings.Split(media.MalId.String, "@")
					if len(splits) == 2 {
						ty := splits[0]
						id := splits[1]

						switch ty {
						case "anime":
							providers[myanimelist.AnimeProviderName] = id
							providerUsed = myanimelist.AnimeProviderName
						default:
							fmt.Println("WARNING: Invalid MAL-ID", media.MalId.String)
						}
					} else {
						fmt.Println("WARNING: Invalid MAL-ID", media.MalId.String)
					}
				}

				_, err := db.CreateMedia(ctx, database.CreateMediaParams{
					Id:           media.Id,
					Type:         media.Type,
					Title:        media.Title,
					Description:  media.Description,
					Score:        media.Score,
					Status:       media.Status,
					Rating:       media.Rating,
					AiringSeason: media.AiringSeason,
					StartDate:    media.StartDate,
					EndDate:      media.EndDate,
					CoverFile:    media.CoverFile,
					LogoFile:     media.LogoFile,
					BannerFile:   media.BannerFile,
					DefaultProvider: sql.NullString{
						String: providerUsed,
						Valid:  providerUsed != "",
					},
					Providers:    providers,
					Created:      media.Created,
					Updated:      media.Updated,
				})
				if err != nil {
					logger.Fatal("failed", "err", err)
				}

				for _, tag := range media.Tags.Data {
					err := db.AddTagToMedia(ctx, media.Id, tag)
					if err != nil {
						logger.Fatal("failed", "err", err)
					}
				}

				for _, tag := range media.Creators.Data {
					err := db.AddCreatorToMedia(ctx, media.Id, tag)
					if err != nil {
						logger.Fatal("failed", "err", err)
					}
				}

			}
		}

		parts, err := oldDb.GetAllMediaParts(ctx)
		if err != nil {
			logger.Fatal("failed", "err", err)
		}

		for _, part := range parts {
			err := db.CreateMediaPart(ctx, database.CreateMediaPartParams{
				Index:   part.Index,
				MediaId: part.MediaId,
				Name:    part.Name,
				Created: part.Created,
				Updated: part.Updated,
			})
			if err != nil {
				logger.Fatal("failed", "err", err)
			}
		}

		cols, err := oldDb.GetAllCollections(ctx)
		if err != nil {
			logger.Fatal("failed", "err", err)
		}

		for _, col := range cols {
			_, err := db.CreateCollection(ctx, database.CreateCollectionParams{
				Id:         col.Id,
				Type:       col.Type,
				Name:       col.Name,
				CoverFile:  col.CoverFile,
				LogoFile:   col.LogoFile,
				BannerFile: col.BannerFile,
				Created:    col.Created,
				Updated:    col.Updated,
			})
			if err != nil {
				logger.Fatal("failed", "err", err)
			}
		}

		items, err := oldDb.GetAllCollectionMediaItems(ctx)
		if err != nil {
			logger.Fatal("failed", "err", err)
		}

		for _, item := range items {
			err := db.CreateCollectionMediaItem(ctx, database.CreateCollectionMediaItemParams{
				CollectionId: item.CollectionId,
				MediaId:      item.MediaId,
				Name:         item.Name,
				SearchSlug:   item.SearchSlug,
				Position:     int(item.OrderNumber),
				Created:      item.Created,
				Updated:      item.Updated,
			})
			if err != nil {
				logger.Fatal("failed", "err", err)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(oldCmd)
}
