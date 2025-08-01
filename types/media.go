package types

import (
	"errors"
	"strconv"
	"time"
)

const MediaDateLayout = "2006-01-02"

func GetAiringSeason(d string) string {
	t, err := time.Parse(MediaDateLayout, d)
	if err != nil {
		return ""
	}

	year := t.Year()

	switch t.Month() {
	case time.January, time.February, time.March:
		return "winter-"+strconv.Itoa(year)
	case time.April, time.May, time.June:
		return "spring-"+strconv.Itoa(year)
	case time.July, time.August, time.September:
		return "summer-"+strconv.Itoa(year)
	case time.October, time.November, time.December:
		return "winter-"+strconv.Itoa(year)
	}

	return ""
}

func IsReleased(d string) bool {
	t, err := time.Parse(MediaDateLayout, d)
	if err != nil {
		return false
	}

	newT := time.Now().Sub(t)
	return newT.Seconds() > 0
}

type MediaType string

const (
	MediaTypeUnknown     MediaType = "unknown"
	MediaTypeSeason      MediaType = "season"
	MediaTypeMovie       MediaType = "movie"
	MediaTypeAnimeSeason MediaType = "anime-season"
	MediaTypeAnimeMovie  MediaType = "anime-movie"
)

func (t MediaType) IsMovie() bool {
	return t == MediaTypeMovie || t == MediaTypeAnimeMovie
}

func IsValidMediaType(t MediaType) bool {
	switch t {
	case MediaTypeUnknown,
		MediaTypeSeason,
		MediaTypeMovie,
		MediaTypeAnimeSeason,
		MediaTypeAnimeMovie:
		return true
	}

	return false
}

func ValidateMediaType(val any) error {
	if s, ok := val.(string); ok {
		if s == "" {
			return nil
		}

		t := MediaType(s)
		if !IsValidMediaType(t) {
			return errors.New("invalid type")
		}
	} else if p, ok := val.(*string); ok {
		if p == nil {
			return nil
		}

		s := *p
		if s == "" {
			return nil
		}

		t := MediaType(s)
		if !IsValidMediaType(t) {
			return errors.New("invalid type")
		}
	} else {
		return errors.New("expected string")
	}

	return nil
}

// type MediaType string
//
// const (
// 	MediaTypeUnknown   MediaType = "unknown"
// 	MediaTypeTV        MediaType = "tv"
// 	MediaTypeOVA       MediaType = "original-video-anime"
// 	MediaTypeMovie     MediaType = "movie"
// 	MediaTypeSpecial   MediaType = "special"
// 	MediaTypeONA       MediaType = "original-network-anime"
// 	MediaTypeMusic     MediaType = "music"
// 	MediaTypeCM        MediaType = "commercial"
// 	MediaTypePV        MediaType = "promotional-video"
// 	MediaTypeTVSpecial MediaType = "tv-special"
// )

type MediaStatus string

const (
	MediaStatusUnknown  MediaStatus = "unknown"
	MediaStatusAiring   MediaStatus = "airing"
	MediaStatusFinished MediaStatus = "finished"
	MediaStatusNotAired MediaStatus = "not-aired"
)

func IsValidMediaStatus(s MediaStatus) bool {
	switch s {
	case MediaStatusUnknown,
		MediaStatusAiring,
		MediaStatusFinished,
		MediaStatusNotAired:
		return true
	}

	return false
}

func ValidateMediaStatus(val any) error {
	if s, ok := val.(string); ok {
		if s == "" {
			return nil
		}

		t := MediaStatus(s)
		if !IsValidMediaStatus(t) {
			return errors.New("invalid status")
		}
	} else if p, ok := val.(*string); ok {
		if p == nil {
			return nil
		}

		s := *p
		if s == "" {
			return nil
		}

		t := MediaStatus(s)
		if !IsValidMediaStatus(t) {
			return errors.New("invalid status")
		}
	} else {
		return errors.New("expected string")
	}

	return nil
}

type MediaRating string

const (
	MediaRatingUnknown     MediaRating = "unknown"
	MediaRatingAllAges     MediaRating = "all-ages"
	MediaRatingPG          MediaRating = "pg"
	MediaRatingPG13        MediaRating = "pg-13"
	MediaRatingR17         MediaRating = "r-17"
	MediaRatingRMildNudity MediaRating = "r-mild-nudity"
	MediaRatingRHentai     MediaRating = "r-hentai"
)

func IsValidMediaRating(r MediaRating) bool {
	switch r {
	case MediaRatingUnknown,
		MediaRatingAllAges,
		MediaRatingPG,
		MediaRatingPG13,
		MediaRatingR17,
		MediaRatingRMildNudity,
		MediaRatingRHentai:
		return true
	}

	return false
}

func ValidateMediaRating(val any) error {
	if s, ok := val.(string); ok {
		if s == "" {
			return nil
		}

		t := MediaRating(s)
		if !IsValidMediaRating(t) {
			return errors.New("invalid rating")
		}
	} else if p, ok := val.(*string); ok {
		if p == nil {
			return nil
		}

		s := *p
		if s == "" {
			return nil
		}

		t := MediaRating(s)
		if !IsValidMediaRating(t) {
			return errors.New("invalid rating")
		}
	} else {
		return errors.New("expected string")
	}

	return nil
}

type MediaUserList string

const (
	MediaUserListWatching    MediaUserList = "watching"
	MediaUserListCompleted   MediaUserList = "completed"
	MediaUserListOnHold      MediaUserList = "on-hold"
	MediaUserListDropped     MediaUserList = "dropped"
	MediaUserListPlanToWatch MediaUserList = "plan-to-watch"
)

func IsValidMediaUserList(l MediaUserList) bool {
	switch l {
	case MediaUserListWatching,
		MediaUserListCompleted,
		MediaUserListOnHold,
		MediaUserListDropped,
		MediaUserListPlanToWatch:
		return true
	}

	return false
}

type MediaImageType string

const (
	MediaImageTypeUnknown MediaImageType = "unknown"
	MediaImageTypeCover   MediaImageType = "cover"
	MediaImageTypeBanner  MediaImageType = "banner"
	MediaImageTypeLogo    MediaImageType = "logo"
)

func IsValidMediaImageType(l MediaImageType) bool {
	switch l {
	case MediaImageTypeUnknown,
		MediaImageTypeCover,
		MediaImageTypeBanner,
		MediaImageTypeLogo:
		return true
	}

	return false
}

func ValidateMediaImageType(val any) error {
	if s, ok := val.(string); ok {
		if s == "" {
			return nil
		}

		t := MediaImageType(s)
		if !IsValidMediaImageType(t) {
			return errors.New("invalid media image type")
		}
	} else if p, ok := val.(*string); ok {
		if p == nil {
			return nil
		}

		s := *p
		if s == "" {
			return nil
		}

		t := MediaImageType(s)
		if !IsValidMediaImageType(t) {
			return errors.New("invalid media image type")
		}
	} else {
		return errors.New("expected string")
	}

	return nil
}
