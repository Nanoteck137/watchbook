package mal

import (
	"fmt"

	"github.com/nanoteck137/watchbook/types"
)

func ConvertAnimeType(typ string) types.AnimeType {
	switch typ {
	case "TV":
		return types.AnimeTypeTV
	case "OVA":
		return types.AnimeTypeOVA
	case "Movie":
		return types.AnimeTypeMovie
	case "Special":
		return types.AnimeTypeSpecial
	case "ONA":
		return types.AnimeTypeONA
	case "Music":
		return types.AnimeTypeMusic
	case "CM":
		return types.AnimeTypeCM
	case "PV":
		return types.AnimeTypePV
	case "TV Special":
		return types.AnimeTypeTVSpecial
	case "":
	default:
		// TODO(patrik): Better logging
		fmt.Printf("WARN: Unknown anime type \"%s\"\n", typ)
	}

	return types.AnimeTypeUnknown
}

func ConvertAnimeStatus(status string) types.AnimeStatus {
	switch status {
	case "Currently Airing":
		return types.AnimeStatusAiring
	case "Finished Airing":
		return types.AnimeStatusFinished
	case "Not yet aired":
		return types.AnimeStatusNotAired
	case "":
	default:
		// TODO(patrik): Better logging
		fmt.Printf("WARN: Unknown anime status \"%s\"\n", status)
	}

	return types.AnimeStatusUnknown
}

func ConvertAnimeRating(rating string) types.AnimeRating {
	switch rating {
	case "G - All Ages":
		return types.AnimeRatingAllAges
	case "PG - Children":
		return types.AnimeRatingPG
	case "PG-13 - Teens 13 or older":
		return types.AnimeRatingPG13
	case "R - 17+ (violence & profanity)":
		return types.AnimeRatingR17
	case "R+ - Mild Nudity":
		return types.AnimeRatingRMildNudity
	case "Rx - Hentai":
		return types.AnimeRatingRHentai
	case "":
	default:
		// TODO(patrik): Better logging
		fmt.Printf("WARN: Unknown anime rating \"%s\"\n", rating)
	}

	return types.AnimeRatingUnknown
}

func ConvertThemeSongType(typ ThemeSongType) types.ThemeSongType {
	switch(typ) {
	case ThemeSongOpening:
		return types.ThemeSongTypeOpening
	case ThemeSongEnding:
		return types.ThemeSongTypeEnding
	default:
		// TODO(patrik): Better logging
		fmt.Printf("WARN: Unknown theme song type \"%s\"\n", typ)
	}

	return types.ThemeSongTypeUnknown
}
