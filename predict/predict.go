package predict

import (
	"context"
	"time"

	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/types"
)

func UpdateRelease(db *database.Database, m database.MediaPartRelease, t time.Time) error {
	// next, err := time.Parse(time.RFC3339, m.NextAiring)
	// if err != nil {
	// 	return err
	// }
	//
	// if t.Before(next) {
	// 	if m.Status == types.MediaPartReleaseStatusUnknown {
	// 		err := db.UpdateMediaPartRelease(context.Background(), m.MediaId, database.MediaPartReleaseChanges{
	// 			Status: database.Change[types.MediaPartReleaseStatus]{
	// 				Value:   types.MediaPartReleaseStatusWaiting,
	// 				Changed: true,
	// 			},
	// 		})
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }
	//
	// if t.After(next) {
	// 	newCurrentPart := m.CurrentPart + 1
	// 	newNext := next.Add(time.Duration(m.IntervalDays) * 24 * time.Hour)
	//
	// 	newStatus := types.MediaPartReleaseStatusRunning
	//
	// 	if m.NumExpectedParts > 0 {
	// 		boolToInt := func(val bool) int {
	// 			if val {
	// 				return 1
	// 			}
	//
	// 			return 0
	// 		}
	//
	// 		if newCurrentPart >= m.NumExpectedParts {
	// 			newStatus = types.MediaPartReleaseStatusCompleted
	// 			newNext = time.Time{}
	// 			newCurrentPart = m.NumExpectedParts
	// 		}
	//
	// 		err := db.UpdateMediaPartRelease(context.Background(), m.MediaId, database.MediaPartReleaseChanges{
	// 			Status: database.Change[types.MediaPartReleaseStatus]{
	// 				Value:   newStatus,
	// 				Changed: true,
	// 			},
	// 			CurrentPart: database.Change[int]{
	// 				Value:   newCurrentPart,
	// 				Changed: true,
	// 			},
	// 			// NextAiring: database.Change[string]{
	// 			// 	Value:   newNext.Format(time.RFC3339),
	// 			// 	Changed: true,
	// 			// },
	// 			// IsActive: database.Change[int]{
	// 			// 	Value:   boolToInt(newCurrentPart < m.NumExpectedParts),
	// 			// 	Changed: true,
	// 			// },
	// 			// Created:      database.Change[int64]{},
	// 		})
	// 		if err != nil {
	// 			return err
	// 		}
	// 	} else {
	//
	// 		newStatus := types.MediaPartReleaseStatusRunning
	//
	// 		err := db.UpdateMediaPartRelease(context.Background(), m.MediaId, database.MediaPartReleaseChanges{
	// 			Status: database.Change[types.MediaPartReleaseStatus]{
	// 				Value:   newStatus,
	// 				Changed: true,
	// 			},
	// 			CurrentPart: database.Change[int]{
	// 				Value:   newCurrentPart,
	// 				Changed: true,
	// 			},
	// 			// NextAiring: database.Change[string]{
	// 			// 	Value:   newNext.Format(types.MediaDateLayout),
	// 			// 	Changed: true,
	// 			// },
	// 		})
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }
	//
	// return nil

	return nil
}

func RunPredict(ctx context.Context, app core.App) error {
	media, err := app.DB().GetAllMediaPartReleases(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, m := range media {
		if m.Status == types.MediaPartReleaseStatusCompleted {
			continue
		}

		err := UpdateRelease(app.DB(), m, now)
		if err != nil {
			return err
		}
	}

	return nil
}
