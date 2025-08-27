package cmd

import (
	"context"
	"time"

	"github.com/kr/pretty"
	"github.com/nanoteck137/watchbook/apis"
	"github.com/nanoteck137/watchbook/config"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/predict"
	"github.com/nanoteck137/watchbook/types"
	"github.com/spf13/cobra"
)

func test(app core.App) error {
	testId := "cka7s522p9zt"

	{
		t := time.Now().Truncate(24 * time.Hour)

		t = t.AddDate(0, 0, -1)
		// t.Round()

		app.DB().RemoveMediaPartRelease(context.Background(), testId)

		err := app.DB().CreateMediaPartRelease(context.Background(), database.CreateMediaPartReleaseParams{
			MediaId:          testId,
			NumExpectedParts: 12,
			CurrentPart:      11,
			NextAiring:       t.Format(types.MediaDateLayout),
			IntervalDays:     7,
			IsActive:         1,
		})
		if err != nil {
			return err
		}
	}

	{
		m, err := app.DB().GetMediaById(context.Background(), nil, testId)
		if err != nil {
			return err
		}

		pretty.Println(m)
	}

	err := predict.RunPredict(context.Background(), app)
	if err != nil {
		return err
	}

	{
		m, err := app.DB().GetMediaById(context.Background(), nil, testId)
		if err != nil {
			return err
		}

		pretty.Println(m)
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
