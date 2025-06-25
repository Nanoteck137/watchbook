package mal

import (
	"fmt"

	"github.com/nanoteck137/watchbook/types"
)

func ConvertAnimeType(typ string) types.MediaType {
	// switch typ {
	// case "TV":
	// 	return types.AnimeTypeTV
	// case "OVA":
	// 	return types.AnimeTypeOVA
	// case "Movie":
	// 	return types.AnimeTypeMovie
	// case "Special":
	// 	return types.AnimeTypeSpecial
	// case "ONA":
	// 	return types.AnimeTypeONA
	// case "Music":
	// 	return types.AnimeTypeMusic
	// case "CM":
	// 	return types.AnimeTypeCM
	// case "PV":
	// 	return types.AnimeTypePV
	// case "TV Special":
	// 	return types.AnimeTypeTVSpecial
	// case "":
	// default:
	// 	// TODO(patrik): Better logging
	// 	fmt.Printf("WARN: Unknown anime type \"%s\"\n", typ)
	// }

	return types.MediaTypeUnknown
}

func ConvertAnimeStatus(status string) types.MediaStatus {
	switch status {
	case "Currently Airing":
		return types.MediaStatusAiring
	case "Finished Airing":
		return types.MediaStatusFinished
	case "Not yet aired":
		return types.MediaStatusNotAired
	case "":
	default:
		// TODO(patrik): Better logging
		fmt.Printf("WARN: Unknown anime status \"%s\"\n", status)
	}

	return types.MediaStatusUnknown
}

func ConvertAnimeRating(rating string) types.MediaRating {
	switch rating {
	case "G - All Ages":
		return types.MediaRatingAllAges
	case "PG - Children":
		return types.MediaRatingPG
	case "PG-13 - Teens 13 or older":
		return types.MediaRatingPG13
	case "R - 17+ (violence & profanity)":
		return types.MediaRatingR17
	case "R+ - Mild Nudity":
		return types.MediaRatingRMildNudity
	case "Rx - Hentai":
		return types.MediaRatingRHentai
	case "":
	default:
		// TODO(patrik): Better logging
		fmt.Printf("WARN: Unknown anime rating \"%s\"\n", rating)
	}

	return types.MediaRatingUnknown
}
