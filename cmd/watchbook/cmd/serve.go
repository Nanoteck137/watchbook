package cmd

import (
	"context"
	"time"

	"github.com/nanoteck137/watchbook/apis"
	"github.com/nanoteck137/watchbook/config"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/predict"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
)

func test(app core.App) error {
	testId := "cka7s522p9zt"

	{
		// t := time.Now()
		// t = t.AddDate(0, 0, 0)

		// t, _ := time.Parse(time.RFC3339, "2025-07-03T16:00:00Z")
		t, _ := time.Parse(time.RFC3339, "2025-08-30T16:00:00Z")

		app.DB().RemoveMediaPartRelease(context.Background(), testId)

		err := app.DB().CreateMediaPartRelease(context.Background(), database.CreateMediaPartReleaseParams{
			MediaId:          testId,
			StartDate:        t,
			NumExpectedParts: 12,
			CurrentPart:      9,
			NextAiring:       t,
			IntervalDays:     7,
			DelayDays:        7,
		})
		if err != nil {
			return err
		}

		m, err := app.DB().GetMediaPartReleaseById(context.Background(), testId)
		if err != nil {
			return err
		}

		err = predict.UpdateRelease(app.DB(), m, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		app := core.NewBaseApp(&config.LoadedConfig)

		err := app.Bootstrap()
		if err != nil {
			app.Logger().Fatal("Failed to bootstrap app", "err", err)
		}

		err = test(app)
		if err != nil {
			app.Logger().Fatal("failed to run test", "err", err)
		}

		crun := cron.New()

		// @hourly
		crun.AddFunc("*/30 * * * *", func() {
			app.Logger().Info("running part prediction")

			err := predict.RunPredict(context.Background(), app)
			if err != nil {
				app.Logger().Error("prediction failed", "err", err)
			}
		})

		crun.Start()

		e, err := apis.Server(app)
		if err != nil {
			app.Logger().Fatal("Failed to create server", "err", err)
		}

		err = e.Start(app.Config().ListenAddr)
		if err != nil {
			app.Logger().Fatal("Failed to start server", "err", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
