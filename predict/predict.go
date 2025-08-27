package predict

import (
	"context"
	"fmt"
	"time"

	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/kvstore"
	"github.com/nanoteck137/watchbook/types"
)

func RunPredict(ctx context.Context, app core.App) error {
	media, err := app.DB().GetAllMediaPartReleases(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, m := range media {
		next, err := time.Parse(types.MediaDateLayout, m.NextAiring)
		if err != nil {
			return err
		}

		app.Logger().Info("Next Release", "release", next)

		if now.After(next) {
			newCurrentPart := m.CurrentPart + 1
			newNext := next.Add(time.Duration(m.IntervalDays) * 24 * time.Hour)

			app.Logger().Info("Hello World", "part", newCurrentPart, "next-airing", newNext)

			if m.NumExpectedParts > 0 {
				boolToInt := func(val bool) int {
					if val {
						return 1
					}

					return 0
				}

				err := app.DB().UpdateMediaPartRelease(context.Background(), m.MediaId, database.MediaPartReleaseChanges{
					// NumExpectedParts: database.Change[int]{},
					CurrentPart: database.Change[int]{
						Value:   newCurrentPart,
						Changed: true,
					},
					NextAiring: database.Change[string]{
						Value:   newNext.Format(types.MediaDateLayout),
						Changed: true,
					},
					// IntervalDays: database.Change[int]{},
					IsActive: database.Change[int]{
						Value:   boolToInt(newCurrentPart < m.NumExpectedParts),
						Changed: true,
					},
					// Created:      database.Change[int64]{},
				})
				if err != nil {
					return err
				}

				fullMedia, err := app.DB().GetMediaById(ctx, nil, m.MediaId)
				if err != nil {
					return err
				}

				users, err := app.DB().GetAllUsers(ctx)
				if err != nil {
					return err
				}

				for _, user := range users {
					_, err := app.DB().CreateNotification(ctx, database.CreateNotificationParams{
						UserId:  user.Id,
						Type:    types.NotificationTypePartRelease,
						Title:   fmt.Sprintf("New Release of '%s' (Part %d)", fullMedia.Title, newCurrentPart),
						Message: "",
						Metadata: kvstore.Store{
							"mediaId":    fullMedia.Id,
							"mediaTitle": fullMedia.Title,
						},
						DedupKey: fmt.Sprintf("part-release-%s-%d", m.MediaId, newCurrentPart),
					})
					if err != nil {
						return err
					}
				}

			} else {
				err := app.DB().UpdateMediaPartRelease(context.Background(), m.MediaId, database.MediaPartReleaseChanges{
					// NumExpectedParts: database.Change[int]{},
					CurrentPart: database.Change[int]{
						Value:   newCurrentPart,
						Changed: true,
					},
					NextAiring: database.Change[string]{
						Value:   newNext.Format(types.MediaDateLayout),
						Changed: true,
					},
					// IntervalDays: database.Change[int]{},
					// IsActive: database.Change[int]{
					// 	Value:   boolToInt(newCurrentPart < m.NumExpectedParts),
					// 	Changed: true,
					// },
					// Created:      database.Change[int64]{},
				})
				if err != nil {
					return err
				}
			}

			//
			// // Notify users who haven't been notified yet
			// bookmarks := getBookmarksForShow(db, show.ID)
			// for _, bm := range bookmarks {
			// 	if bm.LastNotifiedEpisode < newEpisode {
			// 		sendNotification(bm.UserID, show.Title, newEpisode)
			// 		updateBookmarkLastNotified(db, bm.UserID, show.ID, newEpisode)
			// 	}
			// }
		}
	}

	return nil
}
