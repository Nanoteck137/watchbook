package mal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/nanoteck137/watchbook/downloader"
)

var ErrCheckFailed = errors.New("entry check failed")

const perPage = 100
const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:136.0) Gecko/20100101 Firefox/136.0"

var episodeRegex = regexp.MustCompile(`\((.+)\/(.+)\)`)

func readEpisodeCount(p string) (int, error) {
	f, err := os.Open(p)
	if err != nil {
		return 0, err
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return 0, err
	}

	episodeCount := doc.Find("tbody > tr > td > div > .di-ib").Text()
	captures := episodeRegex.FindStringSubmatch(episodeCount)
	if captures == nil {
		return 0, errors.New("Episode count not found")
	}

	currentRaw := captures[1]
	currentRaw = strings.ReplaceAll(currentRaw, ",", "")

	current, err := strconv.ParseInt(currentRaw, 10, 32)
	if err != nil {
		return 0, err
	}

	return int(current), nil
}

// TODO(patrik): Move this?
func checkValidEntry(dir string) error {
	err := filepath.Walk(dir, func(p string, info fs.FileInfo, err error) error {
		name := info.Name()
		ext := path.Ext(name)

		if strings.HasPrefix(name, ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		if info.IsDir() {
			return nil
		}

		if ext == ".html" {
			f, err := os.Open(p)
			if err != nil {
				return err
			}

			doc, err := goquery.NewDocumentFromReader(f)
			if err != nil {
				return err
			}

			title := doc.Find(".title-name")
			if len(title.Nodes) == 0 {
				return ErrCheckFailed
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

type EntryData struct {
	Id string `json:"id"`

	Anime       Anime    `json:"anime"`
	PictureUrls []string `json:"pictureUrls"`

	DownloadDate time.Time `json:"downloadDate"`
}

func process(id, p string) error {
	anime, err := ExtractAnimeData(path.Join(p, "root.html"))
	if err != nil {
		return fmt.Errorf("failed to extract anime data: %w", err)
	}

	pictures, err := ExtractPictures(path.Join(p, "pictures.html"))
	if err != nil {
		return fmt.Errorf("failed to extract anime pictures: %w", err)
	}

	data := EntryData{
		Id:           id,
		Anime:        anime,
		PictureUrls:  pictures,
		DownloadDate: time.Now(),
	}

	d, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal entry data: %w", err)
	}

	err = os.WriteFile(path.Join(p, "data.json"), d, 0644)
	if err != nil {
		return fmt.Errorf("failed to write entry data: %w", err)
	}

	return nil
}

// func DownloadEntry(dl *downloader.Downloader, workDir types.WorkDir, id string) (string, error) {
// 	p := path.Join(workDir.RawCurrentDir(), id)
//
// 	err := os.Mkdir(p, 0755)
// 	if err != nil {
// 		if os.IsExist(err) {
// 			// TODO(patrik): Make this a config
// 			newName := fmt.Sprintf("%s-%d", id, time.Now().UnixMilli())
// 			err = os.Rename(p, path.Join(workDir.RawOldDir(), newName))
// 			if err != nil {
// 				return "", fmt.Errorf("failed to rename old entry: %w", err)
// 			}
//
// 			err = os.Mkdir(p, 0755)
// 			if err != nil {
// 				return "", fmt.Errorf("failed to create entry dir after rename: %w", err)
// 			}
// 		} else {
// 			return "", fmt.Errorf("failed to create entry dir: %w", err)
// 		}
// 	}
//
// 	episodesDir := path.Join(p, "episodes")
// 	err = os.Mkdir(episodesDir, 0755)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to create entry episodes dir: %w", err)
// 	}
//
// 	baseUrl := fmt.Sprintf("https://myanimelist.net/anime/%s/random", id)
//
// 	err = dl.DownloadToFile(baseUrl, path.Join(p, "root.html"))
// 	if err != nil {
// 		return "", fmt.Errorf("failed to download entry root page: %w", err)
// 	}
//
// 	hasEpisodes := true
//
// 	firstEpisodePath := path.Join(episodesDir, "0.html")
// 	err = dl.DownloadToFile(baseUrl+"/episode", firstEpisodePath)
// 	if err != nil {
// 		if errors.Is(err, downloader.NotFound) {
// 			hasEpisodes = false
// 		} else {
// 			return "", fmt.Errorf("failed to download entry initial episode page: %w", err)
// 		}
// 	}
//
// 	if hasEpisodes {
// 		count, err := readEpisodeCount(firstEpisodePath)
// 		if err != nil {
// 			return "", fmt.Errorf("failed to read entry episode count: %w", err)
// 		}
//
// 		pageCount := utils.TotalPages(perPage, count)
//
// 		for i := 1; i < pageCount; i++ {
// 			name := fmt.Sprintf("%d.html", i)
//
// 			p := path.Join(episodesDir, name)
//
// 			err = dl.DownloadToFile(baseUrl+fmt.Sprintf("/episode?offset=%d", i*perPage), p)
// 			if err != nil && !errors.Is(err, downloader.NotFound) {
// 				return "", fmt.Errorf("failed to download entry episode page (%d): %w", i, err)
// 			}
// 		}
// 	}
//
// 	err = dl.DownloadToFile(baseUrl+"/pics", path.Join(p, "pictures.html"))
// 	if err != nil {
// 		return "", fmt.Errorf("failed to download entry pictures page: %w", err)
// 	}
//
// 	err = checkValidEntry(p)
// 	if err != nil {
// 		return "", fmt.Errorf("entry failed check: %w", err)
// 	}
//
// 	err = process(id, p)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to process entry: %w", err)
// 	}
//
// 	return p, nil
// }

func FetchAnimeData(dl *downloader.Downloader, id string, fetchPictures bool) (*Anime, error) {
	p, err := os.MkdirTemp("", "anime*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(p)

	err = os.Mkdir(p, 0755)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("failed to create entry dir: %w", err)
	}

	baseUrl := fmt.Sprintf("https://myanimelist.net/anime/%s/random", id)

	err = dl.DownloadToFile(baseUrl, path.Join(p, "root.html"))
	if err != nil {
		return nil, fmt.Errorf("failed to download entry root page: %w", err)
	}

	// TODO(patrik): Remove?
	err = checkValidEntry(p)
	if err != nil {
		return nil, fmt.Errorf("entry failed check: %w", err)
	}

	anime, err := ExtractAnimeData(path.Join(p, "root.html"))
	if err != nil {
		return nil, fmt.Errorf("failed to extract anime data: %w", err)
	}

	if fetchPictures {
		baseUrl := fmt.Sprintf("https://myanimelist.net/anime/%s/random", id)

		picturesDst := path.Join(p, "pictures.html")
		err = dl.DownloadToFile(baseUrl+"/pics", picturesDst)
		if err != nil {
			return nil, fmt.Errorf("failed to download entry pictures page: %w", err)
		}

		pictures, err := ExtractPictures(picturesDst)
		if err != nil {
			return nil, err
		}

		anime.Pictures = pictures
	}

	return &anime, nil
}
