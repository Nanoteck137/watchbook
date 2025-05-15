package types

type ThemeSongType string

const (
	ThemeSongTypeUnknown ThemeSongType = "unknown"
	ThemeSongTypeOpening ThemeSongType = "opening"
	ThemeSongTypeEnding  ThemeSongType = "ending"
)

type AnimeType string

const (
	AnimeTypeUnknown   AnimeType = "unknown"
	AnimeTypeTV        AnimeType = "tv"
	AnimeTypeOVA       AnimeType = "original-video-anime"
	AnimeTypeMovie     AnimeType = "movie"
	AnimeTypeSpecial   AnimeType = "special"
	AnimeTypeONA       AnimeType = "original-network-anime"
	AnimeTypeMusic     AnimeType = "music"
	AnimeTypeCM        AnimeType = "commercial"
	AnimeTypePV        AnimeType = "promotional-video"
	AnimeTypeTVSpecial AnimeType = "tv-special"
)

type AnimeStatus string

const (
	AnimeStatusUnknown  AnimeStatus = "unknown"
	AnimeStatusAiring   AnimeStatus = "airing"
	AnimeStatusFinished AnimeStatus = "finished"
	AnimeStatusNotAired AnimeStatus = "not-aired"
)

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

type AnimeUserList string

const (
	AnimeUserListWatching    AnimeUserList = "watching"
	AnimeUserListCompleted   AnimeUserList = "completed"
	AnimeUserListOnHold      AnimeUserList = "on-hold"
	AnimeUserListDropped     AnimeUserList = "dropped"
	AnimeUserListPlanToWatch AnimeUserList = "plan-to-watch"
)
