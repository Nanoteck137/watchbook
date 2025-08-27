package predict

import (
	"context"
	"time"

	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
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

		if now.After(next) {
			newCurrentPart := m.CurrentPart + 1
			newNext := next.Add(time.Duration(m.IntervalDays) * 24 * time.Hour)

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
		}
	}

	return nil
}
