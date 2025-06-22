package types

import "errors"

type AnimeDownloadType string

const (
	AnimeDownloadTypeNone    AnimeDownloadType = "none"
	AnimeDownloadTypeMal     AnimeDownloadType = "myanimelist"
	AnimeDownloadTypeAnilist AnimeDownloadType = "anilist"
)

type AnimeThemeSongType string

const (
	AnimeThemeSongTypeUnknown AnimeThemeSongType = "unknown"
	AnimeThemeSongTypeOpening AnimeThemeSongType = "opening"
	AnimeThemeSongTypeEnding  AnimeThemeSongType = "ending"
)

type AnimeType string

const (
	AnimeTypeUnknown     AnimeType = "unknown"
	AnimeTypeSeason      AnimeType = "season"
	AnimeTypeMovie       AnimeType = "movie"
	AnimeTypeAnimeSeason AnimeType = "anime-season"
	AnimeTypeAnimeMovie  AnimeType = "anime-movie"
)

func (t AnimeType) IsMovie() bool {
	return t == AnimeTypeMovie || t == AnimeTypeAnimeMovie
}

func IsValidAnimeType(t AnimeType) bool {
	switch t {
	case AnimeTypeUnknown,
		AnimeTypeSeason,
		AnimeTypeMovie,
		AnimeTypeAnimeSeason,
		AnimeTypeAnimeMovie:
		return true
	}

	return false
}

func ValidateAnimeType(val any) error {
	if s, ok := val.(string); ok {
		if s == "" {
			return nil
		}

		t := AnimeType(s)
		if !IsValidAnimeType(t) {
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

		t := AnimeType(s)
		if !IsValidAnimeType(t) {
			return errors.New("invalid type")
		}
	} else {
		return errors.New("expected string")
	}

	return nil
}

// type AnimeType string
//
// const (
// 	AnimeTypeUnknown   AnimeType = "unknown"
// 	AnimeTypeTV        AnimeType = "tv"
// 	AnimeTypeOVA       AnimeType = "original-video-anime"
// 	AnimeTypeMovie     AnimeType = "movie"
// 	AnimeTypeSpecial   AnimeType = "special"
// 	AnimeTypeONA       AnimeType = "original-network-anime"
// 	AnimeTypeMusic     AnimeType = "music"
// 	AnimeTypeCM        AnimeType = "commercial"
// 	AnimeTypePV        AnimeType = "promotional-video"
// 	AnimeTypeTVSpecial AnimeType = "tv-special"
// )

type AnimeStatus string

const (
	AnimeStatusUnknown  AnimeStatus = "unknown"
	AnimeStatusAiring   AnimeStatus = "airing"
	AnimeStatusFinished AnimeStatus = "finished"
	AnimeStatusNotAired AnimeStatus = "not-aired"
)

func IsValidAnimeStatus(s AnimeStatus) bool {
	switch s {
	case AnimeStatusUnknown,
		AnimeStatusAiring,
		AnimeStatusFinished,
		AnimeStatusNotAired:
		return true
	}

	return false
}

func ValidateAnimeStatus(val any) error {
	if s, ok := val.(string); ok {
		if s == "" {
			return nil
		}

		t := AnimeStatus(s)
		if !IsValidAnimeStatus(t) {
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

		t := AnimeStatus(s)
		if !IsValidAnimeStatus(t) {
			return errors.New("invalid status")
		}
	} else {
		return errors.New("expected string")
	}

	return nil
}

type AnimeRating string

const (
	AnimeRatingUnknown     AnimeRating = "unknown"
	AnimeRatingAllAges     AnimeRating = "all-ages"
	AnimeRatingPG          AnimeRating = "pg"
	AnimeRatingPG13        AnimeRating = "pg-13"
	AnimeRatingR17         AnimeRating = "r-17"
	AnimeRatingRMildNudity AnimeRating = "r-mild-nudity"
	AnimeRatingRHentai     AnimeRating = "r-hentai"
)

func IsValidAnimeRating(r AnimeRating) bool {
	switch r {
	case AnimeRatingUnknown,
		AnimeRatingAllAges,
		AnimeRatingPG,
		AnimeRatingPG13,
		AnimeRatingR17,
		AnimeRatingRMildNudity,
		AnimeRatingRHentai:
		return true
	}

	return false
}

func ValidateAnimeRating(val any) error {
	if s, ok := val.(string); ok {
		if s == "" {
			return nil
		}

		t := AnimeRating(s)
		if !IsValidAnimeRating(t) {
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

		t := AnimeRating(s)
		if !IsValidAnimeRating(t) {
			return errors.New("invalid rating")
		}
	} else {
		return errors.New("expected string")
	}

	return nil
}

type AnimeUserList string

const (
	AnimeUserListWatching    AnimeUserList = "watching"
	AnimeUserListCompleted   AnimeUserList = "completed"
	AnimeUserListOnHold      AnimeUserList = "on-hold"
	AnimeUserListDropped     AnimeUserList = "dropped"
	AnimeUserListPlanToWatch AnimeUserList = "plan-to-watch"
)

func IsValidAnimeUserList(l AnimeUserList) bool {
	switch l {
	case AnimeUserListWatching,
		AnimeUserListCompleted,
		AnimeUserListOnHold,
		AnimeUserListDropped,
		AnimeUserListPlanToWatch:
		return true
	}

	return false
}
