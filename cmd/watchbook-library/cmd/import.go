package cmd

import (
	"os"

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

func init() {
	rootCmd.AddCommand(initCmd)
}
