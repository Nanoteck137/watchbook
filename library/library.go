package library

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/nanoteck137/watchbook/types"
	"github.com/pelletier/go-toml/v2"
)

type MediaImages struct {
	Cover  string `toml:"cover"`
	Banner string `toml:"banner"`
	Logo   string `toml:"logo"`
}

type MediaGeneral struct {
	Title       string `toml:"title"`

	Score        float64           `toml:"score"`
	Status       types.MediaStatus `toml:"status"`
	Rating       types.MediaRating `toml:"rating"`
	AiringSeason string            `toml:"airingSeason"`

	StartDate string `toml:"startDate"`
	EndDate   string `toml:"endDate"`

	Studios []string `toml:"studios"`
	Tags    []string `toml:"tags"`
}

type MediaIds struct {
	TheMovieDB  string `toml:"theMovieDB"`
	MyAnimeList string `toml:"myAnimeList"`
	Anilist     string `toml:"anilist"`
}

type Media struct {
	Id        string          `toml:"id"`
	MediaType types.MediaType `toml:"mediaType"`

	Ids     MediaIds     `toml:"ids"`
	General MediaGeneral `toml:"general"`
	Images  MediaImages  `toml:"images"`

	Path string `toml:"-"`
}

type MediaSearch struct {
	Media  []Media
	Errors map[string]error
}

func readMedia(p string) (Media, error) {
	metadataPath := path.Join(p, "media.toml")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return Media{}, err
	}

	var media Media
	err = toml.Unmarshal(data, &media)
	if err != nil {
		return Media{}, err
	}

	// if metadata.General.Cover != "" {
	// 	metadata.General.Cover = path.Join(p, metadata.General.Cover)
	// }
	//
	// for i, t := range metadata.Tracks {
	// 	metadata.Tracks[i].File = path.Join(p, t.File)
	// }

	media.Path = p

	return media, nil
}

func FindMedia(p string) (*MediaSearch, error) {
	var media []string

	err := filepath.WalkDir(p, func(p string, d fs.DirEntry, err error) error {
		if d == nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		name := d.Name()

		if strings.HasPrefix(name, ".") {
			return nil
		}

		if name == "media.toml" {
			media = append(media, path.Dir(p))
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	errors := map[string]error{}
	res := make([]Media, 0, len(media))

	for _, p := range media {
		media, err := readMedia(p)
		if err != nil {
			errors[p] = err
			continue
		}

		res = append(res, media)
	}

	return &MediaSearch{
		Media:  res,
		Errors: errors,
	}, nil
}
