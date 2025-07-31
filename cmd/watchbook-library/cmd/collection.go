package cmd

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/nanoteck137/watchbook/cmd/watchbook-library/ui/colui"
	"github.com/nanoteck137/watchbook/library"
	"github.com/nanoteck137/watchbook/utils"
	"github.com/spf13/cobra"
)

// TODO(patrik):
//  - We need to save after every edit

var collectionCmd = &cobra.Command{
	Use: "collection",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")

		data, err := library.ReadCollection(dir)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				name := ""
				if a, err := filepath.Abs(dir); err == nil {
					name = path.Base(a)
				}

				res, err := colui.FormNewCollection(name)
				if err != nil {
					logger.Fatal("failed to run new collection form", "err", err)
				}

				data = library.Collection{
					Type: res.Type,
					Id:   utils.CreateCollectionId(),
					General: library.CollectionGeneral{
						Name: name,
					},
					Images: library.Images{},
					Groups: []library.CollectionGroup{
						{
							Name:    "Default",
							Order:   0,
							Entries: []library.CollectionEntry{},
						},
					},
					Path: dir,
				}

				err = library.WriteCollection(dir, data)
				if err != nil {
					logger.Fatal("failed to write collection to disk", "err", err)
				}
			} else {
				logger.Fatal("failed to read collection", "err", err)
			}
		}

		err = colui.EditMainMenu(data, dir)
		if err != nil {
			logger.Fatal("failed to run main menu", "err", err)
		}
	},
}

func init() {
	collectionCmd.Flags().StringP("dir", "d", ".", "")

	rootCmd.AddCommand(collectionCmd)
}
