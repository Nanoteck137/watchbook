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

type Images struct {
	Cover  string `toml:"cover"`
	Banner string `toml:"banner"`
	Logo   string `toml:"logo"`
}

type MediaGeneral struct {
	Title       string `toml:"title"`
	Description string `toml:"description"`

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
	Imdb        string `toml:"IMDb"`
	MyAnimeList string `toml:"myAnimeList"`
	Anilist     string `toml:"anilist"`
}

type MediaPart struct {
	Name string `toml:"name"`
}

type Media struct {
	Id        string          `toml:"id"`
	MediaType types.MediaType `toml:"mediaType"`

	Ids     MediaIds     `toml:"ids"`
	General MediaGeneral `toml:"general"`
	Images  Images       `toml:"images"`

	Parts []MediaPart `toml:"parts"`

	Path string `toml:"-"`
}

func (m Media) GetCoverPath() string {
	if m.Images.Cover == "" {
		return ""
	}

	return path.Join(m.Path, m.Images.Cover)
}

func (m Media) GetLogoPath() string {
	if m.Images.Logo == "" {
		return ""
	}

	return path.Join(m.Path, m.Images.Logo)
}

func (m Media) GetBannerPath() string {
	if m.Images.Banner == "" {
		return ""
	}

	return path.Join(m.Path, m.Images.Banner)
}

type CollectionGeneral struct {
	Name string `toml:"name"`
}

type CollectionEntry struct {
	Path       string `toml:"path"`
	SearchSlug string `toml:"searchSlug"`
	Order      int    `toml:"order"`
	SubOrder   int    `toml:"subOrder"`
}

type Collection struct {
	Id string `toml:"id"`

	General CollectionGeneral `toml:"general"`
	Images  Images            `toml:"images"`

	Entries []CollectionEntry `toml:"entries"`

	Path string `toml:"-"`
}

func (c Collection) GetCoverPath() string {
	if c.Images.Cover == "" {
		return ""
	}

	return path.Join(c.Path, c.Images.Cover)
}

func (c Collection) GetLogoPath() string {
	if c.Images.Logo == "" {
		return ""
	}

	return path.Join(c.Path, c.Images.Logo)
}

func (c Collection) GetBannerPath() string {
	if c.Images.Banner == "" {
		return ""
	}

	return path.Join(c.Path, c.Images.Banner)
}

type LibrarySearch struct {
	Media       []Media
	Collections []Collection
	Errors      map[string]error
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

	media.Path = p

	return media, nil
}

func readCollection(p string) (Collection, error) {
	metadataPath := path.Join(p, "collection.toml")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return Collection{}, err
	}

	var collection Collection
	err = toml.Unmarshal(data, &collection)
	if err != nil {
		return Collection{}, err
	}

	collection.Path = p

	return collection, nil
}

func SearchLibrary(p string) (*LibrarySearch, error) {
	var mediaPaths []string
	var collectionPaths []string

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
			mediaPaths = append(mediaPaths, path.Dir(p))
		}

		if name == "collection.toml" {
			collectionPaths = append(collectionPaths, path.Dir(p))
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	errors := map[string]error{}
	media := make([]Media, 0, len(mediaPaths))
	collections := make([]Collection, 0, len(collectionPaths))

	for _, p := range mediaPaths {
		m, err := readMedia(p)
		if err != nil {
			errors[p] = err
			continue
		}

		media = append(media, m)
	}

	for _, p := range collectionPaths {
		collection, err := readCollection(p)
		if err != nil {
			errors[p] = err
			continue
		}

		collections = append(collections, collection)
	}

	return &LibrarySearch{
		Media:       media,
		Errors:      errors,
		Collections: collections,
	}, nil
}
