package cmd

import (
	"context"
	"time"

	"github.com/nanoteck137/watchbook/apis"
	"github.com/nanoteck137/watchbook/config"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/spf13/cobra"
)

func test(app core.App) error {
	testId := "cka7s522p9zt"

	{
		// t := time.Now()
		// t = t.AddDate(0, 0, 0)

		// t, _ := time.Parse(time.RFC3339, "2025-07-03T16:00:00Z")
		// t, _ := time.Parse(time.RFC3339, "2025-09-03T2:00:00Z")
		t, _ := time.Parse(time.RFC3339, "2025-09-03T00:00:00+02:00")

		err := app.DB().SetMediaPartRelease(context.Background(), testId, database.SetMediaPartRelease{
			StartDate:        t.Format(time.RFC3339),
			NumExpectedParts: 12,
			PartOffset:       11,
			IntervalDays:     7,
			DelayDays:        0,
		})
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
